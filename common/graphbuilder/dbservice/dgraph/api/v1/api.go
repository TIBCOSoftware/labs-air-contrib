/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/dgraph-io/dgraph/x"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	apiInterface "github.com/SteveNY-Tibco/labs-air-contrib/common/graphbuilder/dbservice/dgraph/api"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("dgraph-service")

func NewAPI(
	url string,
	user string,
	password string,
	tlsEnabled bool,
	tlsUserCfg map[string]interface{},
) (apiInterface.API, error) {
	api := &API{
		_url:        url,
		_user:       user,
		_password:   password,
		_tlsEnabled: tlsEnabled,
		_tlsUserCfg: tlsUserCfg,
	}
	return api, nil
}

type API struct {
	_url        string
	_user       string
	_password   string
	_tlsEnabled bool
	_tlsUserCfg map[string]interface{}

	_dgConnection *grpc.ClientConn
	_dgraphClient *dgo.Dgraph

	_mux sync.Mutex
}

func (this *API) EnsureConnection() error {
	var connErr error

	if nil == this._dgConnection {
		this._mux.Lock()
		defer this._mux.Unlock()
		if nil == this._dgConnection {
			if nil != this._dgConnection {
				this._dgConnection.Close()
			}

			fmt.Println("[DgraphService::ensureConnection] Will try to connect ..........")
			fmt.Println("[DgraphService::ensureConnection] url = " + this._url)
			fmt.Println("[DgraphService::ensureConnection] user = " + this._user)
			fmt.Println("[DgraphService::ensureConnection] password = " + this._password)
			fmt.Println("[DgraphService::ensureConnection] tlsUserCfg = ", this._tlsUserCfg)

			if !this._tlsEnabled {
				this._dgConnection, connErr = grpc.Dial(this._url, grpc.WithInsecure())
			} else {
				helperConfig := &x.TLSHelperConfig{}
				if nil != this._tlsUserCfg["tlsCertDir"] {
					helperConfig.CertDir = this._tlsUserCfg["tlsCertDir"].(string)
				}
				if nil != this._tlsUserCfg["tlsCertRequired"] {
					helperConfig.CertRequired = this._tlsUserCfg["tlsCertRequired"].(bool)
				}
				if nil != this._tlsUserCfg["tlsCert"] {
					helperConfig.Cert = this._tlsUserCfg["tlsCert"].(string)
				}
				if nil != this._tlsUserCfg["tlsKey"] {
					helperConfig.Key = this._tlsUserCfg["tlsKey"].(string)
				}

				if nil != this._tlsUserCfg["tlsServerName"] {
					helperConfig.ServerName = this._tlsUserCfg["tlsServerName"].(string)
				}
				if nil != this._tlsUserCfg["tlsRootCACert"] {
					helperConfig.RootCACert = this._tlsUserCfg["tlsRootCACert"].(string)
				}
				if nil != this._tlsUserCfg["tlsClientAuth"] {
					helperConfig.ClientAuth = this._tlsUserCfg["tlsClientAuth"].(string)
				}
				if nil != this._tlsUserCfg["tlsUseSystemCACerts"] {
					helperConfig.UseSystemCACerts = this._tlsUserCfg["tlsUseSystemCACerts"].(bool)
				}

				tlsCfg, err := x.GenerateClientTLSConfig(helperConfig)
				if nil != err {
					log.Error("[DgraphService::ensureConnection] Unable to configure TLS connection !!! Will not connect ......")
					this._dgConnection = nil
					return connErr
				}

				this._dgConnection, connErr = grpc.Dial(this._url, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
			}

			if nil != connErr {
				log.Error("[DgraphService::ensureConnection] Unable to create connection !!! Will not connect ......")
				this._dgConnection = nil
				return connErr
			}

			this._dgraphClient = dgo.NewDgraphClient(api.NewDgraphClient(this._dgConnection))
		}
	}

	return nil
}

func (this *API) NewTransaction() apiInterface.Transaction {
	txn := &Transaction{
		_api: this,
		_txn: this._dgraphClient.NewTxn(),
	}
	return txn
}

func (this *API) BuildSchema(schema string) error {
	fmt.Println("***************** schema query ********************")
	fmt.Println(schema)
	fmt.Println("***************************************************")

	err := this._dgraphClient.Alter(
		context.Background(),
		&api.Operation{Schema: schema},
	)

	return err
}

func (this *API) DropGraph() {
	this._dgraphClient.Alter(context.Background(), &api.Operation{DropAll: true})
}

func (this *API) Destroy() {
	this._dgConnection.Close()
}

type Transaction struct {
	_api *API
	_txn *dgo.Txn
}

func (this *Transaction) QueryWithVars(query string, vars map[string]string) (map[string]interface{}, error) {
	this._api.EnsureConnection()

	data := make(map[string]interface{})
	res, err := this._txn.QueryWithVars(context.Background(), query, vars)
	if nil != err {
		return data, err
	}

	err = json.Unmarshal(res.GetJson(), &data)
	if nil != err {
		return nil, err
	}

	return data, nil
}

func (this *Transaction) Mutate(nquads string) error {
	_, err := this._txn.Mutate(
		context.Background(),
		&api.Mutation{
			SetNquads: []byte(nquads),
		},
	)
	return err
}

func (this *Transaction) Commit() error {
	return this._txn.Commit(context.Background())
}

func (this *Transaction) Discard() {
	this._txn.Discard(context.Background())
}

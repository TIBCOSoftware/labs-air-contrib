/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package factory

import (
	"sync"

	"github.com/SteveNY-Tibco/labs-air-contrib/common/graphbuilder/dbservice"
	"github.com/SteveNY-Tibco/labs-air-contrib/common/graphbuilder/dbservice/dgraph/services"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("dgraph-service")

type DgraphServiceFactory struct {
	dbservice.BaseDBServiceFactory
	mux sync.Mutex
}

func (this *DgraphServiceFactory) CreateImportService(serviceId string, properties map[string]interface{}) (dbservice.ImportService, error) {
	this.mux.Lock()
	defer this.mux.Unlock()

	dgraphService := this.DBServices[serviceId]
	var err error
	if nil == dgraphService {
		dgraphService, err = services.NewDgraphImportRDF(properties)
		log.Info("(DgraphServiceFactory.CreateImportService) imortService : ", dgraphService)
		if nil != err {
			log.Info("(DgraphServiceFactory.CreateImportService) err : ", err)
			return nil, err
		}
		this.DBServices[serviceId] = dgraphService.(*services.DgraphImportRDF)
	}
	return dgraphService.(dbservice.ImportService), nil
}

func (this *DgraphServiceFactory) CreateUpsertService(serviceId string, properties map[string]interface{}) (dbservice.UpsertService, error) {
	this.mux.Lock()
	defer this.mux.Unlock()

	dgraphService := this.DBServices[serviceId]
	var err error
	if nil == dgraphService {
		dgraphService, err = services.NewDgraphService(properties)
		log.Info("(DgraphServiceFactory.CreateUpsertService) upsertService : ", dgraphService)
		if nil != err {
			log.Info("(DgraphServiceFactory.CreateUpsertService) err : ", err)
			return nil, err
		}
		this.DBServices[serviceId] = dgraphService.(*services.DgraphService)
	}
	return dgraphService.(dbservice.UpsertService), nil
}

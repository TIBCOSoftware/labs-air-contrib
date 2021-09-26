/*
 * Copyright © 2020. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
package api

import (
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

var log = logger.GetLogger("dgraph-service")

type API interface {
	EnsureConnection() error
	NewTransaction() Transaction
	BuildSchema(schema string) error
	DropGraph()
	Destroy()
}

type Transaction interface {
	QueryWithVars(query string, vars map[string]string) (map[string]interface{}, error)
	Mutate(nquads string) error
	Commit() error
	Discard()
}

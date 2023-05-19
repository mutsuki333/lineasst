/*
	init.go
	Purpose: The Entry point to the main program,
	it is also responsible for setting up default configs
	and define basic informations for the program.

	@author Evan Chen
	@version 1.0 2023/02/22
*/

package main

import (
	"app/core/config"
	"app/core/property"
)

var Version = "v0.0.0"
var ID = "dev"
var Build = "+%Y-%m-%d_%H:%M:%S"
var BinName = "lineasst"

const Usage = "Personal Line Assistant Server"

func init() {
	config.SetEnvPrefix("APP_")

	config.SetDefault(property.DB, "sqlite")
	config.SetDefault(property.DSN, "app.db?cache=shared&_fk=1")

	config.SetDefault(property.PORT, "80")
	config.SetDefault(property.ADDR, "0.0.0.0")
	config.SetDefault(property.GRPC_PORT, "8080")
	config.SetDefault(property.GRPC_ADDR, "localhost")

	// config.SetDefault(property.CUSTOM, "custom")

	config.SetDefault(property.LOG_LEVEL, "info")
	config.SetDefault(property.LOG_FORMAT, "text")
	config.SetDefault(property.LOG_STD, true)
	config.SetDefault(property.LOG_FILE, false)
	config.SetDefault(property.LOG_FILE_SIZE, 10)
	config.SetDefault(property.LOG_FILE_AGE, 365)
}

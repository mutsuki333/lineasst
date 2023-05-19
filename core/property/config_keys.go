/*
	config_keys.go
	Purpose: A place to define common config keys
	for modules to reference throughout the applicaion.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/02/22  v1.0.0 Evan Chen   Initial release

*/

package property

import "app/core/config"

//-------------------------------------------------
//- database related configs                      -
//-------------------------------------------------

const (
	DB  config.Key = "DB"  // config key for db types. ex: db2, sqlite.
	DSN config.Key = "DSN" // config key for db connection strings.
)

//-------------------------------------------------
//- Server related configs                        -
//-------------------------------------------------

const (
	PORT config.Key = "PORT" // config key for which port the http server listens on.
	ADDR config.Key = "ADDR" // config key for which ip address to bind on.

	GRPC_PORT config.Key = "GRPC_PORT" // config key for which port the grpc server listens on.
	GRPC_ADDR config.Key = "GRPC_ADDR" // config key for which ip address for grpc server to bind on.
)

//-------------------------------------------------
//- System related configs                        -
//-------------------------------------------------

const (
	DEBUG      config.Key = "DEBUG"      // config key to set run mode to debug
	CUSTOM     config.Key = "CUST"       // config key to set where the custom folder is located
	SECRET     config.Key = "SECRET"     // config key to set the key to sign jwt tokens
	NO_CRON    config.Key = "NO_CRON"    // config key to disable cron jobs
	AUTO_LOGIN config.Key = "AUTO_LOGIN" // config key to set the authentication mechanism to always identify requests as given user
)

//-------------------------------------------------
//- Logging related configs                       -
//-------------------------------------------------

const (
	LOG_LEVEL     config.Key = "LOG_LEVEL"     // config key to set the log level
	LOG_FORMAT    config.Key = "LOG_FORMAT"    // config key to set the logging format
	LOG_STD       config.Key = "LOG_STD"       // config key to set the logging output to stdout
	LOG_FILE      config.Key = "LOG_FILE"      // config key to set the logging output to file
	LOG_FILE_SIZE config.Key = "LOG_FILE_SIZE" // config key to set the logging output file rotating max size (mb)
	LOG_FILE_AGE  config.Key = "LOG_FILE_AGE"  // config key to set the logging output file max age to keep
	LOG_SQL       config.Key = "LOG_SQL"       // config key to set to log all sql outputs
)

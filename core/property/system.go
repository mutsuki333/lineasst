package property

import (
	"golang.org/x/exp/slog"
)

//-------------------------------------------------
//- Application run mode                          -
//-------------------------------------------------

type Mode int

const (
	DEBUG_MODE Mode = iota
	PROD_MODE
)

// RUN_MODE stores the current runmode
var RUN_MODE Mode = PROD_MODE

func IsDebug() bool {
	return RUN_MODE == DEBUG_MODE
}

//-------------------------------------------------
//- Application state                             -
//-------------------------------------------------

type State int

const (
	STATE_NULL     State = iota
	STATE_CONN           // connecting to database
	STATE_DBM            // migrating database
	STATE_PREPARE        // prepare rpc/restful services
	STATE_INIT           // lifecycle init
	STATE_LOAD           // lifecycle load
	STATE_PRESTART       // before app start
	STATE_STARTED        // application started
	STATE_TERM           // application in termination state
)

var APP_STATE = STATE_NULL

func SetState(state State) {
	APP_STATE = state
}

//-------------------------------------------------
//- Application logging                           -
//-------------------------------------------------

var LogLevel = new(slog.LevelVar)

func SetLogLevel(level string) {
	switch level {
	case "debug":
		LogLevel.Set(slog.LevelDebug)
	case "warn":
		LogLevel.Set(slog.LevelWarn)
	case "error":
		LogLevel.Set(slog.LevelError)
	default:
		LogLevel.Set(slog.LevelInfo)
	}
}

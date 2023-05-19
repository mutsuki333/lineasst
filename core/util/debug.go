/*
	debug.go
	Purpose: debug utiliy functions.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/02  v1.0.0 Evan Chen   Initial release
*/

package util

import (
	"fmt"
	"runtime"

	"golang.org/x/exp/slog"
)

// Stack gets the current call stack for debug
func Stack() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}

// Caller returns the caller file
func Caller(depth int) string {
	_, file, line, _ := runtime.Caller(depth)
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return fmt.Sprintf("%s:%d", short, line)
}

type SlogIface interface {
	Attr() slog.Attr
}

func ErrAtrr(err error) slog.Attr {
	if e, ok := err.(SlogIface); ok {
		return e.Attr()
	}
	return slog.String("err", err.Error())
}

func RpcAtrr(methodName string) slog.Attr {
	return slog.String("rpc", methodName)
}

var (
	C_attr = slog.String("act", "create")   // act: create
	R_attr = slog.String("act", "retrieve") // act: retrieve
	U_attr = slog.String("act", "update")   // act: update
	D_attr = slog.String("act", "delete")   // act: delete
)

/*
	grpc.go
	Purpose: grpc utiliy functions.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/02  v1.0.0 Evan Chen   Initial release
*/

package server

import (
	"net/http"
	"strings"

	"app/core/errors"
	"app/core/util"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// GetHttpLocale gets the locale from http header
func GetHttpLocale(r *http.Request) string {
	return strings.ToLower(r.Header.Get("accept-language"))
}

// GetHttpAuthToken gets the auth token from the http header
func GetHttpAuthToken(r *http.Request) string {
	return r.Header.Get("authorization")
}

func HttpAbort(w http.ResponseWriter, r *http.Request, err *errors.Error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(runtime.HTTPStatusFromCode(err.Status))
	util.PbMarshaler.Marshal(w, err.Exec(GetHttpAuthToken(r), nil).GRPCStatus().Proto())
}

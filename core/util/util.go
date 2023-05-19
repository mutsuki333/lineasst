/*
	util.go
	Purpose: Utiliy functions.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/02  v1.0.0 Evan Chen   Initial release
*/

package util

import "github.com/golang/protobuf/jsonpb"

// Ref returns a pointer to src
func Ref[T any](src T) *T {
	return &src
}

// DeRef retruns the value of src, if src is a nil pointer, returns the default value
func DeRef[T any](src *T) T {
	if src == nil {
		var s T
		return s
	}
	return *src
}

var PbMarshaler = &jsonpb.Marshaler{
	OrigName: true,
}

var PbUnmarshaler = &jsonpb.Unmarshaler{
	AllowUnknownFields: true,
}

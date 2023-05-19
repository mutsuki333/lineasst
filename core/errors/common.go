/*
	common.go
	Purpose: Define some common errors.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/06  v1.0.0 Evan Chen   Initial release

*/

package errors

import "google.golang.org/grpc/codes"

// ErrBadRequest 400
var ErrBadRequest = New("ECMN-03-0", codes.InvalidArgument, nil)

// ErrUnauthorized 401
var ErrUnauthorized = New("ECMN-16-0", codes.Unauthenticated, nil)

// ErrForbidden 403
var ErrForbidden = New("ECMN-07-0", codes.PermissionDenied, nil)

// ErrNotFound 404
var ErrNotFound = New("ECMN-05-0", codes.NotFound, nil)

// ErrConflict 409
var ErrConflict = New("ECMN-06-0", codes.AlreadyExists, nil)

// ErrInternal 500
var ErrInternal = New("ECMN-02-0", codes.Internal, nil)

// ErrServiceUnavailable 503
var ErrServiceUnavailable = New("ECMN-14-0", codes.Unavailable, nil)

var ErrResourceInUse = New("ECMN-09-0", codes.FailedPrecondition, nil)

func ECodeFromGCode(c codes.Code) string {
	switch c {
	case codes.Unknown, codes.DataLoss:
		return ErrInternal.Code
	case codes.InvalidArgument, codes.OutOfRange:
		return ErrBadRequest.Code
	case codes.FailedPrecondition:
		return ErrResourceInUse.Code
	case codes.NotFound:
		return ErrNotFound.Code
	case codes.PermissionDenied:
		return ErrForbidden.Code
	case codes.Unauthenticated:
		return ErrUnauthorized.Code
	case codes.AlreadyExists, codes.Aborted:
		return ErrConflict.Code
	case codes.Unavailable:
		return ErrServiceUnavailable.Code

	default:
		return ErrInternal.Code
	}
}

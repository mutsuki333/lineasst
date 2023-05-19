/*
	errors.go
	Purpose: Define Error structs.

	@author Evan Chen

	MODIFICATION HISTORY
	   Date        Ver    Name     Description
	---------- ------- ----------- -------------------------------------------
	2023/03/06  v1.0.0 Evan Chen   Initial release

*/

package errors

import (
	"errors"
	"strings"

	"app/core/msg"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Error struct {
	// Code is the error code,
	Code string `json:"code,omitempty"`
	// Msg is the main message of this error, it is set by calling `Error()` or `Exec()`
	Msg string `json:"msg,omitempty"`
	// GrpcCode is the grpc status code for services to return
	Status codes.Code `json:"status"`

	// data is the payload that will pass to `Exec()` to construct the Msg
	data any `json:"-"`
}

// Error implements the error interface.
// If the Error has not been [Exec], it will execute with the default locale and the current data,
// which is most likely to be null.
func (e *Error) Error() string {
	if e.Msg != "" {
		return e.Msg
	}
	return e.Exec("", e.data).Error()
}

// String implements the Stringer interface.
func (e *Error) String() string {
	return e.Error()
}

// Exec copies the parent Error and executes the message with the given @locale and @Data
func (e *Error) Exec(locale string, data any) *Error {
	if data == nil {
		data = e.data
	}
	return &Error{
		Code: e.Code, Status: e.Status, // copy
		data: data,
		Msg:  e.Code + ": " + msg.T(e.Code, locale, data), // execute new message
	}
}

// SetData copies the Error and sets the new data to @data,
// passing the pointer to the new Error
func (e *Error) SetData(data any) *Error {
	return &Error{
		Code: e.Code, Status: e.Status, data: data,
	}
}

// SetInfo is a short hand to set data for simple messages that has a {{.Info}} in the template.
// It is the same as callint SetData(msg.Plain(str))
func (e *Error) SetInfo(str string) *Error {
	return e.SetData(msg.Plain(str))
}

// SetData sets the current Code to @Code,
// passing the original pointer back for easier method chaining
func (e *Error) SetCode(Code string) *Error {
	e.Code = Code
	return e
}

// SetData sets the current Status to @status,
// passing the original pointer back for easier method chaining
func (e *Error) SetStatus(status codes.Code) *Error {
	e.Status = status
	return e
}

func (e *Error) Copy() *Error {
	return &Error{
		Code: e.Code, Status: e.Status, data: e.data,
	}
}

// GRPCStatus implements the interface for custom errors
// to convert to standard grpc errors.
func (e *Error) GRPCStatus() *status.Status {
	err := status.New(e.Status, e.Error())
	err, _ = err.WithDetails(&ErrorDetial{
		Code:    e.Code,
		Message: strings.TrimPrefix(e.Error(), e.Code+": "),
		Emited:  timestamppb.Now(),
	})
	return err
}

func (e *Error) Attr() slog.Attr {
	return slog.String("err", e.Error())
}

// New creates a new Error struct with the given @Code and @MsgKey and an optional grpc failure code.
//
// While the system uses grpc as primary protocol, it also supports restful proxying,
// where the http status codes will map according to
// https://github.com/grpc-ecosystem/grpc-gateway/blob/main/runtime/errors.go#L36
func New(Code string, status codes.Code, data any) *Error {
	e := &Error{
		Code:   Code,
		Status: status,
		data:   data,
	}
	return e
}

// Is is a wrapper of [errors/Is] from the stdlib
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// As is a wrapper of [errors/As] from the stdlib
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Convert converts err to *Error of this package.
//
//   - If @err is already *Error, return as-is
func Convert(err error) *Error {
	if err == nil {
		return nil
	}

	e, ok := err.(*Error)
	if ok {
		return e
	}
	stat := status.Convert(err)

	return New(ECodeFromGCode(stat.Code()), stat.Code(), msg.Plain(err.Error()))
}

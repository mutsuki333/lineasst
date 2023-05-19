/*
	middleware.go
	Purpose: Middleware such as logging, panic recovery, etc. for gin or grpc.

	@author Evan Chen
*/

package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"app/core/auth"
	"app/core/errors"
	"app/core/msg"
	"app/core/property"
	"app/core/util"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const grpc_method = "GRPC"

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	start := time.Now()
	// content := fmt.Sprintf("%s %s", grpc_method, info.FullMethod)

	defer func() {
		pan := recover()
		args := []any{
			slog.String("method", grpc_method),
			slog.String("ip", GetGrpcClientIP(ctx)),
			slog.Duration("duration", time.Since(start)),
			slog.String("rpc", info.FullMethod),
		}
		if usr, ok := auth.GetUser(ctx); ok {
			args = append(args, slog.String("usr", usr.Username))
		}

		if property.IsDebug() {
			args = append(args, slog.Any("request", req))
		}
		if pan != nil {
			if err == nil {
				err = errors.ErrInternal.Exec("", msg.Plain(fmt.Sprint(pan)))
			}
			args = append(args,
				slog.Int("status", int(codes.Internal)),
				slog.String("err", err.Error()),
			)
			if property.IsDebug() {
				args = append(args,
					slog.String("caller", util.Caller(2)),
					slog.String("stack", util.Stack()),
				)
			}
			err = errors.ErrInternal.Exec("", msg.Plain(fmt.Sprint(pan)))
			log.Error("panic", args...)
		} else {
			if err != nil {
				if ce, ok := err.(*errors.Error); ok {
					args = append(args,
						slog.Int("status", int(ce.Status)),
						slog.String("code", ce.Code),
						ce.Attr(),
					)
				} else {
					args = append(args, slog.String("err", err.Error()))
				}
				log.Error(err.Error(), args...)
			} else { // success
				args = append(args, slog.Int("status", int(codes.OK)))
				log.Info("ok", args...)
			}
		}
	}()

	switch property.APP_STATE {
	case property.STATE_LOAD, property.STATE_TERM:
		err = errors.ErrServiceUnavailable
		return nil, err
	}

	//-------------------------------------------------
	//- Authorization interceptor                     -
	//-------------------------------------------------
	ctx = auth.SetUser(ctx, GetGrpcAuthToken(ctx))
	if err = auth.Authenticate(ctx, info.FullMethod); err != nil {
		return nil, err
	}

	ctx = OptimizeGrpcLocale(ctx)

	resp, err = handler(ctx, req)
	if err != nil {
		err = errors.Convert(err).Exec(GetGrpcLocale(ctx), nil)
	}

	return resp, err
}

const HttpHeaderPrefix = "http-"

func OutGoingHeaderMatcher(key string) (string, bool) {
	if strings.HasPrefix(key, HttpHeaderPrefix) {
		return key[len(HttpHeaderPrefix):], true
	}
	return fmt.Sprintf("%s%s", runtime.MetadataHeaderPrefix, key), true
}

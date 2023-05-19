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
	"context"
	"strings"

	"app/core/property"
	"app/core/util"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// GetGrpcClientIP gets the client ip from metadata,
// it will try to resolve proxy ips if presented.
func GetGrpcClientIP(ctx context.Context) string {

	p, _ := peer.FromContext(ctx)
	peerIP := util.StripAfterChar(p.Addr.String(), ':')

	// first look for proxy
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return peerIP
	}
	var results []string
	if results = md.Get("x-forwarded-for"); len(results) == 0 {
		return peerIP
	}
	return results[0]
}

// GetGrpcLocale gets the locale from metadata,
// it will use grpcgateway values if presented.
func GetGrpcLocale(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	if r := md.Get("x-accept-language"); len(r) > 0 {
		return r[0]
	}
	var results []string
	if results = md.Get("accept-language"); len(results) == 0 {
		// Try grpc-gateway
		if results = md.Get("grpcgateway-accept-language"); len(results) == 0 {
			return ""
		}
		ip, _, _ := strings.Cut(results[0], ";")
		return strings.ToLower(ip)
	}
	ip, _, _ := strings.Cut(results[0], ";")
	return strings.ToLower(ip)

}

// GetGrpcAuthToken gets the auth token from the metadata.
func GetGrpcAuthToken(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	var results []string
	if results = md.Get("authorization"); len(results) == 0 {
		// Try grpc-gateway
		if results = md.Get("grpcgateway-authorization"); len(results) == 0 {
			return ""
		}
		return results[0]
	}
	return results[0]
}

// OptimizeGrpcLocale speeds up the GetGrpcLocale if called multiple times.
//
// Getting locale from metadata is expected to be called very often, so we extract it in advance.
func OptimizeGrpcLocale(ctx context.Context) context.Context {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{"x-accept-language": property.DefaultLocale})
	}

	md.Set("x-accept-language", GetGrpcLocale(ctx))

	return metadata.NewIncomingContext(ctx, md)
}

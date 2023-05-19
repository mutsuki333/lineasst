/*
	server.go
	Purpose: Operaions for .

	@author Evan Chen
	@version 1.0 2023/02/22
*/

package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"app/core/property"
	"app/core/util"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

var log *slog.Logger

func SetLogger(l *slog.Logger) {
	log = l
}

type RegisterService func(gsrv *grpc.Server)
type RegisterProxy func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption)
type RegisterRouting func(http.Handler) http.Handler

type Option struct {
	Service  RegisterService
	Proxy    RegisterProxy
	Routing  RegisterRouting
	GrpcAddr string
	GrpcPort string
	HttpAddr string
	HttpPort string
	CORS     bool
}

type Server struct {
	grpc_ln       net.Listener
	grpc_endpoint string
	http_endpoint string
	ctx           context.Context
	cancel        context.CancelFunc

	gsrv *grpc.Server
	gate *http.Server
	opt  *Option
}

func NewServer(opt *Option) (*Server, error) {

	grpc_endpoint := fmt.Sprintf("%s:%s", opt.GrpcAddr, opt.GrpcPort)
	ln, err := net.Listen("tcp", grpc_endpoint)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	s := &Server{
		grpc_endpoint: grpc_endpoint,
		http_endpoint: fmt.Sprintf("%s:%s", opt.HttpAddr, opt.HttpPort),
		grpc_ln:       ln,
		ctx:           ctx, cancel: cancel,
		opt: opt,
	}

	return s, nil
}

// Start starts the combined server
//
// The start process will follow these steps
//   - initiate a grpc server
//   - register grpc services
//   - start the grpc server in a go routine
//   - initiate a grpc gateway
//   - register gateways
//   - setup additional routing rules
//   - start the gateway server in a go routine
func (s *Server) Start() {

	//-------------------------------------------------
	//- Initiate grpc server                          -
	//-------------------------------------------------
	s.gsrv = grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	//-------------------------------------------------
	//- Register gRPC Services                        -
	//-------------------------------------------------
	if s.opt.Service != nil {
		s.opt.Service(s.gsrv)
	}

	reflection.Register(s.gsrv) // enable gRPC reflection

	// Start grpc server in go routine
	go func() {
		if err := s.gsrv.Serve(s.grpc_ln); err != nil {
			slog.Error("start grpc server failed", util.ErrAtrr(err))
			os.Exit(1)
		}
	}()

	//-------------------------------------------------
	//- Initiate gRPC gateway proxy                   -
	//-------------------------------------------------
	gateway := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					EmitUnpopulated: true,
					UseProtoNames:   true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
		runtime.WithOutgoingHeaderMatcher(OutGoingHeaderMatcher),
	)

	//-------------------------------------------------
	//- Register gRPC gateway proxy                   -
	//-------------------------------------------------
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if s.opt.Proxy != nil {
		s.opt.Proxy(s.ctx, gateway, s.grpc_endpoint, opts)
	}

	//-------------------------------------------------
	//- Setup Router                                  -
	//-------------------------------------------------

	var mux http.Handler
	if s.opt.Routing != nil {
		if property.IsDebug() || s.opt.CORS {
			mux = s.opt.Routing(allowCORS(gatewayInterceptor(gateway)))
		} else {
			mux = s.opt.Routing(gatewayInterceptor(gateway))
		}
	} else {
		mux = gatewayInterceptor(gateway)
	}

	//-------------------------------------------------
	//- Start HTTP server in goroutine                -
	//-------------------------------------------------
	s.gate = &http.Server{
		Addr:    s.http_endpoint,
		Handler: mux,
	}
	go func() {
		if err := s.gate.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				slog.Debug("gateway server closed")
			} else {
				slog.Error("start http server failed", err)
				os.Exit(1)
			}
		}
	}()
}

func (s *Server) Stop(timeout time.Duration) {
	termGW, doneGW := context.WithTimeout(context.Background(), timeout)
	term, done := context.WithTimeout(context.Background(), timeout)

	// Try GracefulStop first
	go func() {
		s.gate.Shutdown(termGW)
		doneGW()
	}()
	go func() {
		s.gsrv.GracefulStop()
		done()
	}()

	<-termGW.Done()
	<-term.Done()
	if errors.Is(term.Err(), context.DeadlineExceeded) {
		// force stop grpc server if GracefulStop reached timeout.
		slog.Error("server termination timeout reached!", term.Err())
		s.gsrv.Stop()
	}
	if errors.Is(termGW.Err(), context.DeadlineExceeded) {
		// canceling the app context signals the grpc-gateway to close its connection
		s.cancel()
	}
}

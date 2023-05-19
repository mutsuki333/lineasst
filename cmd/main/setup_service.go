package main

import (
	"app/core/config"
	"app/core/property"
	"app/core/server"
	"app/src"
	"context"
	_ "embed"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"kumoly.io/lib/swaggerui"
)

func setup_service() (*server.Server, error) {

	//-------------------------------------------------
	//- Setup Server                                  -
	//-------------------------------------------------
	property.SetState(property.STATE_PREPARE)

	//-------------------------------------------------
	//- Initiate and Register gRPC Services           -
	//-------------------------------------------------

	return server.NewServer(&server.Option{
		GrpcAddr: config.GetString(property.GRPC_ADDR), GrpcPort: config.GetString(property.GRPC_PORT),
		HttpAddr: config.GetString(property.ADDR), HttpPort: config.GetString(property.PORT),

		//-------------------------------------------------
		//- Register services to gRPC server              -
		//-------------------------------------------------
		Service: func(gsrv *grpc.Server) {
			// service.RegisterCoreServiceServer(gsrv, coreSvc)
		},

		//-------------------------------------------------
		//- Register gRPC gateway proxy                   -
		//-------------------------------------------------
		Proxy: func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
			// service.RegisterCoreServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)

			// add http only handlers
			// mux.HandlePath("POST", "/api/insp_item/import", pmmSvc.ImportInspection)
		},

		//-------------------------------------------------
		//- Setup Router                                  -
		//-------------------------------------------------
		Routing: func(sm http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case strings.HasPrefix(r.URL.Path, "/api"):
					// gRPC gateway
					sm.ServeHTTP(w, r)

				case strings.HasPrefix(r.URL.Path, "/swagger"):
					switch r.URL.Path {
					case "/swagger":
						http.Redirect(w, r, "/swagger/", http.StatusTemporaryRedirect)
					case "/swagger/apidocs.swagger.json":
						w.Write(src.ApiDoc)
					default:
						// swagger ui
						http.StripPrefix("/swagger/", http.FileServer(http.FS(swaggerui.FS))).ServeHTTP(w, r)
					}

				default:
					// default spa website
					// stat, err := fs.Stat(webopi.FS, strings.TrimPrefix(r.URL.Path, "/"))
					// if err != nil || stat.IsDir() {
					// 	index, err := webopi.FS.Open("index.html")
					// 	if err != nil {
					// 		slog.Error("failed to read index.html", util.ErrAtrr(err))
					// 	}
					// 	io.Copy(w, index)
					// 	return
					// }
					// fileserver.ServeHTTP(w, r)
				}
			})
		},
	})

}

package main

import (
	"context"
	"flag"
	"fmt"
	"mime"
	"net"
	"net/http"
	"strings"

	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/log"
	{{cookiecutter.app_name|lower}} "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/proto"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/version"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/service"
	"{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"

	_ "{{cookiecutter.source_path}}/{{cookiecutter.app_name}}/statik"
)

// getOpenAPIHandler serves an OpenAPI UI.
// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")

	statikFS, err := fs.New()
	if err != nil {
		panic("creating OpenAPI filesystem: " + err.Error())
	}

	return http.FileServer(statikFS)
}

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

func runHTTP() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := {{cookiecutter.app_name|lower}}.Register{{cookiecutter.service_name}}ServiceHandlerFromEndpoint(ctx, mux,  *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	gatewayAddr := fmt.Sprintf("0.0.0.0:%d", config.Get().HTTPPort)
	gwServer := &http.Server{
		Addr: gatewayAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/swagger/") {
				http.StripPrefix("/swagger/", getOpenAPIHandler()).ServeHTTP(w, r)
				return
			}
			mux.ServeHTTP(w, r)
		}),
	}
	log.Info("Starting HTTP server on ", gatewayAddr)
	return gwServer.ListenAndServe()

}

func runGRPC() error {
	grpcServerEndpoint := fmt.Sprintf("0.0.0.0:%d", config.Get().GRPCPort)
	lis, err := net.Listen("tcp", grpcServerEndpoint)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
/*
	if *tls {
		if *certFile == "" {
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = data.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}
	*/
	log.Info("Starting GRPC server on ", grpcServerEndpoint)
	grpcServer := grpc.NewServer(opts...)
	{{cookiecutter.app_name|lower}}.Register{{cookiecutter.service_name}}ServiceServer(grpcServer, service.New())
	return grpcServer.Serve(lis)
}

func main() {

	versionFlag := flag.Bool("version", false, "Version")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Build Date:", version.BuildDate)
		fmt.Println("Git Commit:", version.GitCommit)
		fmt.Println("Version:", version.Version)
		fmt.Println("Go Version:", version.GoVersion)
		fmt.Println("OS / Arch:", version.OsArch)
		return
	}

	errChan := make(chan error, 0)
	go func() {
		errChan <- runHTTP()
	}()
	go func() {
		errChan <- runGRPC()
	}()
	log.Error(<-errChan)
}

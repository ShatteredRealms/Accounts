package main

import (
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/internal/option"
	"github.com/ShatteredRealms/Accounts/internal/srv"
	"github.com/ShatteredRealms/Accounts/pkg/repository"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	utilRepository "github.com/ShatteredRealms/GoUtils/pkg/repository"
	utilService "github.com/ShatteredRealms/GoUtils/pkg/service"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

func main() {
	config := option.NewConfig()

	var logger log.LoggerService
	if config.IsRelease() {
		logger = log.NewLogger(log.Info, "")
	} else {
		logger = log.NewLogger(log.Debug, "")
	}

	file, err := ioutil.ReadFile(*config.DBFile.Value)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("reading db file: %v", err))
		os.Exit(1)
	}

	c := &utilRepository.DBConnections{}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("parsing db file: %v", err))
		os.Exit(1)
	}

	db, err := utilRepository.DBConnect(*c)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("db: %v", err))
		os.Exit(1)
	}

	userRepository := repository.NewUserRepository(db)
	if err := userRepository.Migrate(); err != nil {
		logger.Log(log.Error, fmt.Sprintf("user repo: %v", err))
		os.Exit(1)
	}

	userService := service.NewUserService(userRepository)
	jwtService, err := utilService.NewJWTService(*config.KeyDir.Value)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("jwt service: %v", err))
		os.Exit(1)
	}

	grpcServer, gwmux, err := srv.NewServer(userService, jwtService, logger, config)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("server creation: %v", err))
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", config.Address())
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("listen: %v", err))
		os.Exit(1)
	}

	server := &http.Server{
		Addr:    config.Address(),
		Handler: grpcHandlerFunc(grpcServer, gwmux),
	}

	err = server.Serve(lis)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("listen: %v", err))
		os.Exit(1)
	}
}

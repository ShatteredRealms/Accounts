package main

import (
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/internal/option"
	v1 "github.com/ShatteredRealms/Accounts/internal/router/v1"
	"github.com/ShatteredRealms/Accounts/pkg/repository"
	"os"
)

func main() {
	config := option.NewConfig()

	var logger log.LoggerService
	if config.IsRelease() {
		logger = log.NewLogger(log.Info, "")
	} else {
		logger = log.NewLogger(log.Debug, "")
	}

	db, err := repository.DBConnect(*config.DBFile.Value)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("db: %v", err))
		os.Exit(1)
	}

	r, err := v1.InitRouter(db, config, logger)

	if config.IsRelease() {
		logger.Log(log.Info, "Service running")
	}

	err = r.Run(config.Address())
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("error: server crash: %v", err))
	}
}

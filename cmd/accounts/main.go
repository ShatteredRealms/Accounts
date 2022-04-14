package main

import (
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/internal/option"
	v1 "github.com/ShatteredRealms/Accounts/internal/router/v1"
	"github.com/ShatteredRealms/Accounts/pkg/repository"
	"gopkg.in/yaml.v3"
	"io/ioutil"
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

	file, err := ioutil.ReadFile(*config.DBFile.Value)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("reading db file: %v", err))
		os.Exit(1)
	}

	c := &repository.DBConnections{}
	err = yaml.Unmarshal(file, c)
	if err != nil {
		logger.Log(log.Error, fmt.Sprintf("parsing db file: %v", err))
		os.Exit(1)
	}

	db, err := repository.DBConnect(*c)
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
		logger.Log(log.Error, fmt.Sprintf("server: %v", err))
	}
}

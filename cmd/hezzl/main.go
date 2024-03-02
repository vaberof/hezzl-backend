package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/vaberof/hezzl-backend/pkg/database/postgres"
	"github.com/vaberof/hezzl-backend/pkg/database/redis"
)

var appConfigPaths = flag.String("config.files", "not-found.yaml", "List of application config files separated by comma")
var environmentVariablesPath = flag.String("env.vars.file", "not-found.env", "Path to environment variables file")

func main() {
	flag.Parse()
	if err := loadEnvironmentVariables(); err != nil {
		panic(err)
	}

	appConfig := mustGetAppConfig(*appConfigPaths)

	fmt.Printf("%+v\n", appConfig)

	_, err := postgres.New(&appConfig.Postgres)
	if err != nil {
		panic(err)
	}

	_, err = redis.New(&appConfig.Redis)
	if err != nil {
		panic(err)
	}
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}

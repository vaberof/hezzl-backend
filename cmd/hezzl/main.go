package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/vaberof/hezzl-backend/internal/app/entrypoint/http"
	"github.com/vaberof/hezzl-backend/internal/domain/good"
	"github.com/vaberof/hezzl-backend/internal/infra/storage/postgres/pggood"
	redisstorage "github.com/vaberof/hezzl-backend/internal/infra/storage/redis"
	"github.com/vaberof/hezzl-backend/pkg/database/postgres"
	"github.com/vaberof/hezzl-backend/pkg/database/redis"
	"github.com/vaberof/hezzl-backend/pkg/http/httpserver"
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

	postgresManagedDb, err := postgres.New(&appConfig.Postgres)
	if err != nil {
		panic(err)
	}

	redisManagedDb, err := redis.New(&appConfig.Redis)
	if err != nil {
		panic(err)
	}

	pgGoodStorage := pggood.NewPgGoodStorage(postgresManagedDb.PostgresDb)
	redisStorage := redisstorage.NewRedisStorage(redisManagedDb.RedisDb)

	domainGoodService := good.NewGoodService(pgGoodStorage, redisStorage)

	httpHandler := http.NewHandler(domainGoodService)

	appServer := httpserver.New(&appConfig.Server)

	httpHandler.InitRoutes(appServer.Server)

	<-appServer.StartAsync()
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/vaberof/hezzl-backend/internal/app/entrypoint/http"
	"github.com/vaberof/hezzl-backend/internal/domain/good"
	"github.com/vaberof/hezzl-backend/internal/infra/messagebroker/nats/publisher"
	"github.com/vaberof/hezzl-backend/internal/infra/messagebroker/nats/subscriber"
	"github.com/vaberof/hezzl-backend/internal/infra/storage/clickhouse/chgoodlog"
	"github.com/vaberof/hezzl-backend/internal/infra/storage/postgres/pggood"
	redisstorage "github.com/vaberof/hezzl-backend/internal/infra/storage/redis"
	"github.com/vaberof/hezzl-backend/pkg/database/clickhouse"
	"github.com/vaberof/hezzl-backend/pkg/database/postgres"
	"github.com/vaberof/hezzl-backend/pkg/database/redis"
	"github.com/vaberof/hezzl-backend/pkg/http/httpserver"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	clickHouseManagedDb, err := clickhouse.New(&appConfig.ClickHouse)
	if err != nil {
		panic(err)
	}

	goodLogPublisher, err := publisher.New(&appConfig.NatsPublisher)
	if err != nil {
		panic(err)
	}

	goodLogSubscriber, err := subscriber.New(&appConfig.NatsSubscriber)
	if err != nil {
		panic(err)
	}

	pgGoodStorage := pggood.NewPgGoodStorage(postgresManagedDb.PostgresDb, goodLogPublisher)
	redisStorage := redisstorage.NewRedisStorage(redisManagedDb.RedisDb)
	chGoodStorage := chgoodlog.NewCHGoodLogStorage(clickHouseManagedDb.ClickHouseDb)

	goodLogSubscriber.SubscribeOnGoodLogsSubject(chGoodStorage)

	domainGoodService := good.NewGoodService(pgGoodStorage, redisStorage)

	httpHandler := http.NewHandler(domainGoodService)

	appServer := httpserver.New(&appConfig.Server)

	httpHandler.InitRoutes(appServer.ChiRouter)

	serverExitChannel := appServer.StartAsync()

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)

	select {
	case signalValue := <-quitCh:
		log.Println("stopping application", "signal", signalValue.String())

		gracefulShutdown(appServer, postgresManagedDb, redisManagedDb, clickHouseManagedDb)
	case err := <-serverExitChannel:
		log.Println("stopping application", "err", err.Error())

		gracefulShutdown(appServer, postgresManagedDb, redisManagedDb, clickHouseManagedDb)
	}
}

func gracefulShutdown(server *httpserver.AppServer, postgresManagedDb *postgres.ManagedDatabase, redisManagedDb *redis.ManagedDatabase, clickHouseManagedDb *clickhouse.ManagedDatabase) {
	if err := server.Server.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server Shutdown: %v\n", err)
	}

	if err := postgresManagedDb.Disconnect(); err != nil {
		log.Printf("Postgres database Shutdown: %v\n", err)
	}

	if err := redisManagedDb.RedisDb.Close(); err != nil {
		log.Printf("Redis database Shutdown: %v\n", err)
	}

	if err := clickHouseManagedDb.ClickHouseDb.Close(); err != nil {
		log.Printf("ClickHouse database Shutdown: %v\n", err)
	}

	log.Println("Server successfully shutdown")
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}

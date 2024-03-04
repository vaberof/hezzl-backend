package main

import (
	"errors"
	"github.com/vaberof/hezzl-backend/internal/infra/messagebroker/nats/publisher"
	"github.com/vaberof/hezzl-backend/internal/infra/messagebroker/nats/subscriber"
	"github.com/vaberof/hezzl-backend/pkg/config"
	"github.com/vaberof/hezzl-backend/pkg/database/clickhouse"
	"github.com/vaberof/hezzl-backend/pkg/database/postgres"
	"github.com/vaberof/hezzl-backend/pkg/database/redis"
	"github.com/vaberof/hezzl-backend/pkg/http/httpserver"
	"os"
)

type AppConfig struct {
	Server         httpserver.ServerConfig
	Postgres       postgres.Config
	Redis          redis.Config
	ClickHouse     clickhouse.Config
	NatsPublisher  publisher.Config
	NatsSubscriber subscriber.Config
}

func mustGetAppConfig(sources ...string) AppConfig {
	config, err := tryGetAppConfig(sources...)
	if err != nil {
		panic(err)
	}

	if config == nil {
		panic(errors.New("config cannot be nil"))
	}

	return *config
}

func tryGetAppConfig(sources ...string) (*AppConfig, error) {
	if len(sources) == 0 {
		return nil, errors.New("at least 1 source must be set for app config")
	}

	provider := config.MergeConfigs(sources)

	var serverConfig httpserver.ServerConfig
	err := config.ParseConfig(provider, "app.http.server", &serverConfig)
	if err != nil {
		return nil, err
	}

	var postgresConfig postgres.Config
	err = config.ParseConfig(provider, "app.postgres", &postgresConfig)
	if err != nil {
		return nil, err
	}
	postgresConfig.User = os.Getenv("POSTGRES_USER")
	postgresConfig.Password = os.Getenv("POSTGRES_PASSWORD")

	var redisConfig redis.Config
	err = config.ParseConfig(provider, "app.redis", &redisConfig)
	if err != nil {
		return nil, err
	}

	var clickHouseConfig clickhouse.Config
	err = config.ParseConfig(provider, "app.clickhouse", &clickHouseConfig)
	if err != nil {
		return nil, err
	}
	clickHouseConfig.User = os.Getenv("CLICKHOUSE_USER")
	clickHouseConfig.Password = os.Getenv("CLICKHOUSE_PASSWORD")

	var natsPublisher publisher.Config
	err = config.ParseConfig(provider, "app.nats.publisher", &natsPublisher)
	if err != nil {
		return nil, err
	}

	var natsSubscriber subscriber.Config
	err = config.ParseConfig(provider, "app.nats.subscriber", &natsSubscriber)
	if err != nil {
		return nil, err
	}

	appConfig := AppConfig{
		Server:         serverConfig,
		Postgres:       postgresConfig,
		Redis:          redisConfig,
		ClickHouse:     clickHouseConfig,
		NatsPublisher:  natsPublisher,
		NatsSubscriber: natsSubscriber,
	}

	return &appConfig, nil
}

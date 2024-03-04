package clickhouse

import (
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Config struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

type ManagedDatabase struct {
	ClickHouseDb driver.Conn
}

func New(config *Config) (*ManagedDatabase, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", config.Host, config.Port)},
		Auth: clickhouse.Auth{
			Database: config.Database,
			Username: config.User,
			Password: config.Password,
		},
	})
	if err != nil {
		return nil, err
	}

	managedDatabase := &ManagedDatabase{
		ClickHouseDb: conn,
	}

	return managedDatabase, nil
}

func (db *ManagedDatabase) Disconnect() error {
	return db.ClickHouseDb.Close()
}

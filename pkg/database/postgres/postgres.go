package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

type ManagedDatabase struct {
	PostgresDb *sqlx.DB
}

func New(config *Config) (*ManagedDatabase, error) {
	psqlUrl := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Database)

	psqlDb, err := sqlx.Connect("postgres", psqlUrl)
	if err != nil {
		return nil, err
	}

	managedDatabase := &ManagedDatabase{
		PostgresDb: psqlDb,
	}

	return managedDatabase, nil
}

func (db *ManagedDatabase) Disconnect() error {
	return db.PostgresDb.Close()
}

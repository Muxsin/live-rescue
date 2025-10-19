package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDatabase struct {
	Db *gorm.DB
}

var instance *PostgresDatabase

func NewPostgres() (*PostgresDatabase, error) {
	if instance == nil {
		db, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES")), &gorm.Config{})
		if err != nil {
			return nil, err
		}

		instance = &PostgresDatabase{Db: db}
	}

	return instance, nil
}

package config

import (
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(cfg Config) (*gorm.DB, error) {
	if strings.HasPrefix(strings.TrimSpace(cfg.Database.DSN), "jdbc:") {
		return nil, fmt.Errorf("DATABASE_DSN must be a Postgres DSN (keyword/value), not a JDBC URL. Example: host=localhost user=postgres password=... dbname=... port=5432 sslmode=disable")
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

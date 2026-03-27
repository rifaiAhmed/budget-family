package config

import (
	"database/sql"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func RunMigrations(cfg Config, logger *zap.Logger) error {
	if !cfg.Migrate.Auto {
		return nil
	}

	db, err := sql.Open("postgres", cfg.Database.DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(cfg.Migrate.Path, "postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("migrations: no change")
			return nil
		}
		return err
	}

	logger.Info("migrations: applied")
	return nil
}

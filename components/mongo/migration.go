package mongo

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jucardi/go-titan/logx"
)

func MigrateMongo(cfg *Config) {
	if cfg == nil || cfg.MigrationSource == "" {
		return
	}
	migrator, err := migrate.New(cfg.MigrationSource, cfg.url())

	logx.WithObj(err).Fatal("Could not initialize migrator")

	logx.Info(fmt.Sprintf("Applying migrations from source %s", cfg.MigrationSource))

	migErr := migrator.Up()

	if migErr != nil && migErr != migrate.ErrNoChange {
		logx.WithObj(migErr).Fatal("Could not run migrations")
	}

	if migErr == migrate.ErrNoChange {
		logx.Info("No migrations to apply")
	}
}

package mongo

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jucardi/go-titan/logx"
)

func Migrate(cfg ...*Config) {
	var c *Config
	if len(cfg) > 0 && cfg[0] != nil {
		c = cfg[0]
	} else {
		c = getConfig()
	}

	if c == nil || c.MigrationSource == "" {
		return
	}
	migrator, err := migrate.New(c.MigrationSource, c.url())

	logx.WithObj(err).Fatal("Could not initialize migrator")

	logx.Info(fmt.Sprintf("Applying migrations from source %s", c.MigrationSource))

	migErr := migrator.Up()

	if migErr != nil && migErr != migrate.ErrNoChange {
		logx.WithObj(migErr).Fatal("Could not run migrations")
	}

	if migErr == migrate.ErrNoChange {
		logx.Info("No migrations to apply")
	}
}

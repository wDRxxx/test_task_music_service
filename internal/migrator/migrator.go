package migrator

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(connectionURL string, migrationsPath string) error {
	m, err := migrate.New("file://"+migrationsPath, connectionURL)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	log.Println("migrations were applied successfully!")
	return nil
}

package config

import (
	"fmt"
	"os"
	"time"
)

type PostgresConfig struct {
	host           string
	port           string
	user           string
	password       string
	database       string
	migrationsPath string
	timeout        time.Duration
}

func (c *PostgresConfig) Timeout() time.Duration {
	return c.timeout
}

func (c *PostgresConfig) MigrationsPath() string {
	return c.migrationsPath
}

func NewPostgresConfig() *PostgresConfig {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		panic("POSTGRES_HOST environment variable is empty")
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		panic("POSTGRES_PORT environment variable is empty")
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		panic("POSTGRES_USER environment variable is empty")
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		panic("POSTGRES_PASSWORD environment variable is empty")
	}

	database := os.Getenv("POSTGRES_DB")
	if database == "" {
		panic("POSTGRES_DB environment variable is empty")
	}

	migrationsPath := os.Getenv("POSTGRES_MIGRATIONS_PATH")
	if migrationsPath == "" {
		panic("POSTGRES_MIGRATIONS_PATH environment variable is empty")
	}

	strTimeout := os.Getenv("POSTGRES_TIMEOUT")
	if strTimeout == "" {
		panic("POSTGRES_TIMEOUT environment variable is empty")
	}
	timeout, err := time.ParseDuration(strTimeout)
	if err != nil {
		panic("POSTGRES_TIMEOUT environment variable is invalid")
	}

	return &PostgresConfig{
		host:           host,
		port:           port,
		user:           user,
		password:       password,
		database:       database,
		migrationsPath: migrationsPath,
		timeout:        timeout,
	}
}

func (c *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.host,
		c.port,
		c.database,
		c.user,
		c.password,
	)
}

func (c *PostgresConfig) ConnectionURL() string {
	return fmt.Sprintf(
		"pgx5://%s:%s@%s:%s/%s",
		c.user,
		c.password,
		c.host,
		c.port,
		c.database,
	)
}

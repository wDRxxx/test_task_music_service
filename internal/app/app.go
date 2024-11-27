package app

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/golang-migrate/migrate/v4"

	"github.com/wDRxxx/test-task/internal/closer"
	"github.com/wDRxxx/test-task/internal/config"
	"github.com/wDRxxx/test-task/internal/migrator"
)

type App struct {
	wg sync.WaitGroup

	serviceProver *serviceProvider

	httpServer       *http.Server
	prometheusServer *http.Server
}

func NewApp(ctx context.Context, envPath string) (*App, error) {
	err := config.Load(envPath)
	if err != nil {
		return nil, err
	}

	app := &App{wg: sync.WaitGroup{}}
	err = app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) initDeps(ctx context.Context) error {
	a.serviceProver = newServiceProvider()

	a.initHttpServer(ctx)

	return nil
}

func (a *App) initHttpServer(ctx context.Context) {
	s := a.serviceProver.HTTPServer(ctx)
	a.httpServer = &http.Server{
		Addr:    a.serviceProver.HttpConfig().Address(),
		Handler: s.Handler(),
	}
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	err := migrator.Migrate(
		a.serviceProver.PostgresConfig().ConnectionURL(),
		a.serviceProver.PostgresConfig().MigrationsPath(),
	)
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		err = a.runHttpServer()
		if err != nil {
			log.Fatalf("error running http server: %v", err)
		}
	}()
	a.wg.Wait()

	return nil
}

func (a *App) runHttpServer() error {
	slog.Info("starting http server...")

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"log"
	"log/slog"
	"os"

	"spy-cat-agency/internal/cats"
	"spy-cat-agency/internal/env"
	"spy-cat-agency/internal/missions"
	"spy-cat-agency/internal/storage"

	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "spy-cat-agency/docs"
)

// @title Spy Cat Agency API
// @version 1.0
// @description This is a sample server for a spy cat agency.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:7777
// @BasePath /
type config struct {
	port      string
	breedsApi string
	db        storage.Config
}

type application struct {
	config
	cats     *cats.Service
	missions *missions.Service
}

func main() {
	cfg := config{
		port:      env.GetString("SPY_CAT_AGENCY_PORT", ":7777"),
		breedsApi: env.GetString("CATS_BREEDS_API", "https://api.thecatapi.com/v1/breeds"),
		db: storage.Config{
			Dsn:          env.GetString("DB_DSN", ""),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		},
	}

	valid := validator.New()

	db, err := storage.ConnectSQL(cfg.db)
	if err != nil {
		panic(err)
	}

	slog.Info("Applying Migrations")
	pm, err := migrate.New("file://internal/storage/migrations", cfg.db.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := pm.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}
	slog.Info("MIgrations applied")

	breeds, err := cats.NewBreeds(cfg.breedsApi)
	if err != nil {
		panic(err)
	}

	catsRepo := cats.NewRepository(db)
	catsService := cats.NewService(catsRepo, valid, breeds)

	missionsRepo := missions.NewRepository(db)
	missionsService := missions.NewService(missionsRepo, valid)

	valid.RegisterValidation("validBreed", catsService.Exists)
	catsService.Validate = valid

	app := &application{
		config:   cfg,
		cats:     catsService,
		missions: missionsService,
	}

	slog.Info("Listening on", "port", cfg.port)
	if err := app.routes().Run(app.config.port); err != nil {
		slog.Error("Server error, shutting down...", "error", err)
		os.Exit(1)
	}
}

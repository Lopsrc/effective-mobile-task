package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"effectiveM-test-task/internal/config"
	delcar "effectiveM-test-task/internal/handlers/car/delete"
	getcar "effectiveM-test-task/internal/handlers/car/get"
	getallcars "effectiveM-test-task/internal/handlers/car/getall"
	reccar "effectiveM-test-task/internal/handlers/car/recover"
	savecar "effectiveM-test-task/internal/handlers/car/save"
	upcar "effectiveM-test-task/internal/handlers/car/update"
	delperson "effectiveM-test-task/internal/handlers/person/delete"
	getperson "effectiveM-test-task/internal/handlers/person/get"
	recperson "effectiveM-test-task/internal/handlers/person/recover"
	saveperson "effectiveM-test-task/internal/handlers/person/save"
	upperson "effectiveM-test-task/internal/handlers/person/update"
	carstore "effectiveM-test-task/internal/storage/postgresql/car"
	personstore "effectiveM-test-task/internal/storage/postgresql/person"
	psql "effectiveM-test-task/pkg/client/postgresql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	pathConfig  = "config/local.yaml"
	envLocal    = "local"
	envDev      = "dev"
	envProd     = "prod"
	maxAttempts = 3
)

func main() {
	// init config
	cfg := config.GetConfig(pathConfig)
	adress := getAdress(cfg.Listen.BindIP, cfg.Listen.Port)
	// init logger
	log := setupLogger(cfg.Env)
	log.Info(
		"starting car",
		slog.String("env", cfg.Env),
		slog.String("version", "1"),
	)

	log.Debug("debug messages are enabled ")
	// Init client of the postgresql
	postgreSQLClient, err := psql.NewClient(context.TODO(), maxAttempts, cfg.Storage)
	if err != nil {
		log.Error("failed to init storage %v", err)
		os.Exit(1)
	}
	defer postgreSQLClient.Close()
	// create a new repositories
	carRepository := carstore.NewRepository(postgreSQLClient, log)
	personRepository := personstore.NewRepository(postgreSQLClient, log)
	// create a new router
	router := chi.NewRouter()
	// init middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	// Init endpoints
	router.Route("/car", func(r chi.Router) {
		r.Post("/save/",	savecar.New(log, carRepository))
		r.Get("/get/", getcar.New(log, carRepository))
		r.Get("/getall", getallcars.New(log, carRepository))
		r.Post("/up/", upcar.New(log, carRepository))
		r.Delete("/del/", delcar.New(log, carRepository))
		r.Post("/rec/", reccar.New(log, carRepository))
	})

	router.Route("/person", func(r chi.Router) {
		r.Post("/save/", saveperson.New(log, personRepository))
		r.Get("/get/", getperson.New(log, personRepository))
		r.Post("/up/", upperson.New(log, personRepository))
		r.Delete("/del/", delperson.New(log, personRepository))
		r.Post("/rec/", recperson.New(log, personRepository))
	})

	log.Info("starting server", slog.String("address", adress))

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	serv := &http.Server{
		Addr:         adress,
		Handler:      router,
		ReadTimeout:  cfg.Listen.Timeout,
		WriteTimeout: cfg.Listen.Timeout,
		IdleTimeout:  cfg.Listen.IdleTimeout,
	}

	go func() {
		if err := serv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")
	// Init context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := serv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server %v", err)
		return
	}

	
	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func getAdress(ip, port string) string {
	return fmt.Sprintf("%s:%s", ip, port)
}

package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"app/config"
	"app/conn"
	"app/logger"
)

var env = flag.String("environment", "dev", "set the application environment")

func main() {
	flag.Parse()

	if *env != "dev" && *env != "prod" && *env != "staging" {
		panic(fmt.Sprintf("the app only has 'dev', 'prod' & 'staging' environemnt and no %s", *env))
	}

	appConfig, err := config.LoadConfig(fmt.Sprintf(".env.%s", *env))
	if err != nil {
		panic(err)
	}

	appLogger := logger.New(appConfig.LogFileLoc)
	appLogger.Info("application started", zap.String("ENVIRONMENT_NAME", *env))

	pool, err := conn.ConnectToPG(&appConfig.Database, appLogger)
	if err != nil {
		panic(err)
	}

	redis, err := conn.NewRedisClient(&appConfig.Redis, appLogger)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	Route(r, pool, appLogger, &appConfig, &redis)
	hub := newHub()

	go hub.run()

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWebSocket(hub, w, r)
	})

	server := http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	appLogger.Info("server started on successfully", zap.Int("SERVER_PORT", appConfig.Server.Port))

	if err := server.ListenAndServe(); err != nil {
		appLogger.Panic("server did not start", zap.Error(err))
	}
}

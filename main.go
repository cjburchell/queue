package main

import (
	"context"
	"fmt"
	queueroute "github.com/cjburchell/queue/routes/queue"
	"github.com/cjburchell/queue/routes/status"
	"github.com/cjburchell/queue/serivce/data"
	"github.com/cjburchell/queue/serivce/queue"
	"github.com/cjburchell/queue/settings"
	config "github.com/cjburchell/settings-go"
	"github.com/cjburchell/tools-go/env"
	"github.com/cjburchell/uatu-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	configFile := config.Get(env.Get("SettingsFile", ""))
	logger := log.Create(configFile)

	appConfig, err := settings.Get(logger, configFile)
	if err != nil{
		logger.Fatal(err, "Unable to verify settings")
	}

	dataService, err := data.NewService(appConfig.MongoURL)
	if err != nil {
		logger.Fatalf(err, "Unable to Connect to mongo %s", appConfig.MongoURL)
	}

	srv := startHTTPServer(appConfig.Port, dataService, logger)
	defer stopHTTPServer(srv, logger)

	workers := queue.StartWorkers(*appConfig, dataService, logger)
	defer workers.Stop()

	// wait for app shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	logger.Print("shutting down")
	os.Exit(0)
}

func stopHTTPServer(srv *http.Server, logger log.ILog) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Error(err)
	}
}

func startHTTPServer(port int, dataService data.IService, logger log.ILog) *http.Server {
	r := mux.NewRouter()
	status.Setup(r, logger)
	queueroute.Setup(r, dataService, logger)

	loggedRouter := handlers.LoggingHandler(logger.GetWriter(log.DEBUG), r)

	logger.Printf("Starting Server at port %d", port)
	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(err)
		}
	}()

	return srv
}

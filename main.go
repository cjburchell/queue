package main

import (
	"context"
	"fmt"
	log "github.com/cjburchell/go-uatu"
	logSettings "github.com/cjburchell/go-uatu/settings"
	queueRoute "github.com/cjburchell/queue/routes/queue-route"
	"github.com/cjburchell/queue/routes/status-route"
	"github.com/cjburchell/queue/serivce/data"
	"github.com/cjburchell/queue/serivce/queue"
	"github.com/cjburchell/queue/settings"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	err := logSettings.SetupLogger()
	if err != nil {
		log.Warn(err, "Unable to Connect to logger")
	}

	config, err := settings.Get()
	if err != nil{
		log.Fatal(err, "Unable to verify settings")
	}

	dataService, err := data.NewService(config.MongoUrl)
	if err != nil {
		log.Fatalf(err, "Unable to Connect to mongo %s", config.MongoUrl)
	}

	srv := startHttpServer(config.Port, dataService)
	defer stopHttpServer(srv)

	workers := queue.StartWorkers(*config, dataService)
	defer workers.Stop()

	// wait for app shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	log.Print("shutting down")
	os.Exit(0)
}

func stopHttpServer(srv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Error(err)
	}
}

func startHttpServer(port int, dataService data.IService) *http.Server {
	r := mux.NewRouter()
	status_route.Setup(r)
	queueRoute.Setup(r, dataService)

	loggedRouter := handlers.LoggingHandler(log.Writer{Level: log.DEBUG}, r)

	log.Printf("Starting Server at port %d", port)
	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err)
		}
	}()

	return srv
}

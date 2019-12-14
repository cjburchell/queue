package queue_route

import (
	"encoding/json"
	log "github.com/cjburchell/go-uatu"
	"github.com/cjburchell/queue/routes/contract"
	"github.com/cjburchell/queue/routes/token"
	"github.com/cjburchell/queue/serivce/data"
	"github.com/gorilla/mux"
	"net/http"
)

func Setup(r *mux.Router, dataService data.IService) {
	r.HandleFunc("/queue/job", token.ValidateMiddleware(func(writer http.ResponseWriter, request *http.Request) {
		handlePostJob(writer, request, dataService)
	})).Methods("POST")
}

func handlePostJob(writer http.ResponseWriter, request *http.Request, dataService data.IService) {
	decoder := json.NewDecoder(request.Body)
	var job contract.Job
	err := decoder.Decode(&job)
	if err != nil {
		log.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = dataService.AddJob(job.Call, job.Repeat, job.Delay, job.Retries, job.Priority)
	if err != nil {
		log.Error(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}


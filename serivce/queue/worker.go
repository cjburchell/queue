package queue

import (
	"bytes"
	"github.com/cjburchell/queue/serivce/data/models"
	"github.com/cjburchell/uatu-go"
	"net/http"
)

type queueItem struct {
	models.Job
	Completed bool
}

type worker struct {
	workerQueue chan worker
	quitChan    chan bool
	Job         chan queueItem
	logger      log.ILog
}

func newWorker(workerQueue chan worker, logger log.ILog) worker {

	worker := worker{
		Job:         make(chan queueItem),
		workerQueue: workerQueue,
		quitChan:    make(chan bool),
		logger: logger,
	}

	return worker

}

func (w *worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.workerQueue <- *w

			select {
			case work := <-w.Job:
				work.Completed = false
				// Receive a work request.
				work = w.Process(work)
				//Return work
				w.Job <- work

			case <-w.quitChan:
				// We have been asked to stop.
				return
			}
		}
	}()
}

func (w worker) Process(job queueItem) queueItem {
	if w.doHTTPRequest(job.Call){
		job.Completed = true
	}

	return job
}

func (w worker) doHTTPRequest(call models.Call) bool {

	restClient := &http.Client{}

	var body *bytes.Buffer
	if len(call.StringBody) == 0 {
		body= bytes.NewBuffer(call.ResponseBody)
	} else {
		body= bytes.NewBuffer([]byte(call.StringBody))
	}


	req, err := http.NewRequest(call.Method, call.Path, body)
	if err != nil {
		w.logger.Error(err)
		return false
	}

	req.Header.Add("Content-Type", call.ContentType)

	resp, err := restClient.Do(req)
	if err != nil {
		w.logger.Error(err)
		return false
	}

	if resp.StatusCode != http.StatusOK {
		w.logger.Warnf("Unable to %s %s(%d)",call.Method, call.Path, resp.StatusCode)
		return false
	}

	return true
}

func (w worker) Stop() {
	go func() {
		w.quitChan <- true
	}()
}
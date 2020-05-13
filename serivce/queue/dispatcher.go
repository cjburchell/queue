package queue

import (
	"github.com/cjburchell/queue/serivce/data"
	"github.com/cjburchell/queue/settings"
	"github.com/cjburchell/uatu-go"
	"time"
)


func getJob(maxJobTimeMilliseconds int64, dataService data.IService) (*queueItem, error) {

	job, err := dataService.GetNextJob(maxJobTimeMilliseconds)
	if err != nil || job == nil {
		return nil, err
	}

	return &queueItem{Job: *job}, err
}

// Dispatcher object
type Dispatcher struct {
	quitChan chan bool
	workerQueue chan worker
	workers []worker
	configuration settings.Configuration
	dataService data.IService
	logger log.ILog
}

// Stop the dispatcher
func (d Dispatcher)Stop()  {
	go func() {
		d.quitChan <- true
	}()
}

// StartWorkers starts all workers
func StartWorkers(configuration settings.Configuration, dataService data.IService, logger log.ILog) Dispatcher {
	workerQueue := make(chan worker, configuration.MaxWorkers)
	workers := make([]worker, configuration.MaxWorkers)
	for i := 0; i < configuration.MaxWorkers; i++ {
		workers[i] = newWorker(workerQueue, logger)
		logger.Printf("Starting worker %d", i+1)
		workers[i].Start()
	}

	dispatcher := Dispatcher{ make(chan bool), workerQueue, workers, configuration, dataService, logger}

	go dispatcher.dispatch()

	return dispatcher
}

func (d Dispatcher)dispatch() {
	for {
		select {
		case worker := <-d.workerQueue:
			job, err := getJob(d.configuration.MaxJobTimeMilliseconds, d.dataService)
			if err != nil || job == nil {
				if err != nil {
					d.logger.Error(err, "Failed to Get job")
				}

				d.workerQueue <- worker
				time.Sleep(time.Duration(d.configuration.SleepMilliseconds) * time.Millisecond)
			} else {
				go process(worker, job, d.dataService, d.configuration, d.logger)
			}
		case <-d.quitChan:
			for i := 0; i < d.configuration.MaxWorkers; i++ {
				d.logger.Printf("Stopping worker %d", i+1)
				d.workers[i].Stop()
			}
			return
		}
	}
}

func process(worker worker, job *queueItem, dataService data.IService, configuration settings.Configuration, logger log.ILog) {
	worker.Job <- *job
	workDone := <-worker.Job

	if workDone.Completed == true {
		err := stopJob(job, dataService)
		if err != nil {
			logger.Error(err)
		}
	} else {
		if workDone.Tries >= workDone.MaxRetries {
			logger.Warnf("Maximum number of retries for a job reached (%d), removing job: %+v", workDone.MaxRetries, workDone)
			err := stopJob(job, dataService)
			if err != nil {
				logger.Error(err)
			}
		} else {
			err := dataService.DelayJob(job.ID, configuration.RetryDelay*int64(workDone.Tries), job.Priority)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func stopJob(job *queueItem, dataService data.IService) error {
	if job.Repeat {
		// requeue the job
		err := dataService.DelayJob(job.ID, job.Delay, job.InitialPriority)
		if err != nil {
			return err
		}
	} else {
		err := dataService.DeleteJob(job.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
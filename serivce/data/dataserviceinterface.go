package data

import "github.com/cjburchell/queue/serivce/data/models"

// IService interface
type IService interface {
	AddJob(call models.Call, repeat bool, delay int64, retries int, priority int) error
	DeleteJob(jobID string) error
	GetNextJob(maxJobTime int64) (*models.Job, error)
	DelayJob(jobID string, delay int64, priority int) error
}

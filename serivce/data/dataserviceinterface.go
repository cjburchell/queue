package data

import "github.com/cjburchell/queue/serivce/data/models"

type IService interface {
	AddJob(call models.Call, repeat bool, delay int64, retries int, priority int) error
	DeleteJob(jobId string) error
	GetNextJob(maxJobTime int64) (*models.Job, error)
	DelayJob(jobId string, delay int64, priority int) error
}

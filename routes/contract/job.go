package contract

import "github.com/cjburchell/queue/serivce/data/models"

type Job struct {
	Call     models.Call `json:"call"`
	Repeat   bool        `json:"repeat"`
	Delay    int64       `json:"delay"`
	Retries  int         `json:"retries"`
	Priority int         `json:"priority"`
}

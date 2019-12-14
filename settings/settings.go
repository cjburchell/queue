package settings

import (
	"fmt"
	"github.com/cjburchell/queue/log"
	"github.com/cjburchell/queue/tools/env"
)

const defaultMongoUrl ="localhost"

const defaultPort  = 8091
const maxWorkers = 10
const defaultRetryDelay int64 = 1000
const sleepMilliseconds int64 = 1000
const maxJobTime int64 = 60000

type Configuration struct {
	MongoUrl               string
	Port                   int
	MaxWorkers             int
	RetryDelay             int64
	SleepMilliseconds      int64
	MaxJobTimeMilliseconds int64
}

func Get(logger log.ILog) (*Configuration, error) {
	err := verify(logger)
	if err != nil {
		return nil, err
	}

	return &Configuration{
		MongoUrl:               env.Get("MONGO_URL", defaultMongoUrl),
		Port:                   env.GetInt("PORT", defaultPort),
		MaxWorkers:             env.GetInt("maxWorkers", maxWorkers),
		RetryDelay:             env.GetInt64("retryDelay", defaultRetryDelay),
		SleepMilliseconds:      env.GetInt64("sleep", sleepMilliseconds),
		MaxJobTimeMilliseconds: env.GetInt64("maxJobTime", maxJobTime),
	}, nil
}


func verify(logger log.ILog) error {

	warningMessage := ""
	if env.Get("MONGO_URL", defaultMongoUrl) == defaultMongoUrl {
		warningMessage += fmt.Sprintf("\nMONGO_URL set to default value (%s)", defaultMongoUrl)
	}

	if warningMessage != "" {
		logger.Warn("Warning: " + warningMessage)
	}

	return nil
}
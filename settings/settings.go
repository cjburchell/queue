package settings

import (
	"fmt"
	"github.com/cjburchell/settings-go"
	"github.com/cjburchell/uatu-go"
)

const defaultMongoURL ="localhost"

const defaultPort  = 8091
const maxWorkers = 10
const defaultRetryDelay int64 = 1000
const sleepMilliseconds int64 = 1000
const maxJobTime int64 = 60000

// Configuration of the application
type Configuration struct {
	MongoURL               string
	Port                   int
	MaxWorkers             int
	RetryDelay             int64
	SleepMilliseconds      int64
	MaxJobTimeMilliseconds int64
}

// Get the application settings
func Get(logger log.ILog, settings settings.ISettings) (*Configuration, error) {
	err := verify(logger, settings)
	if err != nil {
		return nil, err
	}

	return &Configuration{
		MongoURL:               settings.Get("MONGO_URL", defaultMongoURL),
		Port:                   settings.GetInt("PORT", defaultPort),
		MaxWorkers:             settings.GetInt("maxWorkers", maxWorkers),
		RetryDelay:             settings.GetInt64("retryDelay", defaultRetryDelay),
		SleepMilliseconds:      settings.GetInt64("sleep", sleepMilliseconds),
		MaxJobTimeMilliseconds: settings.GetInt64("maxJobTime", maxJobTime),
	}, nil
}


func verify(logger log.ILog, settings settings.ISettings) error {
	warningMessage := ""
	if settings.Get("MONGO_URL", defaultMongoURL) == defaultMongoURL {
		warningMessage += fmt.Sprintf("\nMONGO_URL set to default value (%s)", defaultMongoURL)
	}

	if warningMessage != "" {
		logger.Warn("Warning: " + warningMessage)
	}

	return nil
}
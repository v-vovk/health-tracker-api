package logger

import (
	"go.uber.org/zap"
	"os"
)

var Log *zap.Logger

func InitLogger(env string) {
	environment := env
	if environment == "" {
		environment = "PROD"
	}

	// Ensure logs directory exists
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
		if err := os.Mkdir("./logs", 0755); err != nil {
			panic("Failed to create logs directory: " + err.Error())
		}
	}

	var config zap.Config
	if environment == "DEV" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	// Add logs to both stdout and file
	config.OutputPaths = []string{"stdout", "./logs/app.log"}

	// Add initial fields
	config.InitialFields = map[string]interface{}{
		"service":     "health-tracker-api",
		"environment": environment,
	}

	var err error
	Log, err = config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}

// Sync flushes any buffered log entries
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

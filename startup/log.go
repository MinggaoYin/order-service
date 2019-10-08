package startup

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	appEnv := os.Getenv("APP_ENV")

	switch appEnv {
	case "testing":
		fallthrough
	case "localhost":
		// log as text form easier for human to read
		logrus.SetFormatter(&logrus.TextFormatter{})

		// log to stdout
		logrus.SetOutput(os.Stdout)

		// lower log level for testing and debugging
		logrus.SetLevel(logrus.DebugLevel)
	case "development":
		fallthrough
	case "staging":
		fallthrough
	case "production":
		setupLogger()
	default:
		fmt.Printf("Failed to start, since APP_ENV = %s is not registered\n", appEnv)
		os.Exit(1)
	}
}

func setupLogger() {
	// log as JSON format
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// log to stdout
	logrus.SetOutput(os.Stdout)

	// only log the info severity or above.
	logrus.SetLevel(logrus.InfoLevel)
}

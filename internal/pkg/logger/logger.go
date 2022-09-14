package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func NewLogger() (*logrus.Logger, error) {
	log := logrus.New()
	lvl, err := logrus.ParseLevel("debug")
	if err != nil {
		return nil, fmt.Errorf("error parsing log level: %w", err)
	}
	log.SetLevel(lvl)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(true)

	return log, nil
}

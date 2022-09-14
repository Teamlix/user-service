package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/teamlix/user-service/internal/pkg/config"
	log "github.com/teamlix/user-service/internal/pkg/logger"
)

func Run(configPath string) error {

	var cfg config.Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return err
	}

	logger, err := log.NewLogger()
	if err != nil {
		return nil
	}

	// run grpc server

	// listen to os signals
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-s:
		logger.Infof("OS signal: %s", sig.String())
	}

	return nil
}

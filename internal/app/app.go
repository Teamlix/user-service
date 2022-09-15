package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/teamlix/user-service/internal/cache"
	"github.com/teamlix/user-service/internal/pkg/bcrypt"
	"github.com/teamlix/user-service/internal/pkg/config"
	grpc_server "github.com/teamlix/user-service/internal/pkg/grpc/server"
	log "github.com/teamlix/user-service/internal/pkg/logger"
	"github.com/teamlix/user-service/internal/pkg/mongo"
	"github.com/teamlix/user-service/internal/pkg/redis"
	"github.com/teamlix/user-service/internal/pkg/validator"
	"github.com/teamlix/user-service/internal/repository"
	"github.com/teamlix/user-service/internal/service"
)

func Run(configPath string) error {
	errCh := make(chan error, 1)
	ctx := context.Background()
	var grpcsrv grpc_server.Server

	var cfg config.Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return err
	}

	logger, err := log.NewLogger()
	if err != nil {
		return nil
	}

	mCon, err := mongo.NewMongo(ctx, cfg.MongoDB.URL)
	if err != nil {
		return err
	}

	repo := repository.NewRepository(mCon)

	rCon, err := redis.NewRedis(ctx, cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return err
	}

	c := cache.NewCache(rCon)

	b := bcrypt.NewBcrypt(10)

	v := validator.NewValidator()

	s := service.NewService(repo, c, b, v)

	// run grpc server
	go func() {
		grpcsrv = grpc_server.NewServer(
			cfg.Grpc.Server.Host,
			cfg.Grpc.Server.Port,
			s,
			logger,
		)
		err := grpcsrv.Serve()
		if err != nil {
			errCh <- err
		}
	}()

	logger.Infof("gRPC server started on: %s:%s", cfg.Grpc.Server.Host, cfg.Grpc.Server.Port)

	// listen to os signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		logger.Infof("OS signal: %s", sig.String())

	case err := <-errCh:
		logger.Errorf("Got error: %v", err)
	}

	grpcsrv.Stop()
	logger.Info("gRPC server closed")

	err = rCon.Disconnect()
	if err != nil {
		return nil
	}

	logger.Info("Redis connection closed")

	err = mCon.Disconnect(ctx)
	if err != nil {
		return nil
	}

	logger.Info("Mongo connection closed")

	return nil
}

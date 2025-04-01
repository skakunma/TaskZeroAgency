package config

import (
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/logrusadapter"
	"github.com/sirupsen/logrus"
	"github.com/skakunma/TaskZeroAgency/internal/storage"
	"os"
)

type Config struct {
	Store     storage.Storage
	Logger    *logrusadapter.Logger
	SecretKEY string
}

func CreateConfig() *Config {
	conf := &Config{}

	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	reformLogger := logrusadapter.NewLogger(logger)
	conf.Logger = reformLogger

	dsn := os.Getenv("DATABASE_URL")

	if dsn == "" {
		reformLogger.Log(pgx.LogLevelError, "Отсутсвует DATABASE_URL в env", nil)
		return nil
	}

	secretkey := os.Getenv("SECRET_KEY")
	if secretkey == "" {
		reformLogger.Log(pgx.LogLevelWarn, "SECRET_KEY не был найден будет использованно деволтное значение ", nil)
		secretkey = "very hard secret key"
	}

	conf.SecretKEY = secretkey

	store, err := storage.CreatePostgreStorage(dsn)
	if err != nil {
		reformLogger.Log(pgx.LogLevelError, err.Error(), nil)
		return nil
	}

	conf.Store = store
	return conf
}

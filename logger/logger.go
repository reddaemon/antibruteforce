package logger

import (
	"github.com/reddaemon/antibruteforce/config"
	"go.uber.org/zap"
)

func GetLogger(cfg *config.Config) (*zap.Logger, error) {
	var err error

	var l *zap.Logger

	if !cfg.Debug {
		l = zap.NewNop()
		return l, nil
	}

	switch cfg.Environment {
	case "production":
		l, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	case "dev":
		l, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	default:
		l = zap.NewExample()
	}

	return l, nil
}

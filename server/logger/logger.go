package logger

import (
	"golty/api"

	"go.uber.org/zap"
)

func Init(environment string) (*zap.Logger, error) {
	if environment == api.LOCAL_ENV || environment == api.STAGING_ENV {
		logger, err := zap.NewDevelopment()
		if err != nil {
			return nil, err
		}

		return logger, nil
	}

	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

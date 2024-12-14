package logger

import "go.uber.org/zap"

func New() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	return sugar, nil
}

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	LOG = getLogger()
)

func getLogger() *zap.SugaredLogger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		panic("fail to init logger, Err:" + err.Error())
	}
	//defer func(logger *zap.Logger) {
	//	err := logger.Sync()
	//	if err != nil {
	//		panic(err)
	//	}
	//}(logger) // flushes buffer, if any
	sugar := logger.Sugar()
	return sugar
}

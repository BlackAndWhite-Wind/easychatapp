package initialize

import (
	"go.uber.org/zap"
	"log"
)

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("日志初始化失败", err.Error())
	}
	zap.ReplaceGlobals(logger)
}

package initialize

import (
	"EasyChatApp/global"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {
	v := viper.New()                               //实例化对象
	v.SetConfigFile("../HiChat/config-debug.yaml") //读取配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := v.Unmarshal(&global.ServiceConfig); err != nil {
		panic(err)
	}
	zap.S().Info("配置信息", global.ServiceConfig)
}

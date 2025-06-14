package initialize

import (
	"github.com/spf13/viper"
	"os"
	"path"
	"srv/userSrv/global"
)

func GetEnv(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig() {
	// 通过环境变量读取是否debug
	// 然后读取相对应的配置文件
	// 然后解析到serverConfig中
	cur, _ := os.Getwd()
	fileName := "config.yaml"
	if !GetEnv("GoProjectDebug") {
		fileName = "config_debug.yaml"
	}

	v := viper.New()
	v.SetConfigFile(path.Join(cur, fileName))
	err := v.ReadInConfig()
	if err != nil {
		panic("initialize config file error:" + err.Error())
	}

	err = v.Unmarshal(&global.SrvConfig)
	if err != nil {
		panic("initialize config file error:" + err.Error())
	}
}

package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"goodsSrv/global"
	"goodsSrv/structs"
	"os"
	"path"
	"path/filepath"
)

var (
	nacos = &structs.NacosCnf{}
)

func GetEnv(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

// 通过环境变量以及文件读取nacos文件，然后读取nacos文件的配置
// 然后通过解析到global.SrvConfig

func InitConfig() {
	curPath, _ := os.Getwd()
	configName := "config.json"
	debug := GetEnv("GoProjectDebug")
	if debug {
		configName = "config_debug.json"
	}

	file, err := os.ReadFile(filepath.Join(curPath, configName))
	if err != nil {
		zap.L().Fatal("[InitConfig] 读取nacos配置文件失败:", zap.Error(err))
	}

	err = json.Unmarshal(file, nacos)
	if err != nil {
		zap.L().Fatal("[InitConfig] 解析nacos配置文件失败:", zap.Error(err))
	}

	// 读取nacos的配置文件之后，打印以下看看
	fmt.Printf("nacos配置文件：%#v\n", nacos)

	// 连接nacos获取配置
	cnf, err := getConfig()
	if err != nil {
		zap.L().Fatal("[InitConfig].[getConfig] with error:", zap.Error(err))
	}

	err = json.Unmarshal([]byte(cnf), &global.SrvConfig)
	if err != nil {
		zap.L().Fatal("[InitConfig].[json.Unmarshal] with error:", zap.Error(err))
	}
	zap.L().Info("[InitConfig] <UNK>", zap.Any("cnf", cnf))
}

func getConfig() (string, error) {
	clientConfig := constant.ClientConfig{
		NamespaceId:         nacos.Namespace, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
		Username:            nacos.Username, // Nacos服务端的API鉴权Username
		Password:            nacos.Password, // Nacos服务端的API鉴权Password
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      nacos.Host,
			ContextPath: "/nacos",
			Port:        nacos.Port,
			Scheme:      "http",
		},
	}

	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		return "", err
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: nacos.DataID,
		Group:  nacos.Group,
	})
	if err != nil {
		return "", err
	}

	return content, nil
}

func InitConfig2() {
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

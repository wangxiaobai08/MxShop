package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"user_srv/config"
	"user_srv/global"
)

// InitConfig
// @Description:利用viper和nacos初始化文件配置
//
//param:null
//return:null
func InitConfig() {
	configfilepath := fmt.Sprintf(global.FilePath.ConfigFile)
	viper.SetConfigFile(configfilepath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	// 定义配置结构体
	global.NacosConfig = &config.NacosConfig{}

	// 将读取到的配置数据反序列化到结构体
	if err := viper.Unmarshal(global.NacosConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
		panic(err)
	}

	// Nacos 配置
	serverConfig := []constant.ServerConfig{{
		IpAddr: global.NacosConfig.Host,
		Port:   uint64(global.NacosConfig.Port),
	},
	}

	nacosLogDir := fmt.Sprintf("%s/%s/%s", global.FilePath.LogFile, "nacos", "log")
	nacosCacheDir := fmt.Sprintf("%s/%s/%s", global.FilePath.LogFile, "nacos", "cache")
	clientConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果使用的是默认命名空间，保持为 public
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              nacosLogDir,
		CacheDir:            nacosCacheDir,
		LogLevel:            "debug",
	}

	// 创建 Nacos 客户端
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfig,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		log.Fatalf("Error creating Nacos client: %v", err)
		panic(err)
	}

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Dataid,
		Group:  global.NacosConfig.Group,
	})
	if err != nil {
		zap.S().Errorw("client.GetConfig读取文件失败", "err", err.Error())
		return
	}
	global.ServiceConfig = &config.ServiceConfig{}
	err = json.Unmarshal([]byte(content), global.ServiceConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println("nacos配置拉取成功")
	global.ServiceConfig.Host = "192.168.8.1"

}

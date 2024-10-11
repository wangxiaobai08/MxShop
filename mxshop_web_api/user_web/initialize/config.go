package initialize

import (
	"encoding/json"                                       // 用于处理JSON数据
	"fmt"                                                 // 格式化字符串
	"github.com/nacos-group/nacos-sdk-go/clients"         // Nacos客户端库，用于与Nacos服务通信
	"github.com/nacos-group/nacos-sdk-go/common/constant" // Nacos常量库，定义客户端和服务器的配置
	"github.com/nacos-group/nacos-sdk-go/vo"              // Nacos的参数定义，用于构造请求
	"github.com/spf13/viper"                              // Viper库，用于读取本地配置文件
	"go.uber.org/zap"                                     // Zap日志库，用于记录日志
	"mxshop_web_api/config"                               // 项目的配置结构体
	"mxshop_web_api/global"                               // 项目的全局变量
)

// InitConfig
// @Description: 初始化配置，加载Nacos远程配置并保存到全局配置中
func InitConfig() {
	// 获取本地配置文件名
	configFileName := fmt.Sprintf("%s", global.FileConfig.ConfigFile)

	v := viper.New()                // 新建一个Viper实例
	v.SetConfigFile(configFileName) // 设置要读取的配置文件
	err := v.ReadInConfig()         // 读取配置文件
	if err != nil {
		// 如果读取失败，记录错误并返回
		zap.S().Errorw("viper.ReadInConfig失败", "err", err.Error())
		return
	}

	// 将读取到的本地配置文件内容解析到global.NacosConfig中
	global.NacosConfig = &config.NacosConfig{}
	err = v.Unmarshal(global.NacosConfig) // 将配置文件映射到结构体
	if err != nil {
		zap.S().Errorw("viper unmarshal失败", "err", err.Error())
		return
	}
	// 输出Nacos的配置内容
	zap.S().Infof("%#v", global.NacosConfig)

	// 创建Nacos的服务器配置
	sConfig := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,         // Nacos服务器IP地址
			Port:   uint64(global.NacosConfig.Port), // Nacos服务器端口号
		},
	}

	// 设置Nacos日志和缓存目录
	nacosLogDir := fmt.Sprintf("%s/%s/%s", global.FileConfig.LogFile, "nacos", "log")     // 日志目录
	nacosCacheDir := fmt.Sprintf("%s/%s/%s", global.FileConfig.LogFile, "nacos", "cache") // 缓存目录

	// 创建Nacos的客户端配置
	cConfig := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // Nacos命名空间
		TimeoutMs:           5000,                         // 超时时间
		NotLoadCacheAtStart: true,                         // 不使用缓存启动
		LogDir:              nacosLogDir,                  // 日志文件路径
		CacheDir:            nacosCacheDir,                // 缓存文件路径
		LogLevel:            "debug",                      // 日志级别
	}

	// 创建Nacos配置客户端
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sConfig, // 服务器配置
		"clientConfig":  cConfig, // 客户端配置
	})
	if err != nil {
		// 如果客户端创建失败，记录错误并返回
		zap.S().Errorw("客户端连接失败", "err", err.Error())
		return
	}

	// 从Nacos服务器获取配置信息
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.Dataid, // 配置的DataId
		Group:  global.NacosConfig.Group,  // 配置的Group
	})
	if err != nil {
		// 如果读取失败，记录错误并返回
		zap.S().Errorw("client.GetConfig读取文件失败", "err", err.Error())
		return
	}

	// 将Nacos中获取的配置内容解析为global.WebServiceConfig结构体
	global.WebServiceConfig = &config.WebServiceConfig{}
	err = json.Unmarshal([]byte(content), global.WebServiceConfig) // 将JSON字符串解析为结构体
	if err != nil {
		// 如果解析失败，记录错误并返回
		zap.S().Errorw("读取的配置content解析到global.serviceConfig失败", "err", err.Error())
		return
	}

	// 成功拉取并解析配置后，记录配置信息
	zap.S().Infof("nacos配置拉取成功 %#v", global.WebServiceConfig)
}

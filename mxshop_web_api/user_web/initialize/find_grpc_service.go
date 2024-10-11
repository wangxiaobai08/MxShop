package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api" // 引入Consul客户端API，用于和Consul交互
	"go.uber.org/zap"                 // 引入Zap日志库，用于日志记录
	"mxshop_web_api/global"           // 项目全局配置包，包含Consul和服务信息
)

// InitUserService 用于初始化与user-service的连接，返回其地址
func InitUserService() string {
	// 获取Consul默认配置
	cfg := api.DefaultConfig()

	// 从全局配置中获取Consul的主机和端口信息
	consulConfig := global.WebServiceConfig.ConsulInfo
	cfg.Address = fmt.Sprintf("%s:%d", consulConfig.Host, consulConfig.Port) // 设置Consul地址

	// 定义变量用于存储user-service的地址和端口
	var userServiceHost string
	var userServicePort int

	// 创建一个新的Consul客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		// 如果创建客户端失败，记录错误并返回空字符串
		zap.S().Errorw("连接注册中心失败", "err", err.Error())
		return ""
	}

	// 使用过滤器查询服务名称为"user-service"的服务实例
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.WebServiceConfig.UserServiceInfo.Name))
	if err != nil {
		// 如果查询失败，记录错误并返回空字符串
		zap.S().Errorw("查询 user-service失败", "err", err.Error())
		return ""
	}

	// 遍历查询结果，取出第一个匹配的服务的地址和端口
	for _, value := range data {
		userServiceHost = value.Address
		userServicePort = value.Port
		break // 找到一个实例就退出循环
	}

	// 如果查询结果为空，说明服务未找到，记录错误并终止程序
	if userServiceHost == "" || userServicePort == 0 {
		zap.S().Fatal("InitRPC失败")
		return ""
	}

	// 记录查询到的服务地址和端口
	zap.S().Infof("查询到user-service %s:%d", userServiceHost, userServicePort)

	// 将地址和端口组合成服务目标地址
	target := fmt.Sprintf("%s:%d", userServiceHost, userServicePort)
	return target // 返回目标地址
}

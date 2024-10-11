package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop_web_api/config"
	"mxshop_web_api/proto"
)

var (
	WebServiceConfig *config.WebServiceConfig
	UserClient       proto.UserClient   // grpc客户端
	Translator       ut.Translator      // 翻译器
	FileConfig       *config.FileConfig // 文件配置
	NacosConfig      *config.NacosConfig
	Port             int
)

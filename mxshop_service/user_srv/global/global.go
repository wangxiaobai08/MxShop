package global

import (
	"gorm.io/gorm"
	"user_srv/config"
)

var (
	DB            *gorm.DB
	ServiceConfig *config.ServiceConfig
	FilePath      *config.FilePathConfig
	NacosConfig   *config.NacosConfig
)

package initialize

import (
	"fmt"
	"path"
	"runtime"
	"user_srv/config"
	"user_srv/global"
)

// InitFileAbsPath
// @Description: 初始化文件路径
//
//param:null
//return:null
func InitFileAbsPath() {
	basepath := GetCurrentAbsolutePath()
	global.FilePath = &config.FilePathConfig{
		ConfigFile: basepath + "/config-debug.yaml",
		LogFile:    basepath + "/log",
	}
	fmt.Println("文件路径初始化成功", basepath)
}

// GetCurrentAbsolutePath
// @Description: 获取调用方绝对路径地址
//
//param:null
//return:string
func GetCurrentAbsolutePath() string {
	var abspath string
	_, fileName, _, ok := runtime.Caller(2)
	if ok {
		abspath = path.Dir(fileName)
	}
	return abspath
}

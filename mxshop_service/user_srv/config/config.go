package config

// ConsulConfig
// @Description: consul配置
type ConsulConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// ServiceConfig
// @Description: 服务配置
type ServiceConfig struct {
	Name       string `json:"name"`
	Host       string
	ConsulInfo ConsulConfig `json:"consul"`
}

// FilePathConfig
// @Description: 文件路径配置
type FilePathConfig struct {
	ConfigFile string
	LogFile    string
}

// NacosConfig
// @Description: Nacos连接配置
type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Dataid    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

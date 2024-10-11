package config

type WebServiceConfig struct {
	Name            string            `json:"name"`
	JWTInfo         JwtConfig         `json:"jwt"`
	UserServiceInfo UserServiceConfig `json:"user_service_info"`
	JaegerInfo      JaegerConfig      `json:"jaeger_info"`
	RedisInfo       RedisConfig       `json:"redis"`
	AliSmsInfo      AliSmsConfig      `json:"aliyun_message"`
	ConsulInfo      ConsulConfig      `json:"consul"`
}

type JwtConfig struct {
	Signingkey string `json:"key"`
}

type UserServiceConfig struct {
	Name string `json:"name"`
}

type JaegerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

type AliSmsConfig struct {
	ApiKey       string `json:"key"`
	ApiSecret    string `json:"secret"`
	SignName     string `json:"signName"`
	TemplateCode string `json:"template_code"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type FileConfig struct {
	ConfigFile string
	LogFile    string
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	Dataid    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

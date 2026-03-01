package config

import (
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func setDefault(v *viper.Viper) {
	v.SetDefault(LEAVE_CONTAINERNAME, uuid.New().String())
	v.SetDefault(LEAVE_HERTZ_LISTENADDR, "0.0.0.0")
	v.SetDefault(LEAVE_HERTZ_LISTENPORT, "8080")
	v.SetDefault(LEAVE_HERTZ_MONITORPORT, "8081")
	v.SetDefault(LEAVE_PGSQL_ADDR, "pgsql")
	v.SetDefault(LEAVE_PGSQL_PORT, "5432")
	v.SetDefault(LEAVE_PGSQL_USER, "root")
	v.SetDefault(LEAVE_PGSQL_PASSWD, "")
	v.SetDefault(LEAVE_PGSQL_DBNAME, "leave")
	v.SetDefault(LEAVE_OAUTH_CLIENTID, "")
	v.SetDefault(LEAVE_OAUTH_CLIENTSECRET, "")
	v.SetDefault(LEAVE_OAUTH_REDIRECTURL, "")
}

func LoadConfig() (*Config, error) {
	// 初始化结构体指针
	conf := new(Config)
	// 初始化Viper
	v := viper.New()
	//这样环境变量需要以 LEAVE_ 开头，如 LEAVE_HERTZ_LISTENADDR
	v.SetEnvPrefix("LEAVE")
	// 加载默认配置
	setDefault(v)
	// 加载环境变量
	v.AutomaticEnv()
	//设置键名转换器（将环境变量中的 _ 映射到结构体的嵌套字段）
	//例如：LEAVE_HERTZ_LISTENADDR -> Hertz.ListenAddr
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 解析配置到结构体
	if err := v.Unmarshal(conf); err != nil {
		return nil, err
	}

	return conf, nil
}

package config

import (
	"github.com/spf13/viper"
	"github.com/teanft/ethscan/util"
	"time"
)

var Cfg = &Config{}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	//Timeout    time.Duration `mapstructure:"timeout"`
	//MaxRetries int           `mapstructure:"maxRetries"`
}

type ClientConfig struct {
	URL        string        `mapstructure:"url"`
	Timeout    time.Duration `mapstructure:"timeout"`
	MaxRetries int           `mapstructure:"maxRetries"`
	Sleeper    time.Duration `mapstructure:"sleeper"`
}

type AccountConfig struct {
	PrivateKey string `mapstructure:"privateKey"`
	To         string `mapstructure:"toAddress"`
}

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Client  ClientConfig  `mapstructure:"client"`
	Account AccountConfig `mapstructure:"account"`
}

func InitConfig(name, typ, path string) (*Config, error) {
	viper.SetConfigName(name)
	viper.SetConfigType(typ)
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, util.NewErr("读取config失败", err)
	}

	cfg := Cfg
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, util.NewErr("无法从struct中解码", err)
	}
	return cfg, nil
}

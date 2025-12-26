package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	WeChat   WeChatConfig   `mapstructure:"wechat"`
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	TimeZone string `mapstructure:"timezone"`
}

type WeChatConfig struct {
	AppID      string `mapstructure:"appid"`
	AppSecret  string `mapstructure:"appsecret"`
	TemplateID string `mapstructure:"template_id"`
	MchID      string `mapstructure:"mch_id"`
	MchKey     string `mapstructure:"mch_key"`
	NotifyURL  string `mapstructure:"notify_url"`
}

func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// TODO: MCH_ID, MCH_KEY, NOTIFY_URL
	// Manually bind specific environment variables
	_ = viper.BindEnv("wechat.appid", "APP_ID")
	_ = viper.BindEnv("wechat.appsecret", "APP_SECRET")
	_ = viper.BindEnv("wechat.template_id", "TEMPLATE_ID")
	_ = viper.BindEnv("wechat.mch_id", "MCH_ID")
	_ = viper.BindEnv("wechat.mch_key", "MCH_KEY")
	_ = viper.BindEnv("wechat.notify_url", "NOTIFY_URL")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

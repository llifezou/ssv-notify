package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type Config struct {
	Network        string         `json:"network"`
	LarkConfig     LarkConfig     `json:"larkconfig"`
	TelegramConfig TelegramConfig `json:"telegramconfig"`
	Aim            string         `json:"aim"`
	ClusterOwner   []string       `json:"clusterowner"`
}

type LarkConfig struct {
	WebHook string `json:"webhook"`
}

type TelegramConfig struct {
	AccessToken string `json:"accesstoken"`
	ChatId      string `json:"chatid"`
}

var conf Config

func GetConfig() Config {
	return conf
}

func Init(p string) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	if p != "" {
		v.AddConfigPath(filepath.Dir(p))
	} else {
		v.AddConfigPath("./config")
		v.AddConfigPath("../config")
	}
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = v.Unmarshal(&conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if conf.LarkConfig.WebHook == "" && conf.TelegramConfig.AccessToken == "" || conf.TelegramConfig.ChatId == "" {
		fmt.Println("At least configure lark or telegram")
		os.Exit(1)
	}
}

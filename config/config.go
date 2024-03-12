package config

import (
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var log = logging.Logger("config")

type Config struct {
	Network            string                   `json:"network"`
	EthRpc             string                   `json:"ethrpc"`
	LarkConfig         LarkConfig               `json:"larkconfig"`
	TelegramConfig     TelegramConfig           `json:"telegramconfig"`
	GmailConfig        GmailConfig              `json:"gmailconfig"`
	DiscordConfig      DiscordConfig            `json:"discordconfig"`
	OperatorMonitor    OperatorMonitorConfig    `json:"operatormonitor"`
	LiquidationMonitor LiquidationMonitorConfig `json:"liquidationmonitor"`
}

type OperatorMonitorConfig struct {
	Aim          string   `json:"aim"`
	ClusterOwner []string `json:"clusterowner"`
}

type LiquidationMonitorConfig struct {
	Threshold    uint64   `json:"threshold"`
	ClusterOwner []string `json:"clusterowner"`
}

type LarkConfig struct {
	WebHook string `json:"webhook"`
}

type TelegramConfig struct {
	AccessToken string `json:"accesstoken"`
	ChatId      string `json:"chatid"`
}

type GmailConfig struct {
	From     string `json:"from"`
	Password string `json:"password"`
	To       string `json:"to"`
}

type DiscordConfig struct {
	WebHook string `json:"webhook"`
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

	if conf.LarkConfig.WebHook == "" &&
		(conf.TelegramConfig.AccessToken == "" || conf.TelegramConfig.ChatId == "") &&
		conf.GmailConfig.Password == "" &&
		conf.DiscordConfig.WebHook == "" {
		log.Warn("At least configure lark or telegram or gmail or discord")
		os.Exit(1)
	}

	if conf.LiquidationMonitor.Threshold < 10 {
		log.Warnw("config LiquidationMonitor.Threshold is too low, use the default", "default threshold", "10 day")
		conf.LiquidationMonitor.Threshold = 10
	}
}

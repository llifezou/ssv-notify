package ssv

import (
	"fmt"
	"github.com/llifezou/ssv-notify/config"
	"github.com/llifezou/ssv-notify/notify"
	"github.com/llifezou/ssv-notify/notify/lark"
	"github.com/llifezou/ssv-notify/notify/telegram"
	"os"
	"testing"
)

func TestMonitor(t *testing.T) {
	config.Init("../config/config.yaml")
	conf := config.GetConfig()
	var senders []notify.Sender
	if conf.LarkConfig.WebHook != "" {
		senders = append(senders, lark.NewLarkClient(conf.LarkConfig.WebHook))
	}
	if conf.TelegramConfig.AccessToken != "" && conf.TelegramConfig.ChatId != "" {
		senders = append(senders, telegram.NewTelegramClient(conf.TelegramConfig.AccessToken, conf.TelegramConfig.ChatId))
	}

	notifyClient, err := notify.NewNotify(senders...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	monitor(notifyClient, "", "holesky", conf.ClusterOwner)
}

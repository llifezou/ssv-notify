package notify

import (
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/llifezou/ssv-notify/config"
	"github.com/llifezou/ssv-notify/notify/discord"
	"github.com/llifezou/ssv-notify/notify/gmail"
	"github.com/llifezou/ssv-notify/notify/lark"
	"github.com/llifezou/ssv-notify/notify/telegram"
)

var log = logging.Logger("notify")

type Sender interface {
	Send(msg string) error
	Platform() string
}

type Notify struct {
	Senders []Sender
}

func NewSender() ([]Sender, error) {
	conf := config.GetConfig()
	var senders []Sender
	if conf.LarkConfig.WebHook != "" {
		senders = append(senders, lark.NewLarkClient(conf.LarkConfig.WebHook))
	}
	if conf.TelegramConfig.AccessToken != "" && conf.TelegramConfig.ChatId != "" {
		senders = append(senders, telegram.NewTelegramClient(conf.TelegramConfig.AccessToken, conf.TelegramConfig.ChatId))
	}
	if conf.GmailConfig.Password != "" && conf.GmailConfig.From != "" && conf.GmailConfig.To != "" {
		senders = append(senders, gmail.NewGmailClient(conf.GmailConfig.From, conf.GmailConfig.Password, conf.GmailConfig.To))
	}
	if conf.DiscordConfig.WebHook != "" {
		senders = append(senders, discord.NewDiscordClient(conf.DiscordConfig.WebHook))
	}

	if len(senders) == 0 {
		return nil, fmt.Errorf("no alarm configuration")
	}
	return senders, nil
}

func NewNotify() (*Notify, error) {
	senders, err := NewSender()
	if err != nil {
		return nil, err
	}

	n := &Notify{
		Senders: senders,
	}

	return n, nil
}

func (n *Notify) Send(msg string) []error {
	var errs []error
	for _, sender := range n.Senders {
		name := sender.Platform()
		err := sender.Send(msg)
		if err != nil {
			werr := fmt.Errorf("Notification failed! platform: %s, err: %s \n", name, err)
			errs = append(errs, werr)
			log.Warn(werr.Error())
		} else {
			log.Infow("Alarm sent successfully", "platform", name)
		}
	}
	return errs
}

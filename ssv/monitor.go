package ssv

import (
	"fmt"
	"github.com/llifezou/ssv-notify/config"
	"github.com/llifezou/ssv-notify/notify"
	"github.com/llifezou/ssv-notify/notify/lark"
	"github.com/llifezou/ssv-notify/notify/telegram"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func StartMonitor(shutdown <-chan struct{}) {
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

	ticker := time.NewTicker(32 * 12 * time.Second) // 1 epoch
	for {
		select {
		case <-ticker.C:
			log.Println("Monitoring trigger")
			monitor(notifyClient, conf.Aim, conf.Network, conf.ClusterOwner)
		case <-shutdown:
			return
		}
	}
}

func monitor(notify *notify.Notify, aim string, network string, clusterOwner []string) {
	var reportedOperatorId = make(map[int]struct{})

	for _, owner := range clusterOwner {
		clusterValidators, err := GetClusterValidators(network, owner)
		if err != nil {
			msg := fmt.Sprintf("ssv api request failed, GetClusterValidators: %s", err.Error())
			log.Println(msg)
			notify.Send(msg)
			continue
		}

		// check from ssv
		for _, validator := range clusterValidators.Validators {
			validatorDuties, err := GetValidatorDuties(network, validator.PublicKey)
			if err != nil {
				msg := fmt.Sprintf("ssv api request failed, GetValidatorDuties: %s", err.Error())
				log.Println(msg)
				notify.Send(msg)
				continue
			}

			badOperator, name := CheckDuty(validatorDuties.Duties)
			for i, opId := range badOperator {
				if aim != "" && aim != "all" {
					if !strings.Contains(aim, strconv.Itoa(opId)) {
						continue
					}
				}
				if _, ok := reportedOperatorId[opId]; ok {
					continue
				}
				reportedOperatorId[opId] = struct{}{}
				msg := fmt.Sprintf("[Data From SSV API]: OperatorId: %d (name: %s) inactive in epech: %d !!!", opId, name[i], validatorDuties.Duties[0].Epoch)
				log.Println(msg)
				notify.Send(msg)
			}
		}

		// check from ssvscan
		if len(clusterValidators.Validators) > 0 {
			baseMsg := "[Data From SSVScan API]: "
			willReport := make(map[int]string)
			for _, operator := range clusterValidators.Validators[0].Operators {
				opId := operator.ID
				msg := ""
				status, err := GetOperatorStatusFromSSVScan(network, opId)
				if err != nil {
					msg = fmt.Sprintf("ssvscan api request failed, err: %s", err.Error())
					log.Println(msg)
					notify.Send(msg)
					continue
				}

				if !status {
					msg = baseMsg + fmt.Sprintf("OperatorId: %d (name: %s) inactive", opId, operator.Name)
				}

				if msg == "" {
					log.Println(baseMsg + fmt.Sprintf("OperatorId: %d (name: %s) active", opId, operator.Name))
					continue
				}

				if aim != "" && aim != "all" {
					if !strings.Contains(aim, strconv.Itoa(opId)) {
						continue
					}
				}

				if _, ok := reportedOperatorId[opId]; ok {
					continue
				}

				reportedOperatorId[opId] = struct{}{}
				willReport[opId] = msg
				log.Println(msg)
			}

			// If they are all inactive, ssvscan data may be broken.
			if len(willReport) != len(clusterValidators.Validators[0].Operators) {
				for _, msg := range willReport {
					notify.Send(msg)
				}
			} else {
				log.Println("All operators will be reported")
			}
		}
	}
}

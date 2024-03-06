package operator

import (
	"fmt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/llifezou/ssv-notify/config"
	"github.com/llifezou/ssv-notify/notify"
	"os"
	"strconv"
	"strings"
	"time"
)

var log = logging.Logger("operator-monitor")

func StartOperatorMonitor(shutdown <-chan struct{}) {
	conf := config.GetConfig()

	notifyClient, err := notify.NewNotify()
	if err != nil {
		log.Warn(err)
		os.Exit(1)
	}

	ticker := time.NewTicker(32 * 12 * time.Second) // 1 epoch
	for {
		select {
		case <-ticker.C:
			log.Info("Monitoring trigger")
			monitor(notifyClient, conf.OperatorMonitor.Aim, conf.Network, conf.OperatorMonitor.ClusterOwner)
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
			log.Warn(msg)
			continue
		}

		// check from ssv
		for _, validator := range clusterValidators.Validators {
			validatorDuties, err := GetValidatorDuties(network, validator.PublicKey)
			if err != nil {
				msg := fmt.Sprintf("ssv api request failed, GetValidatorDuties: %s", err.Error())
				log.Warn(msg)
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
				log.Warn(msg)
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
					log.Warn(fmt.Sprintf("ssvscan api request failed, err: %s", err.Error()))
					continue
				}

				if !status {
					msg = baseMsg + fmt.Sprintf("OperatorId: %d (name: %s) inactive", opId, operator.Name)
				}

				if msg == "" {
					log.Info(baseMsg + fmt.Sprintf("OperatorId: %d (name: %s) active", opId, operator.Name))
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
				log.Warn(msg)
			}

			// If they are all inactive, ssvscan data may be broken.
			if len(willReport) != len(clusterValidators.Validators[0].Operators) {
				for _, msg := range willReport {
					notify.Send(msg)
				}
			} else {
				log.Warn("All operators will be reported, Maybe it’s a third-party data error")
			}
		}
	}
}

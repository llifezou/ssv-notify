package operator

import (
	"encoding/json"
	"fmt"
	"github.com/llifezou/ssv-notify/ssv/utils"
)

var ssvScanUrl = "https://%sssvscan.io/api/v1/operator/status?id=%d"

type OperatorStatus struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status           string `json:"status"`
		SuccessfulDuties string `json:"successful_duties"`
		FailedDuties     int    `json:"failed_duties"`
		TotalDuties      string `json:"total_duties"`
	} `json:"data"`
}

func GetOperatorStatusFromSSVScan(network string, operatorId int) (bool, error) {
	var operatorStatus *OperatorStatus
	if network == "mainnet" {
		network = ""
	} else {
		network = network + "."
	}

	url := fmt.Sprintf(ssvScanUrl, network, operatorId)
	b, err := utils.HttpGet(url)
	if err != nil {
		return false, err
	}

	if err = json.Unmarshal(b, &operatorStatus); err != nil {
		return false, err
	}

	return operatorStatus.Data.Status == "active", nil
}

package ssv

import (
	"encoding/json"
	"fmt"
	"log"
)

var (
	ssvClusterValidatorsUrl = "https://api.ssv.network/api/v4/%s/validators?page=1&perPage=10&ownerAddress=%s"
	ssvValidatorDutiesUrl   = "https://api.ssv.network/api/v4/%s/duties/%s?page=1&perPage=1"
)

type ClusterValidatorsInfo struct {
	Pagination struct {
		Total   int `json:"total"`
		Page    int `json:"page"`
		Pages   int `json:"pages"`
		PerPage int `json:"per_page"`
	} `json:"pagination"`
	Validators []struct {
		PublicKey        string `json:"public_key"`
		Cluster          string `json:"cluster"`
		OwnerAddress     string `json:"owner_address"`
		Status           string `json:"status"`
		IsValid          bool   `json:"is_valid"`
		IsDeleted        bool   `json:"is_deleted"`
		IsPublicKeyValid bool   `json:"is_public_key_valid"`
		IsSharesValid    bool   `json:"is_shares_valid"`
		IsOperatorsValid bool   `json:"is_operators_valid"`
		Operators        []struct {
			ID               int    `json:"id"`
			IDStr            string `json:"id_str"`
			DeclaredFee      string `json:"declared_fee"`
			PreviousFee      string `json:"previous_fee"`
			Fee              string `json:"fee"`
			PublicKey        string `json:"public_key"`
			OwnerAddress     string `json:"owner_address"`
			AddressWhitelist string `json:"address_whitelist"`
			Location         string `json:"location"`
			SetupProvider    string `json:"setup_provider"`
			Eth1NodeClient   string `json:"eth1_node_client"`
			Eth2NodeClient   string `json:"eth2_node_client"`
			MevRelays        string `json:"mev_relays"`
			Description      string `json:"description"`
			WebsiteURL       string `json:"website_url"`
			TwitterURL       string `json:"twitter_url"`
			LinkedinURL      string `json:"linkedin_url"`
			DkgAddress       string `json:"dkg_address"`
			Logo             string `json:"logo"`
			Type             string `json:"type"`
			Name             string `json:"name"`
			Performance      any    `json:"performance"`
			IsValid          bool   `json:"is_valid"`
			IsDeleted        bool   `json:"is_deleted"`
			IsActive         int    `json:"is_active"`
			Status           string `json:"status"`
			ValidatorsCount  int    `json:"validators_count"`
			Version          string `json:"version"`
			Network          string `json:"network"`
		} `json:"operators"`
		ValidatorInfo struct {
		} `json:"validator_info"`
		Version string `json:"version"`
		Network string `json:"network"`
	} `json:"validators"`
}

func GetClusterValidators(network, clusterOwner string) (*ClusterValidatorsInfo, error) {
	var clusterValidators *ClusterValidatorsInfo

	url := fmt.Sprintf(ssvClusterValidatorsUrl, network, clusterOwner)
	b, err := httpGet(url)

	if err = json.Unmarshal(b, &clusterValidators); err != nil {
		return nil, err
	}

	return clusterValidators, nil
}

type ValidatorDutiesInfo struct {
	Pagination struct {
		Total   int `json:"total"`
		Page    int `json:"page"`
		Pages   int `json:"pages"`
		PerPage int `json:"per_page"`
	} `json:"pagination"`
	Duties []Duty `json:"duties"`
}
type Duty struct {
	PublicKey string `json:"publicKey"`
	Operators []struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"operators"`
	MissingOperators []any  `json:"missing_operators"`
	Slot             any    `json:"slot"`
	Epoch            int    `json:"epoch"`
	Duty             any    `json:"duty"`
	Status           string `json:"status"`
	Sequence         any    `json:"sequence"`
}

func GetValidatorDuties(network, pubKey string) (*ValidatorDutiesInfo, error) {
	var validatorDutiesInfo *ValidatorDutiesInfo

	url := fmt.Sprintf(ssvValidatorDutiesUrl, network, pubKey)
	b, err := httpGet(url)

	if err = json.Unmarshal(b, &validatorDutiesInfo); err != nil {
		return nil, err
	}

	return validatorDutiesInfo, nil
}

func CheckDuty(duties []Duty) ([]int, []string) {
	var badOperator []int
	var name []string
	for _, duty := range duties {
		for _, operator := range duty.Operators {
			if operator.Status != "success" {
				badOperator = append(badOperator, operator.ID)
				name = append(name, operator.Name)
			} else {
				log.Println(fmt.Sprintf("[Data From SSV API]: OperatorId: %d active", operator.ID))
			}
		}
	}
	return badOperator, name
}

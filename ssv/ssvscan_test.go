package ssv

import "testing"

func TestGetOperatorStatusFromSSVScan(t *testing.T) {
	status, err := GetOperatorStatusFromSSVScan("mainnet", 23)
	if err != nil {
		t.Fatal(err)
	}
	status2, err := GetOperatorStatusFromSSVScan("holesky", 23)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(status)
	t.Log(status2)
}

package operator

import "testing"

func TestGetClusterValidators(t *testing.T) {
	info, err := GetClusterValidators("holesky", "0x344152eD7110694B004962CD61ddA876559Fd8a4")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
}

func TestGetValidatorDuties(t *testing.T) {
	info, err := GetValidatorDuties("holesky", "af22f5801dd7bf1ee421d45962089e8997a2c2fe38370299831b2ba9a94eb5a61570ee972593c41c2304d238faa7f352")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(info)
	_, names := CheckDuty(info.Duties)
	for _, b := range names {
		t.Log(b)
	}

}

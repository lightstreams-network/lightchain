package web3

import (
	"testing"
)

func TestParseWei(t *testing.T) {
	var values = []struct {
		wei string
	}{
		// 0.01 PHT
		{"10000000000000000"},
		// 0.1 PHT
		{"100000000000000000"},
		// 1 PHT
		{"1000000000000000000"},
		// 105M PHT
		{"105000000000000000000000000"},
		// 300M PHT
		{"300000000000000000000000000"},
	}

	for _, v := range values {
		weiBN, err := ParseWei(v.wei)
		if err != nil {
			t.Fatal(err)
		}

		if weiBN.String() != v.wei {
			t.Errorf("unable to parse %v Wei. Actual: %s", v.wei, weiBN.String())
		}
	}
}
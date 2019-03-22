package distribution

import (
	"testing"
	"os"
	"path"
	"github.com/ethereum/go-ethereum/common"
)

const csvFixtureFileName = "csv_distribution_test.csv"
const csvFixtureEmptyFileName = "csv_distribution_empty_test.csv"
const csvFixtureWrongHeaderFileName = "csv_distribution_wrong_header_test.csv"
const csvFixtureWrongBodyFileName = "csv_distribution_wrong_body_test.csv"

func TestParseCsvDistributions(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	address2distributionsCount, distributions, err := parseCsvDistributions(path.Join(cwd, csvFixtureFileName))
	if err != nil {
		t.Fatal(err)
	}

	if len(distributions) != 4 {
		t.Fatal(err)
	}

	if address2distributionsCount[common.HexToAddress("0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e")] != 3 {
		t.Fatal("0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e should have 3 distribution to perform")
	}

	if address2distributionsCount[common.HexToAddress("0x336f959c88b6c66f952859a0a53be2f5e0b152cf")] != 1 {
		t.Fatal("0x336f959c88b6c66f952859a0a53be2f5e0b152cf should have 1 distribution to perform")
	}
}

func TestParseEmptyCsvDistributions(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	address2distributionsCount, distributions, err := parseCsvDistributions(path.Join(cwd, csvFixtureEmptyFileName))
	if err != nil {
		t.Fatal(err)
	}

	if len(distributions) != 0 {
		t.Fatal(err)
	}

	if len(address2distributionsCount) != 0 {
		t.Fatal("no addresses should have been found")
	}
}

func TestCsvHasWrongHeader(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	_, _, err = parseCsvDistributions(path.Join(cwd, csvFixtureWrongHeaderFileName))
	if err == nil {
		t.Fatal("an err should be returned")
	}
}

func TestCsvHasWrongBody(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	_, _, err = parseCsvDistributions(path.Join(cwd, csvFixtureWrongBodyFileName))
	if err == nil {
		t.Fatal("an err should be returned")
	}
}
package distribution

import (
	"testing"
	"os"
	"path"
)

const csvFixtureFileName = "csv_distribution_test.csv"
const csvWrongHeaderFileName = "csv_distribution_wrong_header_test.csv"
const csvWrongBodyFileName = "csv_distribution_wrong_body_test.csv"

func TestParseCsvDistributions(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	distributions, err := parseCsvDistributions(path.Join(cwd, csvFixtureFileName))
	if err != nil {
		t.Fatal(err)
	}

	if len(distributions) != 4 {
		t.Fatal(err)
	}
}

func TestCsvHasWrongHeader(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	_, err = parseCsvDistributions(path.Join(cwd, csvWrongHeaderFileName))
	if err == nil {
		t.Fatal(err)
	}
}

func TestCsvHasWrongBody(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	_, err = parseCsvDistributions(path.Join(cwd, csvWrongBodyFileName))
	if err == nil {
		t.Fatal(err)
	}
}
package distribution

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"encoding/csv"
	"path/filepath"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/lightstreams-network/lightchain/fs"
)

const csvHeaderColumnsCount = 3
const csvHeaderColumnFrom = "from"
const csvHeaderColumnTo = "to"
const csvHeaderColumnAmountWei = "amount_wei"

func parseCsvDistributions(csvFilePath string) (map[common.Address]int, []distribution, error) {
	records, err := parseCsvRecords(csvFilePath)
	if err != nil {
		return nil, nil, err
	}

	address2distributionsCount := make(map[common.Address]int)
	distributions := make([]distribution, len(records)-1)
	for i := 0; i < len(records); i++ {
		if i == 0 {
			if err := isValidCsvHeader(records[i]); err != nil {
				return nil, nil, err
			}

			continue
		}

		distribution, err := parseCsvRow(records[i])
		if err != nil {
			return nil, nil, err
		}

		address2distributionsCount[distribution.from]++
		distributions[i-1] = distribution
	}

	return address2distributionsCount, distributions, nil
}

func parseCsvRecords(csvFilePath string) ([][]string, error) {
	if err := isValidCsv(csvFilePath); err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(csvFilePath)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(bytes.NewReader(content))
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV distribution file is empty")
	}

	return records, nil
}

func isValidCsv(csvFilePath string) error {
	if !fs.FileExist(csvFilePath) {
		return fmt.Errorf("unable to find the distribution CSV file")
	}

	if filepath.Ext(csvFilePath) != ".csv" {
		return fmt.Errorf("the tokens distribution source file must be a CSV file")
	}

	return nil
}

func isValidCsvHeader(header []string) error {
	if len(header) != 3 {
		return fmt.Errorf("expected %v header columns", csvHeaderColumnsCount)
	}

	if header[0] != csvHeaderColumnFrom {
		return fmt.Errorf("first header column should be '%s'", csvHeaderColumnFrom)
	}

	if header[1] != csvHeaderColumnTo {
		return fmt.Errorf("second header column should be '%s'", csvHeaderColumnTo)
	}

	if header[2] != csvHeaderColumnAmountWei {
		return fmt.Errorf("third header column should be '%s'", csvHeaderColumnAmountWei)
	}

	return nil
}

func parseCsvRow(row []string) (distribution, error) {
	if len(row) != 3 {
		return distribution{}, fmt.Errorf("expected %v row columns", csvHeaderColumnsCount)
	}

	from := row[0]
	to := row[1]
	amountWei := row[2]

	if !common.IsHexAddress(from) {
		return distribution{}, fmt.Errorf("'%s' is invalid hex address", from)
	}

	if !common.IsHexAddress(to) {
		return distribution{}, fmt.Errorf("'%s' is invalid hex address", to)
	}

	amount, ok := math.ParseBig256(amountWei)
	if !ok {
		return distribution{}, fmt.Errorf("unable to parse '%s' amount in wei to big int", amountWei)
	}

	if amount.String() != amountWei {
		return distribution{}, fmt.Errorf("the CSV amount '%s' was incorrectly parsed to '%s'", amountWei, amount.String())
	}

	return newDistribution(common.HexToAddress(from), common.HexToAddress(to), amount), nil
}
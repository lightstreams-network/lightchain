package whitelist
//
//import (
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/lightstreams-network/lightchain/consensus"
//)
//
//type AccountWhitelist struct {
//	accounts map[common.Address]common.Address
//}
//
//func NewAccountWhitelist(cfg consensus.Config) AccountWhitelist {
//	accounts := make(map[common.Address]common.Address, len(cfg.DeploymentWhitelist))
//
//	for _, account := range cfg.DeploymentWhitelist {
//		address := common.HexToAddress(account)
//		accounts[address] = address
//	}
//
//	return AccountWhitelist{accounts: accounts}
//}
//
//func (aw AccountWhitelist) HasAccount(address common.Address) bool {
//	_, found := aw.accounts[address];
//
//	return found
//}

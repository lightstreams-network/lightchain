package deployment
//
//import (
//	"github.com/ethereum/go-ethereum/core/types"
//	"github.com/lightstreams-network/lightchain/consensus/transaction/validator/deployment/whitelist"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/lightstreams-network/lightchain/node"
//)
//
//type Validator struct {
//	accountWhitelist whitelist.AccountWhitelist
//}
//
//func New(cfg node.Config) (Validator) {
//	return Validator{whitelist.NewAccountWhitelist(cfg)}
//}
//
//func (v Validator) IsValid(tx types.Transaction) (bool, error) {
//	sender, err := extractSender(&tx)
//	if err != nil {
//		return false, err
//	}
//
//	return v.accountWhitelist.HasAccount(sender), nil
//}
//
//// Extracts sender from ethereum transaction
//func extractSender(tx *types.Transaction) (common.Address, error) {
//	var signer types.Signer = types.FrontierSigner{}
//	if tx.Protected() {
//		signer = types.NewEIP155Signer(tx.ChainId())
//	}
// 	// Make sure the transaction is signed properly
//	return types.Sender(signer, tx)
//}

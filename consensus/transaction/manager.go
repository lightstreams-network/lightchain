package transaction
//
//import (
//	"github.com/lightstreams-network/lightchain/consensus/transaction/validator"
//	"github.com/ethereum/go-ethereum/core/types"
//	"github.com/lightstreams-network/lightchain/node"
//)
//
//type TxHandler interface {
//	IsValid(tx types.Transaction) (bool, error)
//}
//
//type TxManager struct {
//	validator validator.Validator
//}
//
//func NewManager(cfgPath string) (TxManager, error) {
//	cfg, err :=  node.NewLightchainConfig(cfgPath)
//	if err != nil {
//		return TxManager{}, err
//	}
//
//	return TxManager{validator.New(cfg)}, nil
//}
//
//func (txm TxManager) IsValid(tx types.Transaction) (bool, error) {
//	return txm.validator.IsValid(tx)
//}

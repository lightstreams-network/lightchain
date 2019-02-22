package authy

import (
	"crypto/ecdsa"
	"io/ioutil"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/accounts"
)

// EthAccount is a wrapper around ETH common.Address.
type EthAccount struct {
	address common.Address
}

// User responsible for signing TXs.
type User struct {
	account EthAccount
	privKey *ecdsa.PrivateKey
	pubKey  *ecdsa.PublicKey
}

// NewEthAccount creates a new instance of EthAccount which is a wrapper around common.Address.
func NewEthAccount(address common.Address) EthAccount {
	return EthAccount{address}
}

// NewEthAccountFromHex creates a new instance of EthAccount which is a wrapper around common.Address, from hex format.
func NewEthAccountFromHex(address string) EthAccount {
	return EthAccount{common.HexToAddress(address)}
}

// GenerateEthAccount creates a new ethereum address and stores it encrypted in the given keystoreDir.
func GenerateEthAccount(keystoreDir string, pwd string) (EthAccount, error) {
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.NewAccount(pwd)
	if err != nil {
		return EthAccount{}, err
	}

	return NewEthAccount(acc.Address), nil
}

// NewFromKeystoreFile constructs a new instance of User from GETH keystore, populates the ETH addr and unlocks priv key with PWD.
func NewFromKeystoreFile(keystoreFilePath string, pwd string) (User, error) {
	keyJSON, err := ioutil.ReadFile(keystoreFilePath)
	if err != nil {
		return User{}, err
	}

	key, err := keystore.DecryptKey(keyJSON, pwd)
	if err != nil {
		return User{}, err
	}

	pubKey := key.PrivateKey.Public()
	pubKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return User{}, err
	}

	return User{
		NewEthAccount(key.Address),
		key.PrivateKey,
		pubKeyECDSA,
	}, nil
}

// FindInKeystoreDir searches the given `keystoreDir` for `account`, decrypts it with the given `pwd` and returns a User.
func FindInKeystoreDir(keystoreDir string, account EthAccount, pwd string) (User, error) {
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	ksAcc := accounts.Account{Address: account.Addr()}

	userKeystore, err := ks.Find(ksAcc)
	if err != nil {
		return User{}, fmt.Errorf("unable to generate token due to an error locating account's keystore file. %s", err.Error())
	}

	user, err := NewFromKeystoreFile(userKeystore.URL.Path, pwd)
	if err != nil {
		return User{}, fmt.Errorf("unable to generate token due to an error constructing leth user from keystore. %s", err.Error())
	}

	return user, nil
}

// Addr returns the user's eth address in common.Address format.
func (a EthAccount) Addr() common.Address {
	return a.address
}

// String returns the user's eth address (common.Address) in string, hex format.
func (a EthAccount) String() string {
	return a.address.Hex()
}

// PrivKey returns user ecdsa.PrivateKey.
func (u User) PrivKey() *ecdsa.PrivateKey {
	return u.privKey
}

// EthAccount returns user's EthAccount.
func (u User) EthAccount() EthAccount {
	return u.account
}

// EthAccountAddress returns user's EthAccount addr.
func (u User) EthAccountAddress() common.Address {
	return u.account.address
}
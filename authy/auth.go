package authy

import (
	"crypto/ecdsa"
	"io/ioutil"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/accounts"
)

// Auth responsible for signing TXs.
type Auth struct {
	address common.Address
	privKey *ecdsa.PrivateKey
	pubKey  *ecdsa.PublicKey
}

// NewFromKeystoreFile constructs a new instance of Auth from GETH keystore, populates the ETH addr and unlocks priv key with PWD.
func NewFromKeystoreFile(keystoreFilePath string, pwd string) (Auth, error) {
	keyJSON, err := ioutil.ReadFile(keystoreFilePath)
	if err != nil {
		return Auth{}, err
	}

	key, err := keystore.DecryptKey(keyJSON, pwd)
	if err != nil {
		return Auth{}, err
	}

	pubKey := key.PrivateKey.Public()
	pubKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return Auth{}, err
	}

	return Auth{
		key.Address,
		key.PrivateKey,
		pubKeyECDSA,
	}, nil
}

// FindInKeystoreDir searches the given `keystoreDir` for `address`, decrypts it with the given `pwd` and returns a Auth.
func FindInKeystoreDir(keystoreDir string, account common.Address, pwd string) (Auth, error) {
	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
	ksAcc := accounts.Account{Address: account}

	userKeystore, err := ks.Find(ksAcc)
	if err != nil {
		return Auth{}, fmt.Errorf("unable to generate token due to an error locating address's keystore file. %s", err.Error())
	}

	auth, err := NewFromKeystoreFile(userKeystore.URL.Path, pwd)
	if err != nil {
		return Auth{}, fmt.Errorf("unable to generate token due to an error constructing auth from keystore. %s", err.Error())
	}

	return auth, nil
}

func (u Auth) PrivKey() *ecdsa.PrivateKey {
	return u.privKey
}

func (u Auth) Address() common.Address {
	return u.address
}
package solana

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"

	"github.com/anyproto/go-slip10"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	"github.com/sirupsen/logrus"
)

type SolanaAdapter struct {
	logger *logrus.Logger
}

func NewSolanaAdapter(logger *logrus.Logger) *SolanaAdapter {
	return &SolanaAdapter{logger: logger}
}

func (adapter *SolanaAdapter) CanDo(coinType uint) bool {
	return coinType == 501
}

func (adapter *SolanaAdapter) deriveKeysForPath(seed []byte, derivationPath string) (ed25519.PrivateKey, ed25519.PublicKey, error) {
	node, err := slip10.DeriveForPath(derivationPath, seed)
	if err != nil {
		return nil, nil, err
	}

	pub, prv := node.Keypair()
	return prv, pub, nil
}

func (adapter *SolanaAdapter) DerivePrivateKey(seed []byte, derivationPath string, isDev bool) (string, error) {
	privKey, _, err := adapter.deriveKeysForPath(seed, derivationPath)
	if err != nil {
		return "", err
	}
	return base58.Encode(privKey), nil
}

func (adapter *SolanaAdapter) DerivePublicKey(seed []byte, derivationPath string, isDev bool) (string, error) {
	_, pubKey, err := adapter.deriveKeysForPath(seed, derivationPath)
	if err != nil {
		return "", err
	}
	return base58.Encode(pubKey), nil
}

func (adapter *SolanaAdapter) DeriveAddress(seed []byte, derivationPath string, isDev bool) (string, error) {
	_, pubKey, err := adapter.deriveKeysForPath(seed, derivationPath)
	if err != nil {
		return "", err
	}
	return base58.Encode(pubKey), nil
}

func (adapter *SolanaAdapter) CreateSignedTransaction(seed []byte, derivationPath string, payload string) (string, error) {
	privKey, err := adapter.DerivePrivateKey(seed, derivationPath, false)
	if err != nil {
		return "", err
	}
	senderPrvkey, err := solana.PrivateKeyFromBase58(privKey)
	if err != nil {
		return "", err
	}

	message, err := hex.DecodeString(payload)
	if err != nil {
		return "", err
	}

	tx, err := solana.TransactionFromBytes(message)
	if err != nil {
		return "", err
	}

	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			return &senderPrvkey
		},
	)

	signedTxBytes, err := tx.MarshalBinary()
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signedTxBytes), nil
}

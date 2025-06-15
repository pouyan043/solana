package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/sirupsen/logrus"

	soladapter "solana/solanaadapter"
)

func createTransferInstruction(from, to solana.PublicKey, lamports uint64) solana.Instruction {
	data := make([]byte, 9)
	data[0] = 2 
	binary.LittleEndian.PutUint64(data[1:], lamports)

	return solana.NewInstruction(
		solana.SystemProgramID,
		solana.AccountMetaSlice{
			{PublicKey: from, IsSigner: true, IsWritable: true},
			{PublicKey: to, IsSigner: false, IsWritable: true},
		},
		data,
	)
}

func main() {
	logger := logrus.New()
	adapter := soladapter.NewSolanaAdapter(logger)

	seed := []byte("12345678901234567890123456789012")
	derivationPath := "m/44'/501'/0'/0'"

	privKey, err := adapter.DerivePrivateKey(seed, derivationPath, false)
	if err != nil {
		log.Fatal(err)
	}
	pubKey, err := adapter.DerivePublicKey(seed, derivationPath, false)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Private Key (base58):", privKey)
	fmt.Println("Public Key (base58):", pubKey)

	from, err := solana.PublicKeyFromBase58(pubKey)
	if err != nil {
		log.Fatal(err)
	}
	to := from

	var recentBlockhash solana.Hash

	instruction := createTransferInstruction(from, to, 1)

	tx, err := solana.NewTransaction(
		[]solana.Instruction{instruction},
		recentBlockhash,
		solana.TransactionPayer(from),
	)
	if err != nil {
		log.Fatal(err)
	}

	privKeyObj, err := solana.PrivateKeyFromBase58(privKey)
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		return &privKeyObj
	})
	if err != nil {
		log.Fatal(err)
	}

	txBytes, err := tx.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}

	payloadHex := hex.EncodeToString(txBytes)
	fmt.Println("Payload (hex):", payloadHex)
}

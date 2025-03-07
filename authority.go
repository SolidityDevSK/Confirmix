package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// Authority represents a validator in the PoA system
type Authority struct {
	Address    string
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

// NewAuthority creates a new authority with a keypair
func NewAuthority() (*Authority, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	publicKey := &privateKey.PublicKey
	address := fmt.Sprintf("%x", sha256.Sum256(elliptic.Marshal(elliptic.P256(), publicKey.X, publicKey.Y)))

	return &Authority{
		Address:    address,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

// Sign creates a signature for a block
func (a *Authority) Sign(data []byte) ([]byte, error) {
	hash := sha256.Sum256(data)
	r, s, err := ecdsa.Sign(rand.Reader, a.PrivateKey, hash[:])
	if err != nil {
		return nil, err
	}

	signature := append(r.Bytes(), s.Bytes()...)
	return signature, nil
}

// Verify verifies a signature
func (a *Authority) Verify(data, signature []byte) bool {
	hash := sha256.Sum256(data)
	
	// Signature'ı r ve s değerlerine ayır
	r := new(big.Int).SetBytes(signature[:len(signature)/2])
	s := new(big.Int).SetBytes(signature[len(signature)/2:])

	return ecdsa.Verify(a.PublicKey, hash[:], r, s)
} 
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

const (
	filename = "key"
	keySize  = 4096
)

func haveKeys(dir string) bool {
	_, err := os.Stat(dir+"/"+filename+".rsa")
	// If there's any error, assume the files don't exist
	// If this causes problems, the plan is to refactor to look more like:
	// https://stackoverflow.com/a/12527546
	if err != nil {
		return false
	}

	_, err = os.Stat(dir+"/"+filename+".rsa.pub")
	if err != nil {
		return false
	}

	return true
}

func genKeyPair(dir string) error {
	key, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return err
	}

	pub := key.Public()

	keyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	pubPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)

	err = os.WriteFile(dir+"/"+filename+".rsa", keyPem, 0700)
	if err != nil {
		return err
	}

	err = os.WriteFile(dir+"/"+filename+".rsa.pub", pubPem, 0755)
	if err != nil {
		return err
	}

	return nil
}

type keyPair struct {
	key *rsa.PrivateKey
	pubKey *rsa.PublicKey
}

func getKeyPair(dir string) (*keyPair, error) {
	privBytes, err := os.ReadFile(dir+"/"+filename+".rsa")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(privBytes)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pubBytes, err := os.ReadFile(dir+"/"+filename+".rsa.pub")
	if err != nil {
		return nil, err
	}

	pBlock, _ := pem.Decode(pubBytes)
	pubKey, err := x509.ParsePKCS1PublicKey(pBlock.Bytes)
	if err != nil {
		return nil, err
	}
	
	return &keyPair{key, pubKey}, nil
}

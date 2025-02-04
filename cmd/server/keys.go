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
	return false
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

func getKeyPair(dir string) (rsa.PrivateKey, rsa.PublicKey) {
	return rsa.PrivateKey{}, rsa.PublicKey{}
}

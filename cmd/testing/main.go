package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Queueue0/qpass/internal/crypto"
)

func main() {
	text, passwd := os.Args[1], os.Args[2]

	salt, err := crypto.GenSalt(16)
	if err != nil {
		log.Fatal(err)
	}

	key := crypto.GetKey(passwd, salt)
	encrypted, err := crypto.Encrypt(text, key)
	if err != nil {
		log.Fatal(err)
	}

	decrypted, err := crypto.Decrypt(encrypted, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Salt: %s\nKey: %s\nEncrypted: %s\nDecrypted: %s\n", salt, key, encrypted, decrypted)
}

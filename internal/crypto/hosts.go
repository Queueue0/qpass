package crypto

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/Queueue0/qpass/internal/dbman"
)

func getKnownHosts() (map[string]*rsa.PublicKey, error) {
	home, err := dbman.GetQpassHome()
	if err != nil {
		return nil, err
	}

	fName := home + "/known_hosts"
	hostsFile, err := os.OpenFile(fName, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	defer hostsFile.Close()

	hosts := make(map[string]*rsa.PublicKey)

	scanner := bufio.NewScanner(hostsFile)
	for scanner.Scan() {
		host, keyString, ok := strings.Cut(scanner.Text(), " ")
		if !ok {
			// Something went wrong, skip current line
			continue
		}

		keyBytes, err := base64.RawStdEncoding.DecodeString(keyString)
		if err != nil {
			return nil, err
		}

		key, err := x509.ParsePKCS1PublicKey(keyBytes)
		if err != nil {
			return nil, err
		}

		hosts[host] = key
	}

	return hosts, nil
}

func addHost(addr string, key *rsa.PublicKey) error {
	home, err := dbman.GetQpassHome()
	if err != nil {
		return err
	}

	fName := home + "/known_hosts"
	f, err := os.OpenFile(fName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	keyBytes := x509.MarshalPKCS1PublicKey(key)
	keyString := base64.RawStdEncoding.EncodeToString(keyBytes)

	f.WriteString(fmt.Sprintf("%s %s\n", addr, keyString))

	return nil
}

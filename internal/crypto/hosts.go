package crypto

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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
	hostsFile, err := os.Open(fName)
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

		block, _ := pem.Decode([]byte(keyString))
		if block == nil {
			// Same as above
			continue
		}

		key, err := x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		hosts[host] = key
	}

	return hosts, nil
}

func addHost(addr string, key *rsa.PublicKey) error {
	return nil
}

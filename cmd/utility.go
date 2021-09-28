package cmd

import (
	"encoding/hex"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/nknorg/nkn-sdk-go"
	"io/ioutil"
	"os"
	"reflect"
	"time"
)

// invalidKeyPairErr is raised when a key pair is invalid.
var invalidKeyPairErr = errors.New("key pair does not match or is invalid")

// formatBinary formats a binary key to a hex encoded string for readability.
func formatBinary(data []byte) string {
	return hex.EncodeToString(data)
}

// readPrivateKey reads a stored private key from a path.
func readPrivateKey(path string) ([]byte, error) {
	byteKey, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	var privateKey = make([]byte, len(byteKey))

	length, err := hex.Decode(privateKey, byteKey)

	return privateKey[0:length], nil

}

// newNKNPrivateKey returns a new NKN account with a random seed.
func newNKNPrivateKey() ([]byte, error) {
	account, err := nkn.NewAccount(nil)
	if err != nil {
		return nil, err
	}

	secret := account.Seed()

	return secret, nil
}

// newMultiClient returns a new MultiClient with the given private key.
func newMultiClient(private []byte, numSubClients int) (*nkn.MultiClient, error) {
	// Create an account using the private key.
	account, err := nkn.NewAccount(private)
	if err != nil {
		return nil, err
	}

	// Create a new MultiClient using the private key.
	client, err := nkn.NewMultiClient(account, "", numSubClients, true, nil)
	if err != nil {
		return nil, err
	}

	return client, err
}

// writeSecretKey writes a new secret key to the desired path and returns the public key.
func writeSecretKey(keyPath string) (string, error) {
	var err error

	// Create a new key pair.
	newKey, err := newNKNPrivateKey()
	if err != nil {
		return "", nil
	}
	client, err := newMultiClient(newKey, defaultNumSubClients)
	if err != nil {
		return "", nil
	}

	// Convert the key to a hex and write it to the desired path.
	err = ioutil.WriteFile(keyPath, []byte(formatBinary(newKey)), os.ModePerm)
	if err != nil {
		return "", err
	}

	// Return the public key.
	return formatBinary(client.PubKey()), nil
}

// validateCfgKeyPair validates that the provided public key is expected of the created client.
func validateCfgKeyPair(cfgPublic string, client *nkn.MultiClient) error {
	publicBytes, err := hex.DecodeString(cfgPublic)
	if err != nil {
		return err
	}
	expectedPublic := client.PubKey()
	if !reflect.DeepEqual(publicBytes, expectedPublic) {
		return invalidKeyPairErr
	}

	return nil
}

// unixSeconds gets the current time in unix seconds.
func unixSeconds() int64 {
	return time.Now().UTC().UnixNano()
}

// newUuid returns a new UUID.
func newUuid() (string, error) {
	newUuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return newUuid.String(), nil
}
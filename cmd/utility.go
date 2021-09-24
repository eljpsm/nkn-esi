package cmd

import (
	"encoding/hex"
	"fmt"
	"github.com/nknorg/nkn-sdk-go"
	"os"
)

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

// formatBinary formats a binary key to a hex encoded string for readability.
func formatBinary(data []byte) string {
	return hex.EncodeToString(data)
}

// printPublicPrivateKeys prints the provided private and public keys with additional info.
func printPublicPrivateKeys(private []byte, public []byte) {
	fmt.Println(fmt.Sprintf("Private Key: %s", formatBinary(private)))
	fmt.Println(fmt.Sprintf("Public Key: %s", formatBinary(public)))
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

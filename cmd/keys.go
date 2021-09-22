/*
Copyright © 2021 Ecogy Energy

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"github.com/nknorg/nkn-sdk-go"
	"github.com/spf13/cobra"
)

// keysCmd represents the key command
var keysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Create a new NKN hex encoded keypair",
	Long: `Create a new NKN hex encoded keypair.

Be sure to write down the private key for later!
`,
	Run: newKeyPair,
}

func init() {
	rootCmd.AddCommand(keysCmd)
}

func newKeyPair(cmd *cobra.Command, args []string) {
	private, err := newNKNPrivateKey()
	if err != nil {
		fmt.Println(err.Error())
	}
	client, err := openMulticlient(private, defaultNumSubClients)
	if err != nil {
		fmt.Println(err.Error())
	}
	public := client.PubKey()

	printPublicPrivateKeys(private, public)
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

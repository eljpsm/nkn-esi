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
	"encoding/hex"
	"fmt"
	"github.com/nknorg/nkn-sdk-go"
	"github.com/spf13/cobra"
)

// registryNumSubClients is the number of subclients in the Registry Multiclient.
var registryNumSubClients int

// registryStartCmd represents the start command
var registryStartCmd = &cobra.Command{
	Use:   "start <registry-config.json>",
	Short: "Start the registry on the desired address and port",
	Long:  `Start the registry on the desired address and port.`,
	Args:  cobra.ExactArgs(1),
	RunE:   registryStart,
}

// init initializes registry_list.go.
func init() {
	registryCmd.AddCommand(registryStartCmd)

	registryStartCmd.Flags().IntVarP(&registryNumSubClients, "subclients", "s", defaultNumSubClients, "number of subclients to use in multiclient")
}

// registryStart is the function run by registryStartCmd.
func registryStart(cmd *cobra.Command, args []string) error {
	if verboseFlag {
		fmt.Printf("Starting Registry ...\n")
	}

	// The path to the registry config should be the first and only argument.
	registryPath := args[0]

	var private []byte
	var public []byte
	var err error

	// Get the registry config located at registryPath.
	registry, err := openRegistryConfig(registryPath)
	if err != nil {
		return err
	}

	private, err = hex.DecodeString(registry.PrivateKey)
	if err != nil {
		return err
	}
	// Open a new multiclient with the private key.
	if verboseFlag {
		fmt.Printf("Opening Multiclient with private key: %s\n", hex.EncodeToString(private))
	}
	client, err := openMulticlient(private, registryNumSubClients)
	if err != nil {
		return err
	}
	public = client.PubKey()

	// Print the key information.
	printPublicPrivateKeys(private, public)

	<-client.OnConnect.C
	if verboseFlag {
		fmt.Println("Connection opened on Registry")
	}

	err = registryLoop(client)
	if err != nil {
		return err
	}

	return nil
}

// registryLoop is the main loop of Registry.
func registryLoop(client *nkn.MultiClient) error {
	for {
	}
}

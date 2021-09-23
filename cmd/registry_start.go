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

var registryClient *nkn.MultiClient

// registryStartCmd represents the start command
var registryStartCmd = &cobra.Command{
	Use:   "start <registry-config.json>",
	Short: "Start a Registry instance",
	Long:  `Start a Registry instance.`,
	Args:  cobra.ExactArgs(2),
	RunE:  registryStart,
}

// init initializes registry_list.go.
func init() {
	registryCmd.AddCommand(registryStartCmd)

	registryStartCmd.Flags().IntVarP(&numSubClients, "subclients", "s", defaultNumSubClients, "number of subclients to use in multiclient")
}

// registryStart is the function run by registryStartCmd.
func registryStart(cmd *cobra.Command, args []string) error {
	var err error

	if verboseFlag {
		fmt.Printf("Starting Registry instance ...\n")
	}

	// The path to the registry config should be the first argument.
	registryPath := args[0]
	// The private key associated with the Registry.
	registryPrivateKey, err := hex.DecodeString(args[1])
	if err != nil {
		return err
	}

	// Get the registry config located at registryPath.
	err = openRegistryConfig(registryPath)
	if err != nil {
		return err
	}

	// Open a Multiclient with the private key and the desired number of subclients.
	registryClient, err = openMulticlient(registryPrivateKey, numSubClients)
	if err != nil {
		return err
	}

	<-registryClient.OnConnect.C
	fmt.Println(fmt.Sprintf("\nConnection opened on Registry '%s'\n"), registryInfo.Name)
	fmt.Println(registryInfo)

	// Enter the Registry shell.
	err = registryLoop()
	if err != nil {
		return err
	}

	return nil
}

// registryLoop is the main loop of a Registry.
func registryLoop() error {
	for {
		msg := <- registryClient.OnMessage.C

		switch msg.Data {
		default:
			fmt.Println(fmt.Sprintf("%s attempted to send data", msg.Src))
		}
	}
}

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
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var registryNumSubClients int

type Registry struct {
	Name       string   `json:"name"`
	PrivateKey string   `json:"privateKey"`
	Peers      []string `json:"peers"`
}

// registryStartCmd represents the start command
var registryStartCmd = &cobra.Command{
	Use:   "start <registry-config.json>",
	Short: "Start the registry on the desired address and port",
	Long:  `Start the registry on the desired address and port.`,
	Args:  cobra.ExactArgs(1),
	Run:   registryStart,
}

// init initializes registry_list.go.
func init() {
	registryCmd.AddCommand(registryStartCmd)

	registryStartCmd.Flags().IntVarP(&registryNumSubClients, "subclients", "s", defaultNumSubClients, "number of subclients to use in multiclient")
}

// registryStart is the function run by registryStartCmd.
func registryStart(cmd *cobra.Command, args []string) {
	if verboseFlag {
		fmt.Printf("Starting Registry ...\n")
	}

	registryPath := args[0]
	var registry Registry

	var private []byte
	var public []byte
	var err error

	// Open and unmarshal registry file.
	registryFile, err := os.Open(registryPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer registryFile.Close()

	byteValue, err := ioutil.ReadAll(registryFile)
	if err != nil {
		fmt.Println(err.Error())
	}
	json.Unmarshal(byteValue, &registry)

	private, err = hex.DecodeString(registry.PrivateKey)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Open a new multiclient with the private key.
	if verboseFlag {
		fmt.Printf("Opening Multiclient with private key: %s\n", hex.EncodeToString(private))
	}
	client, err := openMulticlient(private, registryNumSubClients)
	if err != nil {
		fmt.Println(err.Error())
	}
	public = client.PubKey()

	// Print the key information.
	printPublicPrivateKeys(private, public)

	// Upon successfully connecting, print a message.
	<-client.OnConnect.C
	fmt.Println("Connection opened on Registry")

	registryLoop()
}

func registryLoop() {
	for {
	}
}

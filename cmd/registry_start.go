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
	"github.com/spf13/cobra"
)

// registryStartCmd represents the start command
var registryStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the registry on the desired address and port",
	Long:  `Start the registry on the desired address and port.`,
	Args:  cobra.MaximumNArgs(0),
	Run:   registryStart,
}

// init initializes registry_list.go.
func init() {
	registryCmd.AddCommand(registryStartCmd)
}

// registryStart is the function run by registryStartCmd.
func registryStart(cmd *cobra.Command, args []string) {
	if verboseFlag {
		fmt.Printf("Starting Registry ...\n")
	}
	//publicKey := viper.Get(registryPublicKeyCfgName).([]byte)
	//privateKey := viper.Get(registryPrivateKeyCfgName).([]byte)

	//client := newNKNMulticlient()
	// Create a new Registry Multiclient.
	//client := newNKNMulticlient("registry", defaultNumSubClients)

	// Upon successfully connecting, print a message.
	//<- client.OnConnect.C
	fmt.Println("Connection opened on Registry")
}

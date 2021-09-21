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

const (
)

// registryStartCmd represents the start command
var registryStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the registry on the desired address and port",
	Long:  `Start the registry on the desired address and port.`,
	Args: cobra.MaximumNArgs(0),
	Run:   registryStart,
}

// init initializes registry_list.go.
func init() {
	registryCmd.AddCommand(registryStartCmd)

	registryStartCmd.Flags().StringP("address", "a", defaultRegistryAddress, "the address the registry will listen on")
	registryStartCmd.Flags().IntP("port", "p", defaultRegistryPort, "the port the registry will listen on")
}

// registryStart is the function run by registryStartCmd.
func registryStart(cmd *cobra.Command, args []string) {
	address, _ := cmd.Flags().GetString("address")
	port, _ := cmd.Flags().GetInt("port")

	if verboseFlag {
		fmt.Printf("Starting Registry: %s:%d\n", address, port)
	}
}

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
	"github.com/nknorg/nkn-sdk-go"
	"github.com/spf13/cobra"
)

var (
	// coordinationNodeClient is the Multiclient opened representing the Facility.
	coordinationNodeClient *nkn.MultiClient
	// coordinationNodePath is the name of what to initialize the new coordination node as.
	coordinationNodePath   string
)

// coordinationNodeStartCmd represents the start command
var coordinationNodeStartCmd = &cobra.Command{
	Use:   "start <coordination-node-config.json> <coordination-node-key.secret>",
	Short: "Start a coordination node instance",
	Long:  `Start a coordination node instance.`,
	Args:  cobra.ExactArgs(2),
	RunE:  coordinationNodeStart,
}

// init initializes coordination_node_start.go.
func init() {
	coordinationNodeCmd.AddCommand(coordinationNodeStartCmd)

	coordinationNodeStartCmd.Flags().IntVarP(&numSubClients, "subclients", "s", defaultNumSubClients, "number of subclients to use in multiclient")
}

// coordinationNodeStart is the function run by coordinationNodeStartCmd.
func coordinationNodeStart(cmd *cobra.Command, args []string) error {
	var err error

	// The path to the coordination-node-config config should be the first argument.
	coordinationNodePath = args[0]
	// The private key associated with the Facility.
	privateKey, err := readPrivateKey(args[1])
	if err != nil {
		return err
	}

	// Get the coordination-node-config config located at coordinationNodePath.
	err = readCoordinationNodeConfig(coordinationNodePath)

	// Open a Multiclient with the private key and the desired number of subclients.
	coordinationNodeClient, err = newMultiClient(privateKey, numSubClients)
	if err != nil {
		return err
	}

	// Validate the key pair.
	err = validateCfgKeyPair(coordinationNodeInfo.PublicKey, coordinationNodeClient)
	if err != nil {
		return err
	}

	// Enter the Facility shell.
	coordinationNodeShell()

	return nil
}

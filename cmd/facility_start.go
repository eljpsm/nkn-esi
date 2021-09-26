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

var (
	// facilityClient is the Multiclient opened representing the Facility.
	facilityClient *nkn.MultiClient
)

// facilityStartCmd represents the start command
var facilityStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Facility instance",
	Long:  `Start a Facility instance.`,
	Args:  cobra.ExactArgs(2),
	RunE:  facilityStart,
}

// init initializes facility_start.go.
func init() {
	facilityCmd.AddCommand(facilityStartCmd)

	facilityStartCmd.Flags().IntVarP(&numSubClients, "subclients", "s", defaultNumSubClients, "number of subclients to use in multiclient")
}

// facilityStart is the function run by facilityStartCmd.
func facilityStart(cmd *cobra.Command, args []string) error {
	var err error

	// The path to the facility-config config should be the first argument.
	facilityPath := args[0]
	// The private key associated with the Facility.
	facilityPrivateKey, err := readPrivateKey(args[1])
	if err != nil {
		return err
	}

	// Get the facility-config config located at facilityPath.
	err = readFacilityConfig(facilityPath)

	// Open a Multiclient with the private key and the desired number of subclients.
	facilityClient, err = newMultiClient(facilityPrivateKey, numSubClients)
	if err != nil {
		return err
	}

	// Validate the key pair.
	err = validateCfgKeyPair(facilityInfo.FacilityPublicKey, facilityClient)
	if err != nil {
		return err
	}

	<-facilityClient.OnConnect.C
	infoMsgColor.Println(fmt.Sprintf("\nConnection opened on Facility '%s'\n", noteMsgColorFunc(facilityInfo.Name)))
	fmt.Printf("Public Key: %s\n", formatBinary(facilityClient.PubKey()))

	// Enter the Facility shell.
	facilityShell()

	return nil
}

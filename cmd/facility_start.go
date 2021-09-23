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
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/nknorg/nkn-sdk-go"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	userFacility      esi.DerFacilityExchangeInfo
	userFacilityClient *nkn.MultiClient
	facilityPrivateKey []byte
	facilityPublicKey []byte
	UnknownCommandErr = errors.New("unknown command")
)

// facilityNumSubClients is the number of subclients in the Facility Multiclient.
var facilityNumSubClients int

// facilityStartCmd represents the start command
var facilityStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Facility instance",
	Long:  `Start a Facility instance.`,
	Args: cobra.ExactArgs(2),
	RunE:  facilityStart,
}

// init initializes facility_start.go.
func init() {
	facilityCmd.AddCommand(facilityStartCmd)

	facilityStartCmd.Flags().IntVarP(&facilityNumSubClients, "subclients", "s", defaultNumSubClients, "number of subclients to use in multiclient")
}

// facilityStart is the function run by facilityStartCmd.
func facilityStart(cmd *cobra.Command, args []string) error {
	if verboseFlag {
		fmt.Printf("Starting Facility instance ...\n")
	}

	var err error

	// The path to the facility config should be the first and only argument.
	facilityPath := args[0]
	// The private key associated with the Facility.
	facilityPrivateKey, err = hex.DecodeString(args[1])

	// Get the facility config located at facilityPath.
	userFacility, err = openFacilityConfig(facilityPath)
	if err != nil {
		return err
	}

	userFacilityClient, err = openMulticlient(facilityPrivateKey, facilityNumSubClients)
	if err != nil {
		return err
	}

	facilityPublicKey = userFacilityClient.PubKey()

	// Print the key information.
	printPublicPrivateKeys(facilityPrivateKey, facilityPublicKey)

	<-userFacilityClient.OnConnect.C
	if verboseFlag {
		fmt.Println("Connection opened on Facility")
	}

	err = facilityLoop()
	if err != nil {
		return err
	}

	return nil
}

// facilityLoop is the main loop of a Facility.
func facilityLoop() error {
	//defer handleExit()
	var input string

	if verboseFlag {
		fmt.Printf("Entering facilityLoop ...\n")
	}
	for {
		// Prompt the user for input.
		input = prompt.Input("> ", facilityCompleter)

		// Execute the input and receive a message and error.
		message, err := facilityExecutor(input)

		// If the execution results in an error, alert the user.
		if err != nil {
			fmt.Println(err.Error())
		}

		// If a message was sent back, then show the user.
		if message != "" {
			fmt.Println(message)
		}
	}
}

// facilityCompleter is the completer for a Facility.
func facilityCompleter(d prompt.Document) []prompt.Suggest {
	// Useful prompts that the user can use in the shell.
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit out of Facility instance"},
		{Text: "info", Description: "Print info on Facility"},
		{Text: "discover", Description: "Discover and send Facility info to Registry"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// facilityExecutor is the function which executes user input.
func facilityExecutor(input string) (string, error) {
	fields := strings.Fields(input)
	fmt.Println(fields)

	// If there is no input, simply return.
	if len(fields) == 0 {
		return "", nil
	}

	// Evaluate the first string.
	switch fields[0] {
	default:
		return "", UnknownCommandErr
	case "exit":
		// Exit out of the program.
		os.Exit(0)
	case "info":
		fmt.Println(userFacility)
	case "discover":
	case "register":
	}

	return "", nil
}

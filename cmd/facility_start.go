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
	"github.com/nknorg/nkn-sdk-go"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	UnknownCommandErr = errors.New("unknown command")
)

// facilityStartCmd represents the start command
var facilityStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Facility instance",
	Long:  `Start a Facility instance.`,
	RunE:  facilityStart,
}

// init initializes facility_start.go.
func init() {
	facilityCmd.AddCommand(facilityStartCmd)
}

// facilityStart is the function run by facilityStartCmd.
func facilityStart(cmd *cobra.Command, args []string) error {
	if verboseFlag {
		fmt.Printf("Starting Facility instance ...\n")
	}

	// The path to the facility config should be the first and only argument.
	facilityPath := args[0]

	var private []byte
	var public []byte
	var err error

	// Get the facility config located at facilityPath.
	registry, err := openFacilityConfig(facilityPath)
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

	err = facilityLoop(client)
	if err != nil {
		return err
	}

	return nil
}

// facilityLoop is the main loop of a Facility.
func facilityLoop(client *nkn.MultiClient) error {
	//defer handleExit()
	var input string

	if verboseFlag {
		fmt.Printf("Entering facilityLoop ...\n")
	}
	for {
		// Prompt the user for input.
		input = prompt.Input("> ", facilityCompleter)

		// Execute the input and receive a message and error.
		message, err := facilityExecutor(input, client)

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
		{Text: "register", Description: "Register your facility with a Registry"},
		{Text: "exit", Description: "Exit out of Facility instance"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// facilityExecutor is the function which executes user input.
func facilityExecutor(input string, client *nkn.MultiClient) (string, error) {
	fields := strings.Fields(input)
	fmt.Println(fields)

	// If there is no input, simply return.
	if len(fields) == 0 {
		return "", nil
	}

	// If there is only one split input, treat it as a single command with no arguments.
	if len(fields) == 1 {
		switch fields[0] {
		default:
			return "", UnknownCommandErr
		case "exit":
			os.Exit(0)
		}
	}

	// Evaluate the first command.
	switch fields[0] {
	default:
		return "", UnknownCommandErr
	case "register":
		msg, err := client.Send(nkn.NewStringArray(fields[1]), []byte("Hello, World!"), nil)
		if err != nil {
			return "", err
		}
		fmt.Println(msg.C)
	}

	return "", nil
}

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
	"github.com/golang/protobuf/proto"
	"github.com/nknorg/nkn-sdk-go"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	// facilityClient is the Multiclient opened representing the Facility.
	facilityClient *nkn.MultiClient
	// unknownCommandErr is the error returned for any unknown input.
	unknownCommandErr = errors.New("unknown command")
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

	if verboseFlag {
		fmt.Printf("Starting Facility instance ...\n")
	}

	// The path to the facility config should be the first argument.
	facilityPath := args[0]
	// The private key associated with the Facility.
	facilityPrivateKey, err := hex.DecodeString(args[1])
	if err != nil {
		return err
	}

	// Get the facility config located at facilityPath.
	err = openFacilityConfig(facilityPath)
	if err != nil {
		return err
	}

	// Open a Multiclient with the private key and the desired number of subclients.
	facilityClient, err = openMulticlient(facilityPrivateKey, numSubClients)
	if err != nil {
		return err
	}

	<-facilityClient.OnConnect.C
	infoMsgColor.Println(fmt.Sprintf("\nConnection opened on Facility '%s'\n", noteMsgColorFunc(facilityInfo.Name)))

	// Enter the Facility shell.
	err = facilityShell()
	if err != nil {
		return err
	}

	return nil
}

// facilityLoop is the main loop of a Facility.
func facilityShell() error {
	// Create two channels, one for incoming messages and another for outgoing inputs.
	messages := make(chan string)
	inputs := make(chan string)
	go messageReceiver(messages)
	go inputReceiver(inputs)

	for {
		select {

		case message, ok := <-messages:
			// Evaluate incoming messages.
			if !ok {
				break
			}

			// Print any new message received from the receiver.
			fmt.Println(message)

		case input, ok := <-inputs:
			if !ok {
				break
			}

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
}

// messageReceiver receives and returns any incoming messages.
func messageReceiver(messagesCh chan string) {
	message := &esi.RegistryMessage{}

	for {
		msg := <-facilityClient.OnMessage.C
		err := proto.Unmarshal(msg.Data, message)
		if err != nil {
			continue
		}
		messagesCh <- message.String()
	}
}

// inputReceiver receives and returns any inputs.
func inputReceiver(inputCh chan string) {
	for {
		input := prompt.Input("> ", facilityCompleter)
		inputCh <- string(input)
	}
}

// facilityCompleter is the completer for a Facility.
func facilityCompleter(d prompt.Document) []prompt.Suggest {
	// Useful prompts that the user can use in the shell.
	s := []prompt.Suggest{
		{Text: "exit", Description: "Exit out of Facility instance"},
		{Text: "info", Description: "Print info on Facility"},
		{Text: "signup", Description: "Signup and send Facility info to Registry"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// facilityExecutor is the function which executes user input.
func facilityExecutor(input string) (string, error) {
	fields := strings.Fields(input)

	// If there is no input, simply return.
	if len(fields) == 0 {
		return "", nil
	}

	// Evaluate the first string.
	switch fields[0] {
	default:
		return "", unknownCommandErr
	case "exit":
		// Exit out of the program.
		os.Exit(0)
	case "info":
		fmt.Println(facilityInfo)
	case "signup":
		_, err := esi.DiscoverRegistry(facilityClient, fields[1], facilityInfo)
		if err != nil {
			return "", err
		}
	case "register":
	}

	return "", nil
}

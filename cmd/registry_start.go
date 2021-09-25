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
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"github.com/golang/protobuf/proto"
	"github.com/nknorg/nkn-sdk-go"
	"github.com/spf13/cobra"
)

var (
	// registryClient is the Multiclient opened representing the Registry.
	registryClient *nkn.MultiClient
)

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

	// The path to the registry config should be the first argument.
	registryPath := args[0]
	// The private key associated with the Registry.
	registryPrivateKey, err := readPrivateKey(args[1])
	if err != nil {
		return err
	}

	// Get the registry config located at registryPath.
	err = readRegistryConfig(registryPath)
	if err != nil {
		return err
	}

	// Open a Multiclient with the private key and the desired number of subclients.
	registryClient, err = newMultiClient(registryPrivateKey, numSubClients)
	if err != nil {
		return err
	}

	<-registryClient.OnConnect.C
	infoMsgColor.Println(fmt.Sprintf("\nConnection opened on Registry '%s'\n", noteMsgColorFunc(registryInfo.Name)))
	fmt.Printf("Public Key: %s\n", formatBinary(registryClient.PubKey()))

	// Enter the Registry shell.
	err = registryLoop()
	if err != nil {
		return err
	}

	return nil
}

// registryLoop is the main loop of a Registry.
func registryLoop() error {
	fmt.Println("Awaiting messages ...")
	message := &esi.RegistryMessage{}

	// Facilities currently stored in memory.
	facilities := make(map[string]*esi.DerFacilityExchangeInfo)

	for {
		msg := <-registryClient.OnMessage.C
		fmt.Printf("Message received from %s\n", noteMsgColorFunc(msg.Src))
		err := proto.Unmarshal(msg.Data, message)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// Case documentation located at api/esi/der_handler.proto.
		switch x := message.Chunk.(type) {
		case *esi.RegistryMessage_DerFacilityExchangeInfo:
			if _, ok := facilities[x.DerFacilityExchangeInfo.FacilityPublicKey]; !ok {
				infoMsgColor.Printf("Saved Facility public key(s) to known Facilities\n")

				facilities[x.DerFacilityExchangeInfo.FacilityPublicKey] = x.DerFacilityExchangeInfo

				for _, v := range facilities {
					data, err := proto.Marshal(&esi.FacilityMessage{Chunk: &esi.FacilityMessage_DerFacilityExchangeInfo{DerFacilityExchangeInfo: v}})
					if err != nil {
						panic(err)
					}

					_, err = registryClient.Send(nkn.NewStringArray(msg.Src), data, nil)
					if err != nil {
						panic(err)
					}
				}
			}

		case *esi.RegistryMessage_DerFacilityExchangeRequest:
			for _, v := range facilities {
				if v.Location.Country == "New Zealand" {
					data, _ := proto.Marshal(&esi.FacilityMessage{Chunk: &esi.FacilityMessage_DerFacilityExchangeInfo{DerFacilityExchangeInfo: v}})
					fmt.Printf("Send Facility %s to %s\n", infoMsgColorFunc(v.FacilityPublicKey), noteMsgColorFunc(msg.Src))
					_, err = registryClient.Send(nkn.NewStringArray(msg.Src), data, nil)
					if err != nil {
						return err
					}
				}
			}
		}
	}
}

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
	"encoding/json"
	"github.com/elijahjpassmore/nkn-esi/api/esi"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// coordinationNodeInitCmd represents the coordinationNodeInitCmd command
var coordinationNodeInitCmd = &cobra.Command{
	Use:   "init <name>",
	Short: "Quickly create a new coordination node configuration and key pair",
	Long:  `Quickly create a new coordination node configuration and key pair.`,
	Args:  cobra.ExactArgs(1),
	RunE:  coordinationNodeInit,
}

// init initializes coordination_node_init.go.
func init() {
	coordinationNodeCmd.AddCommand(coordinationNodeInitCmd)
}

// coordinationNodeInit is the function run by coordinationNodeInitCmd.
func coordinationNodeInit(cmd *cobra.Command, args []string) error {
	var err error
	coordinationNodePath = args[0]

	publicKey, err := writeSecretKey(coordinationNodePath + secretKeySuffix)
	if err != nil {
		return err
	}

	newLatLng := esi.LatLng{
		Latitude:  90.0,
		Longitude: 180.0,
	}
	newLocation := esi.Location{
		Country:       "DC",
		Region:        "Phoney",
		TimeZone:      "DMT",
		StateProvince: "Unreal",
		PostalCode:    "0000",
		Locality:      "Hoax",
		Sublocality:   "Fraud",
		StreetAddress: []string{
			"30 Fake Street",
		},
		Latlng: &newLatLng,
	}
	newDerFacilityExchangeInfo := esi.DerFacilityExchangeInfo{
		Name:      "Dummy Coordination Node",
		PublicKey: publicKey,
		Location:  &newLocation,
	}

	// Write the new config.
	jsonBytes, err := json.MarshalIndent(&newDerFacilityExchangeInfo, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(coordinationNodePath+interfaceCfgSuffix, jsonBytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

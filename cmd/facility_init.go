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

// facilityInitCmd represents the facilityInitCmd command
var facilityInitCmd = &cobra.Command{
	Use:   "init <keyName> <configName>",
	Short: "Quickly create a new registry configuration and key pair",
	Long:  `Quickly create a new registry configuration and key pair.`,
	Args:  cobra.ExactArgs(2),
	RunE:  facilityInit,
}

// init initializes facility_init.go.
func init() {
	facilityCmd.AddCommand(facilityInitCmd)
}

// registryInit is the function run by registryInitCmd.
func facilityInit(cmd *cobra.Command, args []string) error {
	var err error
	keyName := args[0]
	configName := args[1]

	publicKey, err := writeSecretKey(keyName+secretKeySuffix)
	if err != nil {
		return err
	}

	// new latlng config.
	newLatlng := esi.LatLng{
		Latitude: -36.86397,
		Longitude: 174.72052,
	}
	// New location config.
	newLocation := esi.Location{
		Country: "New Zealand",
		Region: "Auckland",
		TimeZone: "NZT",
		StateProvince: "Auckland",
		PostalCode: "1022",
		Locality: "Auckland",
		Sublocality: "Western Springs",
		StreetAddress: []string{
			"Motions Road",
		},
		Latlng: &newLatlng,
	}
	// New registry config.
	newConfig := esi.DerFacilityExchangeInfo{
		Name: "New Facility",
		FacilityPublicKey: publicKey,
		Location: &newLocation,
	}

	// Write the new config.
	jsonBytes, err := json.MarshalIndent(newConfig, "", "  ")
	if err != nil {
		return err
	}
	ioutil.WriteFile(configName+interfaceCfgSuffix, jsonBytes, os.ModePerm)

	return nil
}

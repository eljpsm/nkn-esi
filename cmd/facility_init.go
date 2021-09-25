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
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

// facilityInitCmd represents the facilityInitCmd command
var facilityInitCmd = &cobra.Command{
	Use:   "init <facilityName>",
	Short: "Quickly create a new facility configuration and key pair",
	Long:  `Quickly create a new facility configuration and key pair.`,
	Args:  cobra.ExactArgs(1),
	RunE:  facilityInit,
}

// init initializes facility_init.go.
func init() {
	facilityCmd.AddCommand(facilityInitCmd)
}

// facilityInit is the function run by facilityInitCmd.
func facilityInit(cmd *cobra.Command, args []string) error {
	var err error
	facilityPath := args[0]

	publicKey, err := writeSecretKey(facilityPath+secretKeySuffix)
	if err != nil {
		return err
	}

	newConfig := dummyDerFacilityExchangeInfo
	newConfig.FacilityPublicKey = publicKey

	// Write the new config.
	jsonBytes, err := json.MarshalIndent(newConfig, "", "  ")
	if err != nil {
		return err
	}
	ioutil.WriteFile(facilityPath+interfaceCfgSuffix, jsonBytes, os.ModePerm)

	return nil
}

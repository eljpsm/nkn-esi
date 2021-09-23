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
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

var facilityInfo esi.DerFacilityExchangeInfo

// facilityCmd represents the facility command
var facilityCmd = &cobra.Command{
	Use:   "facility",
	Short: "Manage Facility instances",
	Long:  `Manage Facility instances.`,
	Args:  cobra.ExactArgs(1),
}

// init initializes facility.go.
func init() {
	rootCmd.AddCommand(facilityCmd)
}

// openFacilityConfig opens and reads the given facility config.
func openFacilityConfig(facilityPath string) error {
	// Open facility file.
	registryFile, err := os.Open(facilityPath)
	if err != nil {
		return err
	}
	defer registryFile.Close()

	byteValue, err := ioutil.ReadAll(registryFile)
	if err != nil {
		return err
	}
	// Unmarshal it.
	json.Unmarshal(byteValue, &facilityInfo)

	return nil
}

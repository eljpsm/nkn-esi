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

var coordinationNodeInfo esi.DerFacilityExchangeInfo

// coordinationNodeCmd represents the facility command
var coordinationNodeCmd = &cobra.Command{
	Use:   "coordination-node",
	Short: "Manage coordination node instances",
	Long: `Manage coordination node instances.

A coordination node combines the capabilities of the DER Facility (DERF) and
Interfacing Party with External Responsibility (IPER) defined in the ESI
server.`,
	Args: cobra.ExactArgs(1),
}

// init initializes coordination_node.go.
func init() {
	rootCmd.AddCommand(coordinationNodeCmd)
}

// readCoordinationNodeConfig opens and reads the given coordination-node-config config into a DerFacilityExchangeInfo
// struct.
func readCoordinationNodeConfig(coordinationNodePath string) error {
	coordinationNodeFile, err := os.Open(coordinationNodePath)
	if err != nil {
		return err
	}
	defer coordinationNodeFile.Close()

	byteValue, err := ioutil.ReadAll(coordinationNodeFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteValue, &coordinationNodeInfo)
	if err != nil {
		return err
	}

	return nil
}

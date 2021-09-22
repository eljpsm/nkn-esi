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
	"errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

type Registry struct {
	Name       string   `json:"name"`
	PrivateKey string   `json:"privateKey"`
	Facilities []string `json:"peers"`
}

// registryCmd represents the registry command
var registryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage registry entries",
	Long:  `Manage registry entries.`,
	Args:  cobra.ExactArgs(1),
}

// init initializes registry.go.
func init() {
	rootCmd.AddCommand(registryCmd)
}

func openRegistry(registryPath string) (Registry, error) {
	var registry Registry

	// Open registry file.
	registryFile, err := os.Open(registryPath)
	if err != nil {
		return registry, err
	}
	defer registryFile.Close()

	byteValue, err := ioutil.ReadAll(registryFile)
	if err != nil {
		return registry, err
	}
	// Unmarshal it.
	json.Unmarshal(byteValue, &registry)
	if registry.PrivateKey == "" {
		return registry, errors.New("field PrivateKey must not be empty")
	}
	// Return the result.
	return registry, nil
}

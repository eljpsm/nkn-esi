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
	"github.com/spf13/cobra"
)

// registryListCmd represents the list command
var registryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the available Facilities in the Registry",
	Long: `List the available Facilities in the Registry.`,
	Args: cobra.ExactArgs(1),
	RunE: registryList,
}

// init initializes registry_list.go.
func init() {
	registryCmd.AddCommand(registryListCmd)
}

// registryList is the function run by registryListCmd.
func registryList(cmd *cobra.Command, args []string) error {
	if verboseFlag {
		fmt.Println("Listing Registry facilities ...")
	}
	registry, err := openRegistryConfig(args[0])
	if err != nil {
		return err
	}
	for _, s := range registry.Facilities {
		fmt.Println(s)
	}

	return nil
}
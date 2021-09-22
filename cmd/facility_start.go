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
	"github.com/spf13/cobra"
)

// facilityStartCmd represents the start command
var facilityStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a Facility instance",
	Long: `Start a Facility instance.`,
	RunE: facilityStart,
}

// init initializes facility_start.go.
func init() {
	facilityCmd.AddCommand(facilityStartCmd)
}

// facilityStart is the function run by facilityStartCmd.
func facilityStart(cmd *cobra.Command, args []string) error {
	return nil
}

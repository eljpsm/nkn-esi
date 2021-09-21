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
	"os"

	"github.com/spf13/viper"
)

var (
	// userHome is the current user's home directory.
	userHome, _ = os.UserHomeDir()

	// defaultCfgFile is the default config file path.
	defaultCfgFile = fmt.Sprintf("%s/.config/nkn-esi/nkn-esi.yaml", userHome)

	// cfgFile is the file path given by the user via the config flag.
	cfgFile string
	// verbose is the bool for the persistent flag verbose.
	verboseFlag bool
)

const (
	// configFlagName is the name of the config flag.
	configFlagName = "config"
	// verboseFlagName is the name of the verbose flag.
	verboseFlagName = "verbose"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nkn-esi",
	Short: "NKN-ESI (or nESI) is an NKN based Energy Services Interface (ESI)",
	Long: `NKN-ESI (or nESI) is an NKN based Energy Services Interface (ESI).

Create and maintain facilities and registries. NKN-ESI can be used to facilitate
services such as load shifting, or the timed increased consumption of energy.
This allows an aggregator, utility, or distribution system operator to easily
and cost effectively maintain a stable and resilient electricity grid.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

// init initializes root.go.
func init() {
	cobra.OnInitialize(initMain)

	rootCmd.PersistentFlags().StringVar(&cfgFile, configFlagName, "", fmt.Sprintf("config file (default is %s", defaultCfgFile))
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, verboseFlagName, "v", false, "make the operation more talkative")
}

// initMain is the main initialization function.
func initMain() {
	initConfig()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		if verboseFlag {
			fmt.Printf("Reading config from: %s\n", cfgFile)
		}
		viper.SetConfigFile(cfgFile)
	}
	viper.SetConfigFile(defaultCfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verboseFlag {
			fmt.Printf("Reading config from: %s\n", defaultCfgFile)
		}
	}
}

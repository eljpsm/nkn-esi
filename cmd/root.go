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
	"encoding/hex"
	"fmt"
	"github.com/nknorg/nkn-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	// userHome is the current user's home directory.
	userHome, _ = os.UserHomeDir()

	// defaultCfgFile is the default config file path.
	defaultCfgFile = fmt.Sprintf("%s/.config/nkn-esi/nkn-esi.yaml", userHome)

	// cfgFile is the file path given by the user via the config flag.
	cfgFile string
	// verboseFlag is the bool for the persistent flag verbose.
	verboseFlag bool
)

const (
	registryPublicKeyCfgName  = "registry_public_key"
	registryPrivateKeyCfgName = "registry_private_key"
	// defaultNumSubClients is the default number of sub clients created using nkn.Multiclient.
	defaultNumSubClients = 3
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, configFlagName, "", fmt.Sprintf("config file (default is %s", defaultCfgFile))
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, verboseFlagName, "v", false, "make the operation more talkative")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		if verboseFlag {
			fmt.Printf("Reading config from: %s\n", cfgFile)
		}
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile(defaultCfgFile)

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			if verboseFlag {
				fmt.Printf("Reading config from: %s\n", defaultCfgFile)
			}
		} else {
			// Else, create a new config.
			writeNewCfgFile()
		}
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.ReadInConfig() // read in config
}

func writeNewCfgFile() {
	registryPublicKey, registryPrivateKey, _ := newNKNAccount()
	viper.Set(registryPublicKeyCfgName, hex.EncodeToString(registryPublicKey))
	viper.Set(registryPrivateKeyCfgName, hex.EncodeToString(registryPrivateKey))

	viper.WriteConfig()
}

// newNKNAccount returns a new NKN account with a random seed.
func newNKNAccount() ([]byte, []byte, error) {
	account, err := nkn.NewAccount(nil)
	if err != nil {
		return nil, nil, err
	}
	return account.Seed(), account.PubKey(), nil
}

func newNKNMulticlient(publicKeyName string, privateKeyName string, baseIdentifier string, numSubClients int) (*nkn.MultiClient, error) {
	if verboseFlag {
		fmt.Println("Creating new Multiclient ...")
		fmt.Printf("Base identfier: %s\n", baseIdentifier)
		fmt.Printf("Number of sub clients: %d\n", numSubClients)
	}
	account, err := nkn.NewAccount(nil)
	if err != nil {
		return nil, err
	}
	client, err := nkn.NewMultiClient(account, baseIdentifier, numSubClients, true, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

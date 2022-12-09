/*
Copyright © 2022 Ciena Corporation

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
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
//
//nolint:exhaustruct
var rootCmd = &cobra.Command{
	Use:   "gns3ctl",
	Short: "Controls an instance of a GNS3 server",
	Long: `Allows a user to manipulate the GNS3 network simulator,
including the ability to create example networks and extract
information about those networks.`,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	baseDirectory := path.Join(os.Getenv("HOME"), "GNS3")

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gns3ctl.yaml)")

	rootCmd.PersistentFlags().StringP("base-directory", "d", baseDirectory, "Default project name to use when performing project specific operations")
	_ = viper.BindPFlag("base-directory", rootCmd.PersistentFlags().Lookup("base-directory"))

	rootCmd.PersistentFlags().StringP("project", "p", "default", "Default project name to use when performing project specific operations")
	_ = viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))

	rootCmd.PersistentFlags().StringP("address", "a", "localhost:3080", "Service and port on which to contact the server")
	_ = viper.BindPFlag("address", rootCmd.PersistentFlags().Lookup("address"))

	rootCmd.PersistentFlags().BoolP("insecure-skip-verify", "k", true, "Skip verifing the TLS cert of the host")
	_ = viper.BindPFlag("insecure-skip-verify", rootCmd.PersistentFlags().Lookup("insecure-skip-verify"))

	rootCmd.PersistentFlags().StringP("compute", "c", "local", "Default compute to use when creating projects")
	_ = viper.BindPFlag("compute", rootCmd.PersistentFlags().Lookup("compute"))

	rootCmd.PersistentFlags().DurationP("timeout", "t", 20*time.Second, "Timeout for http requests")
	_ = viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))

	rootCmd.PersistentFlags().StringP("username", "u", "admin", "Username for basic authentication")
	_ = viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))

	rootCmd.PersistentFlags().StringP("password", "w", "admin", "Password for basic authentication")
	_ = viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))

	rootCmd.PersistentFlags().String("download-buffer-size", "10M", "size of in memory buffer to use for file downloads")
	_ = viper.BindPFlag("download-buffer-size", rootCmd.PersistentFlags().Lookup("download-buffer-size"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".gns3ctl" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".gns3ctl")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if viper.ConfigFileUsed() != "" {
			fmt.Fprintf(os.Stderr, "Error while reading configuration file: %s\n", viper.ConfigFileUsed())
			os.Exit(1)
		}
	}
}

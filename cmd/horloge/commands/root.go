/*
Copyright © 2019 Samori Gorse

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
package commands

import (
	"fmt"
	"os"

	"github.com/shinuza/horloge"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var runner *horloge.Runner
var cfgFile string
var redisPasswd string
var redisAddr string
var redisDB int

var rootCmd = &cobra.Command{
	Use:   "horloge",
	Short: "Horloge let's you schedule repeating task based on a pattern",
	Long:  `Horloge is command line tool`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.horloge.yaml)")

	// sync redis options
	rootCmd.PersistentFlags().StringVar(&redisAddr, "redis-addr", ":6379", "address of the redis server")
	rootCmd.PersistentFlags().StringVar(&redisPasswd, "redis-passwd", "", "password of the redis server")
	rootCmd.PersistentFlags().IntVar(&redisDB, "redis-db", 0, "which database to use")

	viper.BindPFlag("redis.addr", rootCmd.PersistentFlags().Lookup("redis-addr"))
	viper.BindPFlag("redis.password", rootCmd.PersistentFlags().Lookup("redis-passwd"))
	viper.BindPFlag("redis.db", rootCmd.PersistentFlags().Lookup("redis-db"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".horloge" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".horloge")
	}

	viper.SetConfigType("yaml")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

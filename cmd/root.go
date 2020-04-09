// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// builddate is the compilation date for the executable, in %Y%m%d format
var builddate string = "2020/03/21"

// version describe the current git version,
var version string = "v0"

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "mcc",
	Short:   "mcc transforms High-Level Petri nets in PNML format into equivalent Place/Transition nets",
	Version: Version(),
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// Version returns information on the current version
func Version() string {
	return fmt.Sprintf("mcc version %s -- %s -- LAAS/CNRS", version, builddate)
}

// Generated returns information that can be embedded in generated files
func Generated() string {
	return fmt.Sprintf("generated with \"mcc %s\", version: %s, build date: %s, at: %s", strings.Join(os.Args[1:], " "), version, builddate, time.Now().Format("2006-01-02T15:04:05"))
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".mcc")  // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package cmd

import (
	"github.com/dalzilio/mcc/corenet"
	"github.com/dalzilio/mcc/hlnet"
	"github.com/dalzilio/mcc/pnml"

	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// hlnetCmd takes a a PNML model for a high level net, with a .pnml extension,
// and generates a .net or .tpn file with the same basename. The command returns
// an error for  PNML files that contain a core net description.
var hlnetCmd = &cobra.Command{
	Use:   "hlnet -i file.pnml",
	Short: "generates a .net or .tpn file from a PNML file describing a high-level net",
	Run: func(cmd *cobra.Command, args []string) {
		convert(hlnetFileName)
	},
}

var hlnetFileName string
var hlnetOutFileName string
var hlnetUseName bool
var hlnetLogger *log.Logger

func init() {
	RootCmd.AddCommand(hlnetCmd)
	hlnetCmd.Flags().StringVarP(&hlnetFileName, "file", "i", "", "name of the input file (.pnml)")
	hlnetCmd.Flags().StringVarP(&hlnetOutFileName, "out", "o", "", "basename of the output file (without extension, default to input file basename)")
	hlnetCmd.Flags().BoolVar(&hlnetUseName, "name", false, "use PNML (document) name for the output file")

	hlnetLogger = log.New(os.Stderr, "MCC HLNET:", 0)
}

func convert(filename string) {
	// we capture panics
	defer func() {
		if r := recover(); r != nil {
			hlnetLogger.Println("Error in generation: cannot compute")
			os.Exit(1)
		}
	}()

	if filename == "" {
		hlnetLogger.Println("Bad command line! Input file mandatory. Use option -i")
		os.Exit(1)
		return
	}

	if filepath.Ext(filename) != ".pnml" {
		hlnetLogger.Println("Wrong file extension!")
		os.Exit(1)
		return
	}

	xmlFile, err := os.Open(filename)
	if err != nil {
		hlnetLogger.Println("Error opening file:", err)
		os.Exit(1)
		return
	}
	defer xmlFile.Close()

	decoder := pnml.NewDecoder(xmlFile)
	var p = new(pnml.Net)
	err = decoder.Build(p)
	if err != nil {
		hlnetLogger.Println("Error decoding PNML file:", err)
		os.Exit(1)
		return
	}

	// compute name of the output file. There are three possible cases
	_, file := filepath.Split(filename)

	outfile := hlnetOutFileName
	if outfile == "" {
		outfile = file[0 : len(file)-len(".pnml")]
	}
	if hlnetUseName {
		outfile = p.Name
	}

	p.SetVerbose(pnml.MINIMAL)
	p.SetFES(false)
	hl := hlnet.Build(p)

	// We try to build a TPN first
	cn, nbcopies, listlr, listh, err := corenet.BuildTPN(p, hl)
	if err == nil {
		hlnetLogger.Println("file " + outfile + " is a TPN")
		outfile = outfile + ".tpn"
		ioutil.WriteFile(outfile, []byte(cn.PrintTPN(nbcopies, listlr, listh)), 0755)
		return
	}

	cn = corenet.Build(p, hl)
	hlnetLogger.Println("file " + outfile + " is a NET")
	outfile = outfile + ".net"
	ioutil.WriteFile(outfile, []byte(cn.String()), 0755)
}

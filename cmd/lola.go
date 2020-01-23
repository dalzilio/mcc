// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package cmd

import (
	"bufio"

	"github.com/dalzilio/mcc/corenet"
	"github.com/dalzilio/mcc/hlnet"
	"github.com/dalzilio/mcc/pnml"

	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// lolaCmd takes a a PNML model for a high level net, with a .pnml extension,
// and generates a .net or .tpn file with the same basename. The command returns
// an error for  PNML files that contain a core net description.
var lolaCmd = &cobra.Command{
	Use:   "lola -i file.pnml",
	Short: "generates a net file in the LoLa format from a PNML file describing a high-level net",
	Run: func(cmd *cobra.Command, args []string) {
		lolaConvert(lolaFileName)
	},
}

var lolaFileName string
var lolaOutFileName string
var lolaUseName bool
var lolaLogger *log.Logger

func init() {
	RootCmd.AddCommand(lolaCmd)
	lolaCmd.Flags().StringVarP(&lolaFileName, "file", "i", "", "name of the input file (.pnml)")
	lolaCmd.Flags().StringVarP(&lolaOutFileName, "out", "o", "", "basename of the output file (without extension, default to input file basename) or - for stdout")
	lolaCmd.Flags().BoolVar(&lolaUseName, "name", false, "use PNML (document) name for the output file")

	lolaLogger = log.New(os.Stderr, "MCC LOLA:", 0)
}

func lolaConvert(filename string) {
	// we capture panics so that that we don't pollute stdout if something goes wrong
	defer func() {
		if r := recover(); r != nil {
			lolaLogger.Println("Error in generation: cannot compute")
			os.Exit(1)
		}
	}()

	if filename == "" {
		lolaLogger.Println("Bad command line! Input file mandatory. Use option -i")
		os.Exit(1)
		return
	}

	if filepath.Ext(filename) != ".pnml" {
		lolaLogger.Println("Wrong file extension!")
		os.Exit(1)
		return
	}

	xmlFile, err := os.Open(filename)
	if err != nil {
		lolaLogger.Println("Error opening file:", err)
		os.Exit(1)
		return
	}
	defer xmlFile.Close()

	decoder := pnml.NewDecoder(xmlFile)
	var p = new(pnml.Net)
	err = decoder.Build(p)
	if err != nil {
		lolaLogger.Println("Error decoding PNML file:", err)
		os.Exit(1)
		return
	}

	// compute name of the output file. There are three possible cases
	_, file := filepath.Split(filename)

	outfile := lolaOutFileName
	if outfile == "" {
		outfile = file[0 : len(file)-len(".pnml")]
	}
	if lolaUseName {
		outfile = p.Name
	}

	p.SetVerbose(pnml.MINIMAL)
	// set the semantic of "overflowing" enumeration types
	p.SetFES(false)
	hl, err := hlnet.Build(p)
	if err != nil {
		hlnetLogger.Println("Error decoding PNML file:", err)
		os.Exit(1)
		return
	}

	cn := corenet.Build(p, hl)
	var out *os.File
	if outfile == "-" {
		out = os.Stdout
	} else {
		out, err = os.Create(outfile + ".net")
		if err != nil {
			hlnetLogger.Println("Error creating result file:", err)
			os.Exit(1)
			return
		}
		defer out.Close()
	}
	w := bufio.NewWriter(out)
	cn.LolaWrite(w)
	w.Flush()
}

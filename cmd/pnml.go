// Copyright 2020. LAAS-CNRS, Vertics. All rights reserved.
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

// pnmlCmd takes a PNML model for a high level net, with a .pnml extension,
// and generates a P/T net equivalent, also in PNML format.
var pnmlCmd = &cobra.Command{
	Use:   "pnml -i file.pnml",
	Short: "generates a P/T net file in PNML format",
	Run: func(cmd *cobra.Command, args []string) {
		pnmlConvert(pnmlFileName)
	},
}

var pnmlFileName string
var pnmlOutFileName string
var pnmlUseName bool
var pnmlUseComplexPNames bool

var pnmlLogger *log.Logger

func init() {
	RootCmd.AddCommand(pnmlCmd)
	pnmlCmd.Flags().StringVarP(&pnmlFileName, "file", "i", "", "name of the input file (.pnml)")
	pnmlCmd.Flags().StringVarP(&pnmlOutFileName, "out", "o", "", "basename of the output file (without extension, default to input file basename) or - for stdout")
	pnmlCmd.Flags().BoolVar(&pnmlUseName, "name", false, "use PNML (document) name for the output file")
	pnmlCmd.Flags().BoolVar(&pnmlUseComplexPNames, "sliced", false, "use structured naming for places")

	pnmlLogger = log.New(os.Stderr, "MCC PNML:", 0)
}

func pnmlConvert(filename string) {
	// we capture panics so that that we don't pollute stdout if something goes wrong
	defer func() {
		if r := recover(); r != nil {
			pnmlLogger.Println("Error in generation: cannot compute")
			os.Exit(1)
		}
	}()

	if filename == "" {
		pnmlLogger.Println("Bad command line! Input file mandatory. Use option -i")
		os.Exit(1)
		return
	}

	if filepath.Ext(filename) != ".pnml" {
		pnmlLogger.Println("Wrong file extension!")
		os.Exit(1)
		return
	}

	xmlFile, err := os.Open(filename)
	if err != nil {
		pnmlLogger.Println("Error opening file:", err)
		os.Exit(1)
		return
	}
	defer xmlFile.Close()

	decoder := pnml.NewDecoder(xmlFile)
	var p = new(pnml.Net)
	err = decoder.Build(p)
	if err != nil {
		pnmlLogger.Println("Error decoding PNML file:", err)
		os.Exit(1)
		return
	}

	// compute name of the output file. There are three possible cases. We try
	// to avoid overwriting over the input file by adding suffix -PT when not
	// using option -o.
	_, file := filepath.Split(filename)
	outfile := pnmlOutFileName
	if outfile == "" {
		outfile = file[0:len(file)-len(".pnml")] + "-PT"
	}
	if pnmlUseName {
		outfile = p.Name
	}

	p.SetSliced(pnmlUseComplexPNames)
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
		out, err = os.Create(outfile + ".pnml")
		if err != nil {
			hlnetLogger.Println("Error creating result file: ", err)
			os.Exit(1)
			return
		}
		defer out.Close()
	}
	w := bufio.NewWriter(out)
	err = cn.PnmlWrite(w)
	if err != nil {
		hlnetLogger.Println("Error encoding PNML file:", err)
		os.Exit(1)
		return
	}
	w.Flush()
}

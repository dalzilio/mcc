// Copyright 2020. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package cmd

import (
	"bufio"
	"fmt"

	"github.com/dalzilio/mcc/corenet"
	"github.com/dalzilio/mcc/hlnet"
	"github.com/dalzilio/mcc/pnml"

	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// tinaCmd takes a PNML model for a high level net, with a .pnml extension, and
// generates a P/T net equivalent in the .net format used in TINA.
var tinaCmd = &cobra.Command{
	Use:   "tina -i file.pnml",
	Short: "mcc tina generates a P/T net file in .net format",
	Run: func(cmd *cobra.Command, args []string) {
		tinaConvert(tinaFileName)
	},
}

var tinaFileName string
var tinaOutFileName string
var tinaUseName bool
var tinaUseComplexPNames bool
var tinaLogger *log.Logger

func init() {
	RootCmd.AddCommand(tinaCmd)
	tinaCmd.Flags().StringVarP(&tinaFileName, "file", "i", "", "name of the input file (.pnml)")
	tinaCmd.Flags().StringVarP(&tinaOutFileName, "out", "o", "", "basename of the output file (without extension, default to input file basename) or - for stdout")
	tinaCmd.Flags().BoolVar(&tinaUseName, "name", false, "use PNML (document) name for the output file")

	tinaLogger = log.New(os.Stderr, "MCC TINA:", 0)

	defaultusage := tinaCmd.UsageString()
	tinaCmd.SetUsageFunc(func(c *cobra.Command) error {
		fmt.Fprintf(os.Stdout, defaultusage)
		fmt.Fprintf(os.Stdout, "\nFiles:\n")
		fmt.Fprintf(os.Stdout, "   infile:    input file should be specified with option -i\n")
		fmt.Fprintf(os.Stdout, "   outfile:   output is stdout when using option \"-o -\"\n")
		fmt.Fprintf(os.Stdout, "   errorfile: errors are reported on stderr\n")
		fmt.Fprintf(os.Stdout, "   help:      is reported on stdout\n")
		return nil
	})
}

func tinaConvert(filename string) {
	// we capture panics so that that we don't pollute stdout if something goes wrong
	defer func() {
		if r := recover(); r != nil {
			tinaLogger.Println("Error in generation: cannot compute")
			os.Exit(1)
		}
	}()

	if filename == "" {
		tinaLogger.Println("Bad command line! Input file mandatory. Use option -i")
		os.Exit(1)
		return
	}

	if filepath.Ext(filename) != ".pnml" {
		tinaLogger.Println("Wrong file extension!")
		os.Exit(1)
		return
	}

	xmlFile, err := os.Open(filename)
	if err != nil {
		tinaLogger.Println("Error opening file:", err)
		os.Exit(1)
		return
	}
	defer xmlFile.Close()

	decoder := pnml.NewDecoder(xmlFile)
	var p = new(pnml.Net)
	err = decoder.Build(p)
	if err != nil {
		tinaLogger.Println("Error decoding PNML file:", err)
		os.Exit(1)
		return
	}

	// compute name of the output file. There are three possible cases.
	_, file := filepath.Split(filename)
	outfile := tinaOutFileName
	if outfile == "" {
		outfile = file[0 : len(file)-len(".pnml")]
	}
	if tinaUseName {
		outfile = p.Name
	}

	p.SetSliced(true)
	p.SetVerbose(pnml.QUIET)
	// p.SetVerbose(pnml.MINIMAL)
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
			hlnetLogger.Println("Error creating result file: ", err)
			os.Exit(1)
			return
		}
		defer out.Close()
	}
	w := bufio.NewWriter(out)
	fmt.Fprintf(w, "# %s\n", Generated())
	cn.Write(w)
	w.Flush()
}

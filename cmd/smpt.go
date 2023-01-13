// Copyright 2021. LAAS-CNRS, Vertics. All rights reserved.
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

// smptCmd takes a PNML model for a high level net, with a .pnml extension, and
// generates a P/T net equivalent in the .net format for use by the SMPT tool.
// The main differenceu with the tina command is that we add information
// relating places and transitions to their colored equivalent. This is used for
// model-checking reachability properties in the Model-Checking Contest.
var smptCmd = &cobra.Command{
	Use:   "smpt -i file.pnml",
	Short: "Generate a P/T net file for use with SMPT",
	Run: func(cmd *cobra.Command, args []string) {
		smptConvert(smptFileName)
	},
}

var smptFileName string
var smptOutFileName string
var smptUseName bool

// var smptUseComplexPNames bool
var smptLogger *log.Logger

func init() {
	RootCmd.AddCommand(smptCmd)
	smptCmd.Flags().StringVarP(&smptFileName, "file", "i", "", "name of the input file (.pnml)")
	smptCmd.Flags().StringVarP(&smptOutFileName, "out", "o", "", "basename of the output file (without extension, default to input file basename) or - for stdout")
	smptCmd.Flags().BoolVar(&smptUseName, "name", false, "use PNML (document) name for the output file")

	smptLogger = log.New(os.Stderr, "MCC SMPT:", 0)

	defaultusage := smptCmd.UsageString()
	smptCmd.SetUsageFunc(func(c *cobra.Command) error {
		fmt.Fprint(os.Stdout, defaultusage)
		fmt.Fprintf(os.Stdout, "\nFiles:\n")
		fmt.Fprintf(os.Stdout, "   infile:    input file should be specified with option -i\n")
		fmt.Fprintf(os.Stdout, "   outfile:   output is stdout when using option \"-o -\"\n")
		fmt.Fprintf(os.Stdout, "   errorfile: errors are reported on stderr\n")
		fmt.Fprintf(os.Stdout, "   help:      is reported on stdout\n")
		return nil
	})
}

func smptConvert(filename string) {
	// we capture panics so that we don't pollute stdout if something goes wrong
	defer func() {
		if r := recover(); r != nil {
			smptLogger.Println("Error in generation: cannot compute")
			os.Exit(1)
		}
	}()

	if filename == "" {
		smptLogger.Println("Bad command line! Input file mandatory. Use option -i")
		os.Exit(1)
		return
	}

	if filepath.Ext(filename) != ".pnml" {
		smptLogger.Println("Wrong file extension!")
		os.Exit(1)
		return
	}

	xmlFile, err := os.Open(filename)
	if err != nil {
		smptLogger.Println("Error opening file:", err)
		os.Exit(1)
		return
	}
	defer xmlFile.Close()

	decoder := pnml.NewDecoder(xmlFile)
	var p = new(pnml.Net)
	err = decoder.Build(p)
	if err != nil {
		smptLogger.Println("Error decoding PNML file:", err)
		os.Exit(1)
		return
	}

	// compute name of the output file. There are three possible cases.
	_, file := filepath.Split(filename)
	outfile := smptOutFileName
	if outfile == "" {
		outfile = file[0 : len(file)-len(".pnml")]
	}
	if smptUseName {
		outfile = p.Name
	}

	p.SetVerbose(pnml.SMPT)
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

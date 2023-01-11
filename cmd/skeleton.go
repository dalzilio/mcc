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

// skeletonCmd takes a PNML model for a high level net, with a .pnml extension,
// and generates a skeleton net, in the .net format, which means a P/T net
// obtained by erasing all colors from the places and all guards from the
// transitions. This construction can be used, in some cases, to efficiently
// answer reachability problems.
var skeletonCmd = &cobra.Command{
	Use:   "skeleton -i file.pnml",
	Short: "mcc skeleton generates a P/T net file in .net format",
	Run: func(cmd *cobra.Command, args []string) {
		skeletonConvert(skeletonFileName)
	},
}

var skeletonFileName string
var skeletonOutFileName string
var skeletonUseName bool

// var skeletonUseComplexPNames bool
var skeletonLogger *log.Logger

func init() {
	RootCmd.AddCommand(skeletonCmd)
	skeletonCmd.Flags().StringVarP(&skeletonFileName, "file", "i", "", "name of the input file (.pnml)")
	skeletonCmd.Flags().StringVarP(&skeletonOutFileName, "out", "o", "", "basename of the output file (without extension, default to input file basename) or - for stdout")
	skeletonCmd.Flags().BoolVar(&skeletonUseName, "name", false, "use PNML (document) name for the output file")

	skeletonLogger = log.New(os.Stderr, "MCC skeleton:", 0)

	defaultusage := skeletonCmd.UsageString()
	skeletonCmd.SetUsageFunc(func(c *cobra.Command) error {
		fmt.Fprint(os.Stdout, defaultusage)
		fmt.Fprintf(os.Stdout, "\nFiles:\n")
		fmt.Fprintf(os.Stdout, "   infile:    input file should be specified with option -i\n")
		fmt.Fprintf(os.Stdout, "   outfile:   output is stdout when using option \"-o -\"\n")
		fmt.Fprintf(os.Stdout, "   errorfile: errors are reported on stderr\n")
		fmt.Fprintf(os.Stdout, "   help:      is reported on stdout\n")
		return nil
	})
}

func skeletonConvert(filename string) {
	// we capture panics so that we don't pollute stdout if something goes wrong
	defer func() {
		if r := recover(); r != nil {
			skeletonLogger.Println("Error in generation: cannot compute")
			os.Exit(1)
		}
	}()

	if filename == "" {
		skeletonLogger.Println("Bad command line! Input file mandatory. Use option -i")
		os.Exit(1)
		return
	}

	if filepath.Ext(filename) != ".pnml" {
		skeletonLogger.Println("Wrong file extension!")
		os.Exit(1)
		return
	}

	xmlFile, err := os.Open(filename)
	if err != nil {
		skeletonLogger.Println("Error opening file:", err)
		os.Exit(1)
		return
	}
	defer xmlFile.Close()

	decoder := pnml.NewDecoder(xmlFile)
	var p = new(pnml.Net)
	err = decoder.Build(p)
	if err != nil {
		skeletonLogger.Println("Error decoding PNML file:", err)
		os.Exit(1)
		return
	}

	// compute name of the output file. There are three possible cases.
	_, file := filepath.Split(filename)
	outfile := skeletonOutFileName
	if outfile == "" {
		outfile = file[0 : len(file)-len(".pnml")]
	}
	if skeletonUseName {
		outfile = p.Name
	}

	hl, err := hlnet.Build(p)
	if err != nil {
		hlnetLogger.Println("Error decoding PNML file:", err)
		os.Exit(1)
		return
	}

	cn := corenet.BuildSkeleton(p, hl)

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

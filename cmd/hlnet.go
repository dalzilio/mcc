// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package cmd

import (
	"bufio"
	"fmt"
	"time"

	"github.com/dalzilio/mcc/corenet"
	"github.com/dalzilio/mcc/hlnet"
	"github.com/dalzilio/mcc/pnml"

	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// hlnetCmd takes a a PNML model for a high level net, with a .pnml extension,
// and outputs statistics on the size and computation time of the result (option
// `--stats`), or generates a .net file with information about the (option
// `--debug`) structure of the colored net. This replaces the `hlnet` command
// that is deprecated.` The command returns an error for PNML files that contain
// a core net description.
var hlnetCmd = &cobra.Command{
	Use:   "info -i file.pnml",
	Short: "Print statistics or generate textual version for use with NetDraw (nd)",
	Run: func(cmd *cobra.Command, args []string) {
		convert(hlnetFileName)
	},
}

var hlnetFileName string
var hlnetOutFileName string
var hlnetUseName bool
var hlnetDebugMode bool
var hlnetStat bool
var hlnetLogger *log.Logger

func init() {
	RootCmd.AddCommand(hlnetCmd)
	hlnetCmd.Flags().StringVarP(&hlnetFileName, "file", "i", "", "name of the input file (.pnml)")
	hlnetCmd.Flags().StringVarP(&hlnetOutFileName, "out", "o", "", "basename of the output file (without extension, default to input file basename) or - for stdout")
	hlnetCmd.Flags().BoolVar(&hlnetUseName, "name", false, "use PNML (document) name for the output file")
	hlnetCmd.Flags().BoolVar(&hlnetDebugMode, "debug", false, "output a readable version of the colored net for use with Tina's NetDraw (nd)")
	hlnetCmd.Flags().BoolVar(&hlnetStat, "stats", false, "print statistics (nb. of places, trans. and computation time) and quit; do not output the net")

	hlnetLogger = log.New(os.Stderr, "MCC HLNET:", 0)

	defaultusage := hlnetCmd.UsageString()
	hlnetCmd.SetUsageFunc(func(c *cobra.Command) error {
		fmt.Fprint(os.Stdout, defaultusage)
		fmt.Fprintf(os.Stdout, "\nFiles:\n")
		fmt.Fprintf(os.Stdout, "   infile:    input file should be specified with option -i\n")
		fmt.Fprintf(os.Stdout, "   outfile:   output is stdout when using option \"-o -\"\n")
		fmt.Fprintf(os.Stdout, "   errorfile: errors are reported on stderr\n")
		fmt.Fprintf(os.Stdout, "   help:      is reported on stdout\n")
		return nil
	})
}

func convert(filename string) {
	// we capture panics
	defer func() {
		if r := recover(); r != nil {
			hlnetLogger.Println("Error in generation: cannot compute")
			os.Exit(1)
		}
	}()

	start := time.Now()

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

	if hlnetDebugMode {
		hl, err := hlnet.Build(p)
		if err != nil {
			hlnetLogger.Println("Error decoding PNML file:", err)
			os.Exit(1)
			return
		}
		if outfile == "-" {
			os.Stdout.Write([]byte(p.String() + hl.Tina()))

		} else {
			os.WriteFile(outfile+".net", []byte(p.String()+hl.Tina()), 0755)
		}
		os.Exit(0)
	}

	if hlnetStat {
		// if we do not print the resulting net, we do not need to build complex
		// identifiers
		p.SetVerbose(pnml.SMPT)
	} else {
		p.SetVerbose(pnml.INFO)
	}

	p.SetFES(false)

	hl, err := hlnet.Build(p)
	if err != nil {
		hlnetLogger.Println("Error decoding PNML file:", err)
		os.Exit(1)
		return
	}

	// // We try to build a TPN first
	// cn, nbcopies, listlr, listh, err := corenet.BuildTPN(p, hl)
	// if err == nil {
	// 	// hlnetLogger.Println("file " + outfile + " is a TPN")
	// 	if hlnetStat {
	// 		elapsed := time.Since(start)
	// 		npl, ntr, narcs := cn.Statistics()
	// 		fmt.Fprintf(os.Stdout, "%d place(s), %d transition(s), %d arc(s), %.3fs\n", nbcopies*npl, nbcopies*(ntr-len(listlr)), nbcopies*narcs, elapsed.Seconds())
	// 		return
	// 	}
	// 	if outfile == "-" {
	// 		os.Stdout.Write([]byte(cn.PrintTPN(nbcopies, listlr, listh)))
	// 		return
	// 	}
	// 	outfile = outfile + ".tpn"
	// 	os.WriteFile(outfile, []byte(cn.PrintTPN(nbcopies, listlr, listh)), 0755)
	// 	return
	// }

	cn := corenet.Build(p, hl)
	if hlnetStat {
		elapsed := time.Since(start)
		npl, ntr, narcs := cn.Statistics()
		fmt.Fprintf(os.Stdout, "%d place(s), %d transition(s), %d arc(s), %.3fs\n", npl, ntr, narcs, elapsed.Seconds())
		return
	}
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
	fmt.Fprintf(w, "# %s\n", Generated())
	cn.Write(w)
	w.Flush()
}

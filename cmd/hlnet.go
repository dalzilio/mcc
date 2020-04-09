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
	Short: "mcc hlnet generates a .net or .tpn file in Tina format",
	Run: func(cmd *cobra.Command, args []string) {
		convert(hlnetFileName)
	},
}

var hlnetFileName string
var hlnetOutFileName string
var hlnetUseName bool
var hlnetDebugMode bool
var hlnetUseComplexPNames bool
var hlnetVerbose bool
var hlnetStat bool
var hlnetLogger *log.Logger

func init() {
	RootCmd.AddCommand(hlnetCmd)
	hlnetCmd.Flags().StringVarP(&hlnetFileName, "file", "i", "", "name of the input file (.pnml)")
	hlnetCmd.Flags().StringVarP(&hlnetOutFileName, "out", "o", "", "basename of the output file (without extension, default to input file basename) or - for stdout")
	hlnetCmd.Flags().BoolVar(&hlnetUseName, "name", false, "use PNML (document) name for the output file")
	hlnetCmd.Flags().BoolVar(&hlnetDebugMode, "debug", false, "output a readable version in a format that can be displayed by Tina")
	hlnetCmd.Flags().BoolVar(&hlnetUseComplexPNames, "sliced", false, "use structured naming for places")
	hlnetCmd.Flags().BoolVar(&hlnetVerbose, "verbose", false, "add extra information in the labels of the .net file")
	hlnetCmd.Flags().BoolVar(&hlnetStat, "stats", false, "print statistics (nb. of places, trans. and computation time); do not output the net")

	hlnetLogger = log.New(os.Stderr, "MCC HLNET:", 0)

	defaultusage := hlnetCmd.UsageString()
	hlnetCmd.SetUsageFunc(func(c *cobra.Command) error {
		fmt.Fprintf(os.Stdout, defaultusage)
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
			ioutil.WriteFile(outfile+".net", []byte(p.String()+hl.Tina()), 0755)
		}
		os.Exit(0)
	}

	p.SetSliced(hlnetUseComplexPNames)

	if hlnetVerbose {
		p.SetVerbose(pnml.MINIMAL)
	} else {
		p.SetVerbose(pnml.QUIET)
	}

	p.SetFES(false)
	hl, err := hlnet.Build(p)
	if err != nil {
		hlnetLogger.Println("Error decoding PNML file:", err)
		os.Exit(1)
		return
	}

	// We try to build a TPN first
	cn, nbcopies, listlr, listh, err := corenet.BuildTPN(p, hl)
	if err == nil {
		// hlnetLogger.Println("file " + outfile + " is a TPN")
		if hlnetStat {
			elapsed := time.Since(start)
			npl, ntr, narcs := cn.Statistics()
			fmt.Fprintf(os.Stdout, "%d place(s), %d transition(s), %d arc(s), %.3fs\n", nbcopies*npl, nbcopies*(ntr-len(listlr)), nbcopies*narcs, elapsed.Seconds())
			return
		}
		if outfile == "-" {
			os.Stdout.Write([]byte(cn.PrintTPN(nbcopies, listlr, listh)))
			return
		}
		outfile = outfile + ".tpn"
		ioutil.WriteFile(outfile, []byte(cn.PrintTPN(nbcopies, listlr, listh)), 0755)
		return
	}

	cn = corenet.Build(p, hl)
	// hlnetLogger.Println("file " + outfile + " is a NET")
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

// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

import (
	"fmt"
	"io"
	"sort"

	"github.com/dalzilio/mcc/pnml"
)

// ----------------------------------------------------------------------

func (pl Place) Write(w io.Writer) {
	if pl.init == 0 {
		// it is not necessary to print a place when it is not initially marked
		// fmt.Fprintf(w, "pl %s\n", pl.name)
		return
	}
	fmt.Fprintf(w, "pl %s (%d)\n", pl.name, pl.init)
}

func (pl corep) Write(w io.Writer) {
	if pl.int == 1 {
		fmt.Fprintf(w, " %s", pl.name)
		return
	}
	fmt.Fprintf(w, " %s*%d", pl.name, pl.int)
}

func (tr Trans) Write(w io.Writer, k int, verbosity pnml.VERB) {
	if (verbosity == pnml.QUIET) || (tr.label == "") {
		fmt.Fprintf(w, "tr t%d ", k)
	} else {
		fmt.Fprintf(w, "tr t%d : {%s} ", k, tr.label)
	}
	for _, v := range tr.in {
		v.Write(w)
	}
	fmt.Fprint(w, " ->")
	for _, v := range tr.out {
		v.Write(w)
	}
	fmt.Fprint(w, "\n")
}

// Write outputs the corenet in .net format on an io.Writer.
func (net Net) Write(w io.Writer) {
	fmt.Fprintf(w, "# net %s has %d places and %d transitions\n", net.name, len(net.pl), len(net.tr))
	fmt.Fprintf(w, "net {%s}\n", net.name)

	// we start by sorting the slice of places
	sort.Slice(net.pl, func(i, j int) bool {
		return net.pl[i].name < net.pl[j].name
	})

	for _, v := range net.pl {
		v.Write(w)
	}

	for k, v := range net.tr {
		v.Write(w, k, net.verbose)
	}
}

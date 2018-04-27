// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

import (
	"fmt"
	"io"
)

// ----------------------------------------------------------------------

func (pl Place) Write(w io.Writer) {
	if pl.init == 0 {
		fmt.Fprintf(w, "pl %s\n", pl.name)
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

func (tr Trans) Write(w io.Writer) {
	fmt.Fprintf(w, "tr t%d ", tr.count)
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
	fmt.Fprintf(w, "#net %s has %d places and %d transitions\n", net.name, len(net.pl), len(net.tr))
	fmt.Fprintf(w, "net {%s}\n", net.name)

	for _, v := range net.pl {
		v.Write(w)
	}
	for _, v := range net.tr {
		v.Write(w)
	}
}

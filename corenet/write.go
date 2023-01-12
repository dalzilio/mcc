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

// Write prints a place declaration.
func (pl Place) Write(w io.Writer, issliced bool, verbosity pnml.VERB) {
	//  If the net is sliced, its name already give
	// some info about the colored place it originates from. In this case, we do
	// not need to print the place when it is not initially marked. Otherwise,
	// unless the verbosity is MINIMAL, we use the label to keep the relation
	// between the corenet place and the colored one.
	if (verbosity == pnml.MINIMAL) || issliced {
		if pl.init != 0 {
			fmt.Fprintf(w, "pl %s (%d)\n", pl.name, pl.init)
		}
		return
	}
	if pl.init == 0 {
		fmt.Fprintf(w, "pl %s : {%s}\n", pl.name, pl.label)
		return
	}
	fmt.Fprintf(w, "pl %s : {%s} (%d)\n", pl.name, pl.label, pl.init)
}

func (pl corep) Write(w io.Writer) {
	if pl.int == 1 {
		fmt.Fprintf(w, " %s", pl.name)
		return
	}
	fmt.Fprintf(w, " %s*%d", pl.name, pl.int)
}

func (tr Trans) Write(w io.Writer, k int, verbosity pnml.VERB) {
	if (verbosity == pnml.MINIMAL) || (tr.label == "") {
		fmt.Fprintf(w, "tr t%d ", k)
	} else if verbosity == pnml.SKELETON {
		fmt.Fprintf(w, "tr {%s} ", tr.label)
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

	// we start by sorting the slice of places. In the case where the result is
	// not "sliced", place names are all of the form p_k, with k an integer, and
	// are already sorted (numerically and not lexicographically, which is
	// better). We do nothing in this case.
	if net.sliced {
		sort.Slice(net.pl, func(i, j int) bool {
			return net.pl[i].name < net.pl[j].name
		})
	}

	// // we print out properties if needed. We use the fact that places are sorted
	// // by names. Hence (core) places corresponding to the same colored place are
	// // grouped together. Same for transitions.
	// if net.printprops {
	// 	// output list of places for each colored one
	// 	currentname := ""
	// 	for _, v := range net.pl {
	// 		if v.label != currentname {
	// 			currentname = v.label
	// 			fmt.Fprintf(w, "\n# pl %s", v.label)
	// 		}
	// 		fmt.Fprintf(w, " %s", v.name)
	// 	}
	// 	// output list of transitions for each colored one
	// 	currentname = ""
	// 	for k, v := range net.tr {
	// 		if v.label != currentname {
	// 			currentname = v.label
	// 			fmt.Fprintf(w, "\n# tr %s", currentname)
	// 		}
	// 		fmt.Fprintf(w, " t%d", k)
	// 	}
	// 	// find places that are not used in any transitions, they give rise to
	// 	// simple invariants: (m(p) == m0(p)). We visit transitions and mark
	// 	// places that are being used. We also use a counter to speed things up
	// 	// if we find that all places are used early.
	// 	usedpl := make([]bool, len(net.pl))
	// 	cusedpl := len(net.pl)
	// 	for _, v := range net.tr {
	// 		for _, c := range v.in {
	// 			if index := c.Place.count; !usedpl[index] {
	// 				usedpl[index] = true
	// 				cusedpl--
	// 			}
	// 			for _, c := range v.out {
	// 				if index := c.Place.count; !usedpl[index] {
	// 					usedpl[index] = true
	// 					cusedpl--
	// 				}
	// 			}
	// 		}
	// 		if cusedpl == 0 {
	// 			break
	// 		}
	// 	}
	// 	if cusedpl != 0 {
	// 		for k, isused := range usedpl {
	// 			if !isused {
	// 				if net.pl[k].init == 0 {
	// 					fmt.Fprintf(w, "\n# inv %s == %d", net.pl[k].name, net.pl[k].init)
	// 				}
	// 			}
	// 		}
	// 		for k, isused := range usedpl {
	// 			if !isused && (net.pl[k].init == 0) {
	// 				if net.pl[k].init == 0 {
	// 					fmt.Fprintf(w, "\npl %s", net.pl[k].name)
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	// fmt.Fprint(w, "\n")

	// Finally we output the .net declarations
	for _, v := range net.pl {
		v.Write(w, net.sliced, net.verbose)
	}

	for k, v := range net.tr {
		v.Write(w, k, net.verbose)
	}
}

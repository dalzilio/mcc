// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

import "github.com/dalzilio/mcc/pnml"

// Net is the type of (core) Petri nets.
type Net struct {
	verbose pnml.VERB
	sliced  bool
	name    string
	pl      []*Place
	tr      []*Trans
}

// Place is the type of places in a (core) net.
type Place struct {
	count int
	name  string
	label string
	init  int
}

// corep is a pair of a place and a multiplicity. This is used to build arcs in
// the unfolding of a hlnet.
type corep struct {
	*Place
	int
}

// Trans is the type of transitions in a (core) net.
type Trans struct {
	count   int
	label   string
	in, out []corep
}

// Statistics returns the number of places, transitions and arcs in the net
func (net *Net) Statistics() (int, int, int) {
	narcs := 0
	for _, t := range net.tr {
		narcs = narcs + len(t.in) + len(t.out)
	}
	return len(net.pl), len(net.tr), narcs
}

// Copyright 2023. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

import (
	"sort"

	"github.com/dalzilio/mcc/hlnet"
	"github.com/dalzilio/mcc/pnml"
)

// BuildSkeleton returns a core net corresponding to the skeleton of the colored
// Petri net hlnet. All naming information is kept in the labels.
func BuildSkeleton(pnet *pnml.Net, hl *hlnet.Net) *Net {
	var net Net
	net.name = pnet.Name
	net.verbose = pnml.SKELETON
	net.printprops = false

	// we start by building a map between hlnet places and corenet places
	skelpl := make(map[*hlnet.Place]*Place)

	for plname, p := range hl.Places {
		cp := Place{count: 0, name: "", label: escape2aname(plname), init: 0}
		if p.Init != nil {
			cp.init = p.Init.Skeletonize(pnet)
		}
		skelpl[p] = &cp
	}

	// then we build the slice of transitions and sort them to have a
	// deterministic output
	net.pl = []*Place{}
	for _, p := range skelpl {
		net.pl = append(net.pl, p)
	}
	sort.Slice(net.pl, func(i, j int) bool {
		return net.pl[i].label < net.pl[j].label
	})

	// we also have one transition in the skeleton for each transition in the
	// hlnet. We never test if the guard of a transition in the hlnet is
	// satisfiable (which should be the case in practice) and always consider
	// that the condition is true.
	for k, t := range hl.Trans {
		ct := Trans{count: 0, label: escape2aname(k)}
		for _, e := range t.Arcs {
			cp := skelpl[e.Place]
			if e.Kind == hlnet.IN {
				ct.in = append(ct.in, corep{cp, e.Pattern.Skeletonize(pnet)})
			} else {
				ct.out = append(ct.out, corep{cp, e.Pattern.Skeletonize(pnet)})
			}
		}
		// we sort the input and output places to have a consistent,
		// reproducible output
		sort.Slice(ct.in, func(i, j int) bool {
			return ct.in[i].label < ct.in[j].label
		})
		sort.Slice(ct.out, func(i, j int) bool {
			return ct.out[i].label < ct.out[j].label
		})
		net.tr = append(net.tr, &ct)
	}
	// we also sort the transitions. This is easier than for an unfolded net
	// because we only have one transition in the skeleton for each transition
	// on the hlnet and the names are necessarily distinct.
	sort.Slice(net.tr, func(i, j int) bool {
		return net.tr[i].label < net.tr[j].label
	})
	for k, v := range net.tr {
		v.count = k
	}
	return &net
}

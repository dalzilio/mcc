// Copyright 2023. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

import (
	"sort"
	"strings"

	"github.com/dalzilio/mcc/hlnet"
	"github.com/dalzilio/mcc/pnml"
)

// BuildSkeleton returns a core net corresponding to the skeleton of the colored
// Petri net hlnet.
func BuildSkeleton(pnet *pnml.Net, hl *hlnet.Net) *Net {
	var net Net
	net.name = pnet.Name
	net.verbose = pnml.SKELETON
	net.sliced = false
	net.printprops = false

	// we start by building a map between hlnet Place and corenet Place
	pcount := 0
	skelpl := make(map[*hlnet.Place]*Place)

	for k, p := range hl.Places {
		cp := Place{count: pcount, name: normalize2aname(k), label: "", init: 0}
		// we sum the initial marking over all possible colors
		if p.Init != nil {
			// _, multv := p.Init.Match(pnet, nil)
			// for _, v := range multv {
			// 	cp.init += v
			// }
			cp.init = p.Init.Skeletonize(pnet)
		}
		skelpl[p] = &cp
		pcount++
	}

	// then we build the slice of transitions and sort them to have a
	// deterministic output
	net.pl = []*Place{}
	for _, p := range skelpl {
		net.pl = append(net.pl, p)
	}
	sort.Slice(net.pl, func(i, j int) bool {
		return strings.Compare(net.pl[i].name, net.pl[j].name) < 0
	})

	// we also have one transition in the skeleton for each transition in the
	// hlnet. We never test if the guard of a transition in the hlnet is
	// satisfiable (which should be the case in practice) and always consider
	// that the condition is true.
	tcount := 0
	for k, t := range hl.Trans {
		ct := Trans{count: tcount, label: normalize2aname(k)}
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
			return strings.Compare(ct.in[i].name, ct.in[j].name) < 0
		})
		sort.Slice(ct.out, func(i, j int) bool {
			return strings.Compare(ct.out[i].name, ct.out[j].name) < 0
		})
		net.tr = append(net.tr, &ct)
		tcount++
	}
	// we also sort the transitions. This is easier than for an unfolded net
	// because we only have one transition in the skeleton for each transition
	// on the hlnet and the names are necessarily distinct.
	sort.Slice(net.tr, func(i, j int) bool {
		return strings.Compare(net.tr[i].label, net.tr[j].label) < 0
	})
	for k, v := range net.tr {
		v.count = k
	}
	return &net
}

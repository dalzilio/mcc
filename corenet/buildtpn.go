// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

import (
	"fmt"
	"sort"

	"github.com/dalzilio/mcc/hlnet"
	"github.com/dalzilio/mcc/pnml"
)

// BuildTPN returns a core net using the TPN syntax.
func BuildTPN(pnet *pnml.Net, hl *hlnet.Net) (*Net, int, []int, []int, error) {
	// hl := hlnet.Build(pnet)

	var net Net
	net.name = pnet.Name
	net.verbose = pnet.VERBOSE

	// we find the name of the RING variables
	if r := len(pnet.Env); r != 1 {
		return nil, 0, nil, nil, fmt.Errorf("only one ring variable supported, found %d", r)
	}
	var varname string
	var nbcopies int
	for k := range pnet.Env {
		varname = k
		nbcopies = len(pnet.World[pnet.Env[k]])
	}

	// we build all the places in the final net. They are the same than in the
	// hlnet. We keep an association between hlnet.Place and corenet.Place in a
	// map. We list the places by alphabetical order on their ID.
	cpl := make(map[*hlnet.Place]*Place)
	pcount := 0
	hlp := make([]string, len(hl.Places))
	i := 0
	for k := range hl.Places {
		hlp[i] = k
		i++
	}
	sort.Strings(hlp)
	for _, k := range hlp {
		p := hl.Places[k]
		cp := Place{count: pcount, label: k}
		pcount++
		net.pl = append(net.pl, &cp)
		cpl[p] = &cp
		if p.Init != nil {
			mult, err := p.Init.MatchRing(varname, 0)
			if err != nil {
				return nil, 0, nil, nil, err
			}
			cp.init = mult
		}
	}

	// we go through all the transitions; we stop if one of them as a non-nil
	// condition. We build arcs by matching the current identifier. We build
	// synchronization by looking at the match with predecessor(varname). We
	// keep a slice of transition `count` for the left, right transitions that
	// need to be exported. We list the transitions by alphabetical order on
	// their ID.
	var listlr []int
	var listh []int
	tcount := 0
	hlt := make([]string, len(hl.Trans))
	j := 0
	for k := range hl.Trans {
		hlt[j] = k
		j++
	}
	sort.Strings(hlt)
	for _, k := range hlt {
		t := hl.Trans[k]
		if t.Cond.Op != pnml.NIL {
			return nil, 0, nil, nil, fmt.Errorf("we do not support conditions on transitions")
		}
		var inl, inr, outl, outr []corep
		for _, e := range t.Arcs {
			mr, err := e.Pattern.MatchRing(varname, 0)
			if err != nil {
				return nil, 0, nil, nil, err
			}
			cp := corep{Place: cpl[e.Place], int: mr}
			if mr != 0 {
				if e.Kind == hlnet.IN {
					inr = append(inr, cp)
				} else {
					outr = append(outr, cp)
				}
			}
			ml, _ := e.Pattern.MatchRing(varname, -1)
			cp = corep{Place: cpl[e.Place], int: ml}
			if ml != 0 {
				if e.Kind == hlnet.IN {
					inl = append(inl, cp)
				} else {
					outl = append(outl, cp)
				}
			}
		}
		// if we have an 'here' transition
		if inl == nil && outl == nil {
			listh = append(listh, tcount)
			net.tr = append(net.tr, &Trans{count: tcount, label: "h-" + k, in: inr, out: outr})
			tcount++
		} else {
			listlr = append(listlr, tcount)
			net.tr = append(net.tr, &Trans{count: tcount, label: "l-" + k, in: inl, out: outl})
			tcount++
			net.tr = append(net.tr, &Trans{count: tcount, label: "r-" + k, in: inr, out: outr})
			tcount++
		}
	}

	return &net, nbcopies, listlr, listh, nil
}

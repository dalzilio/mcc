// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dalzilio/mcc/hlnet"
	"github.com/dalzilio/mcc/pnml"
)

// ----------------------------------------------------------------------

// makepname returns a (corenet) place name using either a counter or, if we
// want sliced names, a combination of the name of the colored place, pname, and
// a value of the current type.
func makepname(net *pnml.Net, pname string, count int, val *pnml.Value) string {
	if !net.SLICED {
		return fmt.Sprintf("p_%d", count)
	}
	s := strings.Builder{}
	if net.PrintProperties {
		// When we print properties we keep the original name of the PNML place,
		// which may contain illegal characters. This is the case for example
		// with model neoelection that uses `-` in the id of places.
		pname = normalize2aname(pname)
	}
	s.WriteString(pname)
	if val.Head == 0 {
		return s.String()
	}
	if net.VERBOSE == pnml.QUIET {
		fmt.Fprintf(&s, "_%d", val.Head)
		for v := val.Tail; v != nil; v = v.Tail {
			fmt.Fprintf(&s, "_%d", v.Head)
		}
		return s.String()
	}
	fmt.Fprintf(&s, "_%s", net.Identity[val.Head])
	for v := val.Tail; v != nil; v = v.Tail {
		fmt.Fprintf(&s, "_%s", net.Identity[v.Head])
	}
	return s.String()
}

func makeplabel(net *pnml.Net, name string, val *pnml.Value) string {
	if net.VERBOSE == pnml.QUIET {
		if net.PrintProperties {
			return name
		}
		return ""
	}

	s := fmt.Sprintf("%s_%s", name, net.Identity[val.Head])
	if val.Tail == nil {
		return s
	}
	return makeplabel(net, s+"x", val.Tail)
}

func maketlabel(net *pnml.Net, name string, env pnml.Env) string {
	if net.VERBOSE == pnml.QUIET {
		if net.PrintProperties {
			return name
		}
		return ""
	}

	if net.VERBOSE != pnml.MAXIMAL {
		return normalize2aname(name)
	}

	return fmt.Sprintf("%s %s", name, net.PrintEnv(env))
}

// ----------------------------------------------------------------------

type coreAssoc struct {
	place *hlnet.Place
	val   *pnml.Value
}

// ----------------------------------------------------------------------

func appendCorep(in []corep, c corep) []corep {
	// We assume that the list of corep is sorted by Place.count. We find the
	// smallest index i where c should be inserted
	i := sort.Search(len(in), func(j int) bool {
		return in[j].Place.count >= c.Place.count
	})
	if i < len(in) && in[i].Place == c.Place {
		in[i].int += c.int
		return in
	}
	// special trick for inserting
	in = append(in, corep{})
	copy(in[i+1:], in[i:])
	in[i] = c
	return in
}

// ----------------------------------------------------------------------

// normalize2aname returns an identifier that can be used in a .net file. We do
// not ensure that identifiers names will not clash with each other (but it
// should)
func normalize2aname(s string) string {
	anamize := func(r rune) rune {
		switch {
		case (r >= 'A' && r <= 'z') || (r >= '0' && r <= '9'):
			return r
		default:
			return '_'
		}
	}
	return strings.Map(anamize, s)
}

// Build returns a core net from a colored Petri net  by unfolding the
// corresponding hlnet.
func Build(pnet *pnml.Net, hl *hlnet.Net) *Net {
	var net Net
	net.name = pnet.Name
	net.verbose = pnet.VERBOSE
	net.sliced = pnet.SLICED
	net.printprops = pnet.PrintProperties

	// we build all the places in the final net. They are of the form p x val,
	// where val is one of the possible values from the type of p. We build a
	// map to find the given place from the pair {p val}
	cpl := make(map[coreAssoc]*Place)
	pcount := 0
	pname := ""
	for k, p := range hl.Places {
		if pnet.PrintProperties {
			pname = k
		} else {
			pname = normalize2aname(k)
		}
		if p.Stable {
			// when the place is stable, its rechable "values" are the one in
			// its initial marking (moreover the marking of the place is an
			// invariant). We still keep those places in the net in order to have the
			// right value for "maximal number of tokens in a marking" but we do
			// not need to add the edges (this should speed up the computation
			// of invariants).
			initv, multv := p.Init.Match(pnet, nil)
			for k, v := range initv {
				cp := Place{count: pcount, name: makepname(pnet, pname, pcount, v), label: makeplabel(pnet, pname, v)}
				cp.init = multv[k]
				pcount++
				cpl[coreAssoc{place: p, val: v}] = &cp
				net.pl = append(net.pl, &cp)
			}
		} else {
			// the possible values of p are the one in its type
			for _, v := range pnet.World[p.Type] {
				cp := Place{count: pcount, name: makepname(pnet, pname, pcount, v), label: makeplabel(pnet, pname, v)}
				pcount++
				cpl[coreAssoc{place: p, val: v}] = &cp
				net.pl = append(net.pl, &cp)
			}
			if p.Init != nil {
				initv, multv := p.Init.Match(pnet, nil)
				for k, v := range initv {
					cp := cpl[coreAssoc{place: p, val: v}]
					cp.init = multv[k]
				}
			}
		}
	}

	// we go through all the transitions and build coretrans by enumerating all
	// the possible association of variables and values, testing if the
	// condition is true. iterator[i] gives the value (index) we are currently
	// considering for variable varnames[i].
	tcount := 0
	for k, t := range hl.Trans {
		for iter := mkiter(pnet, cpl, t); iter.hasNext(); {
			if ct, ok := iter.check(); ok {
				ct.count = tcount
				ct.label = maketlabel(pnet, k, iter.env)
				if net.sliced {
					sort.Slice(ct.in, func(i, j int) bool {
						return ct.in[i].name < ct.in[j].name
					})
					sort.Slice(ct.out, func(i, j int) bool {
						return ct.out[i].name < ct.out[j].name
					})
				}
				net.tr = append(net.tr, ct)
				tcount++
			}
		}
	}

	// we also sort the transitions when net.sliced is true
	if net.sliced {
		sort.Slice(net.tr, func(i, j int) bool {
			b := strings.Compare(net.tr[i].label, net.tr[j].label)
			if b != 0 {
				return b < 0
			}
			k := 0
			for k < len(net.tr[i].in) && k < len(net.tr[j].in) {
				if net.tr[i].in[k] == net.tr[j].in[k] {
					k++
					continue
				}
				if net.tr[i].in[k].name == net.tr[j].in[k].name {
					return (net.tr[i].in[k].int < net.tr[j].in[k].int)
				}
				return (net.tr[i].in[k].name < net.tr[j].in[k].name)
			}
			if k < len(net.tr[i].in) {
				return false
			}
			if k < len(net.tr[j].in) {
				return true
			}
			k = 0
			for k < len(net.tr[i].out) && k < len(net.tr[j].out) {
				if net.tr[i].out[k] == net.tr[j].out[k] {
					k++
					continue
				}
				if net.tr[i].out[k].name == net.tr[j].out[k].name {
					return (net.tr[i].out[k].int < net.tr[j].out[k].int)
				}
				return (net.tr[i].out[k].name < net.tr[j].out[k].name)
			}
			if k < len(net.tr[i].out) {
				return false
			}
			return true
		})
		// sort.Slice(net.tr, func(i, j int) bool {
		// 	k := 0
		// 	for k < len(net.tr[i].in) && k < len(net.tr[j].in) {
		// 		if net.tr[i].in[k] == net.tr[j].in[k] {
		// 			k++
		// 			continue
		// 		}
		// 		if net.tr[i].in[k].name == net.tr[j].in[k].name {
		// 			return (net.tr[i].in[k].int < net.tr[j].in[k].int)
		// 		}
		// 		return (net.tr[i].in[k].name < net.tr[j].in[k].name)
		// 	}
		// 	if k < len(net.tr[i].in) {
		// 		return false
		// 	}
		// 	if k < len(net.tr[j].in) {
		// 		return true
		// 	}
		// 	k = 0
		// 	for k < len(net.tr[i].out) && k < len(net.tr[j].out) {
		// 		if net.tr[i].out[k] == net.tr[j].out[k] {
		// 			k++
		// 			continue
		// 		}
		// 		if net.tr[i].out[k].name == net.tr[j].out[k].name {
		// 			return (net.tr[i].out[k].int < net.tr[j].out[k].int)
		// 		}
		// 		return (net.tr[i].out[k].name < net.tr[j].out[k].name)
		// 	}
		// 	if k < len(net.tr[i].out) {
		// 		return false
		// 	}
		// 	return true
		// })
		// and we also reflect the new ordering of transitions in the count field
		for k, v := range net.tr {
			v.count = k
		}

	}

	return &net
}

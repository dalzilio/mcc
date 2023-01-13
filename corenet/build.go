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
func makepname(net *pnml.Net, pname string, val *pnml.Value) string {
	if net.VERBOSE == pnml.SMPT {
		return ""
	}

	s := strings.Builder{}
	if net.VERBOSE == pnml.SLICED {
		s.WriteString(normalize2aname(pname))
	} else {
		// we are in the case pnml.INFO
		s.WriteString(pname)
	}

	if val.Head == 0 {
		return s.String()
	}

	// when verbosity is SLICED, identifiers are of the kind `Id__Name_1_4`,
	// where `Id__Name` is a normalized COL identifier and color constants are
	// encoded using integers.
	if net.VERBOSE == pnml.SLICED {
		fmt.Fprintf(&s, "_%d", val.Head)
		for v := val.Tail; v != nil; v = v.Tail {
			fmt.Fprintf(&s, "_%d", v.Head)
		}
		return s.String()
	}

	// otherwise we are in the case INFO, where identifiers are of the kind
	// `Id_c0_s3` where `c0` and `s3` describe the COL constants.
	fmt.Fprintf(&s, "_%s", net.Identity[val.Head])
	for v := val.Tail; v != nil; v = v.Tail {
		fmt.Fprintf(&s, "_%s", net.Identity[v.Head])
	}
	return escape2aname(s.String())
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
		case (r >= 'A' && r <= 'z') ||
			(r >= '0' && r <= '9') ||
			(r == '\''):
			return r
		default:
			return '_'
		}
	}
	return strings.Map(anamize, s)
}

// escape2aname takes an identifier and escape it (with braces) if it is not a
// valid .net identifier
func escape2aname(s string) string {
	for _, c := range s {
		switch {
		case (c >= 'A' && c <= 'z') ||
			(c >= '0' && c <= '9') ||
			(c == '\'') ||
			(c == '_'):
			continue
		default:
			return fmt.Sprintf("{%s}", s)
		}
	}
	return s
}

// Build returns a core net from a colored Petri net by unfolding the
// corresponding hlnet.
func Build(pnet *pnml.Net, hl *hlnet.Net) *Net {
	var net Net
	net.name = pnet.Name
	net.verbose = pnet.VERBOSE
	net.printprops = pnet.PrintProperties

	// we build all the places in the final net. They are of the form p x val,
	// where val is one of the possible values from the type of p. We build a
	// map to find the given place from the pair {p val}
	cpl := make(map[coreAssoc]*Place)
	for plname, p := range hl.Places {
		if p.Stable {
			// when the place is stable, its reachable "values" are the one in
			// its initial marking (moreover the marking of the place is an
			// invariant). We still keep those places in the net in order to have the
			// right value for "maximal number of tokens in a marking" but we do
			// not need to add the edges (this should speed up the computation
			// of invariants).
			initv, multv := p.Init.Match(pnet, nil)
			for k, v := range initv {
				cp := Place{name: makepname(pnet, plname, v), label: plname}
				cp.init = multv[k]
				cpl[coreAssoc{place: p, val: v}] = &cp
				net.pl = append(net.pl, &cp)
			}
		} else {
			// the possible values of p are the one in its type
			for _, v := range pnet.World[p.Type] {
				cp := Place{name: makepname(pnet, plname, v), label: plname}
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

	// we sort the places and instantiate their count field accordingly. We also
	// reflect the new ordering of transitions in the count field. It should not
	// be possible to call function Build with verbosity SKELETON
	switch net.verbose {
	case pnml.INFO: // we can sort using the name
		sort.Slice(net.pl, func(i, j int) bool {
			return net.pl[i].name < net.pl[j].name
		})
		for k, v := range net.pl {
			v.count = k
		}
	case pnml.SKELETON:
		panic("should not call Build with SKELETON")
	case pnml.SLICED: // similar to INFO
		sort.Slice(net.pl, func(i, j int) bool {
			return net.pl[i].name < net.pl[j].name
		})
		for k, v := range net.pl {
			v.count = k
		}
	case pnml.SMPT: // we use a stable sort on the label and also instantiate the name
		sort.SliceStable(net.pl, func(i, j int) bool {
			return net.pl[i].label < net.pl[j].label
		})
		for k, v := range net.pl {
			v.count = k
			v.name = fmt.Sprintf("p%d", k)
		}
	}

	// we go through all the transitions and build coretrans by enumerating all
	// the possible association of variables and values, testing if the
	// condition is true. iterator[i] gives the value (index) we are currently
	// considering for variable varnames[i].
	for trname, t := range hl.Trans {
		for iter := mkiter(pnet, cpl, t); iter.hasNext(); {
			if ct, ok := iter.check(); ok {
				ct.label = trname
				// we sort the places in the IN and OUT arcs to obtain a
				// deterministic output
				sort.Slice(ct.in, func(i, j int) bool {
					return ct.in[i].count < ct.in[j].count
				})
				sort.Slice(ct.out, func(i, j int) bool {
					return ct.out[i].count < ct.out[j].count
				})
				net.tr = append(net.tr, ct)
			}
		}
	}

	// we sort the transitions using their labels, and a stable sort, and
	// instantiate their count field to have a deterministic output. This is the
	// same algorithm with all the different verbosity level.
	sort.SliceStable(net.tr, func(i, j int) bool {
		if net.tr[i].label == net.tr[j].label {
			// transitions with the same label have the same arc cardinality
			for k := range net.tr[i].in {
				if net.tr[i].in[k].count == net.tr[j].in[k].count {
					continue
				}
				return net.tr[i].in[k].count < net.tr[j].in[k].count
			}
			for k := range net.tr[i].out {
				if net.tr[i].out[k].count == net.tr[j].out[k].count {
					continue
				}
				return net.tr[i].out[k].count < net.tr[j].out[k].count
			}
		}
		return net.tr[i].label < net.tr[j].label
	})
	for k, v := range net.tr {
		v.count = k
	}

	return &net
}

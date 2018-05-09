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

func makepname(net *pnml.Net, pname string, count int, hlpcount int, val *pnml.Value) string {
	if !net.SLICED {
		return fmt.Sprintf("p_%d", count)
	}
	s := strings.Builder{}
	fmt.Fprintf(&s, "%s_%d", pname, val.Head)
	for v := val.Tail; v != nil; v = v.Tail {
		fmt.Fprintf(&s, "_%d", v.Head)
	}
	return s.String()
}

func makeplabel(net *pnml.Net, name string, val *pnml.Value) string {
	if net.VERBOSE == pnml.QUIET {
		return ""
	}

	if net.VERBOSE != pnml.MAXIMAL {
		return fmt.Sprintf("%s", name)
	}

	s := fmt.Sprintf("%s %s", name, net.Identity[val.Head])
	if val.Tail == nil {
		return s
	}
	return makeplabel(net, s+" x", val.Tail)
}

func maketlabel(net *pnml.Net, name string, env pnml.Env) string {
	if net.VERBOSE == pnml.QUIET {
		return ""
	}

	if net.VERBOSE != pnml.MAXIMAL {
		return name
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

	// we build all the places in the final net. They are of the form p x val,
	// where val is one of the possible values from the type of p. We build a
	// map to find the given place from the pair {p val}
	cpl := make(map[coreAssoc]*Place)
	pcount := 0
	hlpcount := 0
	for k, p := range hl.Places {
		for _, v := range pnet.World[p.Type] {
			pname := normalize2aname(k)
			cp := Place{count: pcount, name: makepname(pnet, pname, pcount, hlpcount, v), label: makeplabel(pnet, pname, v)}
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
		hlpcount++
	}

	// we go through all the transitions and build coretrans by enumerating all
	// the possible association of variables and values, testing if the
	// condition is true.
	tcount := 0
	for k, t := range hl.Trans {
		env := t.Env
		varnames, iterator, enums := mkiter(pnet, env)
		for {
			if t.Cond.OK(pnet, env) {
				ct := Trans{count: tcount, label: maketlabel(pnet, k, env)}
				sat := true
				for _, e := range t.Arcs {
					if e.Pattern == nil {
						sat = false
						break
					}
					f, m := e.Pattern.Match(pnet, env)
					if f == nil {
						sat = false
						break
					}
					for i := range f {
						place := cpl[coreAssoc{place: e.Place, val: f[i]}]
						if e.Kind == hlnet.IN {
							ct.in = appendCorep(ct.in, corep{Place: place, int: m[i]})
						} else {
							ct.out = appendCorep(ct.out, corep{Place: place, int: m[i]})
						}
					}
				}
				if sat {
					// if pnet.VERBOSE > pnml.QUIET {
					// 	sort.Slice(ct.in, func(i, j int) bool {
					// 		return ct.in[i].count < ct.in[j].count
					// 	})
					// 	sort.Slice(ct.out, func(i, j int) bool {
					// 		return ct.out[i].count < ct.out[j].count
					// 	})
					// }
					net.tr = append(net.tr, &ct)
					tcount++
				}
			}
			if ok := nextiter(pnet, env, varnames, iterator, enums); !ok {
				break
			}
		}
	}

	return &net
}

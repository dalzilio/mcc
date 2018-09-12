// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package hlnet

import "fmt"
import "github.com/dalzilio/mcc/pnml"

func (net Net) String() string {
	// we start by collecting the names of the Places
	secalp := make(map[*Place]string)
	for k, v := range net.Places {
		secalp[v] = k
	}

	s := fmt.Sprintf("# net %s\n", net.Name)
	for k, v := range net.Places {
		s += fmt.Sprintf("# pl %s %s\n", k, pInit(v.Init))
	}
	for k, v := range net.Trans {
		s += fmt.Sprintf("# tr %s %s %s\n", k, v.Cond, v.Env)
		for _, e := range v.Arcs {
			if e.Kind == IN {
				s += fmt.Sprintf("#\t%s -->( %s )\n", secalp[e.Place], e.Pattern)
			} else {
				s += fmt.Sprintf("#\t%s <--( %s )\n", secalp[e.Place], e.Pattern)
			}
		}
	}
	return s
}

// ----------------------------------------------------------------------

// Tina outputs an hlnet in a format that can be displayed by Tina' s nd tool.
func (net Net) Tina() string {
	// we start by collecting the names of the Places
	secalp := make(map[*Place]string)
	for k, v := range net.Places {
		secalp[v] = k
	}

	s := fmt.Sprintf("net {%s}\n", net.Name)
	for k, v := range net.Places {
		isstable := ""
		if v.Stable {
			isstable = "(stable) "
		}
		if is := pInit(v.Init); is != "-" {
			s += fmt.Sprintf("pl {%s} : {%s%s} (1)\n", k, isstable, is)
		} else {
			s += fmt.Sprintf("pl {%s}\n", k)
		}
	}
	for k, v := range net.Trans {
		if v.Cond.Op == pnml.NIL {
			s += fmt.Sprintf("tr {%s} : {%s} ", k, v.Env)
		} else {
			s += fmt.Sprintf("tr {%s} : {%s %s} ", k, v.Cond, v.Env)
		}
		for _, e := range v.Arcs {
			if e.Kind == IN {
				s += fmt.Sprintf(" {%s}", secalp[e.Place])
			}
		}
		s += " -> "
		for _, e := range v.Arcs {
			if e.Kind == OUT {
				s += fmt.Sprintf(" {%s}", secalp[e.Place])
			}
		}
		s += "\n"
		// we output the Pattern of the edges as a comment because it is not
		// possible to draw them with Tina's nd tool.
		for _, e := range v.Arcs {
			if e.Pattern != nil {
				s += fmt.Sprintf("#\t--( %s )-- %s\n", e.Pattern, secalp[e.Place])
			}
		}
	}
	return s
}

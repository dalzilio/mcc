// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package hlnet

import (
	"fmt"

	"github.com/dalzilio/mcc/pnml"
)

// Net is the concrete type of symmetric nets.
type Net struct {
	Name   string
	Places map[string]*Place
	Trans  map[string]*Transition
}

// Place is the concrete type of symmetric nets places.
type Place struct {
	Init   pnml.Expression
	Type   string
	Stable bool
}

// Transition is the concrete type of symmetric nets transitions.
type Transition struct {
	Cond pnml.Operation
	Env  pnml.Env
	Arcs []*Arcs
}

// Arcs is the concrete type of symmetric nets arcs.
type Arcs struct {
	Kind    ARC
	Pattern pnml.Expression
	Place   *Place
}

// ----------------------------------------------------------------------

func pInit(p pnml.Expression) string {
	if p == nil {
		return "-"
	}
	return p.String()
}

// ----------------------------------------------------------------------

// Build returns an hlnet from a PNML net. This structure is easier to deal
// with.
func Build(n *pnml.Net) (*Net, error) {

	var net = Net{Name: n.Name}

	net.Places = make(map[string]*Place)
	for _, p := range n.Page.Places {
		net.Places[p.ID] = &Place{Init: p.InitialMarking, Type: p.Type.ID}
	}

	net.Trans = make(map[string]*Transition)
	for _, t := range n.Page.Trans {
		env := make(pnml.Env)
		var cond pnml.Operation
		if t.Condition == nil {
			cond = pnml.Operation{Op: pnml.NIL}
		} else {
			cond = t.Condition.(pnml.Operation)
		}
		cond.AddEnv(env)
		for varname := range env {
			if _, ok := n.Env[varname]; !ok {
				return nil, fmt.Errorf("variable \"%s\" used in condition of transition %s was never declared", varname, t.ID)
			}
		}
		net.Trans[t.ID] = &Transition{Cond: cond, Env: env}
	}

	for _, a := range n.Page.Arcs {
		e := Arcs{Pattern: a.Pattern}
		if e.Pattern == nil {
			e.Pattern = pnml.Operation{Op: pnml.NIL}
		}
		env := make(pnml.Env)
		e.Pattern.AddEnv(env)
		for varname := range env {
			if _, ok := n.Env[varname]; !ok {
				return nil, fmt.Errorf("variable \"%s\" used in pattern of arc from %s to %s was never declared", varname, a.Source, a.Target)
			}
		}
		if p, ok := net.Places[a.Source]; ok {
			// arc source is a place, target is a transition. The edge is of
			// kind IN. We add the variables in the pattern to env.
			t := net.Trans[a.Target]
			e.Pattern.AddEnv(t.Env)
			e.Place = p
			e.Kind = IN
			t.Arcs = append(t.Arcs, &e)
		}
		if p, ok := net.Places[a.Target]; ok {
			// arc source is a transition, target is a place. The edge is of
			// kind OUT.
			t := net.Trans[a.Source]
			e.Pattern.AddEnv(t.Env)
			e.Place = p
			e.Kind = OUT
			t.Arcs = append(t.Arcs, &e)
		}
	}

	// if there are less than 5 places we try to find "stable places", that is a
	// place p such that, for every transition t, there is an edge t--exp-->p if
	// and only if there is an edge p--exp-->t. When a place is stable we know
	// that the tokens in place p (the possible pairs p x val) are exactly the
	// same than in the initial marking.
	if len(net.Places) < 6 {
		for _, p := range net.Places {
			p.Stable = net.IsPlaceStable(p)
		}
	}

	return &net, nil
}

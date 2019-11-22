// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package pnml

import (
	"fmt"
)

// MatchRing is used to match a pattern in a net with a RING structure; for use
// in the generation of TPN. It takes as parameter the name of the 'scalar set'
// variable representing the process ID and an increment on this id. So a value
// of 0 means the 'current' process id, -1 means the predecessor, +1 the
// successor, and so on. The result is the multiplicity of the result.
func (p All) MatchRing(name string, i int) (int, error) {
	if i == 0 {
		return 1, nil
	}
	return 0, nil
}

func (p Add) MatchRing(name string, i int) (int, error) {
	var res int
	for j := range p {
		r, err := p[j].MatchRing(name, i)
		if err != nil {
			return 0, err
		}
		res += r
	}
	return res, nil
}

func (p Subtract) MatchRing(name string, i int) (int, error) {
	return 0, fmt.Errorf("element not supported in MatchRing (subtract)")
}

func (p Tuple) MatchRing(name string, i int) (int, error) {
	return 0, fmt.Errorf("element not supported in MatchRing (tuple)")

}

func (p Operation) MatchRing(name string, i int) (int, error) {
	return 0, fmt.Errorf("element not supported in MatchRing (operation)")
}

func (p Constant) MatchRing(name string, i int) (int, error) {
	return 0, fmt.Errorf("element not supported in MatchRing (constant)")
}

// We should be able to treat FIRConstant like All
func (p FIRConstant) MatchRing(name string, i int) (int, error) {
	return 0, fmt.Errorf("element not supported in MatchRing (finiteintrangeconstant)")
}

func (p Var) MatchRing(name string, i int) (int, error) {
	if name != string(p) {
		return 0, fmt.Errorf("only one ring variable supported")
	}
	if i == 0 {
		return 1, nil
	}
	return 0, nil
}

func (p Dot) MatchRing(name string, i int) (int, error) {
	return 0, fmt.Errorf("element not supported in MatchRing (dot)")
}

func (p Successor) MatchRing(name string, i int) (int, error) {
	return p.Var.MatchRing(name, i-p.Incr)
}

func (p Numberof) MatchRing(name string, i int) (int, error) {
	r, err := p.Expression.MatchRing(name, i)
	if err != nil {
		return 0, err
	}
	if r != 0 {
		return p.Mult, nil
	}
	return 0, nil
}

// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package pnml

import (
	"fmt"
)

// Value provides a more efficient representation for values
//
// {0 nil} is a Dot
// {i nil} is Constant(name) where i uniquely identifies name
// {i {j {...}}} is for tuples
type Value struct {
	Head int
	Tail *Value
}

// ----------------------------------------------------------------------

// PrintValue returns a readable description of a Value
func (net *Net) PrintValue(val *Value) string {
	c := net.Identity[val.Head]
	if val.Tail == nil {
		return c
	}
	c = fmt.Sprintf("(%s, %s", c, net.Identity[val.Tail.Head])
	return net.printTupleValue(c, val.Tail.Tail)
}

func (net *Net) printTupleValue(s string, val *Value) string {
	if val == nil {
		return s + ")"
	}
	c := net.Identity[val.Head]
	return net.printTupleValue(s+", "+c, val.Tail)
}

// ----------------------------------------------------------------------

func (net *Net) enumprod(elem []string) []*Value {
	if len(elem) == 1 {
		return net.World[elem[0]]
	}
	head := net.World[elem[0]]
	tail := net.enumprod(elem[1:])

	var list []*Value
	for _, a := range head {
		for _, b := range tail {
			val := Value{Head: a.Head, Tail: b}
			pval, ok := net.unique[val]
			if !ok {
				pval = &val
				net.unique[val] = &val
			}
			list = append(list, pval)
		}
	}
	return list
}

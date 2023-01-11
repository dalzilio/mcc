// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package pnml

import (
	"fmt"
	"log"
	"strconv"
)

// ----------------------------------------------------------------------

// Expression is the interface that wraps the type of PNML expression (we only
// consier symmetric nets).
//
// AddEnv is used to accumulate the free variables in the Expression to an
// existing environment. It can be used to add the free variables of an
// expression to an alreay existing environment.
//
// Match return the set of constant values that match an Expression, also called
// a pattern,  together with their multiplicities. The method returns two slices
// that have equal length.
type Expression interface {
	String() string
	AddEnv(Env)
	Match(*Net, Env) ([]*Value, []int)
	MatchRing(string, int) (int, error)
	Skeletonize(*Net) int
}

// ----------------------------------------------------------------------

func (net *Net) compareExp(p1, p2 *Value, op OP) bool {
	switch op {
	case EQ:
		return p1 == p2
	case GREATTHAN:
		if p1.Head >= 0 {
			return p1.Head > p2.Head
		}
		return p1.Head < p2.Head
	case GREATTHANEQ:
		if p1.Head >= 0 {
			return p1.Head >= p2.Head
		}
		return p1.Head <= p2.Head
	case INEQ:
		return p1 != p2
	case LESSTHAN:
		if p1.Head >= 0 {
			return p1.Head < p2.Head
		}
		return p1.Head > p2.Head
	case LESSTHANEQ:
		if p1.Head >= 0 {
			return p1.Head <= p2.Head
		}
		return p1.Head >= p2.Head
	}
	panic("not reachable in compareExp")
}

// ----------------------------------------------------------------------

// All is the type of all expressions.
type All string

func (p All) String() string {
	return "<ALL:" + string(p) + ">"
}

func (p All) AddEnv(env Env) {}

func (p All) Match(net *Net, env Env) ([]*Value, []int) {
	f := net.World[string(p)]
	m := make([]int, len(f))
	for i := range f {
		m[i] = 1
	}
	return f, m
}

func (p All) Skeletonize(net *Net) int {
	return len(net.World[string(p)])
}

// ----------------------------------------------------------------------

// Add is the type of add expressions.
type Add []Expression

func (p Add) String() string {
	return multstring(p, "", " + ", "")
}

func (p Add) AddEnv(env Env) {
	multaddEnv(p, env)
}

func (p Add) Match(net *Net, env Env) ([]*Value, []int) {
	var res []*Value
	var mult []int
	for i := range p {
		f, m := p[i].Match(net, env)
		res = append(res, f...)
		mult = append(mult, m...)
	}
	return res, mult
}

func (p Add) Skeletonize(net *Net) int {
	res := 0
	for i := range p {
		res += p[i].Skeletonize(net)
	}
	return res
}

// ----------------------------------------------------------------------

// Subtract is the type of subtract expressions. It is an array of two
// expressions denoting the left and right elements of a subtract operation.
type Subtract []Expression

func (p Subtract) String() string {
	return multstring(p, "", " - ", "")
}

func (p Subtract) AddEnv(env Env) {
	multaddEnv(p, env)
}

func (p Subtract) Match(net *Net, env Env) ([]*Value, []int) {
	fa, ma := p[0].Match(net, env)
	if fa == nil || len(p) == 1 {
		return fa, ma
	}
	for i := 1; i < len(p); i++ {
		fb, mb := p[i].Match(net, env)
		if fb == nil {
			continue
		}
		fa, ma = subtract(fa, ma, fb, mb)
	}
	return fa, ma
}

// subtract computes multiset difference, taking into account multiplicities
func subtract(fa []*Value, ma []int, fb []*Value, mb []int) ([]*Value, []int) {
	var f []*Value
	var m []int
OUTER:
	for i, a := range fa {
		for j, b := range fb {
			if a == b {
				if ma[i]-mb[j] <= 0 {
					continue OUTER
				}
				ma[i] = ma[i] - mb[j]
			}
		}
		f = append(f, a)
		m = append(m, ma[i])
	}
	return f, m
}

func (p Subtract) Skeletonize(net *Net) int {
	res := p[0].Skeletonize(net)
	for i := 1; i < len(p); i++ {
		res = res - p[i].Skeletonize(net)
	}
	return res
}

// ----------------------------------------------------------------------

// Tuple is the type of tuple expressions.
type Tuple []Expression

func (p Tuple) String() string {
	return multstring(p, "(", ", ", ")")
}

func (p Tuple) AddEnv(env Env) {
	multaddEnv(p, env)
}

func (p Tuple) Match(net *Net, env Env) ([]*Value, []int) {
	res := []*Value{nil}
	for i := len(p) - 1; i >= 0; i-- {
		f, _ := p[i].Match(net, env)
		foo := []*Value{}
		for _, v1 := range f {
			for _, v2 := range res {
				foo = append(foo, net.unique[Value{Head: v1.Head, Tail: v2}])
			}
		}
		res = foo
	}
	mres := []int{}
	for i := 0; i < len(res); i++ {
		mres = append(mres, 1)
	}
	return res, mres
}

func (p Tuple) Skeletonize(net *Net) int {
	return 1
}

// ----------------------------------------------------------------------

// Operation is the type of expressions that apply an operation to a slice of
// expressions.
type Operation struct {
	Op   OP
	Elem []Expression
}

func (p OP) String() string {
	switch p {
	case AND:
		return " and "
	case OR:
		return " or "
	case EQ:
		return " = "
	case INEQ:
		return " != "
	case LESSTHAN:
		return " < "
	case LESSTHANEQ:
		return " <= "
	case GREATTHAN:
		return " > "
	case GREATTHANEQ:
		return " >= "
	}
	return ""
}

func (p Operation) String() string {
	return multstring(p.Elem, "(", p.Op.String(), ")")
}

func (p Operation) AddEnv(env Env) {
	multaddEnv(p.Elem, env)
}

func (p Operation) Match(net *Net, env Env) ([]*Value, []int) {
	panic("Match not authorized on Operation")
}

// OK returns whether the condition evaluates to true.
func (p Operation) OK(net *Net, env Env) bool {
	switch p.Op {
	case NIL:
		return true
	case AND:
		for _, c := range p.Elem {
			if !c.(Operation).OK(net, env) {
				return false
			}
		}
		return true
	case OR:
		for _, c := range p.Elem {
			if c.(Operation).OK(net, env) {
				return true
			}
		}
		return false
	default:
		v1, _ := p.Elem[0].Match(net, env)
		v2, _ := p.Elem[1].Match(net, env)
		if len(v1) == 0 || len(v2) == 0 {
			return false
		}
		if len(v1) > 1 || len(v2) > 1 {
			panic("problem in conditional, too many results")
		}
		return net.compareExp(v1[0], v2[0], p.Op)
	}
}

func (p Operation) Skeletonize(net *Net) int {
	panic("Match not authorized on Operation")
}

// ----------------------------------------------------------------------

// Constant is the type of constant expressions.
type Constant string

func (p Constant) String() string {
	return string(p)
}

func (p Constant) AddEnv(env Env) {}

func (p Constant) Match(net *Net, env Env) ([]*Value, []int) {
	pval, found := net.order[string(p)]
	if !found {
		f, wfound := net.World[string(p)]
		if !wfound {
			log.Fatalf("identifier %s is not a constant or a known type", string(p))
		}
		m := make([]int, len(f))
		for i := range f {
			m[i] = 1
		}
		return f, m
	}
	return []*Value{pval}, []int{1}
}

func (p Constant) Skeletonize(net *Net) int {
	return 1
}

// ----------------------------------------------------------------------

// FIRConstant is the type of finite int range constant expressions.
type FIRConstant struct {
	value int
	start int
	end   int
}

func (p FIRConstant) stringify() string {
	return fmt.Sprintf("_int{%d,%d}%d", p.start, p.end, p.value)
}

func (p FIRConstant) String() string {
	return strconv.Itoa(p.value)
}

func (p FIRConstant) AddEnv(env Env) {}

func (p FIRConstant) Match(net *Net, env Env) ([]*Value, []int) {
	return []*Value{net.order[p.stringify()]}, []int{1}
}

func (p FIRConstant) Skeletonize(net *Net) int {
	return 1
}

// ----------------------------------------------------------------------

// Var is the type of variables.
type Var string

func (p Var) String() string {
	return string(p)
}

func (p Var) AddEnv(env Env) {
	env[string(p)] = nil
}

func (p Var) Match(net *Net, env Env) ([]*Value, []int) {
	return []*Value{env[string(p)]}, []int{1}
}

func (p Var) Skeletonize(net *Net) int {
	return 1
}

// ----------------------------------------------------------------------

// Dot is the type of dot constants.
type Dot struct{}

func (p Dot) String() string {
	return "o"
}

func (p Dot) AddEnv(env Env) {}

func (p Dot) Match(net *Net, env Env) ([]*Value, []int) {
	return []*Value{net.vdot}, []int{1}
}

func (p Dot) Skeletonize(net *Net) int {
	return 1
}

// ----------------------------------------------------------------------

// Successor is the type of successor and predecessor operations.
type Successor struct {
	Var
	Incr int
}

func (p Successor) String() string {
	var mod string
	if p.Incr > 0 {
		mod = "++" + strconv.Itoa(p.Incr)
	} else {
		mod = "--" + strconv.Itoa(-p.Incr)
	}
	return string(p.Var) + mod
}

func (p Successor) AddEnv(env Env) {
	env[string(p.Var)] = nil
}

func (p Successor) Match(net *Net, env Env) ([]*Value, []int) {
	c := env[string(p.Var)]
	res := net.Next(p.Incr, c)
	if res == nil {
		return nil, nil
	}
	return []*Value{res}, []int{1}
}

func (p Successor) Skeletonize(net *Net) int {
	return 1
}

// ----------------------------------------------------------------------

// Numberof is the type of numberof expressions in PNML. This is used to add a
// multiplicity Mult to an Expression in a multiset.
type Numberof struct {
	Expression
	Mult int
}

func (p Numberof) String() string {
	return strconv.Itoa(p.Mult) + "'" + p.Expression.String()
}

func (p Numberof) AddEnv(env Env) {
	p.Expression.AddEnv(env)
}

func (p Numberof) Match(net *Net, env Env) ([]*Value, []int) {
	f, m := p.Expression.Match(net, env)
	for i := range f {
		m[i] = p.Mult
	}
	return f, m
}

func (p Numberof) Skeletonize(net *Net) int {
	return p.Mult * p.Expression.Skeletonize(net)
}

// ----------------------------------------------------------------------

func multstring(ee []Expression, start, delim, end string) string {
	s := start
	if len(ee) == 0 {
		return s + end
	}
	s = s + ee[0].String()
	if len(ee) == 1 {
		return s + end
	}
	for i := 1; i < len(ee); i++ {
		s = s + delim + ee[i].String()
	}
	return s + end
}

func multaddEnv(ee []Expression, env Env) {
	for _, v := range ee {
		v.AddEnv(env)
	}
}

// ----------------------------------------------------------------------

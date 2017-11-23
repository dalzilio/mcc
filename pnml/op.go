// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package pnml

// OP is the type of operations in a PNML expression
type OP int

const (
	NIL OP = iota
	EQ
	INEQ
	LESSTHAN
	LESSTHANEQ
	GREATTHAN
	GREATTHANEQ
	OR
	AND
)

// ----------------------------------------------------------------------

//go:generate stringer -type=TYP

// TYP describes the possible kind of PNML types in PNML.
type TYP int

const (
	DOT TYP = iota
	CENUM
	FENUM
	PROD
	NUMERIC
)

// ----------------------------------------------------------------------

// VERB describes the level of verbosity in the output
type VERB int

const (
	QUIET VERB = iota
	MINIMAL
	MAXIMAL
)

// ----------------------------------------------------------------------

func getpOP(s string) OP {
	switch s {
	case "equality":
		return EQ
	case "inequality":
		return INEQ
	case "lessthan":
		return LESSTHAN
	case "lessthanorequal":
		return LESSTHANEQ
	case "greaterthan":
		return GREATTHAN
	case "greaterthanorequal":
		return GREATTHANEQ
	case "or":
		return OR
	case "and":
		return AND
	}
	panic("not an operation " + s)
}

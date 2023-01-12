// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

//go:generate stringer -type=TYP

// We may need to install go stringer first: go get golang.org/x/tools/cmd/stringer

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
	FINTRANGE
	PROD
	NUMERIC
)

// ----------------------------------------------------------------------

// VERB describes the level of verbosity in the output. Each value has an impact
// on the level of information carried by place and transition names, and the
// details added to the labels when printing the P/T net. Note that labels will
// always be the same identifiers than the ones used in the corresponding COL
// model.
//
// - INFO: place/trans are identifiers + encoding of colors ; exact COL
// identifiers are kept as labels. PT identifiers are normalized version of
// their COL equivalent (e.g. non aname characters like "-" or spaces are
// replaced by underscores "_"). Colors are encoded using color names. Used by
// command `mcc info`.
//
// - SKELETON: place/trans are exactly those of the COL nets, without
// normalization.
//
// - SLICED: place are identifiers + encoding of colors and transitions are of
// the form t%d, with %d an integer; no labels.  PT identifiers are normalized,
// liked with INFO, but colors are encoded using integers. Used by command `mcc
// pnml` and `mcc tina` by default.
//
// - SMPT: place/trans are p/t + integers ; no labels and no association between
// the COL and PT nets. Equivalent to what is used by command `mcc lola`. Used
// as default by command `mcc smpt`, with the difference that we also provide a
// separate list of association between COL and PT identifiers
//
type VERB int

const (
	INFO VERB = iota
	SKELETON
	SLICED
	SMPT
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

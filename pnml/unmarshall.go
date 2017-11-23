// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package pnml

// ----------------------------------------------------------------------

// pnml is the type of PNML net. We ignore the graphical information contained
// in the net and only consider the first page
type pnml struct {
	Net Net `xml:"net"`
}

// Net is the type of the net element in a PNML file.
type Net struct {
	Name        string      `xml:"name>text"`
	Page        Page        `xml:"page"`
	Declaration Declaration `xml:"declaration>structure>declarations"`
	// Env is an association between a variable name and its type name, found in declaration
	Env map[string]string
	// types tells in which type a given constant belongs
	types map[string]*TypeDecl
	// position tells the position of the constant in its type; used for successor
	position map[string]int
	// order associates a unique Value to every Constant; it is used for
	// encoding Constant into Value
	order map[string]*Value
	// Identity associates a string to a constant index. This is only useful for printing
	// debugging information
	Identity []string
	// unique associates a unique representant for each Value
	unique map[Value]*Value
	// World associates a type (name) with a list of all its possible values
	World map[string][]*Value
	// vdot is the Value for the dot constant
	vdot *Value
	// VERBOSE guides the level of information in the outputs
	VERBOSE VERB
	// MCC tells us whether we should allow duplicate transitions on conditions.
	MCC bool
	// FES tells us whether we should use the FiniteEnumeration semantics (the
	// successor/predecessor of a constant in a finite enumeration may be an
	// unvalid value).
	FES bool
}

// ----------------------------------------------------------------------

// SetVerbose sets the value of the VERBOSE setting. true means more information on
// the output.
func (net *Net) SetVerbose(b VERB) {
	net.VERBOSE = b
}

// ----------------------------------------------------------------------

// SetMCC sets the value of the MCC setting. true means we try to emulate the
// expected behavior of PNML unfolding.
func (net *Net) SetMCC(b bool) {
	net.MCC = b
}

// ----------------------------------------------------------------------

// SetFES sets the value of the FES setting. true means we try to emulate the
// expected behavior of successor/predecessor with an overflow/underflow semantics.
func (net *Net) SetFES(b bool) {
	net.FES = b
}

// ----------------------------------------------------------------------

// Page is the type of the page element in a PNML file.
type Page struct {
	Places []*Place      `xml:"place"`
	Trans  []*Transition `xml:"transition"`
	Arcs   []*Arc        `xml:"arc"`
}

// Declaration is the type of a PNML net declaration. It contains declarations
// for types and variables used in the net.
type Declaration struct {
	Sorts []*TypeDecl `xml:"namedsort"`
	Vars  []*VarDecl  `xml:"variabledecl"`
}

// ----------------------------------------------------------------------

// TypeDecl is the type of  PNML type declarations. Test is the field Dot is not
// nil to check whether this is a dot.
type TypeDecl struct {
	Sort    TYP
	Elem    []string
	ID      string    `xml:"id,attr"`
	CEnum   []Fec     `xml:"cyclicenumeration>feconstant,omitempty"`
	FEnum   []Fec     `xml:"finiteenumeration>feconstant,omitempty"`
	Product []Type    `xml:"productsort>usersort,omitempty"`
	Dot     *struct{} `xml:"dot,omitempty"`
}

// ----------------------------------------------------------------------

// Fec is the type of  PNML enumeration constants.
type Fec struct {
	ID string `xml:"id,attr"`
}

// VarDecl is the type of  PNML variable  declarations.
type VarDecl struct {
	ID   string `xml:"id,attr"`
	Type Type   `xml:"usersort"`
}

// Type is the type of a type declaration.
type Type struct {
	ID string `xml:"declaration,attr"`
}

// ----------------------------------------------------------------------

// Place is the type of a PNML place. It can contain a type and an (optional)
// initial marking.
type Place struct {
	ID             string `xml:"id,attr"`
	Type           Type   `xml:"type>structure>usersort"`
	XML            RawXML `xml:"hlinitialMarking>structure"`
	InitialMarking Expression
}

// RawXML is the type of PNML initial marking expressions, patterns and conditions.
type RawXML struct {
	InnerXML []byte `xml:",innerxml"`
}

// ----------------------------------------------------------------------

// Transition is the type of a PNML transition. It can contain a type and an (optional)
// initial marking.
type Transition struct {
	ID        string `xml:"id,attr"`
	XML       RawXML `xml:"condition>structure"`
	Condition Expression
}

// ----------------------------------------------------------------------

// Arc is the type of edges element in a PNML net.
type Arc struct {
	Source  string `xml:"source,attr"`
	Target  string `xml:"target,attr"`
	XML     RawXML `xml:"hlinscription>structure"`
	Pattern Expression
}

// ----------------------------------------------------------------------

// NumberConstant is used in PNML expressions.
type NumberConstant struct {
	Value int `xml:"value,attr"`
}

// Variable is used in PNML expressions.
type Variable struct {
	RefVariable string `xml:"refvariable,attr"`
}

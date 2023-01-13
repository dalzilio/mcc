// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package pnml

import (
	"encoding/xml"
	"fmt"
	"io"
)

// ----------------------------------------------------------------------

// Next gives the ith successor (i can be negative) of the constant with id
// name.
func (e *Net) Next(i int, val *Value) *Value {
	name := e.Identity[val.Head]
	typ := e.types[name]
	pos := e.position[name] + i
	if e.FES && (pos < 0 || pos >= len(typ.Elem)) {
		return nil
	}
	for {
		if pos >= 0 {
			break
		}
		pos = pos + len(typ.Elem)
	}
	return e.order[typ.Elem[pos%len(typ.Elem)]]
}

// // All returns the set of constant in the type with ID name as a slice of
// // Expressions.
// func (e *Net) All(name string) []Expression {
// 	typ := e.types[name].Elem
// 	res := make([]Expression, len(typ))
// 	for k, v := range typ {
// 		res[k] = Constant(v)
// 	}
// 	return res
// }

// ----------------------------------------------------------------------

// A Decoder represents an XML parser reading a particular input stream that
// should be a valid PNML file. It embeds an xml.Decoder for parsing a PNML
// stream passed as an io.Reader. Like with xml.Decoder, the parser assumes that
// its input is encoded in UTF-8.
type Decoder struct {
	*xml.Decoder
}

// NewDecoder creates a new XML parser reading from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{xml.NewDecoder(r)}
}

// ----------------------------------------------------------------------

// Build creates a pnml.Net object by parsing the content of the provided file.
func (d *Decoder) Build(net *Net) error {
	var pnml pnml
	decoder := d.Decoder
	err := decoder.Decode(&pnml)
	if err != nil {
		return fmt.Errorf("error decoding PNML input: %s", err)
	}

	net.Name = pnml.Net.Name
	net.Page = pnml.Net.Page
	net.Declaration = pnml.Net.Declaration

	// we allocate the enumerators and other useful maps
	net.Env = make(map[string]string)
	net.types = make(map[string]*TypeDecl)
	net.ranges = make(map[IntRange]*TypeDecl)
	net.position = make(map[string]int)
	net.order = make(map[string]*Value)
	net.Identity = make([]string, 1)
	net.unique = make(map[Value]*Value)
	net.World = make(map[string][]*Value)

	// we allocate the constant for dot
	net.vdot = &Value{Head: 0}
	net.Identity[0] = "o"
	net.unique[Value{Head: 0}] = net.vdot
	ccount := 1

	// we make a pass through the constant definitions
	for _, v := range net.Declaration.Sorts {
		switch {
		case v.Dot != nil:
			v.Sort = DOT
			net.World[v.ID] = []*Value{net.vdot}
		case v.FIntRan != nil:
			// we expand a range like a FEnum with constants of the form _intx
			v.Sort = FINTRANGE
			pv, found := net.ranges[*v.FIntRan]
			if !found {
				v.Elem = make([]string, v.FIntRan.End-v.FIntRan.Start+1)
				list := make([]*Value, len(v.Elem))
				fir := FIRConstant{start: v.FIntRan.Start, end: v.FIntRan.End}
				for i := 0; i < len(v.Elem); i++ {
					fir.value = v.FIntRan.Start + i
					c := fir.stringify()
					v.Elem[i] = c
					net.Identity = append(net.Identity, c)
					net.position[c] = i
					val := Value{Head: ccount}
					ccount++
					net.order[c] = &val
					net.unique[val] = &val
					list[i] = &val
					net.types[c] = v
				}
				net.ranges[*v.FIntRan] = v
				net.World[v.ID] = list
			} else {
				net.World[v.ID] = net.World[pv.ID]
			}
		case v.CEnum != nil:
			v.Sort = CENUM
			v.Elem = make([]string, len(v.CEnum))
			for i, c := range v.CEnum {
				v.Elem[i] = c.ID
				net.position[c.ID] = i
			}
		case v.FEnum != nil:
			v.Sort = FENUM
			v.Elem = make([]string, len(v.FEnum))
			for i, c := range v.FEnum {
				v.Elem[i] = c.ID
				net.position[c.ID] = i
			}
		case v.Product != nil:
			v.Sort = PROD
			v.Elem = make([]string, len(v.Product))
			for i, c := range v.Product {
				v.Elem[i] = c.ID
			}
		}
		if v.Sort == CENUM || v.Sort == FENUM {
			list := make([]*Value, len(v.Elem))
			for i, c := range v.Elem {
				val := Value{Head: ccount}
				ccount++
				net.types[c] = v
				net.order[c] = &val
				net.Identity = append(net.Identity, c)
				net.unique[val] = &val
				list[i] = &val
			}
			net.World[v.ID] = list
		}
	}

	for _, v := range net.Declaration.Sorts {
		if v.Sort == PROD {
			net.World[v.ID] = net.enumprod(v.Elem)
		}
	}

	for _, v := range net.Declaration.Vars {
		net.Env[v.ID] = v.Type.ID
	}

	// we associate the list of partition element to their identifiers
	for _, p := range net.Declaration.Partitions {
		for _, pe := range p.Partitions {
			val := []*Value{}
			for _, e := range pe.Elem {
				val = append(val, net.order[e.ID])
			}
			net.World[pe.ID] = val
		}
	}

	// We parse the expressions in the places, transitions and arcs of the first
	// page.
	for _, t := range net.Page.Trans {
		exp, err := parseExpression(t.XML.InnerXML)
		if err != nil {
			return err
		}
		if exp == nil {
			exp = Operation{Op: NIL}
		}
		e, ok := exp.(Operation)
		if !ok {
			return fmt.Errorf("invalid condition in a PNML transition")
		}
		t.Condition = e
	}

	for _, p := range net.Page.Places {
		exp, err := parseExpression(p.XML.InnerXML)
		if err != nil {
			return err
		}
		p.InitialMarking = exp
	}

	for _, a := range net.Page.Arcs {
		exp, err := parseExpression(a.XML.InnerXML)
		if err != nil {
			return err
		}
		a.Pattern = exp
	}

	return nil
}

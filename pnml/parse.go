// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package pnml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
)

func parseExprMult(d *xml.Decoder, acc []Expression) ([]Expression, error) {
	res, err := parseExprElement(d)
	if err != nil {
		return nil, fmt.Errorf("error while parsing a sequence of elements; %s", err)
	}
	if res == nil {
		return acc, nil
	}
	return parseExprMult(d, append(acc, res))
}

func parseExprInner(decoder *xml.Decoder) (Expression, error) {
	res, err := parseExprElement(decoder)
	skipUntilEnd(decoder)
	return res, err
}

func skipUntilEnd(d *xml.Decoder) error {
	t, _ := d.Token()
	if t == nil {
		return fmt.Errorf("I was expecting an xml.EndElement")
	}
	switch t.(type) {
	case xml.CharData:
		skipUntilEnd(d)
	case xml.EndElement:
		return nil
	}
	return fmt.Errorf("I was expecting an xml.EndElement, found (%T) %v", t, t)
}

func parseExprElement(decoder *xml.Decoder) (Expression, error) {
	t, _ := decoder.Token()
	if t == nil {
		return nil, nil
	}

	switch se := t.(type) {
	case xml.CharData:
		return parseExprElement(decoder)
	case xml.EndElement:
		return nil, nil
	case xml.StartElement:
		switch se.Name.Local {
		case "subterm":
			return parseExprInner(decoder)
		case "all":
			res, err := parseExprInner(decoder)
			if err != nil {
				return nil, err
			}
			if s, ok := res.(Constant); ok {
				return All(s), nil
			}
			return nil, errors.New("Malformed PNML, an all element should contain usersort")
		case "usersort":
			var s Type
			decoder.DecodeElement(&s, &se)
			return Constant(s.ID), nil
		case "useroperator":
			var s Type
			decoder.DecodeElement(&s, &se)
			return Constant(s.ID), nil
		case "tuple":
			ee, err := parseExprMult(decoder, nil)
			if err != nil {
				return nil, err
			}
			return Tuple(ee), nil
		case "add":
			ee, err := parseExprMult(decoder, nil)
			if err != nil {
				return nil, err
			}
			return Add(ee), nil
		case "subtract":
			// Subtract is actually like an Add since there can be several sub
			// chained together, see the example in PhiloDyn, trans. Initialize,
			// that has an arc pattern of the form All - p - q, where p and q
			// are variables.
			ee, err := parseExprMult(decoder, nil)
			if err != nil {
				return nil, err
			}
			return Subtract(ee), nil
		case "successor", "predecessor":
			inc := 1
			if se.Name.Local == "predecessor" {
				inc = -1
			}
			res, err := parseExprInner(decoder)
			if err != nil {
				return nil, err
			}
			if r, ok := res.(Successor); ok {
				return Successor{Var: r.Var, Incr: r.Incr + inc}, nil
			}
			if r, ok := res.(Var); ok {
				return Successor{Var: r, Incr: inc}, nil
			}
			return nil, errors.New("We only support successor and predecessor on variables")
		case "or", "and", "equality", "inequality",
			"lessthanorequal", "lessthan",
			"greaterthan", "greaterthanorequal":
			ee, err := parseExprMult(decoder, nil)
			if err != nil {
				return nil, err
			}
			return Operation{Op: getpOP(se.Name.Local), Elem: ee}, nil
		case "dotconstant":
			skipUntilEnd(decoder)
			return Dot{}, nil
		case "numberof":
			// we consider the case where numberof does not have two subelements
			// (the multiplicity is missing). This was found in one PNML model,
			// SimpleLoad, and lately on instance DotAndBoxes-COL-3 found in the
			// MCC website. The same problem appears in files dot2.pnml and
			// dot3.pnml of the benchmarks/simple folder.
			res1, _ := parseExprElement(decoder)
			if numb, ok := res1.(Numberof); ok {
				res2, _ := parseExprElement(decoder)
				skipUntilEnd(decoder)
				numb.Expression = res2
				return numb, nil
			}
			res2, _ := parseExprElement(decoder)
			if res2 == nil {
				// it means that res1 was a value, not a numberconstant, and
				// that the multiplicity was forgotten in the XML numberof
				// element.
				return Numberof{Expression: res1, Mult: 1}, nil
			}
			return nil, errors.New("Malformed PNML in numberof, numberconstant is missing")
		case "numberconstant":
			var val NumberConstant
			decoder.DecodeElement(&val, &se)
			return Numberof{Expression: nil, Mult: val.Value}, nil
		case "finiteintrangeconstant":
			var val FIRangeConstant
			decoder.DecodeElement(&val, &se)
			return FIRConstant{
				value: val.Value,
				start: val.Range.Start,
				end:   val.Range.End,
			}, nil
		case "variable":
			var val Variable
			decoder.DecodeElement(&val, &se)
			return Var(val.RefVariable), nil
		default:
			return nil, errors.New("malformed PNML: unexpected Token <" + se.Name.Local + "> in parseExpr")
		}
	}
	return nil, errors.New("malformed PNML: unexpected Token in parseExpr")
}

// parseExpression returns the PNML expression corresponding to the XML content
// found in the provided byte slice.
func parseExpression(b []byte) (Expression, error) {
	if len(b) == 0 {
		return nil, nil
	}
	buff := bytes.NewBuffer(b)
	decoder := xml.NewDecoder(buff)
	return parseExprElement(decoder)
}

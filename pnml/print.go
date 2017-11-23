// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package pnml

// ----------------------------------------------------------------------

func (net Net) String() string {
	var s string
	for _, v := range net.Declaration.Sorts {
		s += "# " + v.String() + "\n"
	}
	for _, v := range net.Declaration.Vars {
		s += "# " + v.ID + " : " + v.Type.ID + "\n"
	}
	return s
}

// ----------------------------------------------------------------------

func (typ TypeDecl) String() string {
	var s = typ.ID + " :"
	switch {
	case typ.CEnum != nil:
		for _, v := range typ.CEnum {
			s += " " + v.ID
		}
	case typ.FEnum != nil:
		for _, v := range typ.FEnum {
			s += " " + v.ID
		}
	case typ.Product != nil:
		s += "("
		if typ.Product == nil {
			s += ")"
			break
		}
		s += typ.Product[0].ID
		for i := 1; i < len(typ.Product); i++ {
			s += " x " + typ.Product[i].ID
		}
		s += ")"
	default:
		s += " dot"
	}
	return s
}

// ----------------------------------------------------------------------

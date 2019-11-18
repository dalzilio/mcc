// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package pnml

import "fmt"

// ----------------------------------------------------------------------

func (net Net) String() string {
	var s string
	for _, v := range net.Declaration.Sorts {
		s += "# type " + v.String() + "\n"
	}
	for _, v := range net.Declaration.Partitions {
		s += "# partition " + v.ID + " : " + v.Type.ID + "\n"
		for _, p := range v.Partitions {
			s += "# +-- " + p.ID + " : " + p.Elem[0].ID + " ... " + p.Elem[len(p.Elem)-1].ID + "\n"
		}
	}
	for _, v := range net.Declaration.Vars {
		s += "# var  " + v.ID + " : " + v.Type.ID + "\n"
	}
	return s
}

// ----------------------------------------------------------------------

func (typ TypeDecl) String() string {
	var s = typ.ID + " :"
	switch {
	case typ.Sort == FINTRANGE:
		s += fmt.Sprintf(" IntRange %d -- %d", typ.FIntRan.Start, typ.FIntRan.End)
	case (typ.Sort == CENUM) || (typ.Sort == FENUM):
		for _, v := range typ.Elem {
			s += " " + v
		}
	case typ.Sort == PROD:
		s += " ("
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

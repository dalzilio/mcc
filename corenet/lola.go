// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

import (
	"fmt"
	"io"
	"strconv"
)

// ----------------------------------------------------------------------

func multstring(ee []string, start, delim, end string) string {
	s := start
	if len(ee) == 0 {
		return s + end
	}
	s = s + ee[0]
	if len(ee) == 1 {
		return s + end
	}
	for i := 1; i < len(ee); i++ {
		s = s + delim + ee[i]
	}
	return s + end
}

// ----------------------------------------------------------------------

// Lola generates a net with the lola format
func (pl Place) Lola() string {
	return "p" + strconv.Itoa(pl.count)
}

// Lola generates a net with the lola format
func (pl corep) Lola() string {
	return fmt.Sprintf("%s: %d", pl.Place.Lola(), pl.int)
}

// Lola generates a net with the lola format
func (tr Trans) Lola() string {
	s := "TRANSITION t" + strconv.Itoa(tr.count) + "\n"
	s += fmt.Sprintf("\t{-- %s --}\n", tr.label)
	var list []string
	for _, v := range tr.in {
		list = append(list, v.Lola())
	}
	s += multstring(list, "\tCONSUME ", ", ", ";\n")
	list = nil
	for _, v := range tr.out {
		list = append(list, v.Lola())
	}
	s += multstring(list, "\tPRODUCE ", ", ", ";\n")
	return s + "\n"
}

// Lola generates a net with the lola format
func (net Net) Lola() string {
	s := "PLACE\n"
	var list, comments []string
	for _, v := range net.pl {
		comments = append(comments, fmt.Sprintf("\t{-- p%d %s --}", v.count, v.label))
		list = append(list, v.Lola())
	}
	s += multstring(comments, "", "\n", "\n")
	s += multstring(list, "\t", ", ", ";\n")

	s += "MARKING\n"
	list = nil
	for _, v := range net.pl {
		if v.init != 0 {
			list = append(list, v.Lola()+": "+strconv.Itoa(v.init))
		}
	}
	s += multstring(list, "\t", ", ", ";\n")

	for _, v := range net.tr {
		s += v.Lola()
	}
	return s
}

// LolaWrite prints a net with the lola format on an io.Writer
func (net Net) LolaWrite(w io.Writer) {
	fmt.Fprint(w, "PLACE\n")
	fmt.Fprint(w, "\t")
	flag := false
	for _, v := range net.pl {
		if flag {
			fmt.Fprint(w, ", ")
		} else {
			flag = true
		}
		fmt.Fprintf(w, " p%d", v.count)
	}
	fmt.Fprint(w, ";\n")

	fmt.Fprint(w, "MARKING\n")
	fmt.Fprint(w, "\t")
	flag = false
	for _, v := range net.pl {
		if v.init != 0 {
			if flag {
				fmt.Fprint(w, ", ")
			} else {
				flag = true
			}
			fmt.Fprintf(w, "%s: %d", v.Lola(), strconv.Itoa(v.init))
		}
	}
	fmt.Fprint(w, ";\n")

	for _, v := range net.tr {
		fmt.Fprintf(w, "TRANSITION t%d\n", strconv.Itoa(v.count))
		fmt.Fprint(w, "\tCONSUME ")
		flag = false
		for _, val := range v.in {
			if flag {
				fmt.Fprint(w, ", ")
			} else {
				flag = true
			}
			fmt.Fprint(w, val.Lola())
		}
		fmt.Fprint(w, ";\n")

		fmt.Fprint(w, "\tPRODUCE ")
		flag = false
		for _, val := range v.out {
			if flag {
				fmt.Fprint(w, ", ")
			} else {
				flag = true
			}
			fmt.Fprint(w, val.Lola())
		}
		fmt.Fprint(w, ";\n")
	}
}

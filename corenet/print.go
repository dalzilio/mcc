// Copyright 2023. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

// ----------------------------------------------------------------------

// DEPRECATED

// func (pl Place) String() string {
// 	if pl.init == 0 {
// 		return fmt.Sprintf("pl p%d\t: {%s}\n", pl.count, pl.label)
// 	}
// 	return fmt.Sprintf("pl p%d\t: {%s} (%d)\n", pl.count, pl.label, pl.init)
// }

// func (pl corep) String() string {
// 	if pl.int == 1 {
// 		return fmt.Sprintf(" p%d", pl.count)
// 	}
// 	return fmt.Sprintf(" p%d*%d", pl.count, pl.int)
// }

// func (tr Trans) String() string {
// 	s := fmt.Sprintf("tr t%d\t: {%s}", tr.count, tr.label)
// 	// s := fmt.Sprintf("tr t : {%s}", tr.label)
// 	for _, v := range tr.in {
// 		s += v.String()
// 	}
// 	s = s + " ->"
// 	for _, v := range tr.out {
// 		s += v.String()
// 	}
// 	return s + "\n"
// }

// func (net Net) String() string {
// 	s := fmt.Sprintf("#net %s has %d places and %d transitions\n", net.name, len(net.pl), len(net.tr))
// 	s += fmt.Sprintf("net {%s}\n", net.name)

// 	// when VERBOSE mode is QUIET we output the raw model
// 	if net.verbose == pnml.QUIET {
// 		for _, v := range net.pl {
// 			if v.init == 0 {
// 				s += fmt.Sprintf("pl p%d\n", v.count)
// 			} else {
// 				s += fmt.Sprintf("pl p%d (%d)\n", v.count, v.init)
// 			}
// 		}
// 		for _, v := range net.tr {
// 			s += fmt.Sprintf("tr t%d ", v.count)
// 			for _, c := range v.in {
// 				s += c.String()
// 			}
// 			s = s + " ->"
// 			for _, c := range v.out {
// 				s += c.String()
// 			}
// 			s += "\n"
// 		}
// 		return s
// 	}

// 	for _, v := range net.pl {
// 		s += v.String()
// 	}
// 	for _, v := range net.tr {
// 		s += v.String()
// 	}
// 	return s
// }

// ----------------------------------------------------------------------

// DEPRECATED

// // PrintTPN returns a net in TPN format for a ring of Petri nets. The parameter
// // listlr list the 'counter' of the transitions that should be synchroized with
// // the predecessor in the ring.
// func (net Net) PrintTPN(nbcopies int, listlr, listh []int) string {
// 	// we cannot mix quiet and ring modes
// 	net.verbose = pnml.MINIMAL
// 	s := net.String()
// 	s += "\n"

// 	var mixlr string
// 	for _, i := range listlr {
// 		mixlr += fmt.Sprintf("/{%s},{%s} ", net.tr[i].label, net.tr[i+1].label)
// 	}
// 	s += fmt.Sprintf("RING sync %d %s", nbcopies, mixlr)

// 	if len(listh) == 0 {
// 		return s + "\n"
// 	}

// 	var mixh string
// 	for _, i := range listh {
// 		mixh += fmt.Sprintf(" {%s}", net.tr[i].label)
// 	}
// 	return s + mixh + "\n"
// }

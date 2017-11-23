// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

import (
	"github.com/dalzilio/mcc/pnml"
)

// ----------------------------------------------------------------------

func mkiter(net *pnml.Net, env pnml.Env) ([]string, []*int, [][]*pnml.Value) {
	varnames := make([]string, len(env))
	iterator := make([]*int, len(env))
	enums := make([][]*pnml.Value, len(env))
	i := 0
	for k := range env {
		varnames[i] = k
		iterator[i] = new(int)
		enums[i] = net.World[net.Env[k]]
		env[k] = enums[i][0]
		i++
	}
	return varnames, iterator, enums
}

func nextiter(net *pnml.Net, env pnml.Env, varname []string, iterator []*int, enums [][]*pnml.Value) bool {
	for i := range iterator {
		if *iterator[i]+1 < len(enums[i]) {
			*iterator[i]++
			env[varname[i]] = enums[i][*iterator[i]]
			return true
		}
		*iterator[i] = 0
		env[varname[i]] = enums[i][0]
	}
	return false
}

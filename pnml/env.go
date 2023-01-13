// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package pnml

import (
	"fmt"
	"sort"
)

// ----------------------------------------------------------------------

// Env is an environment, that is an association between variable (names) and
// Expressions.
type Env map[string]*Value

func (p Env) String() string {
	s := "["
	start := true
	for k := range p {
		if start {
			s += k
			start = false
		} else {
			s += fmt.Sprintf(", %s", k)
		}
	}
	return s + "]"
}

// PrintEnv returns a readable description of a Value environment
func (net *Net) PrintEnv(env Env) string {
	s := "["
	start := true
	var keys []string
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if start {
			s += fmt.Sprintf("%s : %s", k, net.PrintValue(env[k]))
			start = false
		} else {
			s += fmt.Sprintf(", %s : %s", k, net.PrintValue(env[k]))
		}
	}
	return s + "]"
}

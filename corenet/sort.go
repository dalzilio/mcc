// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

// ----------------------------------------------------------------------

type byPlaces []corep

func (s byPlaces) Len() int {
	return len(s)
}

func (s byPlaces) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byPlaces) Less(i, j int) bool {
	return s[i].count < s[j].count
}

// ----------------------------------------------------------------------

type byPLabel []*Place

func (s byPLabel) Len() int {
	return len(s)
}

func (s byPLabel) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byPLabel) Less(i, j int) bool {
	return s[i].label < s[j].label
}

// ----------------------------------------------------------------------

type byTLabel []*Trans

func (s byTLabel) Len() int {
	return len(s)
}

func (s byTLabel) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byTLabel) Less(i, j int) bool {
	return s[i].label < s[j].label
}


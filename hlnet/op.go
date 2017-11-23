// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package hlnet

//go::generate stringer -type=ARC

// ARC is the type of arcs between places and transitions in a symmetric net.
// There are two possibilities. IN is the kind of arcs from places to
// transitions. OUT is an arc from a transition to a place.
type ARC int

// IN is the kind of input arcs while OUT is for output.
const (
	IN ARC = iota
	OUT
)

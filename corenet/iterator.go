// Copyright 2017. LAAS-CNRS, Vertics. All rights reserved.
// Use of this source code is governed by the CeCILL-B license
// that can be found in the LICENSE file.

package corenet

import (
	"github.com/dalzilio/mcc/hlnet"
	"github.com/dalzilio/mcc/pnml"
)

// envIterator allows us to iterate through all the possible valid environments
// (sequence of typed variables). We build one iterator for each transition in
// the symmetric net.
//
// We add an optimization where we cycle through the subset of variables
// occuring in (the patterns) of each arc so that we can rule out some
// configurations early. This is why we use sub-iterators for each arc. When it
// is an arc to a stable place, we also rely on the fact that we know exactly
// the set of possible values in the place.
type envIterator struct {
	net       *pnml.Net
	trans     *hlnet.Transition
	cpl       map[coreAssoc]*Place
	env       pnml.Env
	partition []*subIterator
	finished  bool
}

type subIterator struct {
	*hlnet.Arcs
	varnames stringList
	// enums gives the possible values for each variables in varnames while
	// position tells us what value we are currently using.
	enums    [][]*pnml.Value
	position []int
}

// mkiter initialize an iterator from a transition
func mkiter(net *pnml.Net, cpl map[coreAssoc]*Place, t *hlnet.Transition) *envIterator {
	iter := &envIterator{
		net:       net,
		trans:     t,
		cpl:       cpl,
		env:       t.Env,
		partition: []*subIterator{},
	}
	// we build a subiterator for each arc in the transition but keep only the
	// IN arc when we have a stable place (because we kn ow that there is the
	// equivalent OUT arc). We also collect the set of
	// variables used in the arcs.
	for _, a := range t.Arcs {
		if a.Place.Stable && a.Kind == hlnet.OUT {
			continue
		}
		subiter := &subIterator{
			Arcs:     a,
			varnames: stringListBuild(a.Pattern),
		}
		if a.Place.Stable {
			iter.partition = append([]*subIterator{subiter}, iter.partition...)
		} else {
			iter.partition = append(iter.partition, subiter)
		}
	}
	// next we go through the partition to make sure that every
	// variable occurs only once.
	acc := stringList{}
	for _, p := range iter.partition {
		p.varnames = p.varnames.remove(acc)
		if len(p.varnames) != 0 {
			acc = acc.union(p.varnames)
			p.enums = make([][]*pnml.Value, len(p.varnames))
			p.position = make([]int, len(p.varnames))
			for k, name := range p.varnames {
				p.enums[k] = net.World[net.Env[name]]
				iter.env[name] = p.enums[k][0]
			}
		}
	}
	return iter
}

// hasNext returns false when we have gone through all the possible iterations.
func (iter *envIterator) hasNext() bool {
	return !iter.finished
}

// step compute the next possible assignment on the variables in partition[k];
// it returns false if we have looked at all possible values
func (iter *envIterator) step(k int) bool {
	p := iter.partition[k]
	if len(p.varnames) == 0 {
		return false
	}
	for i, name := range p.varnames {
		if p.position[i]+1 < len(p.enums[i]) {
			p.position[i]++
			iter.env[name] = p.enums[i][p.position[i]]
			// we continue incrementing until the pattern of partition[k] is
			// satisfied
			f, _ := p.Pattern.Match(iter.net, iter.env)
			if f == nil {
				return iter.step(k)
			}
			if p.Place.Stable {
				for _, fv := range f {
					if _, ok := iter.cpl[coreAssoc{place: p.Place, val: fv}]; !ok {
						return iter.step(k)
					}
				}
			}
			return true
		}
		p.position[i] = 0
		iter.env[name] = p.enums[i][0]
	}
	return false
}

// check returns a new corenet transition if the conditions in the arcs are all
// valid with the current environment; the boolean returned value is false
// if the condition failed
func (iter *envIterator) check() (*Trans, bool) {
	ct := &Trans{}
	for k, p := range iter.partition {
		f, m := p.Pattern.Match(iter.net, iter.env)
		if f == nil {
			// the pattern for the k_th edge is not satisfied, so we can stop
			// and skip all the environment that have the same valuation on the
			// variables from partition[0] to partition[k].
			iter.skip(k)
			return nil, false
		}
		for i := range f {
			place, ok := iter.cpl[coreAssoc{place: p.Place, val: f[i]}]
			if !ok {
				// if we cannot find the place {e.Place, f[i]} it means that the
				// place is stable and that f[i] is not one of its possible
				// values. Therefore this cannot be a transition in the
				// resulting corenet.
				iter.skip(k)
				return nil, false
			}
			if !p.Place.Stable {
				if p.Kind == hlnet.IN {
					ct.in = appendCorep(ct.in, corep{Place: place, int: m[i]})
				} else {
					ct.out = appendCorep(ct.out, corep{Place: place, int: m[i]})
				}
			}
			// else {
			// 	ct.in = appendCorep(ct.in, corep{Place: place, int: m[i]})
			// 	ct.out = appendCorep(ct.out, corep{Place: place, int: m[i]})
			// }
		}
	}
	if !iter.trans.Cond.OK(iter.net, iter.env) {
		iter.next()
		return nil, false
	}
	iter.next()
	return ct, true
}

// next returns the next environment that should be tested in an iteration. We
// enumerate the possible solutions in "lexicographic order" with respect to the
// order of the partition.
func (iter *envIterator) next() {
	iter.skip(len(iter.partition) - 1)
}

func (iter *envIterator) skip(k int) {
	for i := len(iter.partition) - 1; i > k; i-- {
		iter.zero(i)
	}
	for i := k; i >= 0; i-- {
		if iter.step(i) {
			return
		}
	}
	iter.finished = true
}

func (iter *envIterator) zero(k int) {
	p := iter.partition[k]
	for i, name := range p.varnames {
		p.position[i] = 0
		iter.env[name] = p.enums[i][0]
	}
}

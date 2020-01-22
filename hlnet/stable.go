package hlnet

// IsPlaceStable returns true if place p is stable in the hlnet. We also check
// that its initial marking is not empty.
func (hl *Net) IsPlaceStable(p *Place) bool {
	if p.Init == nil {
		return false
	}
	for _, t := range hl.Trans {
		// we range through all the arcs and return true if either: (1) p does
		// not appear; or (2) there is an IN and an OUT arc to p with identicall
		// patterns. (We use equality between string representation instead of
		// real (semantic) equality between patterns as a workaround.) We use
		// the fact that a pattern cannot be empty and that there can be at most
		// one IN or OUT edge to p from any transition in hl.
		var pin, pout string
		for _, a := range t.Arcs {
			if a.Place == p {
				switch a.Kind {
				case IN:
					pin = a.Pattern.String()
				case OUT:
					pout = a.Pattern.String()
				}
			}
		}
		if pin != pout {
			return false
		}
	}
	return true
}

// // IsArcSat returns whether
// func (hl *Net) IsArcSat(e Arcs, env pnml.Env) bool {
// 	return false
// }

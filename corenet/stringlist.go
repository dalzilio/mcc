package corenet

import "github.com/dalzilio/mcc/pnml"

// stringList is a container type for lists of strings (without duplication).
type stringList []string

// stringListBuild returns a slice with all the variables names in the
// expression p.
func stringListBuild(p pnml.Expression) stringList {
	env := make(pnml.Env)
	p.AddEnv(env)
	varnames := make(stringList, len(env))
	i := 0
	for name := range env {
		varnames[i] = name
		i++
	}
	return varnames
}

// add returns a stringList obtained by adding the value v at the end of s, only if it
// is not already present. We return a new slice only if the result is different
// from s.
func (s stringList) add(v string) stringList {
	if len(s) == 0 {
		return stringList{v}
	}
	for i := range s {
		if s[i] == v {
			return s
		}
	}
	res := make(stringList, len(s))
	copy(res, s)
	res = append(res, v)
	return res
}

func (s stringList) union(ns stringList) stringList {
	res := s
	for _, v := range ns {
		res = res.add(v)
	}
	return res
}

// stringListZip collects all the strings in a collection of stringList (with the
// hypothesis that they are all disctinct).
// func stringListZip(ss []stringList) stringList {
// 	if len(ss) == 0 {
// 		return stringList{}
// 	}
// 	res := ss[0]
// 	for i := 1; i < len(ss); i++ {
// 		res = res.union(ss[i])
// 	}
// 	return res
// }

// func (s stringList) addEnv(p pnml.Expression) stringList {
// 	res := stringListBuild(p)
// 	for _, v := range s {
// 		res = res.add(v)
// 	}
// 	return res
// }

// member returns the index in s at which element v occurs, or -1 if v does not
// appear in s.
// func (s stringList) member(v string) int {
// 	for k, i := range s {
// 		if i == v {
// 			return k
// 		}
// 	}
// 	return -1
// }

// delete returns a stringList containing the values s that are different from
// v.
func (s stringList) delete(v string) stringList {
	if s == nil {
		return s
	}
	for i := range s {
		if s[i] == v {
			res := make(stringList, len(s)-1)
			copy(res, s[:i])
			copy(res[i:], s[i+1:])
			return res
		}
	}
	return s
}

// remove returns a stringList containing the values that are in s but not in ns.
func (s stringList) remove(ns stringList) stringList {
	rs := s
	for _, name := range ns {
		rs = rs.delete(name)
	}
	return rs
}

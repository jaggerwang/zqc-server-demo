package util

import (
	"sort"
)

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

func Any(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if f(v) {
			return true
		}
	}
	return false
}

func All(vs []string, f func(string) bool) bool {
	for _, v := range vs {
		if !f(v) {
			return false
		}
	}
	return true
}

func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func Keys(m map[string]string) []string {
	ks := make([]string, 0, len(m))
	for k, _ := range m {
		ks = append(ks, k)
	}
	return ks
}

func Values(m map[string]string) []string {
	vs := make([]string, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

func Items(m map[string]string) [][2]string {
	is := make([][2]string, 0, len(m))
	for k, v := range m {
		is = append(is, [...]string{k, v})
	}
	return is
}

func Sort(vs []string) []string {
	vsm := make([]string, 0, len(vs))
	for _, v := range vs {
		vsm = append(vsm, v)
	}
	sort.Strings(vsm)
	return vsm
}

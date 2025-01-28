package slicex

func ConvertSlice[S ~[]E, E, V comparable](s S, fn func(E) V) []V {
	ts := make([]V, len(s))
	for i, sv := range s {
		ts[i] = fn(sv)
	}

	return ts
}

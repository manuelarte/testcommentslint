package slicesutils

// Map transform array in to array out using f.
func Map[IN, OUT any](in []IN, f func(IN) OUT) []OUT {
	out := make([]OUT, len(in))
	for i, v := range in {
		out[i] = f(v)
	}

	return out
}

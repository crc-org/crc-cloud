package maps

func Convert[X comparable, Y any, Z comparable, V any](source map[X]Y,
	convertX func(x X) Z, convertY func(y Y) V) map[Z]V {
	var result = make(map[Z]V)
	for k, v := range source {
		result[convertX(k)] = convertY(v)
	}
	return result
}

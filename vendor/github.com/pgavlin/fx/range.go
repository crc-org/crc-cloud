package fx

type range_ struct {
	i, max int
	any    bool
}

func (r *range_) Value() int {
	return r.i
}

func (r *range_) Next() bool {
	if !r.any {
		if r.i >= r.max {
			return false
		}
		r.any = true
	} else {
		if r.i >= r.max-1 {
			return false
		}
		r.i++
	}
	return true
}

func Range(min, max int) Iterator[int] {
	return &range_{i: min, max: max}
}

type minRange struct {
	i   int
	any bool
}

func (r *minRange) Value() int {
	return r.i
}

func (r *minRange) Next() bool {
	if !r.any {
		r.any = true
	} else {
		r.i++
	}
	return true
}

func MinRange(min int) Iterator[int] {
	return &minRange{i: min}
}

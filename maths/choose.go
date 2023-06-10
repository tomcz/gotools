package maths

// Choose will return yes when ok is true,
// and no when ok is false.
func Choose[V any](ok bool, yes, no V) V {
	if ok {
		return yes
	}
	return no
}

// ChooseFunc will evaluate and return yes when ok is true,
// and no when ok is false.
func ChooseFunc[V any](ok bool, yes, no func() V) V {
	if ok {
		return yes()
	}
	return no()
}

// ChooseFuncE will evaluate and return yes when ok is true,
// and no when ok is false.
func ChooseFuncE[V any](ok bool, yes, no func() (V, error)) (V, error) {
	if ok {
		return yes()
	}
	return no()
}

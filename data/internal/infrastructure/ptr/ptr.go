package ptr

func To[T any](v T) *T {
	return &v
}

func ToOrNil[T comparable](v T) *T {
	if z, ok := any(v).(interface{ IsZero() bool }); ok {
		if z.IsZero() {
			return nil
		}
		return &v
	}
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

func Get[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

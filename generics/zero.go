package generics

func Zero[T any]() T {
	var zero T
	return zero
}

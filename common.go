package utils

// IfElse == condition? trueVal: falseVal
func IfElse[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func Ptr[T any](e T) *T {
	return &e
}

func Value[T any](a *T) T {
	if a == nil {
		var zero T
		return zero
	}
	return *a
}

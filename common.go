package utils

// IfElse == condition? trueVal: falseVal
func IfElse[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

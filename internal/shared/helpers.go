package shared

import (
	"reflect"
)

func Keys[K comparable, V any, M ~map[K]V](m M) []K {
	keys := make([]K, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}

func FilterMap[K comparable, V any, M ~map[K]V](m M, filter func(key K, val V) bool) M {
	result := make(M, len(m))

	for key, val := range m {
		if filter(key, val) {
			result[key] = val
		}
	}

	return result
}

func Filter[V any, S ~[]V](slice S, filter func(V) bool) S {
	result := make(S, 0, len(slice))

	for _, value := range slice {
		if filter(value) {
			result = append(result, value)
		}
	}

	return result
}

func Map[T, U any, S ~[]T](s S, fn func(T) U) []U {
	result := make([]U, 0, len(s))

	for _, v := range s {
		result = append(result, fn(v))
	}

	return result
}

func Chain[T, U any, S ~[]T](s S, fn func(T) []U) []U {
	result := make([]U, 0, len(s))

	for _, v := range s {
		result = append(result, fn(v)...)
	}

	return result
}

func MapCheck[T, U any, S ~[]T](s S, fn func(T) (U, bool)) []U {
	result := make([]U, 0, len(s))

	for _, v := range s {
		value, ok := fn(v)
		if !ok {
			continue
		}

		result = append(result, value)
	}

	return result
}

func Assign[T any](lhs, rhs *T) T {
	lhsR := reflect.Indirect(reflect.ValueOf(lhs))
	rhsR := reflect.Indirect(reflect.ValueOf(rhs))

	if lhsR.Type() != rhsR.Type() {
		return *lhs
	}

	result := reflect.Indirect(reflect.New(lhsR.Type()))

	fields := reflect.VisibleFields(lhsR.Type())

	for _, field := range fields {
		resultField := result.FieldByName(field.Name)

		lhsField := lhsR.FieldByName(field.Name)
		if !lhsField.IsZero() {
			resultField.Set(lhsField)
			continue
		}

		rhsField := rhsR.FieldByName(field.Name)
		if !rhsField.IsZero() {
			resultField.Set(rhsField)
		}
	}

	return result.Interface().(T)
}

func In[T comparable](entry T, items ...T) bool {
	for _, item := range items {
		if entry == item {
			return true
		}
	}

	return false
}

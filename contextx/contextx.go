package contextx

import (
	"context"
)

// key is a unique type that we can use as a key in a context
type key[T any] struct{}

// With returns a copy of parent that contains the given value which can be retrieved by calling From with the resulting context
// The function uses a generic key type to ensure that the stored value is type-safe and can be uniquely identified and retrieved without
// risk of key collisions
func With[T any](ctx context.Context, v T) context.Context {
	return context.WithValue(ctx, key[T]{}, v)
}

// From returns the value associated with the wanted type from the context
// It performs a type assertion to convert the value to the desired type T
// If the type assertion is successful, it returns the value and true
// If the type assertion fails, it returns the zero value of type T and false
func From[T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(key[T]{}).(T)

	return v, ok
}

// MustFrom is similar to from, except that it panics if the type assertion fails / the value is not in the context
func MustFrom[T any](ctx context.Context) T {
	return ctx.Value(key[T]{}).(T)
}

// FromOr returns the value associated with the wanted type or the given default value if the type is not found
// This function is useful when you want to ensure that a value is always returned from the context, even if the
// context does not contain a value of the desired type. By providing a default value, you can avoid handling
// the case where the value is missing and ensure that your code has a fallback value to use
func FromOr[T any](ctx context.Context, def T) T {
	v, ok := From[T](ctx)
	if !ok {
		return def
	}

	return v
}

// FromOrFunc returns the value associated with the wanted type or the result of the given function if the type is not found
// This function is useful when the default value is expensive to compute or when the default value depends on some runtime conditions
func FromOrFunc[T any](ctx context.Context, f func() T) T {
	v, ok := From[T](ctx)
	if !ok {
		return f()
	}

	return v
}

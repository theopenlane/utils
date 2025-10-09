package contextx

import (
	"context"
)

// key is a unique type that we can use as a key in a context
type key[T any] struct{}

// stringKey is a type-constrained key allowing multiple distinct string-based values in context without collisions
// each distinct type T creates a unique context key
//
// Before:
//
//	type orgKey struct{}
//	ctx = context.WithValue(ctx, orgKey{}, "org-123")
//	orgID := ctx.Value(orgKey{}).(string)  // "org-123"
//
// After:
//
//	type OrganizationID string
//	ctx = contextx.WithString(ctx, OrganizationID("org-123"))
//	orgID := contextx.MustStringFrom[OrganizationID](ctx) // "org-123"

// Before: Need separate key structs + string values, unsafe type assertions
// After: Just define typed strings, type-safe retrieval
type stringKey[T ~string] struct{}

// With returns a copy of parent that contains the given value which can be retrieved by calling From with the resulting context
// The function uses a generic key type to ensure that the stored value is type-safe and can be uniquely identified and retrieved without
// risk of key collisions
func With[T any](ctx context.Context, v T) context.Context {
	return context.WithValue(ctx, key[T]{}, v)
}

// WithString stores a string-based value in the context
// The type parameter T must have an underlying type of string
// Each distinct type T creates a unique context key
func WithString[T ~string](ctx context.Context, value T) context.Context {
	return context.WithValue(ctx, stringKey[T]{}, value)
}

// From returns the value associated with the wanted type from the context
// It performs a type assertion to convert the value to the desired type T
// If the type assertion is successful, it returns the value and true
// If the type assertion fails, it returns the zero value of type T and false
func From[T any](ctx context.Context) (T, bool) {
	v, ok := ctx.Value(key[T]{}).(T)

	return v, ok
}

// StringFrom retrieves a string-based value from the context
// Returns the value and true if found, zero value and false otherwise
func StringFrom[T ~string](ctx context.Context) (T, bool) {
	val, ok := ctx.Value(stringKey[T]{}).(T)

	return val, ok
}

// MustFrom is similar to from, except that it panics if the type assertion fails / the value is not in the context
func MustFrom[T any](ctx context.Context) T {
	return ctx.Value(key[T]{}).(T)
}

// MustStringFrom retrieves a string-based value from the context or panics if not found
// Use this when the value must exist for the operation to proceed
func MustStringFrom[T ~string](ctx context.Context) T {
	return ctx.Value(stringKey[T]{}).(T)
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

// StringFromOr retrieves a string-based value from the context or returns the default value if not found
func StringFromOr[T ~string](ctx context.Context, def T) T {
	val, ok := StringFrom[T](ctx)
	if !ok {
		return def
	}

	return val
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

// StringFromOrFunc retrieves a string-based value from the context or returns the result of calling fn if not found
// The function is only called if the value is not present
func StringFromOrFunc[T ~string](ctx context.Context, fn func() T) T {
	val, ok := StringFrom[T](ctx)
	if !ok {
		return fn()
	}

	return val
}

// Has checks if a value of type T exists in the context
// Useful for boolean flags using empty struct types as markers e.g. "OrganizationCreationContextKey"
func Has[T any](ctx context.Context) bool {
	_, ok := From[T](ctx)

	return ok
}

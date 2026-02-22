package contextx

import "context"

// keyID is the internal identity of a Key. Each new(keyID) call produces
// a distinct pointer so two keys over the same type T are always independent
// (the _ field ensures each allocation is unique)
type keyID struct{ _ bool }

// Key stores and retrieves a value of type T in a context; each Key variable
// is independent (so two Key[string] variables never share or overwrite each other's values)
type Key[T any] struct {
	id *keyID
}

// NewKey returns a new Key for values of type T
func NewKey[T any]() Key[T] {
	return Key[T]{id: new(keyID)}
}

// Set returns a copy of ctx with v stored in this key
func (k Key[T]) Set(ctx context.Context, v T) context.Context {
	return context.WithValue(ctx, k.id, v)
}

// Get returns the value stored in this key and true, or the zero value and false if not set
func (k Key[T]) Get(ctx context.Context) (T, bool) {
	v, ok := ctx.Value(k.id).(T)

	return v, ok
}

// MustGet returns the value stored in this key, panicking if not set
func (k Key[T]) MustGet(ctx context.Context) T {
	v, ok := k.Get(ctx)
	if !ok {
		panic("contextx: value not present in context")
	}

	return v
}

// GetOr returns the value stored in this key, or def if not set
func (k Key[T]) GetOr(ctx context.Context, def T) T {
	if v, ok := k.Get(ctx); ok {
		return v
	}

	return def
}

// GetOrFunc returns the value stored in this key, or the result of fn if not set
func (k Key[T]) GetOrFunc(ctx context.Context, fn func() T) T {
	if v, ok := k.Get(ctx); ok {
		return v
	}

	return fn()
}

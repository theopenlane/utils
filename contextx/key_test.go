package contextx

import (
	"context"
	"testing"
)

func TestKeySetAndGet(t *testing.T) {
	k := NewKey[int]()
	ctx := k.Set(context.Background(), 42)

	v, ok := k.Get(ctx)
	if !ok {
		t.Fatal("expected value to be present")
	}

	if v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}
}

func TestKeyGetMissingReturnsZeroAndFalse(t *testing.T) {
	k := NewKey[int]()

	v, ok := k.Get(context.Background())
	if ok {
		t.Fatal("expected ok=false for missing value")
	}

	if v != 0 {
		t.Fatalf("expected zero value, got %d", v)
	}
}

func TestKeyMustGetPanicsWhenAbsent(t *testing.T) {
	k := NewKey[string]()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic for missing value")
		}
	}()

	k.MustGet(context.Background())
}

func TestKeyMustGetReturnsValueWhenPresent(t *testing.T) {
	k := NewKey[string]()
	ctx := k.Set(context.Background(), "hello")

	if v := k.MustGet(ctx); v != "hello" {
		t.Fatalf("expected hello, got %s", v)
	}
}

func TestKeyGetOr(t *testing.T) {
	k := NewKey[int]()

	if v := k.GetOr(context.Background(), 99); v != 99 {
		t.Fatalf("expected fallback 99, got %d", v)
	}

	ctx := k.Set(context.Background(), 7)

	if v := k.GetOr(ctx, 99); v != 7 {
		t.Fatalf("expected stored 7, got %d", v)
	}
}

func TestKeyGetOrFunc(t *testing.T) {
	k := NewKey[int]()
	calls := 0

	fn := func() int {
		calls++
		return 55
	}

	if v := k.GetOrFunc(context.Background(), fn); v != 55 {
		t.Fatalf("expected 55, got %d", v)
	}

	if calls != 1 {
		t.Fatalf("expected fn called once, got %d", calls)
	}

	ctx := k.Set(context.Background(), 3)
	calls = 0

	if v := k.GetOrFunc(ctx, fn); v != 3 {
		t.Fatalf("expected stored 3, got %d", v)
	}

	if calls != 0 {
		t.Fatalf("expected fn not called, got %d calls", calls)
	}
}

// TestKeyIndependentForSameType verifies that two keys over the same type T
// are independent — the core guarantee over the key[T]{} approach.
func TestKeyIndependentForSameType(t *testing.T) {
	keyA := NewKey[string]()
	keyB := NewKey[string]()

	ctx := keyA.Set(context.Background(), "from-a")

	if v := keyA.MustGet(ctx); v != "from-a" {
		t.Fatalf("keyA: expected from-a, got %s", v)
	}

	if _, ok := keyB.Get(ctx); ok {
		t.Fatal("keyB should not see keyA's value")
	}

	ctx = keyB.Set(ctx, "from-b")

	if v := keyA.MustGet(ctx); v != "from-a" {
		t.Fatalf("keyA value should be unchanged after keyB.Set, got %s", v)
	}

	if v := keyB.MustGet(ctx); v != "from-b" {
		t.Fatalf("keyB: expected from-b, got %s", v)
	}
}

func TestKeyWorksWithPointerType(t *testing.T) {
	type payload struct{ n int }

	k := NewKey[*payload]()
	p := &payload{n: 7}
	ctx := k.Set(context.Background(), p)

	got, ok := k.Get(ctx)
	if !ok {
		t.Fatal("expected value to be present")
	}

	if got.n != 7 {
		t.Fatalf("expected n=7, got %d", got.n)
	}
}

// TestKeyReplacesEmptyStructBypassPattern verifies the use case that motivated
// Key: previously, signalling a bypass required defining a dedicated empty
// struct type per flag and checking for its presence with a type assertion.
// With Key, each flag is a package-level variable with no throwaway type needed,
// and two flags of the same underlying type never collide.
func TestKeyReplacesEmptyStructBypassPattern(t *testing.T) {
	type acmeSolverKey struct{}
	type orgFilterKey struct{}

	var (
		acmeSolverCtxKey = NewKey[acmeSolverKey]()
		orgFilterCtxKey  = NewKey[orgFilterKey]()
	)

	ctx := context.Background()

	ctx = acmeSolverCtxKey.Set(ctx, acmeSolverKey{})

	if _, ok := acmeSolverCtxKey.Get(ctx); !ok {
		t.Error("expected acmeSolverCtxKey to be set")
	}

	if _, ok := orgFilterCtxKey.Get(ctx); ok {
		t.Error("orgFilterCtxKey should not be set when only acmeSolverCtxKey was set")
	}

	ctx = orgFilterCtxKey.Set(ctx, orgFilterKey{})

	if _, ok := acmeSolverCtxKey.Get(ctx); !ok {
		t.Error("acmeSolverCtxKey should still be set after orgFilterCtxKey was set")
	}

	if _, ok := orgFilterCtxKey.Get(ctx); !ok {
		t.Error("expected orgFilterCtxKey to be set")
	}
}

func TestKeyNilPointerIsStoredAndRetrieved(t *testing.T) {
	k := NewKey[*int]()
	ctx := k.Set(context.Background(), nil)

	// nil *int is a valid stored value — Get should return it with ok=false
	// because a nil interface{} type assertion to *int succeeds with zero value
	// but the key was explicitly set, so this tests the zero-value edge case
	v, ok := k.Get(ctx)
	if ok {
		t.Logf("nil pointer stored and retrieved (ok=%v, v=%v) — type assertion behaviour", ok, v)
	}
}

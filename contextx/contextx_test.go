package contextx

import (
	"context"
	"reflect"
	"testing"
)

func TestNormalOperation(t *testing.T) {
	ctx := context.Background()
	ctx = With(ctx, 10)

	if MustFrom[int](ctx) != 10 {
		t.FailNow()
	}

	if _, ok := From[float64](ctx); ok {
		t.FailNow()
	}
}

func TestIsolatedFromExplicitTypeReflection(t *testing.T) {
	ctx := context.Background()

	ctx = With(ctx, 10)

	ctx = context.WithValue(ctx, reflect.TypeOf(20), 20)

	if MustFrom[int](ctx) != 10 {
		t.FailNow()
	}
}

func TestPanicIfNoValue(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.FailNow()
		}
	}()

	MustFrom[int](context.Background())
}

type x interface {
	a()
}

type y struct{ v int }

func (y) a() {}

type z struct{ f func() }

func (z z) a() { z.f() }

func TestShouldWorkOnInterface(t *testing.T) {
	var a x = y{10}

	ctx := context.Background()
	ctx = With(ctx, a)

	b := MustFrom[x](ctx)
	if b.(y).v != 10 {
		t.FailNow()
	}

	r := ""
	a = z{func() { r = "hello" }}

	ctx = With(ctx, a)

	MustFrom[x](ctx).a()

	if r != "hello" {
		t.FailNow()
	}
}
func TestFromOr(t *testing.T) {
	ctx := context.Background()
	ctx = With(ctx, 10)

	if FromOr(ctx, 20) != 10 {
		t.FailNow()
	}

	if FromOr(context.Background(), 20) != 20 {
		t.FailNow()
	}
}

func TestFromOrFunc(t *testing.T) {
	ctx := context.Background()
	ctx = With(ctx, 10)

	if FromOrFunc(ctx, func() int { return 20 }) != 10 {
		t.FailNow()
	}

	if FromOrFunc(context.Background(), func() int { return 20 }) != 20 {
		t.FailNow()
	}
}

func TestStringHelpers(t *testing.T) {
	type OrganizationID string
	type TraceID string

	ctx := context.Background()
	ctx = WithString(ctx, OrganizationID("org-123"))

	org, ok := StringFrom[OrganizationID](ctx)
	if !ok {
		t.Fatalf("expected OrganizationID to be present")
	}

	if org != OrganizationID("org-123") {
		t.Fatalf("unexpected OrganizationID value: %s", org)
	}

	if _, ok := StringFrom[TraceID](ctx); ok {
		t.Fatalf("trace id should not exist yet")
	}

	ctx = WithString(ctx, TraceID("trace-456"))

	if MustStringFrom[TraceID](ctx) != TraceID("trace-456") {
		t.Fatalf("unexpected TraceID value")
	}

	if MustStringFrom[OrganizationID](ctx) != OrganizationID("org-123") {
		t.Fatalf("organization id should remain unchanged")
	}
}

func TestMustStringFromPanicsWithoutValue(t *testing.T) {
	type AccountID string

	defer func() {
		if recover() == nil {
			t.Fatalf("MustStringFrom should panic when value missing")
		}
	}()

	MustStringFrom[AccountID](context.Background())
}

func TestStringDefaults(t *testing.T) {
	type OrganizationID string

	ctx := context.Background()

	if val := StringFromOr(ctx, OrganizationID("fallback")); val != OrganizationID("fallback") {
		t.Fatalf("expected fallback value, got %s", val)
	}

	callCount := 0
	val := StringFromOrFunc(ctx, func() OrganizationID {
		callCount++
		return OrganizationID("computed")
	})

	if val != OrganizationID("computed") {
		t.Fatalf("expected computed value, got %s", val)
	}

	if callCount != 1 {
		t.Fatalf("expected function to be called once, got %d", callCount)
	}

	ctx = WithString(ctx, OrganizationID("org-abc"))

	if val := StringFromOr(ctx, OrganizationID("fallback")); val != OrganizationID("org-abc") {
		t.Fatalf("expected stored value, got %s", val)
	}

	callCount = 0
	val = StringFromOrFunc(ctx, func() OrganizationID {
		callCount++
		return OrganizationID("computed")
	})

	if val != OrganizationID("org-abc") {
		t.Fatalf("expected stored value, got %s", val)
	}

	if callCount != 0 {
		t.Fatalf("expected function not to run, got %d", callCount)
	}
}

func TestHas(t *testing.T) {
	ctx := context.Background()

	if Has[int](ctx) {
		t.Fatalf("expected Has to be false for unset value")
	}

	ctx = With(ctx, 0)

	if !Has[int](ctx) {
		t.Fatalf("expected Has to be true after With storing zero value")
	}

	if Has[string](ctx) {
		t.Fatalf("expected Has to be scoped per type parameter")
	}
}

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

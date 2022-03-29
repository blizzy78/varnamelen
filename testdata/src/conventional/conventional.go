package conventional

import (
	"context"
	"testing"
)

func Variable() {
	var (
		ctx context.Context
		b   *testing.B
		f   *testing.F
		m   *testing.M
		pb  *testing.PB
		t   *testing.T
		tb  testing.TB
	)
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = ctx
	_ = b
	_ = f
	_ = m
	_ = pb
	_ = t
	_ = tb
}

func Param(ctx context.Context,
	b *testing.B,
	f *testing.F,
	m *testing.M,
	pb *testing.PB,
	t *testing.T,
	tb testing.TB) {
	// fill
	// fill
	// fill
	// fill
	// fill
	_ = ctx
	_ = b
	_ = f
	_ = m
	_ = pb
	_ = t
	_ = tb
}

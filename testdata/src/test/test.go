package main

import (
	"context"
	"testing"
)

type test struct {
	x int
}

func foo(p int, p2 int, ip int, longParam int) (rv string) { // want "parameter name 'p' is too short for the scope of its usage"
	p2 += longParam

	x := 123 // want "variable name 'x' is too short for the scope of its usage"
	aLongOne := 555
	i := 10
	x++
	x++
	x++
	x++
	x++
	x++
	aLongOne += x + i

	ints := []int{1, 2, 3}
	for _, y := range ints { // want "variable name 'y' is too short for the scope of its usage"
		y++
		y++
		y++
		y++
		y++
		y++
		y := 123 // want "variable name 'y' is too short for the scope of its usage"
		y++
		y++
		y++
		y++
		y++
		y++
	}

	for idx := range ints {
		p := idx // want "variable name 'p' is too short for the scope of its usage"
		p++
		p++
		p++
		p++
		p++
		p++
	}

	p += ip

	rv = "foo"
	return
}

func (t *test) foo() {
	t.x++
	t.x++
	t.x++
	t.x++
	t.x++
	t.x++
}

func Conventionals(ctx context.Context, t *testing.T, b *testing.B, tb testing.TB, m *testing.M, pb *testing.PB) {
	println()
	println()
	println()
	println()
	println()
	println()
	t.Fail()
	b.Fail()
	tb.Fail()
	m.Run()
	pb.Next()
	ctx.Err()
}

package main

import (
	bb "bytes"
	"context"
	"math"
	"strings"
)

const C = 299792458

func foo(c context.Context) {
	println()
	println()
	println()
	println()
	println()
	println()
	println(c)

	println(C)
}

func foo2() (c context.Context) {
	println()
	println()
	println()
	println()
	println()
	println()
	c = context.Background()
	return
}

func foo3() {
	var c context.Context
	println()
	println()
	println()
	println()
	println()
	println()
	c = context.Background()
	println(c)
	return
}

func foo4(i int) {
	println()
	println()
	println()
	println()
	println()
	println()
	println(i)
}

func foo5() (i int) {
	println()
	println()
	println()
	println()
	println()
	println()
	i = 123
	return
}

func foo6() {
	var i int
	println()
	println()
	println()
	println()
	println()
	println()
	i = 123
	println(i)
	return
}

func foo7() {
	var ip *int
	println()
	println()
	println()
	println()
	println()
	println()
	*ip = 123
	println(*ip)
	return
}

func foo8() {
	var b *strings.Builder = &strings.Builder{}
	println()
	println()
	println()
	println()
	println()
	println()
	println(b.String())
	return
}

func foo9() {
	c := context.Background()
	println()
	println()
	println()
	println()
	println()
	println()
	_ = c.Err()
}

func foo10() {
	b := bb.Buffer{}
	println()
	println()
	println()
	println()
	println()
	println()
	println(b.String())

	d := &bb.Buffer{}
	println()
	println()
	println()
	println()
	println()
	println()
	println(d.String())
}

func foo11() {
	f := foo11
	println()
	println()
	println()
	println()
	println()
	println()
	println(f)
}

func foo12() {
	m := map[int]*bb.Buffer{}
	println()
	println()
	println()
	println()
	println()
	println()
	println(m)
}

func foo13() {
	mi := math.MinInt8
	println()
	println()
	println()
	println()
	println()
	println()
	println(mi)
}

func foo14() {
	for i := range []string{"a", "b", "c"} {
		println()
		println()
		println()
		println()
		println()
		println()
		println(i)
	}

	for i, s := range []string{"a", "b", "c"} {
		println()
		println()
		println()
		println()
		println()
		println()
		println(i, s)
	}

	i, s := getTwo()
	println()
	println()
	println()
	println()
	println()
	println()
	println(i, s)
}

func getTwo() (int, string) {
	return 1, "2"
}

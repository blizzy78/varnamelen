package main

import (
	"context"
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

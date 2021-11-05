package main

import "context"

func foo(c context.Context) {
	println()
	println()
	println()
	println()
	println()
	println()
	println(c)
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

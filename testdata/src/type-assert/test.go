package main

func foo() {
	var x interface{}

	str, ok := x.(string)
	println(str)
	println()
	println()
	println()
	println()
	println()
	println(ok)

	var i int
	i, ok = x.(int)
	println(i)
	println()
	println()
	println()
	println()
	println()
	println(ok)
}

func foo2() {
	var x interface{}

	if _, ok := x.(string); ok {
		println()
		println()
		println()
		println()
		println()
		println()
		println(ok)
	}
}

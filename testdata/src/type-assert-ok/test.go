package main

func foo() {
	var something interface{}

	str, ok := something.(string)
	println(str)
	println()
	println()
	println()
	println()
	println()
	println(ok)

	var i int
	i, ok = something.(int)
	println(i)
	println()
	println()
	println()
	println()
	println()
	println(ok)
}

func foo2() {
	var something interface{}

	if _, ok := something.(string); ok {
		println()
		println()
		println()
		println()
		println()
		println()
		println(ok)
	}
}

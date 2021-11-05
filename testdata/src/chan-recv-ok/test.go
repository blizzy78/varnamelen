package main

func foo() {
	x := make(chan int)

	_, ok := <-x
	println()
	println()
	println()
	println()
	println()
	println()
	println(ok)
}

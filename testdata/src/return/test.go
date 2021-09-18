package main

func foo() (r int, r2 int, longRV int, i int) { // want "return value name 'r' is too short for the scope of its usage"
	r2++
	longRV++
	r++
	r++
	r++
	r++
	r++
	r++
	i++
	return
}

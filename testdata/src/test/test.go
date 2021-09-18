package main

func test() {
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
}

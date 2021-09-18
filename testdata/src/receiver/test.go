package main

type test struct {
	x int
}

func (t *test) foo() { // want "parameter name 't' is too short for the scope of its usage"
	t.x++
	t.x++
	t.x++
	t.x++
	t.x++
	t.x++
}

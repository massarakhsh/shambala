package neu

type Test struct {
	name       string
	definition string
	inputs     [][]float64
	outputs    [][]float64
}

func BuildTest(name string, definition string, inputs [][]float64, outputs [][]float64) *Test {
	test := &Test{name: name, definition: definition}
	test.inputs = inputs
	test.outputs = outputs
	return test
}

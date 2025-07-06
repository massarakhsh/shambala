package shambala

func BuildBrainFull(levels ...int) *Brain {
	return buildBrainFull(levels...)
}

func LoadBrainFile(nameFile string) *Brain {
	return loadBrainFile(nameFile)
}

func BuildMentorTest(brain *Brain, test *Test) *Mentor {
	return buildMentor(brain, test)
}

func BuildTest(name string, definition string, inputs [][]float64, outputs [][]float64) *Test {
	return buildTest(name, definition, inputs, outputs)
}

package shambala

import "github.com/massarakhsh/shambala/neu"

func BuildBrainFull(levels ...int) *neu.Brain {
	return neu.BuildBrainFull(levels...)
}

func LoadBrainFile(nameFile string) *neu.Brain {
	return neu.LoadBrainFile(nameFile)
}

func BuildMentorTest(brain *neu.Brain, test *neu.Test) *neu.Mentor {
	return neu.BuildMentor(brain, test)
}

func BuildTest(name string, definition string, inputs [][]float64, outputs [][]float64) *neu.Test {
	return neu.BuildTest(name, definition, inputs, outputs)
}

package shambala

type Neuron struct {
	index   int
	inputs  []*Link
	outputs []*Link
}

func buildNeuron(index int) *Neuron {
	neuron := &Neuron{index: index}
	return neuron
}

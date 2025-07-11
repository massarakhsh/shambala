package shambala

type Network struct {
	countInputs  int
	countOutputs int
	volumeLevels []int
	neurons      []*Neuron
	links        []*Link
}

type Neuron struct {
	index   int
	inputs  []*Link
	outputs []*Link
}

type Link struct {
	index  int
	source *Neuron
	target *Neuron
}

func BuildNetwork() *Network {
	network := &Network{}
	return network
}

func BuildNetworkFull(levels ...int) *Network {
	network := BuildNetwork()
	network.fillNetworkFull(levels...)
	return network
}

func buildNeuron(index int) *Neuron {
	neuron := &Neuron{index: index}
	return neuron
}

func buildLink(index int, source *Neuron, target *Neuron) *Link {
	link := &Link{index: index, source: source, target: target}
	source.outputs = append(source.outputs, link)
	target.inputs = append(target.inputs, link)
	return link
}

func (network *Network) fillNetworkFull(levels ...int) *Network {
	network.volumeLevels = levels
	cLevels := len(levels)
	for level, volume := range levels {
		if level == 0 {
			network.countInputs = volume
		} else if level == cLevels-1 {
			network.countOutputs = volume
		}
		fLevel := len(network.neurons)
		for nu := 0; nu < volume; nu++ {
			neuron := buildNeuron(len(network.neurons))
			network.neurons = append(network.neurons, neuron)
			if level > 0 {
				cInputs := levels[level-1]
				for ni := range cInputs {
					into := network.neurons[fLevel-cInputs+ni]
					link := buildLink(len(network.links), into, neuron)
					network.links = append(network.links, link)
				}
			}
		}
	}
	return network
}

func (network *Network) GetCountNeurons() int {
	return len(network.neurons)
}

func (network *Network) GetCountLinks() int {
	return len(network.links)
}

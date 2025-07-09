package shambala

import (
	"os"

	"github.com/massarakhsh/lik"
)

type Brain struct {
	network *Network
	// zero   *Neurons
	// min    *Neurons
	// max    *Neurons
	// weight *Links
	weights []float64
	values  []float64
}

func buildBrain(network *Network) *Brain {
	brain := &Brain{network: network}
	brain.initializeWeights()
	brain.values = make([]float64, len(brain.network.neurons))
	return brain
}

func buildBrainFull(levels ...int) *Brain {
	network := BuildNetworkFull(levels...)
	return buildBrain(network)
}

func loadBrainFile(fileName string) *Brain {
	if set := lik.SetFromFile(fileName); set != nil {
		return LoadBrainSet(set)
	} else {
		return nil
	}
}

func LoadBrainSet(set lik.Seter) *Brain {
	brain := &Brain{}
	if brain.loadSet(set) {
		return brain
	} else {
		return nil
	}
}

func (brain *Brain) SaveToFile(fileName string) bool {
	if set := brain.SaveToSet(); set == nil {
		return false
	} else {
		data := set.Format("")
		if err := os.WriteFile(fileName, []byte(data), 0644); err != nil {
			return false
		}
		return true
	}
}

func (brain *Brain) SaveToSet() lik.Seter {
	return brain.saveSet()
}

func (brain *Brain) Clone() *Brain {
	dubl := &Brain{network: brain.network}
	dubl.weights = make([]float64, len(brain.weights))
	copy(dubl.weights, brain.weights)
	return dubl
}

func (brain *Brain) CalculateInputs(inputs []float64) (outputs []float64) {
	network := brain.network
	outputs = make([]float64, network.countOutputs)
	fo := len(network.neurons) - network.countOutputs
	for nu, neuron := range network.neurons {
		var value float64
		if nu < network.countInputs {
			value = inputs[nu]
		} else {
			value = brain.CalculateNeuron(neuron)
		}
		brain.values[nu] = value
		if nu >= fo {
			outputs[nu-fo] = value
		}
	}
	return
}

func (brain *Brain) ProbeInputs(inputs []float64, outputs []float64) (quality float64) {
	result := brain.CalculateInputs(inputs)
	delta := 0.0
	for n, val := range outputs {
		res := result[n]
		if res >= val {
			delta += res - val
		} else {
			delta += val - res
		}
	}
	quality = 1 / (1 + delta)
	return
}

func (brain *Brain) ProbeTest(test *Test) (quality float64) {
	total := 0.0
	count := 0
	for ni, inputs := range test.inputs {
		if ni < len(test.outputs) {
			qua := brain.ProbeInputs(inputs, test.outputs[ni])
			total += qua
			count++
		}
	}
	if count > 0 {
		quality = total / float64(count)
	} else {
		quality = 0
	}
	return
}

func (brain *Brain) CalculateNeuron(neuron *Neuron) (result float64) {
	cNeurons := len(brain.network.neurons)
	iNeuron := neuron.index
	result = brain.weights[iNeuron]
	for _, input := range neuron.inputs {
		iLink := input.index
		source := input.source
		value := brain.values[source.index]
		weight := brain.weights[cNeurons*3+iLink]
		result += value * weight
	}
	if len(neuron.inputs) > 0 && len(neuron.outputs) > 0 {
		if min := brain.weights[cNeurons+iNeuron]; result < min {
			result = min
		} else if max := brain.weights[cNeurons*2+iNeuron]; result > max {
			result = max
		}
	}
	return
}

func (brain *Brain) initializeWeights() {
	cNeurons := len(brain.network.neurons)
	cLinks := len(brain.network.links)
	brain.weights = make([]float64, cNeurons*3+cLinks)
	for nw := range brain.weights {
		if nw < cNeurons {
			brain.weights[nw] = 0.5
		} else if nw < cNeurons*2 {
			brain.weights[nw] = 0.0
		} else if nw < cNeurons*3 {
			brain.weights[nw] = 1.0
		} else {
			target := brain.network.links[nw-cNeurons*3].target
			ins := len(target.inputs)
			if ins > 0 {
				brain.weights[nw] = 1.0 / float64(ins)
			} else {
				brain.weights[nw] = 1.0
			}
		}
	}
}

func (brain *Brain) loadSet(set lik.Seter) bool {
	cInputs := int(set.GetInt("count_inputs"))
	cOutputs := int(set.GetInt("count_outputs"))
	levels := set.GetList("volume_levels")
	neurons := set.GetList("neurons")
	links := set.GetList("links")
	weights := set.GetList("weights")
	if cInputs < 1 || cOutputs < 1 || levels == nil || neurons == nil || links == nil || weights == nil {
		return false
	}
	if levels.Count() < 2 || neurons.Count() < cInputs+cOutputs || links.Count() < 1 {
		return false
	}
	network := &Network{}
	brain.network = network
	network.countInputs = cInputs
	network.countOutputs = cOutputs
	cLevels := levels.Count()
	cNeurons := 0
	network.volumeLevels = make([]int, cLevels)
	for lev := 0; lev < cLevels; lev++ {
		volume := int(levels.GetInt(lev))
		if lev == 0 && volume != cInputs {
			return false
		} else if lev == cLevels-1 && volume != cOutputs {
			return false
		}
		network.volumeLevels[lev] = volume
		cNeurons += volume
	}
	if cNeurons != neurons.Count() {
		return false
	}
	network.neurons = make([]*Neuron, neurons.Count())
	for n := 0; n < neurons.Count(); n++ {
		data := neurons.GetSet(n)
		if data == nil {
			return false
		}
		index := int(data.GetInt("index"))
		neuron := buildNeuron(index)
		network.neurons[index] = neuron
	}
	network.links = make([]*Link, links.Count())
	for n := 0; n < links.Count(); n++ {
		data := links.GetSet(n)
		if data == nil {
			return false
		}
		index := int(data.GetInt("index"))
		source := int(data.GetInt("source"))
		target := int(data.GetInt("target"))
		src := network.neurons[source]
		trg := network.neurons[target]
		link := buildLink(index, src, trg)
		network.links[index] = link
	}
	cWeights := weights.Count()
	brain.weights = make([]float64, cWeights)
	for n := 0; n < cWeights; n++ {
		brain.weights[n] = weights.GetFloat(n)
	}
	brain.values = make([]float64, neurons.Count())
	return true
}

func (brain *Brain) saveSet() lik.Seter {
	set := lik.BuildSet()
	network := brain.network
	set.SetValue("count_inputs", network.countInputs)
	set.SetValue("count_outputs", network.countOutputs)
	levels := set.AddList("volume_levels")
	for _, volume := range network.volumeLevels {
		levels.AddItems(volume)
	}
	neurons := set.AddList("neurons")
	for _, neuron := range network.neurons {
		neu := lik.BuildSet("index", neuron.index)
		neurons.AddItems(neu)
	}
	links := set.AddList("links")
	for _, link := range network.links {
		lin := lik.BuildSet("index", link.index)
		lin.SetValue("source", link.source.index)
		lin.SetValue("target", link.target.index)
		links.AddItems(lin)
	}
	weights := set.AddList("weights")
	for _, weight := range brain.weights {
		weights.AddItems(weight)
	}

	return set
}

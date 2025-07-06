package shambala

type Link struct {
	index  int
	source *Neuron
	target *Neuron
}

func buildLink(index int, source *Neuron, target *Neuron) *Link {
	link := &Link{index: index, source: source, target: target}
	source.outputs = append(source.outputs, link)
	target.inputs = append(target.inputs, link)
	return link
}

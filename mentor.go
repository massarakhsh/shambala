package shambala

import (
	"fmt"
	"math/rand"
)

type Mentor struct {
	brain   *Brain
	ograde  int
	origins []float64
	test    *Test

	direction  int
	dirList    [8]int
	dirQual    [8]float64
	dirLeft    int
	dirSuccess int
	isScan     bool

	step    float64
	quality float64
}

func buildMentor(brain *Brain, test *Test) *Mentor {
	mentor := &Mentor{brain: brain, test: test}
	mentor.ograde = len(brain.weights)
	mentor.origins = make([]float64, mentor.ograde)
	copy(mentor.origins, brain.weights)
	mentor.quality = mentor.brain.ProbeTest(mentor.test)
	mentor.step = 1.0
	return mentor
}

func (mentor *Mentor) Probe() (quality float64) {
	quality = mentor.brain.ProbeTest(mentor.test)
	return quality
}

func (mentor *Mentor) GetQuality() (quality float64) {
	quality = mentor.quality
	return
}

func (mentor *Mentor) TrainingStep() float64 {
	if mentor.tryModify() {
		if quali := mentor.brain.ProbeTest(mentor.test); quali > mentor.quality {
			mentor.quality = quali
			mentor.trySuccess()
			if mentor.isScan {
				mentor.tryScan()
			}
		} else {
			mentor.tryFalse()
		}
	}
	return mentor.quality
}

func (mentor *Mentor) tryModify() bool {
	if mentor.isScan {
		mentor.direction++
		if mentor.direction >= mentor.ograde*2 {
			if mentor.dirLeft > 0 {
				//fmt.Printf("Stop scanning: [%d] %d/%.6f\n", mentor.dirLeft, mentor.dirList[0], mentor.dirQual[0])
				mentor.isScan = false
				mentor.dirSuccess = 0
			} else {
				mentor.step /= 2
				if mentor.step < 1e-6 {
					mentor.step = 0.1 + rand.Float64()*0.9
					//fmt.Printf("Return to step: %f\n", mentor.step)
				} else {
					//fmt.Printf("Shift to step: %f\n", mentor.step)
				}
				mentor.direction = 0
			}
		}
	}
	if !mentor.isScan && mentor.dirSuccess == 0 {
		if mentor.dirLeft == 0 {
			//fmt.Printf("Start scanning\n")
			mentor.isScan = true
			mentor.direction = 0
		} else {
			mentor.direction = mentor.dirList[0]
			mentor.dirLeft--
			for d := 0; d < mentor.dirLeft; d++ {
				mentor.dirList[d] = mentor.dirList[d+1]
			}
		}
	}
	dirReal := mentor.direction
	step := mentor.step
	if dirReal >= mentor.ograde {
		dirReal -= mentor.ograde
		step = -step
	}

	brain := mentor.brain
	cNeurons := len(brain.network.neurons)
	cLinks := len(brain.network.links)
	weight := brain.weights[dirReal] + step
	brain.weights[dirReal] = weight
	if dirReal < cNeurons {
	} else if dirReal < cNeurons*2 {
		if weight > brain.weights[dirReal+cNeurons] {
			brain.weights[dirReal+cNeurons] = weight
		}
	} else if dirReal < cNeurons*3 {
		if weight < brain.weights[dirReal-cNeurons] {
			brain.weights[dirReal-cNeurons] = weight
		}
	} else if fLink := cNeurons * 3; dirReal < fLink+cLinks {
		iLink := dirReal - fLink
		link := brain.network.links[iLink]
		target := link.target
		ins := len(target.inputs)
		outs := len(target.outputs)
		if ins > 0 && outs > 0 {
			summa := 0.0
			for _, link := range target.inputs {
				summa += brain.weights[fLink+link.index]
			}
			if delta := summa - 1.0; delta < -1e-6 || delta > 1e-6 {
				delta /= float64(ins)
				for _, link := range target.inputs {
					brain.weights[fLink+link.index] -= delta
				}
			}
		}
	}
	return true
}

func (mentor *Mentor) trySuccess() {
	copy(mentor.origins, mentor.brain.weights)
	mentor.dirSuccess++
	if mentor.dirSuccess >= 10 {
		mentor.dirSuccess = 0
	}
}

func (mentor *Mentor) tryFalse() {
	copy(mentor.brain.weights, mentor.origins)
	mentor.dirSuccess = 0
}

func (mentor *Mentor) tryScan() {
	iLeft := mentor.dirLeft
	for ; iLeft > 0; iLeft-- {
		if mentor.quality <= mentor.dirQual[iLeft-1] {
			break
		}
		if iLeft < len(mentor.dirList) {
			mentor.dirList[iLeft] = mentor.dirList[iLeft-1]
			mentor.dirQual[iLeft] = mentor.dirQual[iLeft-1]
		}
	}
	if iLeft < len(mentor.dirList) {
		mentor.dirList[iLeft] = mentor.direction
		mentor.dirQual[iLeft] = mentor.quality
	}
	if mentor.dirLeft < len(mentor.dirList) {
		mentor.dirLeft++
	}
}

func (mentor *Mentor) PrintInOut(inputs []float64) {
	outputs := mentor.brain.CalculateInputs(inputs)
	fmt.Print(inputs[0])
	fmt.Print(" ")
	fmt.Print(outputs[0])
	fmt.Print("\n")
}

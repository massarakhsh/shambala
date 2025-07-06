package neu

import (
	"fmt"
)

type Mentor struct {
	brain   *Brain
	ograde  int
	origins []float64
	test    *Test

	direction  int
	dirReal    int
	dirList    [8]int
	dirQual    [8]float64
	dirLeft    int
	dirSuccess int
	isScan     bool

	step    float64
	quality float64
	usedMin bool
	usedMax bool
}

func BuildMentor(brain *Brain, test *Test) *Mentor {
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
				//fmt.Printf("Fault scanning, scan=%f\n", mentor.step)
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
	mentor.dirReal = mentor.direction
	step := mentor.step
	if mentor.dirReal >= mentor.ograde {
		mentor.dirReal -= mentor.ograde
		step = -step
	}
	weight := mentor.brain.weights[mentor.dirReal] + step
	mentor.brain.weights[mentor.dirReal] = weight
	mentor.usedMin = false
	mentor.usedMax = false
	if cNeurons := len(mentor.brain.network.neurons); mentor.dirReal < cNeurons {
	} else if mentor.dirReal < cNeurons*2 {
		if weight > mentor.brain.weights[mentor.dirReal+cNeurons] {
			mentor.usedMax = true
			mentor.brain.weights[mentor.dirReal+cNeurons] = weight
		}
	} else if mentor.dirReal < cNeurons*3 {
		if weight < mentor.brain.weights[mentor.dirReal-cNeurons] {
			mentor.usedMin = true
			mentor.brain.weights[mentor.dirReal-cNeurons] = weight
		}
	}
	return true
}

func (mentor *Mentor) trySuccess() {
	cNeurons := len(mentor.brain.network.neurons)
	mentor.origins[mentor.dirReal] = mentor.brain.weights[mentor.dirReal]
	if mentor.usedMax {
		mentor.origins[mentor.dirReal+cNeurons] = mentor.brain.weights[mentor.dirReal+cNeurons]
	}
	if mentor.usedMin {
		mentor.origins[mentor.dirReal-cNeurons] = mentor.brain.weights[mentor.dirReal-cNeurons]
	}
	if mentor.isScan {
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
	mentor.dirSuccess++
	if mentor.dirSuccess >= 10 {
		mentor.dirSuccess = 0
	}
}

func (mentor *Mentor) tryFalse() {
	cNeurons := len(mentor.brain.network.neurons)
	mentor.brain.weights[mentor.dirReal] = mentor.origins[mentor.dirReal]
	if mentor.usedMax {
		mentor.brain.weights[mentor.dirReal+cNeurons] = mentor.origins[mentor.dirReal+cNeurons]
	}
	if mentor.usedMin {
		mentor.brain.weights[mentor.dirReal-cNeurons] = mentor.origins[mentor.dirReal-cNeurons]
	}
	mentor.dirSuccess = 0
}

func (mentor *Mentor) PrintInOut(inputs []float64) {
	outputs := mentor.brain.CalculateInputs(inputs)
	fmt.Print(inputs[0])
	fmt.Print(" ")
	fmt.Print(outputs[0])
	fmt.Print("\n")
}

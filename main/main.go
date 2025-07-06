package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/massarakhsh/shambala"
)

func main() {
	testingSin()
}

func testingSin() {
	brain := shambala.LoadBrainFile("sinus.sha")
	if brain == nil {
		brain = shambala.BuildBrainFull(1, 100, 1)
	}
	var sources [][]float64
	var targets [][]float64
	for n := 0; n < 1000; n++ {
		arg := rand.Float64() * 2 * math.Pi
		sin := math.Sin(arg)
		inputs := []float64{arg}
		outputs := []float64{sin}
		sources = append(sources, inputs)
		targets = append(targets, outputs)
	}
	test := shambala.BuildTest("test", "example", sources, targets)
	mentor := shambala.BuildMentorTest(brain, test)

	lastSay := time.Now()
	quality := mentor.GetQuality()
	fmt.Printf("_start_: %.6f\n", quality)
	for {
		if qua := mentor.TrainingStep(); qua > quality {
			quality = qua
			if time.Since(lastSay) > time.Second {
				fmt.Printf("quality: %.6f\n", quality)
				lastSay = time.Now()
				brain.SaveToFile("sinus.sha")
			}
		}
	}

	// for n := 0; n <= 64; n++ {
	// 	arg := float64(n) * math.Pi * 2 / 64
	// 	mentor.PrintInOut([]float64{arg})
	// }
}

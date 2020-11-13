package custom

import (
	"log"
	"time"

	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
)

func Brake2Arduino(tc *traincontrol.TrainControl, blocks [4]string, start int, target int, step int, dur time.Duration) {
	start_time := time.Now()
	i := start

	log.Println("----------------Brake Ramp started")
	for i >= target {
		for _, block := range blocks { //for each block do action
			tc.SetBlockSpeed(string(block[0]), i)
		}
		time.Sleep(dur * time.Millisecond)
		i = i - step
	}

	end_time := time.Now()
	brake_duration := end_time.Sub(start_time).Seconds()
	log.Println("----------------Brake Ramp done after: ", brake_duration)
}

func Accelerate2Arduino(tc *traincontrol.TrainControl, blocks [4]string, start int, target int, step int, dur time.Duration) {
	start_time := time.Now()

	log.Println("----------------Acceleration Ramp started")
	for i := start; i <= target; i = i + step {
		for _, block := range blocks { //for each block do action
			tc.SetBlockSpeed(string(block[0]), i)
		}
		time.Sleep(dur * time.Millisecond)
	}

	end_time := time.Now()
	acceleration_duration := end_time.Sub(start_time).Seconds()
	log.Println("----------------Acceleration Ramp done after: ", acceleration_duration)
}

func PartialReset2Arduino(tc *traincontrol.TrainControl, block string) {

	log.Println("----------------Track changed. Reset for block: ", block)
	tc.SetBlockSpeed(string(block[0]), 0)
	tc.SetBlockDirection(string(block[0]), "s")
}

func PartialSet2Arduino(tc *traincontrol.TrainControl, block string, direction string, speed int) {

	log.Println("----------------Track changed. Set for block: ", block)
	tc.SetBlockSpeed(string(block[0]), speed)
	tc.SetBlockDirection(string(block[0]), direction)
}

func FullReset2Arduino(tc *traincontrol.TrainControl, blocks [4]string) {

	log.Println("----------------Reset started for blocks: ", blocks)
	for _, block := range blocks {
		tc.SetBlockSpeed(string(block[0]), 0)
		tc.SetBlockDirection(string(block[0]), "s")
	}

	log.Println("----------------Reset done")
}

func Switches2Arduino(tc *traincontrol.TrainControl, block string) {
	if block == "aw" {
		log.Println("----------------Send Switches for Track 1 west in-/outbound to Arduino")
		tc.SetSwitch("a", "0")
		tc.SetSwitch("b", "0")
		tc.SetSwitch("e", "0")
		tc.SetSwitch("f", "0")
	}
	if block == "ae" {
		log.Println("----------------Send Switches for Track 1 to east in-/outbound Arduino")
		tc.SetSwitch("a", "0")
		tc.SetSwitch("b", "0")
		tc.SetSwitch("e", "0")
		tc.SetSwitch("f", "0")
	}
	if block == "bw" {
		log.Println("----------------Send Switches for Track 2 west in-/outbound to Arduino")
		tc.SetSwitch("a", "0")
		tc.SetSwitch("b", "1")
		tc.SetSwitch("d", "0")
		tc.SetSwitch("e", "1")
		tc.SetSwitch("f", "0")
	}
	if block == "be" {
		log.Println("----------------Send Switches for Track 2 east in-/outbound to Arduino")
		tc.SetSwitch("a", "0")
		tc.SetSwitch("b", "1")
		tc.SetSwitch("d", "0")
		tc.SetSwitch("e", "1")
		tc.SetSwitch("f", "0")
	}
	if block == "cw" {
		log.Println("----------------Send Switches for Track 3 west in-/outbound to Arduino")
		tc.SetSwitch("a", "1")
		tc.SetSwitch("c", "0")
		tc.SetSwitch("d", "1")
		tc.SetSwitch("e", "1")
		tc.SetSwitch("f", "0")
	}
	if block == "ce" {
		log.Println("----------------Send Switches for Track 3 east in-/outbound to Arduino")
		tc.SetSwitch("a", "1")
		tc.SetSwitch("c", "0")
		tc.SetSwitch("d", "1")
		tc.SetSwitch("e", "1")
		tc.SetSwitch("f", "0")
	}
	if block == "dw" {
		log.Println("----------------Send Switches for Track 4 west in-/outbound to Arduino")
		tc.SetSwitch("a", "1")
		tc.SetSwitch("c", "1")
		tc.SetSwitch("f", "1")
	}
	if block == "de" {
		log.Println("----------------Send Switches for Track 4 east in-/outbound to Arduino")
		tc.SetSwitch("a", "1")
		tc.SetSwitch("c", "1")
		tc.SetSwitch("f", "1")
	}
}

func Direction2Arduino(tc *traincontrol.TrainControl, block string, direction string) {
	tc.SetBlockDirection(string(block[0]), direction)
}

func EmergencyStop2Arduino(tc *traincontrol.TrainControl) {
	tc.SetBlockDirection("a", "s")
	tc.SetBlockDirection("b", "s")
	tc.SetBlockDirection("c", "s")
	tc.SetBlockDirection("c", "s")
	tc.SetBlockDirection("d", "s")
	tc.SetBlockDirection("f", "s")
	tc.SetBlockDirection("g", "s")

	tc.SetBlockSpeed("a", 0)
	tc.SetBlockSpeed("b", 0)
	tc.SetBlockSpeed("c", 0)
	tc.SetBlockSpeed("d", 0)
	tc.SetBlockSpeed("f", 0)
	tc.SetBlockSpeed("g", 0)
}

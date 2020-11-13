package custom

import (
	"log"
	"time"

	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
)

// Brake2Arduino bla
func Brake2Arduino(tc *traincontrol.TrainControl, blocks [4]string, start int, target int, step int, dur time.Duration) {
	startTime := time.Now()
	i := start

	log.Println("----------------Brake Ramp started")
	for i >= target {
		for _, block := range blocks { //for each block do action
			tc.SetBlockSpeed(string(block[0]), i)
		}
		time.Sleep(dur * time.Millisecond)
		i = i - step
	}

	endTime := time.Now()
	brakeDuration := endTime.Sub(startTime).Seconds()
	log.Println("----------------Brake Ramp done after: ", brakeDuration)
}

// Accelerate2Arduino bla
func Accelerate2Arduino(tc *traincontrol.TrainControl, blocks [4]string, start int, target int, step int, dur time.Duration) {
	startTime := time.Now()

	log.Println("----------------Acceleration Ramp started")
	for i := start; i <= target; i = i + step {
		for _, block := range blocks { //for each block do action
			tc.SetBlockSpeed(string(block[0]), i)
		}
		time.Sleep(dur * time.Millisecond)
	}

	endTime := time.Now()
	accelerationDuration := endTime.Sub(startTime).Seconds()
	log.Println("----------------Acceleration Ramp done after: ", accelerationDuration)
}

// Direction2Arduino sets direction for block
func Direction2Arduino(tc *traincontrol.TrainControl, block byte, direction string) {
	tc.SetBlockDirection(string(block), direction)
}

// Speed2Arduino sets the speed of an arduino with a byte
func Speed2Arduino(tc *traincontrol.TrainControl, block byte, speed int) {
	tc.SetBlockSpeed(string(block), speed)
}

// PartialResetBlock2Arduino resets a block
func PartialResetBlock2Arduino(tc *traincontrol.TrainControl, block byte) {
	Speed2Arduino(tc, block, 0)
	Direction2Arduino(tc, block, "s")
}

// PartialSet2Arduino set block
func PartialSet2Arduino(tc *traincontrol.TrainControl, block byte, direction string, speed int) {
	Speed2Arduino(tc, block, speed)
	Direction2Arduino(tc, block, direction)
}

// Switches2Arduino alters junctions
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

// EmergencyStop2Arduino stops all tracks
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

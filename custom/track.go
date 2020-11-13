package custom

import (
	"log"
	"time"

	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"

)

var EmptyBlock = [4]string{"", "", "f", "g"}

//global variables here
var actualBlocks [4]string = EmptyBlock
var targetBlocks [4]string = EmptyBlock
var direction string = "s"
var actualDirection string = "s"
var targetDirection string = "s"
var actualSpeed int = 0
var targetSpeed int = 0
var previousSpeed int = 0


// ControlRunner performs an arduino loop with controlCycleDuration
func ControlRunner(tc *traincontrol.TrainControl) {

	const controlCycleDuration = 20 * time.Microsecond
	var lastControlCycle = time.Unix(0, 0)

	for {
		var now = time.Now()
		var waitTime = now.Sub(lastControlCycle)
		var sleepTime = controlCycleDuration - waitTime
		if sleepTime < 1 {
			Control(tc)
			lastControlCycle = now
		} else {
			time.Sleep(sleepTime)
		}
	}
}

// Control is run in a short interval
func Control(tc *traincontrol.TrainControl) {
	if IsDriveable() {
		if targetDirection != actualDirection {
			actualDirection = targetDirection
			SetBlocksDirection(tc, actualBlocks, targetDirection)
		}

		if targetSpeed != actualSpeed {
			actualSpeed = targetSpeed
			SetBlocksSpeed(tc, actualBlocks, targetSpeed)
		}
	}

	if targetBlocks != actualBlocks {
		actualBlocks = targetBlocks
		SetSwitches(tc, actualBlocks)
		ResetInactiveBlocks(tc, actualBlocks)
	}

}

// SetDirection sets the direction
func SetDirection(tc *traincontrol.TrainControl, d string) {
	targetDirection = d
	tc.PublishMessage(struct {
		Direction string `json:"direction"`
	}{
		Direction: targetDirection,
	})
	log.Println("----------------Direction set: ", targetDirection)
}

// SetSpeed sets the speed
func SetSpeed(tc *traincontrol.TrainControl, s int) {
	targetSpeed = s
	previousSpeed = actualSpeed
	log.Println("----------------Speed set: ", s)

	tc.PublishMessage(struct {
		Speed int `json:"speed"`
	}{
		Speed: s,
	}) //synchronize all websites with set state
}

// SetTrack sets the track
func SetTrack(tc *traincontrol.TrainControl, track string) {
	var switchLocation = string(GetSwitchLocation(track))
	var block = string(GetBlock(track))
	switch switchLocation {
	case "o":
		targetBlocks[0] = block + switchLocation
	case "w":
		targetBlocks[1] = block + switchLocation
	default:
		targetBlocks[0] = block + "o"
		targetBlocks[1] = block + "w"
	}
	log.Println("----------------setTrack: Blocks set: ", targetBlocks)
	tc.PublishMessage(struct {
		Blocks [4]string `json:"blocks"`
	}{
		Blocks: targetBlocks,
	})
	//synchronize all websites with set state
}

// IsDriveable checks wheather a train can drive
func IsDriveable() bool {
	if targetDirection == "s" {
		return false
	}
	for _, block := range targetBlocks {
		if block == "" {
			return false
		}
	}
	return true
}

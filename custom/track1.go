package custom

import (
	"log"

	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"

)

//global variables here
var blocks [4]string
var direction string = "s"
var previousTrack [4]string
var speed int = 0
var previousSpeed int = 0

// func ManualControl(tc *traincontrol.TrainControl, direction string, speed int, blocks [4]string) {
// 	train := tc.GetActiveTrain()
// 	train = tc.Trains["N700"]

// 	if blocks == [4]string{"", "", "", ""} { //exit manual Control completely when only direction was set (and no blocks set until now)
// 		return
// 	}

// 	actualTrack := blocks
// 	if previousTrack != actualTrack && previousTrack != [4]string{"", "", "", ""} { //Partial reset of tracks (a,b,c,d) in case of track change
// 		log.Println("----------------Send Reset Command for previous Track to Arduino") //Track need to be set to stop and zero in case another track is choosen
// 		if previousTrack[0] != actualTrack[0] {                                         // in case one block changed while speedlock = 0
// 			PartialReset2Arduino(tc, previousTrack[0])               //reset block a,b,c or d (Direction and Speed)
// 			PartialSet2Arduino(tc, actualTrack[0], direction, speed) //Set Direction and Speed directly from user input
// 		}
// 		if previousTrack[1] != actualTrack[1] {
// 			PartialReset2Arduino(tc, previousTrack[1])               //reset block a,b,c or d (Direction and Speed)
// 			PartialSet2Arduino(tc, actualTrack[1], direction, speed) //Set Direction and Speed directly from user input
// 		}
// 	}

// 	previousDirection := string(tc.Blocks[[]rune(blocks[0])[0]].Direction) //gets direction requested from arduino. compare to last input
// 	actualDirection := direction                                           //send Direction
// 	if previousDirection != actualDirection {                              //Execution only by change
// 		log.Println("----------------Manual Control started. (Track & Direction)")
// 		log.Println("----------------Previous Direction was: ", previousDirection)
// 		log.Println("----------------Actual Direction is: ", actualDirection)
// 		for _, block := range blocks {
// 			Direction2Arduino(tc, block, actualDirection)
// 		}
// 	}

// 	if previousTrack != actualTrack { //send Track to Arduino
// 		log.Println("----------------Manual Control started. (Track & Direction)") //Execution only by change
// 		log.Println("----------------Previous Track was: ", previousTrack)
// 		log.Println("----------------Actual Track is: ", actualTrack)
// 		// if flag_driveCircle == 1 { //in case of do circle just ssend command once
// 		// 	Switches2Arduino(tc, blocks[0])
// 		// } else { //iterate through blocks array to set both tracks
// 		for _, block := range blocks { //send command to set junctions to new track
// 			Switches2Arduino(tc, block)
// 		}
// 		// }
// 		previousTrack = blocks //after track was set, store information in previous track for later comparision
// 	}

// 	previousSpeed = tc.Blocks[[]rune(blocks[2])[0]].Speed // compare speed in junctions or open track in case junctions were switched in between. shall prevent intermediate full acceleration
// 	actualSpeed := speed
// 	if flag_direction == 1 && flag_track != 0 && previousSpeed != actualSpeed && flag_speedLock == 0 { //send Speed to Arduino
// 		log.Println("----------------Manual Control started. (Speed)")
// 		log.Println("----------------Send Speed Command to Arduino")
// 		log.Println("----------------Previous Speed was: ", previousSpeed)
// 		log.Println("----------------Actual Speed is: ", actualSpeed)

// 		if previousSpeed != 0 && actualSpeed == 0 { //Execution only by change
// 			flag_speedLock = 1
// 			log.Println("----------------Braking and Full Reset of Blocks ...")
// 			Brake2Arduino(tc, blocks, previousSpeed, actualSpeed, train.Brake.Step, train.Brake.Time)
// 			FullReset2Arduino(tc, blocks)
// 			tc.PublishMessage(struct {
// 				Speed int `json:"speed"`
// 			}{
// 				Speed: actualSpeed,
// 			}) //synchronize all websites with set state. in case sombody plaxed arpund during speedlock
// 			flag_speedLock = 0

// 		} else if previousSpeed < actualSpeed {
// 			flag_speedLock = 1
// 			log.Println("----------------Accelerating ...")
// 			Accelerate2Arduino(tc, blocks, previousSpeed, actualSpeed, train.Accelerate.Step, train.Accelerate.Time)
// 			tc.PublishMessage(struct {
// 				Speed int `json:"speed"`
// 			}{
// 				Speed: actualSpeed,
// 			}) //synchronize all websites with set state. in case sombody plaxed arpund during speedlock
// 			flag_speedLock = 0

// 		} else if previousSpeed > actualSpeed {
// 			flag_speedLock = 1
// 			log.Println("----------------Braking ...")
// 			Brake2Arduino(tc, blocks, previousSpeed, actualSpeed, train.Brake.Step, train.Brake.Time)
// 			tc.PublishMessage(struct {
// 				Speed int `json:"speed"`
// 			}{
// 				Speed: actualSpeed,
// 			}) //synchronize all websites with set state. in case sombody plaxed arpund during speedlock
// 			flag_speedLock = 0

// 		}

// 	}
// }

// SetDirection sets the direction
func SetDirection(tc *traincontrol.TrainControl, d string) {
	if d != direction {
		direction = d
		tc.PublishMessage(struct {
			Direction string `json:"direction"`
		}{
			Direction: d,
		})
	}
	log.Println("----------------Direction set: ", direction)
}

// SetSpeed sets the fucking speed
func SetSpeed(tc *traincontrol.TrainControl, s int) {
	speed = s
	log.Println("----------------Speed set: ", speed)

	tc.PublishMessage(struct {
		Speed int `json:"speed"`
	}{
		Speed: s,
	}) //synchronize all websites with set state

}

// SetTrack sets the track
func SetTrack(tc *traincontrol.TrainControl, t int) {
	if t == 1 {
		blocks = [4]string{"aw", "ao", "f", "g"}
		log.Println("----------------driveCircle: Blocks set: ", blocks)

		publishBlock(tc, blocks)
	}
	if t == 2 {
		blocks = [4]string{"bw", "bo", "f", "g"}
		log.Println("----------------driveCircle: Blocks set: ", blocks)

		publishBlock(tc, blocks)
	}
	if t == 3 {
		blocks = [4]string{"cw", "co", "f", "g"}
		log.Println("----------------driveCircle: Blocks set: ", blocks)

		publishBlock(tc, blocks)
	}
	if t == 4 {
		blocks = [4]string{"dw", "do", "f", "g"}
		log.Println("----------------driveCircle: Blocks set: ", blocks)
		publishBlock(tc, blocks)
	}
}

func publishBlock(tc *traincontrol.TrainControl, blocks [4]string) {
	tc.PublishMessage(struct {
		Blocks [4]string `json:"blocks"`
	}{
		Blocks: blocks,
	})
	//synchronize all websites with set state
}

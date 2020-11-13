package custom

import (
	"log"
	"time"

	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
)

const NO_BLOCK = [4]string{"", "", "", ""}

//global variables here
var blocks [4]string
var direction string = "s"
var previousTrack [4]string
var speed int = 0
var previousSpeed int = 0

func ManualControl(tc *traincontrol.TrainControl, direction string, speed int, blocks [4]string) {
	train := tc.GetActiveTrain()
	train = tc.Trains["N700"]

	if blocks == NO_BLOCK { //exit manual Control completely when only direction was set (and no blocks set until now)
		return
	}

	actualTrack := blocks
	if previousTrack != actualTrack && previousTrack != NO_BLOCK { //Partial reset of tracks (a,b,c,d) in case of track change
		log.Println("----------------Send Reset Command for previous Track to Arduino") //Track need to be set to stop and zero in case another track is choosen
		if previousTrack[0] != actualTrack[0] {                                         // in case one block changed while speedlock = 0
			PartialReset2Arduino(tc, previousTrack[0])               //reset block a,b,c or d (Direction and Speed)
			PartialSet2Arduino(tc, actualTrack[0], direction, speed) //Set Direction and Speed directly from user input
		}
		if previousTrack[1] != actualTrack[1] {
			PartialReset2Arduino(tc, previousTrack[1])               //reset block a,b,c or d (Direction and Speed)
			PartialSet2Arduino(tc, actualTrack[1], direction, speed) //Set Direction and Speed directly from user input
		}
	}

	previousDirection := string(tc.Blocks[[]rune(blocks[0])[0]].Direction) //gets direction requested from arduino. compare to last input
	actualDirection := direction                                           //send Direction
	if previousDirection != actualDirection {                              //Execution only by change
		log.Println("----------------Manual Control started. (Track & Direction)")
		log.Println("----------------Previous Direction was: ", previousDirection)
		log.Println("----------------Actual Direction is: ", actualDirection)
		for _, block := range blocks {
			Direction2Arduino(tc, block, actualDirection)
		}
	}

	if previousTrack != actualTrack { //send Track to Arduino
		log.Println("----------------Manual Control started. (Track & Direction)") //Execution only by change
		log.Println("----------------Previous Track was: ", previousTrack)
		log.Println("----------------Actual Track is: ", actualTrack)
		// if flag_driveCircle == 1 { //in case of do circle just ssend command once
		// 	Switches2Arduino(tc, blocks[0])
		// } else { //iterate through blocks array to set both tracks
		for _, block := range blocks { //send command to set junctions to new track
			Switches2Arduino(tc, block)
		}
		// }
		previousTrack = blocks //after track was set, store information in previous track for later comparision
	}

	previousSpeed = tc.Blocks[[]rune(blocks[2])[0]].Speed // compare speed in junctions or open track in case junctions were switched in between. shall prevent intermediate full acceleration
	actualSpeed := speed
	if flag_direction == 1 && flag_track != 0 && previousSpeed != actualSpeed && flag_speedLock == 0 { //send Speed to Arduino
		log.Println("----------------Manual Control started. (Speed)")
		log.Println("----------------Send Speed Command to Arduino")
		log.Println("----------------Previous Speed was: ", previousSpeed)
		log.Println("----------------Actual Speed is: ", actualSpeed)

		if previousSpeed != 0 && actualSpeed == 0 { //Execution only by change
			flag_speedLock = 1
			log.Println("----------------Braking and Full Reset of Blocks ...")
			Brake2Arduino(tc, blocks, previousSpeed, actualSpeed, train.Brake.Step, train.Brake.Time)
			FullReset2Arduino(tc, blocks)
			tc.PublishMessage(struct {
				Speed int `json:"speed"`
			}{
				Speed: actualSpeed,
			}) //synchronize all websites with set state. in case sombody plaxed arpund during speedlock
			flag_speedLock = 0

		} else if previousSpeed < actualSpeed {
			flag_speedLock = 1
			log.Println("----------------Accelerating ...")
			Accelerate2Arduino(tc, blocks, previousSpeed, actualSpeed, train.Accelerate.Step, train.Accelerate.Time)
			tc.PublishMessage(struct {
				Speed int `json:"speed"`
			}{
				Speed: actualSpeed,
			}) //synchronize all websites with set state. in case sombody plaxed arpund during speedlock
			flag_speedLock = 0

		} else if previousSpeed > actualSpeed {
			flag_speedLock = 1
			log.Println("----------------Braking ...")
			Brake2Arduino(tc, blocks, previousSpeed, actualSpeed, train.Brake.Step, train.Brake.Time)
			tc.PublishMessage(struct {
				Speed int `json:"speed"`
			}{
				Speed: actualSpeed,
			}) //synchronize all websites with set state. in case sombody plaxed arpund during speedlock
			flag_speedLock = 0

		}

	}
}

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

	ManualControl(tc, direction, speed, blocks)
}

func SetSpeed(tc *traincontrol.TrainControl, s int) {
	if flag_speedLock == 0 {
		speed = s
		log.Println("----------------Speed set: ", speed)
	}

	tc.PublishMessage(struct {
		Speed int `json:"speed"`
	}{
		Speed: s,
	}) //synchronize all websites with set state

	ManualControl(tc, direction, speed, blocks)
}

func SetTrack(tc *traincontrol.TrainControl, t int) {
	if flag_driveCircle == 1 && flag_speedLock == 0 { // "a,b,c,d" for Track 1-4, "f" for junctions, "g" for open terrain //Track change only possible if no speed was changed, no ramp is running
		if t == 1 {
			blocks = [4]string{"aw", "ao", "f", "g"}
			flag_track = 11 // flag_track: from track -> to track
			log.Println("----------------driveCircle: Blocks set: ", blocks, flag_track)

			tc.PublishMessage(struct {
				Blocks [4]string `json:"blocks"`
			}{
				Blocks: blocks,
			}) //synchronize all websites with set state
		}
		if t == 2 {
			blocks = [4]string{"bw", "bo", "f", "g"}
			flag_track = 22
			log.Println("----------------driveCircle: Blocks set: ", blocks, flag_track)

			tc.PublishMessage(struct {
				Blocks [4]string `json:"blocks"`
			}{
				Blocks: blocks,
			}) //synchronize all websites with set state
		}
		if t == 3 {
			blocks = [4]string{"cw", "co", "f", "g"}
			flag_track = 33
			log.Println("----------------driveCircle: Blocks set: ", blocks, flag_track)

			tc.PublishMessage(struct {
				Blocks [4]string `json:"blocks"`
			}{
				Blocks: blocks,
			}) //synchronize all websites with set state
		}
		if t == 4 {
			blocks = [4]string{"dw", "do", "f", "g"}
			flag_track = 44
			log.Println("----------------driveCircle: Blocks set: ", blocks, flag_track)

			tc.PublishMessage(struct {
				Blocks [4]string `json:"blocks"`
			}{
				Blocks: blocks,
			}) //synchronize all websites with set state
		}
	}

	ManualControl(tc, direction, speed, blocks)
}

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

// func StateFromArduino(tc *traincontrol.TrainControl) {
// 	tc.GetSensorStates()
// 	tc.GetBlockDirections()
// 	tc.GetBlockSpeeds()
// 	tc.GetSwitchStates()
// 	tc.GetSignalStates()
// }

//==============================================================================================================================
//==============================================================================================================================
//==============================================================================================================================
//============================================================ BELOW TEST CODE =================================================
//==============================================================================================================================
//==============================================================================================================================
//==============================================================================================================================

//==============================================================================================================================
//====================================================== S A M P L E S =========================================================
//==============================================================================================================================

// go get -u github.com/codepuree/tilo-railway-company/pkg/traincontrol			//get latest

//==============================================================================================================================
//available functions===========================================================================================================
//==============================================================================================================================
//to Server

// func ManualControl(tc *traincontrol.TrainControl, direction string, speed int, blocks [4]string)
// func SetDirection(tc *traincontrol.TrainControl, d string)
// func SetSpeed(tc *traincontrol.TrainControl, s int)
// func SetTrack(tc *traincontrol.TrainControl, t int)

//to Arduino
// func Brake2Arduino(tc *traincontrol.TrainControl, blocks [4]string, start int, target int, step int, dur time.Duration)
// func Accelerate2Arduino(tc *traincontrol.TrainControl, blocks [4]string, start int, target int, step int, dur time.Duration)
// func PartialReset2Arduino(tc *traincontrol.TrainControl, block string)
// func PartialSet2Arduino(tc *traincontrol.TrainControl, block string, direction string, speed int)
// func FullReset2Arduino(tc *traincontrol.TrainControl, blocks [4]string)
// func Switches2Arduino(tc *traincontrol.TrainControl, block string)
// func Direction2Arduino(tc *traincontrol.TrainControl, block string, direction string)
// func EmergencyStop2Arduino(tc *traincontrol.TrainControl)
// 	tc.SetBlockDirection("a", "s")	//send command to Arduino
// 	tc.SetBlockSpeed("a", 0)
// 	tc.SetSwitch("a", "1")
//==============================================================================================================================
//available States==============================================================================================================
//==============================================================================================================================
//	speed := tc.Blocks['a'].Speed
//	dir := tc.Blocks['a'].Direction
//	Step := tc.Trains["N700"].Accelerate.Step
//	Step_time := tc.Trains["N700"].Accelerate.Time
//	Max_speed := tc.Trains["N700"].MaxSpeed
//	log.Println("",Step, Step_time, Max_speed)
//
//	if tc.Blocks['a'].Direction == traincontrol.Forward {}	// do action when ...
//	if tc.Blocks['a'].Direction == 'f' {}
//  if tc.Sensors[15].State == false	// do action when state reached
//  tc.Sensors[15].Await(false)  // hold program until state reached
//  tc.Sensors[15].CountTo(10)	// hold program and do action when state reached
// 	tc.GetSensorStates()	//request latest states from arduino. send command to arduino
// 	tc.GetBlockDirections()	//request latest states from arduino. send command to arduino
// 	tc.GetBlockSpeeds()		//request latest states from arduino. send command to arduino
// 	tc.GetSwitchStates()	//request latest states from arduino. send command to arduino
// 	tc.GetSignalStates()	//request latest states from arduino. send command to arduino
// traincontrol.Signal{}.ID		// get Stetes from Server
// traincontrol.Signal{}.State
// traincontrol.Signal{}.Color
// traincontrol.Sensor{}.ID
// traincontrol.Sensor{}.State
// traincontrol.Switch{}.ID
// traincontrol.Switch{}.State
// traincontrol.Block{}.ID
// traincontrol.Block{}.Speed
// traincontrol.Block{}.Direction
// traincontrol.Train{}.Accelerate
// traincontrol.Train{}.Accelerate.Step
// traincontrol.Train{}.Accelerate.Time
// traincontrol.Train{}.Brake.Step
// traincontrol.Train{}.Brake.Time
// traincontrol.Train{}.CrawlSpeed
// traincontrol.Train{}.MaxSpeed
// traincontrol.Train{}.Name
// traincontrol.Train{}.Block
// traincontrol.Track{}.Blocks
//	log.Println("Direction gesetzt")

//to website
// tc.PublishMessage(tc.Trains["N700"])
// // tc.PublishMessage(tc.Trains["beliebiger String"])

// tc.PublishMessage(struct {
// 	Speed int `json:"speed"`
// }{
// 	Speed: actualSpeed,
// }) //synchronize all websites with set state. in case sombody plaxed arpund during speedlock

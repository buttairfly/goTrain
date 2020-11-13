package custom

import (
	"strings"

	"github.com/codepuree/tilo-railway-company/pkg/traincontrol"
)

const allBlocks = "abcdfg"

// SetBlocksDirection sets the direction for all blocks
func SetBlocksDirection(tc *traincontrol.TrainControl, blocks [4]string, direction string) {
	for _, block := range blocks {
		Direction2Arduino(tc, GetBlock(block), direction)
	}
}

// SetBlocksSpeed sets the speed for all blocks
func SetBlocksSpeed(tc *traincontrol.TrainControl, blocks [4]string, speed int) {
	for _, block := range blocks {
		Speed2Arduino(tc, GetBlock(block), speed)
	}
}

// SetSwitches sets all switches
func SetSwitches(tc *traincontrol.TrainControl, blocks [4]string) {
	Switches2Arduino(tc, blocks[0])
	Switches2Arduino(tc, blocks[1])
}

func GetBlock(block string) byte {
	if len(block) > 0 {
		return block[0]
	}
	return '+'
}

func GetSwitchLocation(block string) byte {
	if len(block) > 1 {
		return block[1]
	}
	return '-'
}

func getInactiveBlocks(blocks [4]string) string {
	var currentBlocks = allBlocks
	for _, block := range blocks {
		currentBlocks = strings.Replace(currentBlocks, string(GetBlock(block)), "", 1)
	}
	return currentBlocks
}

func ResetInactiveBlocks(tc *traincontrol.TrainControl, blocks [4]string) {
	var inactiveBlocks = getInactiveBlocks(blocks)
	for _, block := range inactiveBlocks {
		PartialResetBlock2Arduino(tc, byte(block))
	}
}

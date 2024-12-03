package main

import rl "github.com/gen2brain/raylib-go/raylib"

type CustomColors struct {
	primaryTint    rl.Color
	selectionTint  []rl.Color
	projectileTint rl.Color

	increase             bool
	colorFrameCounter    int
	currentSelectionTint int
}

func SelectionColors(baseColor rl.Color) []rl.Color {
	newColors := make([]rl.Color, 0)

	for i := 0; i < 2; i++ {
		newColors = append(newColors, rl.NewColor(baseColor.R, baseColor.G, baseColor.B, uint8(200+(i*50))))
	}
	return newColors
}

// cycle through colors in seconds: cycleTime int in seconds
func (cc *CustomColors) CycleColorsInTime(cycleTime int) {
	cc.colorFrameCounter++
	if ((cc.colorFrameCounter / 120) % cycleTime) == 1 {

		if cc.currentSelectionTint == len(cc.selectionTint)-1 {
			cc.increase = false
		} else if cc.currentSelectionTint == 0 {
			cc.increase = true
		}

		if cc.increase {
			cc.currentSelectionTint++
		} else {
			cc.currentSelectionTint--
		}

		cc.colorFrameCounter = 0
	}
}

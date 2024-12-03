package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Menu struct {
	GameButton Button
	MenuButton Button
	QuitButton Button
}

func (sc *Scene) DrawMenu() {
	currentCamera := sc.camera
	relativeSize := 30 / currentCamera.Zoom
	relativeSpacing := 2 / currentCamera.Zoom
	defaultColor := rl.White
	// endColor := rl.Black
	endString := "  FLATLAND\nGALAXY WARS"

	// boxPos := rl.GetScreenToWorld2D(rl.NewVector2(0, 200), currentCamera)
	textPos := rl.GetScreenToWorld2D(rl.NewVector2(150, 100), currentCamera)
	// rl.DrawRectangleV(boxPos, rl.NewVector2(40*relativeSize, 10*relativeSize), endColor)
	rl.DrawTextEx(rl.GetFontDefault(), endString, textPos, 4*relativeSize, relativeSpacing, defaultColor)

	sc.DrawButton(&sc.QuitButton)
	sc.DrawButton(&sc.GameButton)
}

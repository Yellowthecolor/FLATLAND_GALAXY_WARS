package main

import rl "github.com/gen2brain/raylib-go/raylib"

func (sc *Scene) DrawPauseMenu() {
	sc.DrawCredit()
	currentCamera := sc.camera
	relativeSize := 30 / currentCamera.Zoom
	relativeSpacing := 2 / currentCamera.Zoom
	defaultColor := rl.White
	quitString := "Game Paused"

	boxPos := rl.GetScreenToWorld2D(rl.NewVector2(0, 200), currentCamera)
	textPos := rl.GetScreenToWorld2D(rl.NewVector2(400, 250), currentCamera)
	rl.DrawRectangleV(boxPos, rl.NewVector2(40*relativeSize, 10*relativeSize), rl.DarkGray)
	rl.DrawTextEx(rl.GetFontDefault(), quitString, textPos, 1.5*relativeSize, relativeSpacing, defaultColor)

	sc.DrawButton(&sc.MenuButton)
	sc.DrawButton(&sc.QuitButton)

}

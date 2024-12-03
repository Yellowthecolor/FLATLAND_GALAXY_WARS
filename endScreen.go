package main

import rl "github.com/gen2brain/raylib-go/raylib"

func (sc *Scene) DrawEndScreen() {
	currentCamera := sc.camera
	relativeSize := 30 / currentCamera.Zoom
	relativeSpacing := 2 / currentCamera.Zoom
	defaultColor := rl.White
	endColor := rl.Black
	endString := "Game Over"

	if sc.playerSet.hasWon {
		endColor = sc.playerSet.teamBase.primaryTint
		endString = "You have arrived SQUARE at victory!"
		if !sc.musicChanged {
			sc.PlayMusic(sc.winMusic)
			sc.musicChanged = true
		}
	} else if sc.enemySet.hasWon {
		endColor = sc.enemySet.teamBase.primaryTint
		endString = "They had you CIRCLING like a dog!"
		if !sc.musicChanged {
			sc.PlayMusic(sc.loseMusic)
			sc.musicChanged = true
		}
	}

	boxPos := rl.GetScreenToWorld2D(rl.NewVector2(0, 200), currentCamera)
	textPos := rl.GetScreenToWorld2D(rl.NewVector2(165, 250), currentCamera)
	rl.DrawRectangleV(boxPos, rl.NewVector2(40*relativeSize, 10*relativeSize), endColor)
	rl.DrawTextEx(rl.GetFontDefault(), endString, textPos, 1.5*relativeSize, relativeSpacing, defaultColor)

	sc.DrawButton(&sc.MenuButton)
	sc.DrawButton(&sc.QuitButton)

}

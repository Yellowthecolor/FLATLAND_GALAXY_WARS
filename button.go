package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Button struct {
	currentCamera rl.Camera2D
	Pos           rl.Vector2
	Size          rl.Vector2
	text          string
	textPos       rl.Vector2
	textSize      float32
	textSpacing   float32
	CustomColors

	Collider  rl.Rectangle
	isHovered bool
}

func (sc *Scene) InitButtons() {
	sc.MenuButton.text = "Main Menu"
	sc.MenuButton.Pos = rl.NewVector2(200, 325)

	sc.QuitButton.text = " Quit Game"
	sc.QuitButton.Pos = rl.NewVector2(550, 325)

	sc.GameButton.text = "   Start"
	sc.GameButton.Pos = rl.NewVector2(200, 325)

}

func (sc *Scene) DrawButton(b *Button) {
	currentCamera := sc.camera
	defaultColor := rl.Black

	relativePos := rl.GetScreenToWorld2D(b.Pos, currentCamera)
	relativeSize := 30 / currentCamera.Zoom
	b.Size = rl.NewVector2(relativeSize*10, relativeSize*5)
	b.Collider = rl.NewRectangle(relativePos.X, relativePos.Y, b.Size.X, b.Size.Y)

	relativeSpacing := 2 / currentCamera.Zoom
	b.textSize = relativeSize * 1.8
	relativePos.X = b.Pos.X + 26
	relativePos.Y = b.Pos.Y + 50
	b.textPos = rl.GetScreenToWorld2D(relativePos, currentCamera)

	b.primaryTint = rl.LightGray
	b.selectionTint = SelectionColors(rl.White)

	rl.DrawRectangleRec(b.Collider, b.primaryTint)
	if b.isHovered {
		rl.DrawRectangleLinesEx(b.Collider, b.Size.X*.05, b.selectionTint[b.currentSelectionTint])
		rl.DrawTextEx(rl.GetFontDefault(), b.text, b.textPos, b.textSize, relativeSpacing, rl.RayWhite)
	} else {
		rl.DrawTextEx(rl.GetFontDefault(), b.text, b.textPos, b.textSize, relativeSpacing, defaultColor)
	}
	b.CycleColorsInTime(2)
}

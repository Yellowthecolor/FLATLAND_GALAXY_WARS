package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Selector struct {
	selectorRectangle rl.Rectangle
	recStartPoint     rl.Vector2
}

func (sc *Scene) DrawSelector() {
	rl.DrawRectangleRec(sc.selectorRectangle, rl.NewColor(rl.LightGray.R, rl.LightGray.G, rl.LightGray.B, 50))
	rl.DrawRectangleLinesEx(sc.selectorRectangle, 3, rl.LightGray)
}

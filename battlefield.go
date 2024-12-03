package main

import rl "github.com/gen2brain/raylib-go/raylib"

func DrawBattleField() {
	rl.PushMatrix()
	rl.Translatef(0, 25*50, 0)
	rl.Rotatef(90, 1, 0, 0)
	rl.DrawGrid(150, 100)
	rl.PopMatrix()
}

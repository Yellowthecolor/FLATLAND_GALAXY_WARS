package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(1080, 720, "FLATLAND GALAXY WARS")
	rl.InitAudioDevice()
	defer rl.CloseWindow()
	defer rl.CloseAudioDevice()
	rl.SetExitKey(rl.KeyNull)
	rl.SetTargetFPS(120)

	var sc Scene

	sc.StartGame()

	for !sc.quitGame {
		rl.BeginDrawing()
		rl.BeginMode2D(sc.camera)
		rl.ClearBackground(rl.Black)
		sc.UpdateGameState()
		rl.EndMode2D()
		sc.DrawUI()
		rl.EndDrawing()
	}
}

package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UIItems struct {
	defaultColor rl.Color

	currentCredits int
	creditPos      rl.Vector2
	creditCount    string
	creditString   string
	creditColor    rl.Color

	currentCreatures int
	creaturePos      rl.Vector2
	creatureCount    string
	creatureString   string

	currentDefeats int
	defeatedPos    rl.Vector2
	defeatedCount  string
	defeatedString string

	fpsPos   rl.Vector2
	fpsColor rl.Color

	showSiegeTimer   bool
	opacityCounter   uint8
	siegePos         rl.Vector2
	siegeTimerValue  string
	siegeString      string
	siegeColor       rl.Color
	displayTextCount float32
}

func (sc *Scene) NewUI() {
	ui := UIItems{}

	ui.defaultColor = rl.White

	ui.creditPos = rl.NewVector2(10, 0)
	ui.currentCredits = sc.playerSet.teamBase.creatureCredits
	ui.creditCount = strconv.Itoa(ui.currentCredits)
	ui.creditString = "Creature Credits: " + ui.creditCount
	ui.creditColor = rl.Red

	ui.creaturePos = rl.NewVector2(350, 0)
	ui.currentCreatures = len(sc.playerSet.CreatureSet)
	ui.creatureCount = strconv.Itoa(ui.currentCreatures)
	ui.creatureString = "Active Creatures: " + ui.creatureCount

	ui.defeatedPos = rl.NewVector2(680, 0)
	ui.currentDefeats = sc.playerSet.teamBase.totalEnemiesDefeated
	ui.defeatedCount = strconv.Itoa(ui.currentDefeats)
	ui.defeatedString = "Defeated Enemies: " + ui.defeatedCount

	ui.fpsPos = rl.NewVector2(float32(rl.GetScreenWidth())-80, 5)
	ui.fpsColor = ui.defaultColor

	ui.siegePos = rl.NewVector2(420, 350)
	ui.siegeTimerValue = strconv.Itoa(sc.siegeTimer)
	ui.siegeString = "SIEGE STARTING IN " + ui.siegeTimerValue
	ui.siegeColor = rl.NewColor(rl.Gold.R, rl.Gold.G, rl.Gold.B, 255)
	ui.displayTextCount = 0

	ui.opacityCounter = 255
	sc.UI = &ui
}

func (sc *Scene) DrawUI() {
	if sc.currentScene == MENU {
		return
	}

	// Draw scaling text based on camera zoom

	if sc.UI.currentCredits != sc.playerSet.teamBase.creatureCredits {
		sc.UI.currentCredits = sc.playerSet.teamBase.creatureCredits
		sc.UI.creditCount = strconv.Itoa(sc.UI.currentCredits)
		sc.UI.creditString = "Creature Credits: " + sc.UI.creditCount

		if sc.UI.currentCredits < 1 {
			sc.UI.creditColor = rl.Red
		} else if sc.UI.currentCredits < 5 {
			sc.UI.creditColor = rl.Orange
		} else {
			sc.UI.creditColor = rl.Green
		}
	}

	if sc.UI.currentCreatures != len(sc.playerSet.CreatureSet) {
		sc.UI.currentCreatures = len(sc.playerSet.CreatureSet)
		sc.UI.creatureCount = strconv.Itoa(sc.UI.currentCreatures)
		sc.UI.creatureString = "Active Creatures: " + sc.UI.creatureCount
	}

	if sc.UI.currentDefeats != sc.playerSet.teamBase.totalEnemiesDefeated {
		sc.UI.currentDefeats = sc.playerSet.teamBase.totalEnemiesDefeated
		sc.UI.defeatedCount = strconv.Itoa(sc.UI.currentDefeats)
		sc.UI.defeatedString = "Defeated Enemies: " + sc.UI.defeatedCount
	}

	if rl.GetFPS() < 15 {
		sc.UI.fpsColor = rl.Red
	} else if rl.GetFPS() < 30 {
		sc.UI.fpsColor = rl.Orange
	} else if rl.GetFPS() < 60 {
		sc.UI.fpsColor = rl.Yellow
	} else {
		sc.UI.fpsColor = rl.Green
	}

	if sc.siegeTimer == 120 || sc.siegeTimer == 60 || sc.siegeTimer == 30 || sc.siegeTimer == 10 || sc.siegeTimer == 0 {
		sc.UI.showSiegeTimer = true
		PlaySoundStandAlone(&sc.warningSound)
	}

	if sc.UI.showSiegeTimer {
		if sc.UI.siegeColor.A > 0 && sc.UI.displayTextCount/120 >= .75 && sc.siegeTimer > 10 {
			sc.UI.opacityCounter -= 1
			sc.UI.siegeColor.A = sc.UI.opacityCounter
		}
		if sc.siegeTimer > 0 {
			sc.UI.siegePos = rl.NewVector2(350, 350)
			sc.UI.siegeTimerValue = strconv.Itoa(sc.siegeTimer)
			sc.UI.siegeString = "SIEGE STARTING IN " + sc.UI.siegeTimerValue
		} else {
			sc.UI.siegePos = rl.NewVector2(350, 350)
			sc.UI.siegeString = "SIEGE IS STARTING"
		}
		rl.DrawTextEx(rl.GetFontDefault(), sc.UI.siegeString, sc.UI.siegePos, 40, 1, sc.UI.siegeColor)
	}
	sc.UI.displayTextCount++
	if sc.UI.displayTextCount/120 >= 5 {
		sc.UI.showSiegeTimer = false
		sc.UI.opacityCounter = 255
		sc.UI.siegeColor.A = sc.UI.opacityCounter
		sc.UI.displayTextCount = 0
	}

	rl.DrawTextEx(rl.GetFontDefault(), sc.UI.creditString, sc.UI.creditPos, 20, 1, sc.UI.creditColor)
	rl.DrawTextEx(rl.GetFontDefault(), sc.UI.creatureString, sc.UI.creaturePos, 20, 1, sc.UI.defaultColor)
	rl.DrawTextEx(rl.GetFontDefault(), sc.UI.defeatedString, sc.UI.defeatedPos, 20, 1, sc.UI.defaultColor)
	rl.DrawTextEx(rl.GetFontDefault(), "FPS: "+strconv.Itoa(int(rl.GetFPS())), sc.UI.fpsPos, 20, 1, sc.UI.fpsColor)

}

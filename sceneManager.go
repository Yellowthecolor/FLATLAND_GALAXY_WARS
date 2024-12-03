package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	MENU   = "MenuScene"
	GAME   = "GameScene"
	PLAYER = "PLAYER"
	ENEMY  = "ENEMY"
)

type Scene struct {
	sceneList    []string
	currentScene string
	Menu

	camera rl.Camera2D
	UI     *UIItems
	Selector

	playerSet *CreatureSet
	enemySet  *CreatureSet

	CreditSlice       []*Credit
	maxCreditsSpawned int
	siegeTimer        int
	siegeCounter      int
	enableSiege       bool

	globalClockCounter float32
	playerIdleCount    int
	enemyIdleCount     int

	AutoPlay          bool
	toggleBattleField bool
	pauseWindowOpen   bool
	endScreenOpen     bool
	quitGame          bool

	Audio

	blahhhhh float32
}

func (sc *Scene) StartGame() {
	sc.InitializeGameState()
	sc.currentScene = MENU
	sc.PlayMusic(sc.menuMusic)

}

func (sc *Scene) InitializeGameState() {
	sc.sceneList = []string{MENU, GAME}
	sc.LoadAudio()
	sc.musicChanged = false

	sc.InitButtons()
	sc.siegeTimer = 120

	sc.CreditSlice = make([]*Credit, 0)
	sc.InitializeAllCreatures()
	sc.camera.Zoom = .5
	sc.camera.Target = rl.NewVector2(sc.playerSet.teamBase.Pos.X-200, sc.playerSet.teamBase.Pos.Y-150)
	sc.NewUI()
}

func (sc *Scene) UpdateGameState() {
	rl.UpdateMusicStream(sc.currentMusic.music)
	switch sc.currentScene {
	case MENU:
		{
			sc.MenuInputHandler()
			sc.DrawMenu()
		}
	case GAME:
		{
			sc.InputHandler()
			sc.UpdateGameScene()
			sc.DrawGameScene()
		}
	}
}

func (sc *Scene) UpdateGameScene() {
	if sc.endScreenOpen {
		return
	}

	sc.UpdateGlobalClock()
	sc.UpdateCredits()
	sc.UpadteHealthStatus(sc.playerSet, sc.enemySet)
	sc.UpadteHealthStatus(sc.enemySet, sc.playerSet)
	sc.UpdateCombat(sc.playerSet)
	sc.UpdateCombat(sc.enemySet)
	sc.playerSet.UpdateMovement()
	sc.enemySet.UpdateMovement()

	sc.playerSet.UpdateAIs(sc)
	sc.enemySet.UpdateAIs(sc)

	if sc.playerSet.teamBase.currentHealth <= 0 {
		sc.enemySet.hasWon = true
		sc.endScreenOpen = true
	} else if sc.enemySet.teamBase.currentHealth <= 0 {
		sc.playerSet.hasWon = true
		sc.endScreenOpen = true
	}
}

func (sc *Scene) DrawGameScene() {
	if sc.toggleBattleField {
		DrawBattleField()
	}

	sc.DrawAllCreatures()
	sc.DrawCredit()
	sc.DrawSelector()

	if sc.endScreenOpen {
		sc.DrawEndScreen()
	} else if sc.pauseWindowOpen {
		sc.DrawPauseMenu()
	}

}

func (sc *Scene) UpdateGlobalClock() {
	if sc.playerSet.hasWon || sc.enemySet.hasWon {
		sc.endScreenOpen = true
	}

	sc.siegeCounter++
	if sc.siegeCounter/120 >= 1 {
		sc.siegeTimer--
		if sc.siegeTimer <= 0 {
			sc.enableSiege = true
		}
		sc.siegeCounter = 0
	}

	sc.globalClockCounter++
	if sc.globalClockCounter/120 >= 5 {
		sc.maxCreditsSpawned += 5
		if sc.maxCreditsSpawned > 950 {
			sc.maxCreditsSpawned = 950
		}
		pcs := *sc.playerSet
		ecs := *sc.enemySet

		// Passive Regen and Credits for both teams
		pcs.teamBase.creatureCredits++
		pcs.teamBase.currentHealth += pcs.teamBase.healthRegen
		if pcs.teamBase.currentHealth > pcs.teamBase.maxHealth {
			pcs.teamBase.currentHealth = pcs.teamBase.maxHealth
		}

		ecs.teamBase.creatureCredits++
		ecs.teamBase.currentHealth += ecs.teamBase.healthRegen
		if ecs.teamBase.currentHealth > ecs.teamBase.maxHealth {
			ecs.teamBase.currentHealth = ecs.teamBase.maxHealth
		}

		// Increase enemy count overtime
		sc.enemyIdleCount++
		if ecs.teamBase.creatureCredits >= 5 {
			if sc.enableSiege {
				for ecs.teamBase.creatureCredits > 5 {
					ecs.InitializeNewCreature(ENEMY)
					ecs.teamBase.creatureCredits -= 5
				}
			} else if sc.enemyIdleCount >= 2 {
				for i := 0; i < sc.enemyIdleCount; i++ {
					ecs.InitializeNewCreature(ENEMY)
					ecs.teamBase.creatureCredits -= 5
					if ecs.teamBase.creatureCredits < 5 {
						break
					}
				}
				sc.enemyIdleCount = 0
			}
		}

		// Increase player count automatically if autoplay is on
		if sc.AutoPlay {
			sc.playerIdleCount++
			if sc.playerIdleCount >= 2 && pcs.teamBase.creatureCredits >= 5 {
				for i := 0; i < sc.playerIdleCount; i++ {
					pcs.InitializeNewCreature(PLAYER)
					PlaySoundStandAlone(&sc.spawnSound)
					pcs.teamBase.creatureCredits -= 5
					if pcs.teamBase.creatureCredits < 5 {
						break
					}
				}
				sc.playerIdleCount = 0
			}
		}

		sc.globalClockCounter = 0
	}
}

func (sc *Scene) ResetGameScene() {
	sc.camera.Offset = rl.Vector2Zero()
	sc.playerSet = nil
	sc.enemySet = nil
	sc.CreditSlice = nil
	sc.maxCreditsSpawned = 0
	sc.globalClockCounter = 0
	sc.playerIdleCount = 0
	sc.enemyIdleCount = 0
	sc.enableSiege = false
	sc.AutoPlay = false
	sc.pauseWindowOpen = false
	sc.endScreenOpen = false

}

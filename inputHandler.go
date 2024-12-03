package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (sc *Scene) InputHandler() {
	relativePosition := rl.GetScreenToWorld2D(rl.GetMousePosition(), sc.camera)

	if rl.IsKeyDown(rl.KeyUp) {
		sc.blahhhhh += 1
		fmt.Println(sc.blahhhhh)
	}
	if rl.IsKeyDown(rl.KeyDown) {
		sc.blahhhhh -= 1
		fmt.Println(sc.blahhhhh)

	}

	// Open Exit Menu
	if rl.WindowShouldClose() || rl.IsKeyPressed(rl.KeyEscape) {
		sc.pauseWindowOpen = !sc.pauseWindowOpen
	}

	// Check QuitGame
	if sc.pauseWindowOpen || sc.endScreenOpen {
		if rl.CheckCollisionPointRec(relativePosition, sc.MenuButton.Collider) {
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				sc.PlayMusic(sc.menuMusic)
				sc.currentScene = MENU
			}
			sc.MenuButton.isHovered = true
		} else {
			sc.MenuButton.isHovered = false
		}

		if rl.CheckCollisionPointRec(relativePosition, sc.QuitButton.Collider) {
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				sc.quitGame = true
			}
			sc.QuitButton.isHovered = true
		} else {
			sc.QuitButton.isHovered = false
		}
		return
	}

	pcs := *sc.playerSet
	ecs := *sc.enemySet
	// Zoom in/out
	wheel := rl.GetMouseWheelMove()
	if wheel != 0 || (rl.IsKeyDown(rl.KeyMinus) || rl.IsKeyDown(rl.KeyEqual)) {
		sc.camera.Offset = rl.GetMousePosition()
		sc.camera.Target = relativePosition

		if wheel == 0 {
			if rl.IsKeyDown(rl.KeyMinus) {
				wheel -= .1
			} else {
				wheel += .1
			}
		}
		scaleFactor := 1 + (.25 * math.Abs(float64(wheel)))
		if wheel < 0 {
			scaleFactor = 1 / scaleFactor
		}

		sc.camera.Zoom = rl.Clamp(sc.camera.Zoom*float32(scaleFactor), .0625, 8)

	}

	// drag camera around
	if rl.IsMouseButtonDown(rl.MouseButtonMiddle) || (rl.IsKeyDown(rl.KeyLeftShift) && rl.IsMouseButtonDown(rl.MouseButtonLeft)) {
		delta := rl.GetMouseDelta()
		delta = rl.Vector2Scale(delta, -1/sc.camera.Zoom)
		sc.camera.Target = rl.Vector2Add(sc.camera.Target, delta)
	}

	// Select creatures
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) && !rl.IsKeyDown(rl.KeyLeftShift) { // Select by clicking
		selectedCount := len(pcs.selectedCreatures) // Counter for # of creatures selected
		for _, c := range pcs.CreatureSet {
			if rl.CheckCollisionPointRec(relativePosition, c.Collider) {
				// c.selected = !c.selected <- this doesn't work because of share functionality with rectangle
				if !c.selected {
					c.selected = true
					c.isAI = false

					if c.targetCreature != nil {
						c.targetCreature.enemySelected = true
					}
				} else {
					c.selected = false
					c.isAI = true
				}

			}
			if c.selected {
				pcs.selectedCreatures[c.ID] = c
			} else {
				_, exists := pcs.selectedCreatures[c.ID]
				if exists {
					delete(pcs.selectedCreatures, c.ID)
				}
			}
		}

		// if player clicked anywhere not a creature, clear all selections
		if len(pcs.selectedCreatures) == selectedCount {
			for _, cSelected := range pcs.selectedCreatures {
				cSelected.selected = false
				cSelected.isAI = true
				if cSelected.targetCreature != nil {
					cSelected.targetCreature.enemySelected = false
				}
			}
			clear(pcs.selectedCreatures)
		}
	} else if rl.IsMouseButtonDown(rl.MouseButtonLeft) && !rl.IsKeyDown(rl.KeyLeftShift) { // Select with box
		if sc.selectorRectangle.X == 0 && sc.selectorRectangle.Y == 0 {
			sc.recStartPoint = relativePosition
			sc.selectorRectangle = rl.NewRectangle(relativePosition.X, relativePosition.Y, 0, 0)
		}

		// box direction determined by mouse drag position
		if relativePosition.X < sc.recStartPoint.X {
			sc.selectorRectangle.X = relativePosition.X
			sc.selectorRectangle.Width = sc.recStartPoint.X - sc.selectorRectangle.X
		} else {
			sc.selectorRectangle.Width = relativePosition.X - sc.selectorRectangle.X
		}

		if relativePosition.Y < sc.recStartPoint.Y {
			sc.selectorRectangle.Y = relativePosition.Y
			sc.selectorRectangle.Height = sc.recStartPoint.Y - sc.selectorRectangle.Y
		} else {
			sc.selectorRectangle.Height = relativePosition.Y - sc.selectorRectangle.Y
		}

		// Select ally creatures
		for _, c := range pcs.CreatureSet {
			if rl.CheckCollisionRecs(sc.selectorRectangle, c.Collider) {
				if !rl.Vector2Equals(relativePosition, sc.recStartPoint) {
					c.selected = true
					c.isAI = false
				}
			}

			if c.selected {
				pcs.selectedCreatures[c.ID] = c
			} else {
				_, exists := pcs.selectedCreatures[c.ID]
				if exists {
					delete(pcs.selectedCreatures, c.ID)
				}
			}
		}
	} else {
		sc.selectorRectangle = rl.NewRectangle(0, 0, 0, 0)
	}

	// Move selected creatures
	if rl.IsMouseButtonPressed(rl.MouseButtonRight) {

		if len(pcs.selectedCreatures) > 0 {

			isAttacking := false
			targetCreatureID := 0
			// Select enemy creatures to attack
			for _, ec := range ecs.CreatureSet {
				if rl.CheckCollisionPointRec(relativePosition, ec.Collider) {
					ec.enemySelected = true
					targetCreatureID = ec.ID
				} else {
					ec.enemySelected = false
				}
			}

			// if player clicked anywhere not a creature, clear all selections
			for _, cSelected := range pcs.selectedCreatures {
				if targetCreatureID != 0 {
					cSelected.targetCreature = ecs.CreatureSet[targetCreatureID]
					cSelected.inCombat = true
					isAttacking = true
					_, exists := pcs.inCombatCreatures[cSelected.ID]
					if !exists {
						pcs.inCombatCreatures[cSelected.ID] = cSelected
					}

				} else {
					pcs.ResetCombatCreature(cSelected)
				}
			}

			if isAttacking { // Attack movement
				for _, cSelected := range pcs.selectedCreatures {
					cSelected.targetDestination = cSelected.targetCreature.Center
					cSelected.CalculateStoredMovement(cSelected.AttackRange)
				}
			} else { // Regular Move
				sumDst := rl.Vector2Zero()
				sumCntrs := rl.Vector2Zero()
				for _, cSelected := range pcs.selectedCreatures {
					sumCntrs = rl.Vector2Add(sumCntrs, cSelected.Center)
				}
				averageCntr := rl.Vector2Zero()
				averageCntr.X = sumCntrs.X / float32(len(pcs.selectedCreatures))
				averageCntr.Y = sumCntrs.Y / float32(len(pcs.selectedCreatures))

				for _, cSelected := range pcs.selectedCreatures {
					cSelected.startLocation = cSelected.Center
					cSelected.targetDestination = relativePosition
					subV := rl.Vector2Subtract(cSelected.targetDestination, cSelected.Center)
					sumDst = rl.Vector2Add(sumDst, subV)
					cSelected.targetDistance = rl.Vector2Distance(averageCntr, cSelected.targetDestination)
				}

				for _, cSelected := range pcs.selectedCreatures {
					cSelected.storedMovement.X = sumDst.X / float32(len(pcs.selectedCreatures))
					cSelected.storedMovement.Y = sumDst.Y / float32(len(pcs.selectedCreatures))
				}
			}
		}
	}

	// Spawn Player Units
	if rl.IsKeyPressed(rl.KeySpace) {
		if pcs.teamBase.creatureCredits >= 5 {
			pcs.InitializeNewCreature(PLAYER)
			pcs.teamBase.creatureCredits -= 5
		}
	}

	// Enable AutoPlay
	if rl.IsKeyPressed(rl.KeyA) {
		sc.AutoPlay = !sc.AutoPlay
	}

	// Toggle Battlefield
	if rl.IsKeyPressed(rl.KeyB) {
		sc.toggleBattleField = !sc.toggleBattleField
	}
}

func (sc *Scene) MenuInputHandler() {
	relativePosition := rl.GetScreenToWorld2D(rl.GetMousePosition(), sc.camera)

	if rl.CheckCollisionPointRec(relativePosition, sc.GameButton.Collider) {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			sc.ResetGameScene()
			sc.InitializeGameState()
			sc.PlayMusic(sc.gameMusic)
			sc.currentScene = GAME
		}
		sc.GameButton.isHovered = true
	} else {
		sc.GameButton.isHovered = false
	}

	if rl.CheckCollisionPointRec(relativePosition, sc.QuitButton.Collider) {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			sc.quitGame = true
		}
		sc.QuitButton.isHovered = true
	} else {
		sc.QuitButton.isHovered = false
	}
}

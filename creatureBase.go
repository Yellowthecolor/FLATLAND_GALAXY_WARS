package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Base struct {
	Creature
	healthRegen     float32
	creatureCredits int

	totalEnemiesDefeated int
}

func (cs *CreatureSet) NewBase(team string, scene *Scene) {
	var (
		// Changing Values per team
		startPosX      float32
		startPosY      float32
		primaryTint    rl.Color
		selectionTint  rl.Color
		projectileTint rl.Color
		dangerRange    float32 // Not relevant for player team, just usess collider

		// Default Values all bases
		size        rl.Vector2 = rl.NewVector2(800, 800)
		speed       float32    = 0
		health      float32    = 1000
		damage      float32    = health * .1
		attackRange float32    = size.X * 2
		sightRange  float32    = attackRange
		attackSpeed float32    = 5
		healthRegen float32    = health * .1
	)

	// startPosX = float32(rl.GetRandomValue(-7000, -6500))
	// startPosY = float32(rl.GetRandomValue(-6000, -5500))
	if team == PLAYER {
		startPosX = float32(rl.GetRandomValue(-7000, -5500))
		startPosY = float32(rl.GetRandomValue(-6000, 0))
		primaryTint = rl.Blue
		selectionTint = rl.Green
		projectileTint = rl.Lime
		dangerRange = size.X
	} else if team == ENEMY {
		startPosX = float32(rl.GetRandomValue(4500, 6000))
		startPosY = float32(rl.GetRandomValue(-6000, 0))
		primaryTint = rl.Gold
		selectionTint = rl.Red
		projectileTint = rl.Maroon
		dangerRange = size.X / 2
	}

	newBase := Base{}
	startPos := rl.NewVector2(startPosX, startPosY)
	newBase.Creature = NewCreature(startPos, size, speed, health, primaryTint, selectionTint, projectileTint)
	newBase.Team = team
	newBase.NewCombatStats(damage, attackRange, attackSpeed, dangerRange, sightRange)
	newBase.combatFrameCounter = 120 * attackSpeed
	newBase.healthRegen = healthRegen

	for {
		randomID := int(rl.GetRandomValue(1, 5))
		_, exists := cs.CreatureSet[randomID]
		if !exists {
			newBase.ID = int(randomID)
			cs.CreatureSet[int(randomID)] = &newBase.Creature
			break
		}
	}
	newBase.isAI = true
	newBase.isBase = true
	cs.teamBase = &newBase
}

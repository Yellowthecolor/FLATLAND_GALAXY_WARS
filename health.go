package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Health struct {
	currentHealth float32
	maxHealth     float32
}

func NewHealth(health float32) Health {
	newHealth := Health{
		currentHealth: health,
		maxHealth:     health,
	}
	return newHealth
}

func (c *Creature) DrawHealthBar() {

	var (
		barWidth    = c.Size.X
		barHeight   = c.Size.Y / 4
		offset      = float32(3)
		healthRatio = c.currentHealth / c.maxHealth
		healthWidth = healthRatio * barWidth
	)

	if c.isBase {
		barHeight = c.Size.Y / 10
		offset = float32(10)
	}

	redBar := rl.NewRectangle(c.Pos.X, c.Pos.Y+c.Size.Y+offset, barWidth, barHeight)
	greenBar := rl.NewRectangle(c.Pos.X, c.Pos.Y+c.Size.Y+offset, healthWidth, barHeight)

	rl.DrawRectangleRec(redBar, rl.Red)
	rl.DrawRectangleRec(greenBar, rl.Lime)
	rl.DrawRectangleLinesEx(redBar, 1, rl.Black)

}

func (sc *Scene) UpadteHealthStatus(cs, opponentCS *CreatureSet) {
	for _, c := range cs.CreatureSet {
		if c.currentHealth <= 0 {

			opponentCS.teamBase.totalEnemiesDefeated++
			opponentCS.teamBase.creatureCredits += 2
			if c.targetCreature != nil {
				c.targetCreature.enemySelected = false
			}
			DeleteFromSet(c, cs.inCombatCreatures)
			DeleteFromSet(c, cs.selectedCreatures)
			DeleteFromSet(c, cs.CreatureSet)
		}
	}
}

func DeleteFromSet(c *Creature, set map[int]*Creature) {
	_, exists := set[c.ID]
	if exists {
		delete(set, c.ID)
	}
}

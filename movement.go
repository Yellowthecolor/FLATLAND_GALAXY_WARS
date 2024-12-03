package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Movement struct {
	Speed             float32
	storedMovement    rl.Vector2
	startLocation     rl.Vector2
	targetDestination rl.Vector2
	targetDistance    float32
	distanceTraveled  float32
}

func (c *Creature) CalculateStoredMovement(buffer float32) {
	if c.inAttackRange {
		return
	}
	c.startLocation = c.Center
	subV := rl.Vector2Subtract(c.targetDestination, c.Center)
	c.targetDistance = rl.Vector2Distance(c.Center, c.targetDestination) - buffer
	c.storedMovement.X = subV.X
	c.storedMovement.Y = subV.Y
}

func (cs *CreatureSet) UpdateMovement() {
	for _, c := range cs.CreatureSet {
		c.distanceTraveled = rl.Vector2Distance(c.startLocation, c.Center)

		if c.inCombat {
			if c.Team == ENEMY && rl.CheckCollisionCircleRec(c.Center, c.AttackRange, c.targetCreature.Collider) {
				c.inAttackRange = true
			} else if c.Team == PLAYER && rl.CheckCollisionCircles(c.Center, c.AttackRange, c.targetCreature.Center, c.targetCreature.DangerRange) {
				c.inAttackRange = true
			} else {
				c.targetDestination = c.targetCreature.Center
				c.inAttackRange = false
			}
		}

		if c.inAttackRange || c.distanceTraveled >= c.targetDistance {
			c.distanceTraveled = 0
			c.targetDistance = 0
			c.storedMovement = rl.Vector2Zero()
			c.targetDestination = rl.Vector2Zero()
		} else {
			c.CreatureMove()
		}

	}
}

func (c *Creature) CreatureMove() {
	c.Pos = rl.Vector2Add(c.Pos, rl.Vector2Scale(rl.Vector2Normalize(c.storedMovement), c.Speed*rl.GetFrameTime()))
	c.Center = rl.NewVector2(c.Pos.X+c.Size.X/2, c.Pos.Y+c.Size.Y/2)
	c.Collider.X = c.Pos.X
	c.Collider.Y = c.Pos.Y
}

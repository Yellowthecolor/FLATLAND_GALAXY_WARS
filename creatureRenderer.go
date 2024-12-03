package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type CreatureRenderer struct {
	Pos      rl.Vector2
	Size     rl.Vector2
	Center   rl.Vector2
	Collider rl.Rectangle
}

func (sc *Scene) DrawAllCreatures() {
	sc.playerSet.teamBase.DrawCreature()
	sc.enemySet.teamBase.DrawCreature()
	for _, c := range sc.playerSet.CreatureSet {
		if c.isBase {
			continue
		}
		c.DrawCreature()
	}
	for _, c := range sc.enemySet.CreatureSet {
		if c.isBase {
			continue
		}
		c.DrawCreature()
	}
}

func (c *Creature) DrawCreature() {
	if c.Team == PLAYER {
		rl.DrawRectangleRec(c.Collider, c.primaryTint)
		rl.DrawRectangleLinesEx(c.Collider, 2, rl.Black)
		if c.selected {
			//attack range circle
			rl.DrawCircleLinesV(c.Center, c.AttackRange, rl.White)
			rl.DrawRectangleLinesEx(c.Collider, c.Size.X*.1, c.selectionTint[c.currentSelectionTint])
		}

	} else {
		if c.enemySelected {
			rl.DrawCircleLinesV(c.Center, c.AttackRange, rl.Red)
			rl.DrawCircleV(c.Center, c.Size.X/2, c.selectionTint[c.currentSelectionTint])
		} else {
			rl.DrawCircleV(c.Center, c.Size.X/2, c.primaryTint)
			rl.DrawCircleLinesV(c.Center, c.Size.X/2, rl.Black)
		}
		rl.DrawCircleV(c.Center, c.Size.X/2-(c.Size.X*.1), c.primaryTint)
	}
	c.CycleColorsInTime(2)
	c.DrawProjectiles()
	c.DrawHealthBar()

}

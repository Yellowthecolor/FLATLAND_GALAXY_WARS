package main

import rl "github.com/gen2brain/raylib-go/raylib"

// Used in scouting
func (ai *Creature) FindNearestCredit(sc *Scene) (*Credit, bool) {
	var nearestCredit *Credit = nil
	minDist := ai.SightRange

	for _, cr := range sc.CreditSlice {
		dist := rl.Vector2Distance(ai.Center, cr.Center)
		if dist > ai.SightRange {
			continue
		}
		if dist <= minDist {
			minDist = dist
			nearestCredit = cr
		}
	}

	return nearestCredit, nearestCredit != nil
}

// Used in scouting
func (ai *Creature) FindNearestEnemy(sc *Scene) (*Creature, bool) {
	var nearestEnemy *Creature = nil
	minDist := ai.SightRange

	var cs *CreatureSet
	if ai.Team == PLAYER {
		cs = sc.enemySet
	} else if ai.Team == ENEMY {
		cs = sc.playerSet
	}

	for _, ec := range cs.CreatureSet {
		dist := rl.Vector2Distance(ai.Center, ec.Center)
		if ec.isBase && dist < ai.SightRange*2 {
			nearestEnemy = ec
			break
		}
		if dist > ai.SightRange {
			continue
		}
		if dist < minDist {
			minDist = dist
			nearestEnemy = ec
		}
	}

	return nearestEnemy, nearestEnemy != nil
}

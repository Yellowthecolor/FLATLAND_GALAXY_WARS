package main

import rl "github.com/gen2brain/raylib-go/raylib"

type AIState int

const (
	Idle      = 0
	Scouting  = 1
	Gathering = 2
	Fighting  = 3
)

type CreatureAI struct {
	isAI         bool
	State        AIState
	ChangedState bool
	Timer        float32
	TickCount    int
}

func (cs *CreatureSet) UpdateAIs(sc *Scene) {
	for _, ai := range cs.CreatureSet {
		ai.Tick(sc)
	}
}

func (ai *Creature) SetState(newState AIState) {
	ai.ChangedState = true
	ai.State = newState
}

func (ai *Creature) Tick(sc *Scene) {
	if ai.ChangedState {
		ai.Timer = 0
		ai.TickCount = 0
		ai.ChangedState = false
	}
	switch ai.State {
	case Idle:
		ai.TickIdle(sc)
	case Scouting:
		ai.TickScouting(sc)
	case Gathering:
		ai.TickGather(sc)
	case Fighting:
		ai.TickFighting(sc)
	}
	ai.Timer += rl.GetFrameTime()
	if ai.isAI {
		ai.TickCount += 1
	}

}

func (ai *Creature) TickIdle(sc *Scene) {
	if ai.TickCount == 0 {
	}

	// idle for 3 seconds before doing anything
	if ai.isAI && ai.Timer < 3 {
		return
	}
	if !ai.isAI {
		return
	}
	ai.SetState(Scouting)

}

func (ai *Creature) TickScouting(sc *Scene) {
	if ai.selected {
		ai.SetState(Idle)
		return
	}

	if ai.inAttackRange {
		ai.SetState(Fighting)
	}
	if ai.TickCount == 0 {
		var randPos rl.Vector2
		if sc.enableSiege && ai.Team == ENEMY {
			playerBase := sc.playerSet.teamBase
			siegeDestMinX := playerBase.Pos.X + playerBase.Size.X
			siegeDestMaxX := siegeDestMinX + ai.AttackRange
			siegeDestMinY := playerBase.Pos.Y
			siegeDestMaxY := playerBase.Pos.Y + playerBase.Size.Y
			randPosX := float32(rl.GetRandomValue(int32(siegeDestMinX), int32(siegeDestMaxX)))
			randPosY := float32(rl.GetRandomValue(int32(siegeDestMinY), int32(siegeDestMaxY)))
			randPos = rl.NewVector2(randPosX, randPosY)
			ai.targetCreature = &playerBase.Creature
		} else {
			randPosX := float32(rl.GetRandomValue(-7400, 7400))
			randPosY := float32(rl.GetRandomValue(-6200, 1100))
			randPos = rl.NewVector2(randPosX, randPosY)
		}
		ai.targetDestination = randPos

	}
	var cs *CreatureSet
	if ai.Team == PLAYER {
		cs = sc.enemySet
	} else if ai.Team == ENEMY {
		cs = sc.playerSet
	}
	// Find a credit or an enemy, prioritize enemy
	if ec, found := ai.FindNearestEnemy(sc); found {
		ai.targetCreature = ec
		ai.targetCreature.enemySelected = true
		ai.targetDestination = ai.targetCreature.Center
		ai.inCombat = true
		_, exists := cs.inCombatCreatures[ai.ID]
		if !exists {
			cs.inCombatCreatures[ai.ID] = ai
		}
		ai.SetState(Fighting)

	} else if cr, found := ai.FindNearestCredit(sc); found && !ai.isBase {
		if sc.enableSiege && ai.Team == ENEMY {
			ai.targetCreature = &sc.playerSet.teamBase.Creature
		} else {
			ai.targetDestination = cr.Center
			ai.SetState(Gathering)
		}
	}

	ai.CalculateStoredMovement(0)
	if rl.Vector2Distance(ai.targetDestination, ai.Center) <= 5 {
		ai.SetState(Scouting)
		return
	}
}

func (ai *Creature) TickGather(sc *Scene) {
	if ai.selected {
		ai.SetState(Idle)
		return
	}
	dist := rl.Vector2Distance(ai.Center, ai.targetDestination)

	if dist <= 5 {
		ai.SetState(Scouting)
		return
	}
}

func (ai *Creature) TickFighting(sc *Scene) {
	if ai.selected && !ai.isBase {
		ai.SetState(Idle)
		return
	}

	var dist float32
	var cs *CreatureSet
	if ai.Team == PLAYER {
		cs = sc.enemySet
	} else if ai.Team == ENEMY {
		cs = sc.playerSet
	}
	if ai.TickCount == 0 {
		if ec, ok := ai.FindNearestEnemy(sc); ok {
			ai.targetCreature = ec
			ai.targetCreature.enemySelected = true
			ai.targetDestination = ai.targetCreature.Center
			ai.inCombat = true
			_, exists := cs.inCombatCreatures[ai.ID]
			if !exists {
				cs.inCombatCreatures[ai.ID] = ai
			}
		}
	}
	if ai.targetCreature != nil {
		dist = rl.Vector2Distance(ai.targetCreature.Center, ai.Center)
		ai.targetDestination = ai.targetCreature.Center
		ai.CalculateStoredMovement(0)
		if ai.targetCreature.currentHealth <= 0 {
			cs.ResetCombatCreature(ai)
			ai.SetState(Scouting)
			return
		}
	}

	if dist > ai.SightRange {
		ai.SetState(Scouting)
		return
	}
}

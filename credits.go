package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Credit struct {
	Pos         rl.Vector2
	Center      rl.Vector2
	Radius      float32
	Rotation    float32
	Value       float32
	isCollected bool
	Color       rl.Color
}

func NewCredit(pos rl.Vector2, value float32) *Credit {
	cr := Credit{
		Pos:      pos,
		Radius:   value * 3,
		Value:    value,
		Rotation: float32(rl.GetRandomValue(1, 360)),
		Color:    rl.Gold,
	}
	cr.Center = rl.NewVector2(pos.X+cr.Radius, pos.Y+cr.Radius)
	return &cr
}

func (sc *Scene) DrawCredit() {
	for _, cr := range sc.CreditSlice {
		rl.DrawPoly(cr.Center, 5, cr.Radius, cr.Rotation, cr.Color)
		rl.DrawPolyLinesEx(cr.Center, 5, cr.Radius, cr.Rotation, cr.Radius/5, rl.Orange)
		cr.Rotation++
	}
}

func (sc *Scene) CreditCollection(cr *Credit, cs *CreatureSet) bool {
	for _, c := range cs.CreatureSet {
		if rl.CheckCollisionCircleRec(cr.Center, cr.Radius, c.Collider) {
			cs.teamBase.creatureCredits += int(cr.Value)
			cr.isCollected = true
			if c.Team == PLAYER {
				if cs.teamBase.creatureCredits > 5 {
					PlaySoundStandAlone(&sc.fiveSound)
				} else {
					PlaySoundStandAlone(&sc.collectSound)
				}

			}
			return true
		}
	}
	return false
}

func (sc *Scene) UpdateCredits() {
	if len(sc.CreditSlice) < 50+sc.maxCreditsSpawned {
		randPosX := float32(rl.GetRandomValue(-7400, 7400))
		randPosY := float32(rl.GetRandomValue(-6200, 1100))
		randPos := rl.NewVector2(randPosX, randPosY)
		value := float32(rl.GetRandomValue(1, 5))
		sc.CreditSlice = append(sc.CreditSlice, NewCredit(randPos, value))
	}

	for i, cr := range sc.CreditSlice {
		if sc.CreditCollection(cr, sc.playerSet) || sc.CreditCollection(cr, sc.enemySet) {
			sc.DespawnCredit(i)
			break
		}
	}
}

func (sc *Scene) DespawnCredit(index int) {
	if len(sc.CreditSlice) < 1 {
		return
	}
	if index < len(sc.CreditSlice) {
		sc.CreditSlice = append(sc.CreditSlice[:index], sc.CreditSlice[index+1:]...)
	} else {
		sc.CreditSlice = sc.CreditSlice[:len(sc.CreditSlice)-1]
	}
}

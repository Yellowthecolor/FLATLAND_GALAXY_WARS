package main

import rl "github.com/gen2brain/raylib-go/raylib"

const MAX_CREATURES = 100

type CreatureSet struct {
	teamBase          *Base
	CreatureSet       map[int]*Creature
	selectedCreatures map[int]*Creature
	inCombatCreatures map[int]*Creature

	maxCreatures int

	hasWon bool
}

type Creature struct {
	ID int
	CreatureRenderer
	CustomColors

	Movement
	Health
	Combat

	selected      bool
	enemySelected bool
	isBase        bool

	CreatureAI
	Team string
}

func NewCreatureSet() *CreatureSet {
	return &CreatureSet{
		CreatureSet:       make(map[int]*Creature),
		selectedCreatures: make(map[int]*Creature),
		inCombatCreatures: make(map[int]*Creature),
		maxCreatures:      100,
	}
}

func NewCreature(pos, size rl.Vector2, speed, health float32, tint, selectionTint, projTint rl.Color) Creature {
	newCreature := Creature{
		CreatureRenderer: CreatureRenderer{
			Pos:      pos,
			Size:     size,
			Center:   rl.NewVector2(pos.X+size.X/2, pos.Y+size.Y/2),
			Collider: rl.NewRectangle(pos.X, pos.Y, size.X, size.Y),
		},
		CustomColors: CustomColors{tint, SelectionColors(selectionTint), projTint, true, 0, 0},
		Health:       NewHealth(health),
	}
	newCreature.Speed = speed
	return newCreature
}

func (cs *CreatureSet) InitializeNewCreature(team string) {
	if len(cs.CreatureSet) >= cs.maxCreatures {
		return
	}
	var (
		startPosX      float32
		startPosY      float32
		primaryTint    rl.Color
		selectionTint  rl.Color
		projectileTint rl.Color

		size        rl.Vector2 = rl.NewVector2(50, 50)
		speed       float32    = 500
		health      float32    = 100
		damage      float32    = health * .1
		attackRange float32    = size.X * 4
		sightRange  float32    = size.X * 7
		attackSpeed float32    = .5
		dangerRange float32    = size.X / 3
	)

	if team == PLAYER {
		startPosX = cs.teamBase.Pos.X + cs.teamBase.Size.X + float32(rl.GetRandomValue(100, 200))
		startPosY = cs.teamBase.Pos.Y + float32(rl.GetRandomValue(0, int32(cs.teamBase.Size.Y)))
		primaryTint = rl.DarkBlue
		selectionTint = rl.Green
		projectileTint = rl.Lime
	} else if team == ENEMY {
		startPosX = cs.teamBase.Pos.X - size.X - float32(rl.GetRandomValue(300, 400))
		startPosY = cs.teamBase.Pos.Y + float32(rl.GetRandomValue(0, int32(cs.teamBase.Size.Y)))
		primaryTint = rl.Pink
		selectionTint = rl.Red
		projectileTint = rl.Maroon
	}

	startPos := rl.NewVector2(startPosX, startPosY)
	newCreature := NewCreature(startPos, size, speed, health, primaryTint, selectionTint, projectileTint)
	newCreature.Team = team
	newCreature.NewCombatStats(damage, attackRange, attackSpeed, dangerRange, sightRange)
	newCreature.combatFrameCounter = 120 * attackSpeed
	newCreature.isAI = true
	for {
		randomID := int(rl.GetRandomValue(20, 1000))
		_, exists := cs.CreatureSet[randomID]
		if !exists {
			newCreature.ID = int(randomID)
			cs.CreatureSet[int(randomID)] = &newCreature
			break
		}
	}
}

func (sc *Scene) InitializeAllCreatures() {
	sc.playerSet = NewCreatureSet() // player creature set
	sc.enemySet = NewCreatureSet()  // enemy creature set

	sc.playerSet.NewBase(PLAYER, sc)
	sc.enemySet.NewBase(ENEMY, sc)
	sc.playerSet.InitializeNewCreature(PLAYER)
	sc.enemySet.InitializeNewCreature(ENEMY)
}

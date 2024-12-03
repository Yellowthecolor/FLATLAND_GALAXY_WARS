package main

type Combat struct {
	Damage         float32
	AttackRange    float32
	AttackSpeed    float32
	inCombat       bool
	inAttackRange  bool
	targetCreature *Creature

	DangerRange float32
	SightRange  float32

	combatFrameCounter float32

	projectileSlice []*Projectile
}

func (c *Creature) NewCombatStats(damage, attackRange, attackSpeed, dangerRange, sightRange float32) {
	c.Damage = damage
	c.AttackRange = attackRange
	c.AttackSpeed = attackSpeed
	c.DangerRange = dangerRange
	c.SightRange = sightRange
	c.projectileSlice = make([]*Projectile, 0)
}

func (sc *Scene) UpdateCombat(cs *CreatureSet) {
	for _, c := range cs.CreatureSet {

		if c.combatFrameCounter/120 < c.AttackSpeed {
			c.combatFrameCounter++
		}

		// DEBUGGING COMBAT PRINTS REMOVE LATER
		// if c.selected && c.combatFrameCounter/120 >= c.AttackSpeed {
		// 	fmt.Println(&c.targetCreature)
		// 	fmt.Println(c.inCombat)
		// 	fmt.Println(c.inAttackRange)
		// 	fmt.Println(len(c.projectileSlice))
		// }
	}

	for _, c := range cs.inCombatCreatures {
		if len(c.projectileSlice) > 0 {
			for i, p := range c.projectileSlice {
				sc.ProjectileCollisionCheck(c, cs, p)
				if !p.isActive {
					c.DespawnProjectile(i)
				}
				p.targetDestination = c.targetCreature.Center
			}
		}

		if c.targetCreature.currentHealth <= 0 {
			continue
		}

		if c.inAttackRange && c.combatFrameCounter/120 >= c.AttackSpeed {
			c.ShootProjectile(sc)
			c.combatFrameCounter = 0
		}
	}
}

func (cs *CreatureSet) ResetCombatCreature(creatureToReset *Creature) {
	if creatureToReset.targetCreature != nil {
		creatureToReset.targetCreature.enemySelected = false
	}
	creatureToReset.inCombat = false
	creatureToReset.inAttackRange = false
	_, exists := cs.inCombatCreatures[creatureToReset.ID]
	if exists {
		delete(cs.inCombatCreatures, creatureToReset.ID)
	}

}

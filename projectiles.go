package main

import rl "github.com/gen2brain/raylib-go/raylib"

// import rl "github.com/gen2brain/raylib-go/raylib"

type Projectile struct {
	CreatureRenderer
	CustomColors
	Movement

	isActive   bool
	timeActive int
}

func (c *Creature) NewProjectile() *Projectile {
	newP := Projectile{}
	newP.Pos = c.Center
	if c.isBase {
		newP.Size = rl.NewVector2(c.Size.X/16, c.Size.Y/16)
	} else {
		newP.Size = rl.NewVector2(c.Size.X/8, c.Size.Y/8)
	}
	newP.Center = rl.NewVector2(newP.Pos.X+newP.Size.X/2, newP.Pos.Y+newP.Size.Y/2)
	newP.Speed = 300
	newP.primaryTint = c.projectileTint
	newP.isActive = true
	c.projectileSlice = append(c.projectileSlice, &newP)
	return &newP
}

func (c *Creature) ShootProjectile(sc *Scene) {
	if c.targetCreature == nil {
		return
	}
	projectile := c.NewProjectile()
	PlaySoundOverlap(&sc.pewSound)

	projectile.startLocation = projectile.Center
	projectile.targetDestination = c.targetCreature.Center
	subV := rl.Vector2Subtract(projectile.targetDestination, projectile.Center)
	projectile.storedMovement.X = subV.X
	projectile.storedMovement.Y = subV.Y

}

func (p *Projectile) ProjectileMove(index int, c *Creature) {
	p.startLocation = p.Center
	p.targetDestination = c.targetCreature.Center
	subV := rl.Vector2Subtract(p.targetDestination, p.Center)
	p.storedMovement.X = subV.X
	p.storedMovement.Y = subV.Y

	// p.distanceTraveled = rl.Vector2Distance(p.startLocation, p.Center)
	// if p.distanceTraveled >= p.targetDistance {
	// 	c.DespawnProjectile(index)
	// 	return
	// }

	p.timeActive++
	if p.timeActive/120 >= 5 || rl.Vector2Distance(p.Center, p.targetDestination) < c.targetCreature.Size.X/3 {
		c.DespawnProjectile(index)
		return
	}

	p.Pos = rl.Vector2Add(p.Pos, rl.Vector2Scale(rl.Vector2Normalize(p.storedMovement), p.Speed*rl.GetFrameTime()))
	p.Center = rl.NewVector2(p.Pos.X+p.Size.X/2, p.Pos.Y+p.Size.Y/2)
	p.Collider.X = p.Pos.X
	p.Collider.Y = p.Pos.Y

}

func (sc *Scene) ProjectileCollisionCheck(c *Creature, cs *CreatureSet, p *Projectile) {

	if c.targetCreature.currentHealth <= 0 && len(c.projectileSlice) > 0 {
		cs.ResetCombatCreature(c)
		return
	}
	if c.targetCreature.Team == PLAYER {
		if rl.CheckCollisionCircleRec(p.Center, p.Size.X, c.targetCreature.Collider) {
			if c.isBase {
				PlaySoundOverlap(&sc.impactSound)
			}
			c.targetCreature.currentHealth -= c.Damage
			p.isActive = false
			return
		}
	}
	if c.targetCreature.Team == ENEMY {
		if rl.CheckCollisionCircles(p.Center, p.Size.X, c.targetCreature.Center, c.targetCreature.DangerRange) {
			if c.isBase {
				PlaySoundOverlap(&sc.impactSound)
			}
			c.targetCreature.currentHealth -= c.Damage
			p.isActive = false
			return
		}
	}

}

func (c *Creature) DrawProjectiles() {
	if len(c.projectileSlice) < 1 {
		return
	}
	for i, p := range c.projectileSlice {
		p.ProjectileMove(i, c)
		rl.DrawCircleV(p.Center, p.Size.X, p.primaryTint)
	}
}

func (c *Creature) DespawnProjectile(index int) {
	if len(c.projectileSlice) < 1 {
		return
	}
	if index < len(c.projectileSlice) {
		c.projectileSlice = append(c.projectileSlice[:index], c.projectileSlice[index+1:]...)
	} else {
		c.projectileSlice = c.projectileSlice[:len(c.projectileSlice)-1]
	}
}

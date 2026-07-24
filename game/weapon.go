package game

import (
	"image/color"
	"math"
)

type weaponType int

const (
	weaponLight weaponType = iota
	weaponFlame
	weaponIce
)

var (
	lightColor = color.RGBA{0xff, 0xf0, 0x8c, 0xff}
	flameColor = color.RGBA{0xff, 0x7a, 0x2a, 0xff}
	iceColor   = color.RGBA{0x6a, 0xd0, 0xff, 0xff}
)

func weaponName(w weaponType) string {
	switch w {
	case weaponFlame:
		return "CHAMAS"
	case weaponIce:
		return "GELO"
	default:
		return "LANCA LUZ"
	}
}

func weaponCooldown(w weaponType, level int) int {
	switch w {
	case weaponFlame:
		if level >= 3 {
			return flameCooldownFast
		}
		return flameCooldown
	case weaponIce:
		return iceCooldown
	default:
		return lightCooldown
	}
}

func fireWeapon(w weaponType, level int, x, topY float64) []*Bullet {
	switch w {
	case weaponFlame:
		return fireFlame(level, x, topY)
	case weaponIce:
		return fireIce(level, x, topY)
	default:
		return fireLight(level, x, topY)
	}
}

// Lança de Luz: disparos retos. Sobe de 2 a 4 projéteis conforme o nível.
func fireLight(level int, x, topY float64) []*Bullet {
	speed := lightBulletSpeed
	if level >= 3 {
		speed = lightBulletSpeedFast
	}
	count := 1 + level
	return straightSpread(count, x, topY, speed, lightBulletDamage, 0, lightColor)
}

// Chamas do Dragão: leque de projéteis com dano individual menor.
func fireFlame(level int, x, topY float64) []*Bullet {
	count := 3
	if level >= 2 {
		count = 5
	}
	out := make([]*Bullet, 0, count)
	for _, angle := range fanAngles(count, flameSpread) {
		vx := math.Sin(angle) * flameBulletSpeed
		vy := -math.Cos(angle) * flameBulletSpeed
		out = append(out, newBullet(x-bulletWidth/2, topY, vx, vy, flameBulletDamage, 0, flameColor))
	}
	return out
}

// Lanças de Gelo: lentas, fortes e capazes de atravessar inimigos.
func fireIce(level int, x, topY float64) []*Bullet {
	pierce := icePierce
	if level >= 3 {
		pierce = icePierceMax
	}
	return straightSpread(level, x, topY, iceBulletSpeed, iceBulletDamage, pierce, iceColor)
}

func straightSpread(count int, x, topY, speed float64, damage, pierce int, c color.RGBA) []*Bullet {
	const gap = 6.0
	out := make([]*Bullet, 0, count)
	start := x - gap*float64(count-1)/2
	for i := 0; i < count; i++ {
		bx := start + float64(i)*gap - bulletWidth/2
		out = append(out, newBullet(bx, topY, 0, -speed, damage, pierce, c))
	}
	return out
}

// fanAngles devolve ângulos simétricos ao redor da vertical (0 = para cima).
func fanAngles(count int, spread float64) []float64 {
	if count <= 1 {
		return []float64{0}
	}
	angles := make([]float64, count)
	step := spread / float64(count-1)
	start := -spread / 2
	for i := range angles {
		angles[i] = start + step*float64(i)
	}
	return angles
}

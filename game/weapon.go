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

// weaponColor devolve a cor de identidade da magia (HUD, cargas, efeitos).
func weaponColor(w weaponType) color.RGBA {
	switch w {
	case weaponFlame:
		return flameColor
	case weaponIce:
		return iceColor
	default:
		return lightColor
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

// Lança de Luz: disparos retos e densos, o maior dano focado do jogo.
// Não perfura — perfuração é a identidade do Gelo — e não aplica status.
func fireLight(level int, x, topY float64) []*Bullet {
	speed := lightBulletSpeed
	if level >= 3 {
		speed = lightBulletSpeedFast
	}
	count := 1 + level
	out := straightSpread(count, x, topY, speed, lightBulletDamage, 0, lightColor)
	for _, b := range out {
		b.trail = true
		b.element = weaponLight
	}
	return out
}

// weaponShootSFX escolhe o som de disparo com a assinatura de cada arma.
func weaponShootSFX(w weaponType) soundID {
	switch w {
	case weaponFlame:
		return sfxShootFlame
	case weaponIce:
		return sfxShootIce
	default:
		return sfxShoot
	}
}

// Chamas do Dragão: leque curto para cobertura local (não varre a tela inteira).
// As contagens são sempre ímpares para garantir um tiro central: com 4
// projéteis o leque abria um vão bem no eixo do jogador e o Nv2 rendia menos
// que o Nv1 contra alvos alinhados.
func fireFlame(level int, x, topY float64) []*Bullet {
	count := flameCountBase
	switch {
	case level >= 3:
		count = flameCountMax
	case level >= 2:
		count = flameCountMid
	}
	spread := flameSpread
	if level >= 3 {
		spread = flameSpreadWide
	}
	out := make([]*Bullet, 0, count)
	for _, angle := range fanAngles(count, spread) {
		vx := math.Sin(angle) * flameBulletSpeed
		vy := -math.Cos(angle) * flameBulletSpeed
		b := newBullet(x-bulletWidth/2, topY, vx, vy, flameBulletDamage, 0, flameColor)
		b.trail = true
		b.element = weaponFlame
		out = append(out, b)
	}
	return out
}

// Lanças de Gelo: cadência média, alto dano e perfuração — anti-formação.
func fireIce(level int, x, topY float64) []*Bullet {
	pierce := icePierce
	if level >= 2 {
		pierce = icePierce + 1
	}
	if level >= 3 {
		pierce = icePierceMax
	}
	count := level
	if count < 1 {
		count = 1
	}
	out := straightSpread(count, x, topY, iceBulletSpeed, iceBulletDamage, pierce, iceColor)
	for _, b := range out {
		b.trail = true
		b.element = weaponIce
	}
	return out
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

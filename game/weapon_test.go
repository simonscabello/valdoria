package game

import (
	"math"
	"testing"
)

func TestWeaponSwitchResetsToLevelOne(t *testing.T) {
	p := newPlayer()
	p.gainWeapon(weaponLight)
	p.gainWeapon(weaponLight)
	if p.weaponLevel != 3 {
		t.Fatalf("nível = %d, quero 3", p.weaponLevel)
	}

	p.gainWeapon(weaponIce)
	if p.weapon != weaponIce {
		t.Error("arma deveria ter trocado para gelo")
	}
	if p.weaponLevel != 1 {
		t.Errorf("nível após troca = %d, quero 1", p.weaponLevel)
	}
}

func TestWeaponLevelIncreasesAndCaps(t *testing.T) {
	p := newPlayer()
	for i := 0; i < 5; i++ {
		p.gainWeapon(weaponLight)
	}
	if p.weaponLevel != maxWeaponLevel {
		t.Errorf("nível = %d, quero %d (limite)", p.weaponLevel, maxWeaponLevel)
	}
}

func TestHealDoesNotExceedMax(t *testing.T) {
	p := newPlayer()
	p.health = maxHealth - 1
	p.applyPowerup(powerHeal)
	if p.health != maxHealth {
		t.Errorf("vida = %d, quero %d", p.health, maxHealth)
	}
	p.applyPowerup(powerHeal)
	if p.health != maxHealth {
		t.Errorf("cura não deveria ultrapassar o máximo, vida = %d", p.health)
	}
}

func TestShieldAbsorbsOneHit(t *testing.T) {
	p := newPlayer()
	p.applyPowerup(powerShield)
	full := p.health

	p.hit(1)
	if p.health != full {
		t.Errorf("escudo deveria absorver o primeiro golpe, vida = %d", p.health)
	}
	if p.hasShield() {
		t.Error("escudo deveria sumir após absorver um golpe")
	}

	p.invincible = 0
	p.hit(1)
	if p.health != full-1 {
		t.Errorf("segundo golpe deveria reduzir a vida, vida = %d", p.health)
	}
}

func TestIceBulletPiercesMultipleEnemies(t *testing.T) {
	g := &Game{player: newPlayer()}
	for i := 0; i < 3; i++ {
		e := newCrow(100)
		e.x, e.y = 100, 100
		g.enemies = append(g.enemies, e)
	}
	ice := newBullet(100, 100, 0, -iceBulletSpeed, iceBulletDamage, icePierceMax, iceColor)
	g.bullets = append(g.bullets, ice)

	g.bulletEnemyCollisions()

	for i, e := range g.enemies {
		if !e.dead {
			t.Errorf("inimigo %d deveria ser atravessado e destruído", i)
		}
	}
}

func TestLightBulletStopsAtFirstEnemy(t *testing.T) {
	g := &Game{player: newPlayer()}
	for i := 0; i < 2; i++ {
		e := newCrow(100)
		e.x, e.y = 100, 100
		g.enemies = append(g.enemies, e)
	}
	shot := newBullet(100, 100, 0, -lightBulletSpeed, lightBulletDamage, 0, lightColor)
	g.bullets = append(g.bullets, shot)

	g.bulletEnemyCollisions()

	if !shot.dead {
		t.Error("projétil sem perfuração deveria morrer no primeiro impacto")
	}
	dead := 0
	for _, e := range g.enemies {
		if e.dead {
			dead++
		}
	}
	if dead != 1 {
		t.Errorf("apenas um inimigo deveria morrer, morreram %d", dead)
	}
}

func TestFanAnglesAreSymmetric(t *testing.T) {
	angles := fanAngles(5, 1.0)
	if len(angles) != 5 {
		t.Fatalf("len = %d, quero 5", len(angles))
	}
	if math.Abs(angles[2]) > 1e-9 {
		t.Errorf("ângulo central deveria ser 0, foi %v", angles[2])
	}
	if math.Abs(angles[0]+angles[4]) > 1e-9 {
		t.Errorf("ângulos das pontas deveriam ser simétricos: %v e %v", angles[0], angles[4])
	}
	if math.Abs(angles[0]+0.5) > 1e-9 || math.Abs(angles[4]-0.5) > 1e-9 {
		t.Errorf("pontas deveriam ser -0.5 e 0.5, foram %v e %v", angles[0], angles[4])
	}

	single := fanAngles(1, 1.0)
	if len(single) != 1 || single[0] != 0 {
		t.Errorf("contagem 1 deveria dar [0], deu %v", single)
	}
}

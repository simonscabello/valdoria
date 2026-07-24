package game

import "testing"

func TestPlayerClampToScreen(t *testing.T) {
	p := newPlayer()

	p.x, p.y = -50, ScreenHeight+50
	p.clampToScreen()
	if p.x != 0 {
		t.Errorf("x deveria ser preso em 0, foi %v", p.x)
	}
	if p.y != ScreenHeight-playerSize {
		t.Errorf("y deveria ser preso na base, foi %v", p.y)
	}

	p.x = ScreenWidth + 50
	p.clampToScreen()
	if p.x != ScreenWidth-playerSize {
		t.Errorf("x deveria ser preso na borda direita, foi %v", p.x)
	}
}

func TestPlayerShieldAbsorbsSingleHit(t *testing.T) {
	p := newPlayer()
	p.shieldTimer = shieldDuration
	before := p.health

	if p.hit(1) {
		t.Error("o escudo deveria absorver o dano sem reduzir a vida")
	}
	if p.health != before {
		t.Errorf("a vida não deveria cair com escudo, foi %d", p.health)
	}
	if p.hasShield() {
		t.Error("o escudo deveria ser consumido após absorver um golpe")
	}
}

func TestRespawnReducesWeaponAndRestoresHealth(t *testing.T) {
	p := newPlayer()
	p.weaponLevel = 3
	p.health = 0

	p.respawn()

	if p.health != maxHealth {
		t.Errorf("respawn deveria restaurar a vida ao máximo, foi %d", p.health)
	}
	if p.weaponLevel != 2 {
		t.Errorf("respawn deveria reduzir um nível de arma, foi %d", p.weaponLevel)
	}
	if p.invincible <= 0 {
		t.Error("respawn deveria conceder invencibilidade")
	}
}

func TestBulletOffScreenBounds(t *testing.T) {
	cases := []struct {
		name string
		b    *Bullet
		want bool
	}{
		{"acima", &Bullet{x: 10, y: -bulletHeight - 1, w: bulletWidth, h: bulletHeight}, true},
		{"abaixo", &Bullet{x: 10, y: ScreenHeight + 1, w: bulletWidth, h: bulletHeight}, true},
		{"dentro", &Bullet{x: 10, y: 10, w: bulletWidth, h: bulletHeight}, false},
	}
	for _, c := range cases {
		if got := c.b.offScreen(); got != c.want {
			t.Errorf("%s: offScreen = %v, quero %v", c.name, got, c.want)
		}
	}
}

package game

import (
	"math"
	"testing"
)

// collectRunes entrega n runas de um elemento ao jogador.
func collectRunes(p *Player, w weaponType, n int) {
	for i := 0; i < n; i++ {
		p.gainWeapon(w)
	}
}

func TestWeaponSwitchKeepsLevel(t *testing.T) {
	p := newPlayer()
	collectRunes(p, weaponLight, 2*runesPerLevel) // Nv1 -> Nv3
	if p.weaponLevel != maxWeaponLevel {
		t.Fatalf("nível = %d, quero %d", p.weaponLevel, maxWeaponLevel)
	}

	p.gainWeapon(weaponIce)
	if p.weapon != weaponIce {
		t.Error("arma deveria ter trocado para gelo")
	}
	if p.weaponLevel != maxWeaponLevel {
		t.Errorf("trocar de arma não deveria rebaixar o poder: nível %d, quero %d", p.weaponLevel, maxWeaponLevel)
	}
}

// Trocar de elemento nunca rebaixa o poder e nunca desperdiça a runa: o nível e
// a carga são uma potência compartilhada entre as três magias.
func TestWeaponSwitchAtLowLevelPreservesLevel(t *testing.T) {
	p := newPlayer() // Luz Nv1
	p.gainWeapon(weaponFlame)
	if p.weapon != weaponFlame || p.weaponLevel != 1 {
		t.Fatalf("trocar no nível 1 deveria manter Nv1, foi arma %d nível %d", p.weapon, p.weaponLevel)
	}
	collectRunes(p, weaponFlame, 2*runesPerLevel-1) // completa até o máximo
	if p.weaponLevel != maxWeaponLevel {
		t.Fatalf("as runas deveriam subir até o máximo, foi %d", p.weaponLevel)
	}
	p.gainWeapon(weaponLight) // troca de volta mantendo o nível alto
	if p.weapon != weaponLight || p.weaponLevel != maxWeaponLevel {
		t.Errorf("troca deveria preservar o nível %d na nova arma, foi arma %d nível %d", maxWeaponLevel, p.weapon, p.weaponLevel)
	}
}

// Toda runa elemental avança a carga, mesmo alternando de elemento — caso
// contrário a progressão trava, porque as runas garantidas alternam de propósito.
func TestMixedRunesStillProgress(t *testing.T) {
	p := newPlayer()
	elements := []weaponType{weaponLight, weaponFlame, weaponIce}
	for i := 0; i < runesPerLevel; i++ {
		p.gainWeapon(elements[i%len(elements)])
	}
	if p.weaponLevel != 2 {
		t.Errorf("runas alternadas deveriam subir de nível: nível %d, quero 2", p.weaponLevel)
	}
}

func TestWeaponLevelIncreasesAndCaps(t *testing.T) {
	p := newPlayer()
	collectRunes(p, weaponLight, 5*runesPerLevel)
	if p.weaponLevel != maxWeaponLevel {
		t.Errorf("nível = %d, quero %d (limite)", p.weaponLevel, maxWeaponLevel)
	}
	if p.weaponCharge != 0 {
		t.Errorf("no teto de poder a carga não deveria acumular, foi %d", p.weaponCharge)
	}
}

// Cada runa avança a carga visível no HUD: exigir várias runas por nível sem
// mostrar progresso tiraria recompensa em vez de espalhá-la pela partida.
func TestWeaponChargeAccumulates(t *testing.T) {
	p := newPlayer()
	if p.weaponCharge != 0 {
		t.Fatal("sem runas coletadas a carga deveria ser zero")
	}
	p.gainWeapon(weaponLight)
	if p.weaponCharge != 1 {
		t.Fatalf("uma runa deveria acumular uma carga, foi %d", p.weaponCharge)
	}
	collectRunes(p, weaponLight, runesPerLevel-1)
	if p.weaponLevel != 2 || p.weaponCharge != 0 {
		t.Fatalf("completar a carga deveria subir de nível e zerar: nível %d, carga %d",
			p.weaponLevel, p.weaponCharge)
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
	shot.element = weaponFlame // só testa parar no impacto; cadeia é outro teste
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

func TestWeaponBalanceIdentities(t *testing.T) {
	// Chamas Nv3: leque denso e ainda contido (não varre a tela), com cadência
	// abaixo da Luz.
	flame := fireFlame(3, 120, 200)
	if len(flame) != flameCountMax {
		t.Fatalf("chamas Nv3: %d tiros, quero %d", len(flame), flameCountMax)
	}
	if flameSpreadWide > 0.9 {
		t.Fatalf("leque de chamas ainda largo demais: %v", flameSpreadWide)
	}
	if weaponCooldown(weaponFlame, 3) <= weaponCooldown(weaponLight, 3) {
		t.Fatalf("chamas Nv3 não deveria atirar mais rápido que a luz")
	}

	// Luz: mais dano por tiro e perfuração no Nv3.
	light := fireLight(3, 120, 200)
	if len(light) != 4 {
		t.Fatalf("luz Nv3: %d tiros, quero 4", len(light))
	}
	if light[0].damage <= flameBulletDamage {
		t.Fatalf("luz deveria ferir mais que cada chama (%d vs %d)", light[0].damage, flameBulletDamage)
	}
	// A perfuração é a identidade exclusiva do Gelo. Quando a Luz também
	// perfurava, ela vencia o cenário de coluna além do de foco e não sobrava
	// motivo para escolher outra magia.
	if light[0].pierce != 0 {
		t.Fatalf("a luz não deveria perfurar (identidade do gelo), pierce = %d", light[0].pierce)
	}

	// Gelo: alto dano + perfuração.
	ice := fireIce(3, 120, 200)
	if len(ice) != 3 {
		t.Fatalf("gelo Nv3: %d tiros, quero 3", len(ice))
	}
	if ice[0].damage <= lightBulletDamage {
		t.Fatalf("gelo deveria ser o tiro mais pesado (%d vs luz %d)", ice[0].damage, lightBulletDamage)
	}
	if ice[0].pierce < icePierceMax {
		t.Fatalf("gelo Nv3 pierce=%d, quero %d", ice[0].pierce, icePierceMax)
	}
}

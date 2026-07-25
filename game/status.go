package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// applyWeaponHit aplica o efeito elemental do projétil ao inimigo.
//
// A Lança de Luz **não** aplica status. Ela atordoava a cada acerto
// (stunDuration = 32 frames) com um cooldown de 8 frames: qualquer alvo sob
// fogo contínuo ficava congelado para sempre, inclusive o chefe. Isso apagava
// o combate e tornava as outras duas armas irrelevantes. A identidade da Luz
// passou a ser dano focado puro; o atordoamento fica reservado para uma
// habilidade ativa (o Grito do Grifo).
func (e *Enemy) applyWeaponHit(element weaponType) {
	switch element {
	case weaponIce:
		e.slow = iceSlowDuration
	case weaponFlame:
		e.burn = burnDuration
		e.burnTick = 0
	}
}

func (e *Enemy) tickStatusEffects() {
	if e.slow > 0 {
		e.slow--
	}
	if e.stun > 0 {
		e.stun--
	}
	if e.burn <= 0 {
		return
	}
	e.burn--
	e.burnTick++
	if e.burnTick%burnInterval == 0 && !e.dead {
		e.takeDamage(burnDamage)
	}
}

func (e *Enemy) drawStatusAura(screen *ebiten.Image) {
	x := float32(e.x - 2)
	y := float32(e.y - 2)
	w := float32(e.w + 4)
	h := float32(e.h + 4)
	if e.burn > 0 && (e.animTick/4)%2 == 0 {
		vector.DrawFilledRect(screen, x, y, w, h, withAlpha(flameColor, 70), false)
	}
	if e.slow > 0 {
		vector.StrokeRect(screen, x, y, w, h, 1, withAlpha(iceColor, 180), false)
	}
	if e.stun > 0 && (e.animTick/3)%2 == 0 {
		vector.StrokeRect(screen, x-1, y-1, w+2, h+2, 1, withAlpha(lightColor, 220), false)
	}
}

// --- Status no chefe (mesmas identidades das armas) ---

// applyWeaponHit no chefe segue a mesma regra: nenhum atordoamento passivo.
// Vharak precisa executar seus padrões — é a melhor peça de design do jogo.
func (b *Boss) applyWeaponHit(element weaponType) {
	switch element {
	case weaponIce:
		b.slow = iceSlowDuration
	case weaponFlame:
		b.burn = burnDuration
		b.burnTick = 0
	}
}

func (b *Boss) tickStatusEffects() {
	if b.slow > 0 {
		b.slow--
	}
	if b.stun > 0 {
		b.stun--
	}
	if b.burn <= 0 {
		return
	}
	b.burn--
	b.burnTick++
	if b.burnTick%burnInterval == 0 && b.phase != bossDying && b.phase != bossDead {
		b.takeDamage(burnDamage)
	}
}

func (b *Boss) drawStatusAura(screen *ebiten.Image) {
	x := float32(b.x - 2)
	y := float32(b.y - 2)
	w := float32(b.w + 4)
	h := float32(b.h + 4)
	if b.burn > 0 && (b.tick/4)%2 == 0 {
		vector.DrawFilledRect(screen, x, y, w, h, withAlpha(flameColor, 70), false)
	}
	if b.slow > 0 {
		vector.StrokeRect(screen, x, y, w, h, 1, withAlpha(iceColor, 180), false)
	}
	if b.stun > 0 && (b.tick/3)%2 == 0 {
		vector.StrokeRect(screen, x-1, y-1, w+2, h+2, 1, withAlpha(lightColor, 220), false)
	}
}

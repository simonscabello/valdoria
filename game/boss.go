package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type bossPhase int

const (
	bossEntering bossPhase = iota
	bossPhase1
	bossPhase2
	bossPhase3
	bossDying
	bossDead
)

type bossAction int

const (
	actionWarning bossAction = iota
	actionFiring
	actionCooldown
)

type patternKind int

const (
	patternAimedFire patternKind = iota
	patternCone
	patternArc
	patternSummon
	patternSweep
)

type bossPattern struct {
	kind     patternKind
	duration int
	cooldown int
	interval int
}

var bossPatterns = map[bossPhase][]bossPattern{
	bossPhase1: {
		{kind: patternAimedFire, duration: 90, cooldown: 55, interval: 18},
		{kind: patternCone, duration: 70, cooldown: 60, interval: 30},
	},
	bossPhase2: {
		{kind: patternArc, duration: 90, cooldown: 45, interval: 22},
		{kind: patternSummon, duration: 60, cooldown: 70, interval: 30},
		{kind: patternAimedFire, duration: 80, cooldown: 40, interval: 14},
	},
	bossPhase3: {
		{kind: patternSweep, duration: 100, cooldown: 45, interval: 26},
		{kind: patternCone, duration: 80, cooldown: 35, interval: 16},
		{kind: patternAimedFire, duration: 80, cooldown: 35, interval: 10},
	},
}

type weakSpot struct {
	dx, dy, size float64
}

var bossWeakSpots = []weakSpot{
	{dx: 10, dy: 10, size: 8},
	{dx: bossW/2 - 4, dy: 20, size: 8},
	{dx: bossW - 18, dy: 10, size: 8},
}

// bossContext dá ao chefe acesso ao jogador e às listas que ele alimenta.
type bossContext struct {
	playerX, playerY float64
	bullets          *[]*EnemyBullet
	enemies          *[]*Enemy
}

type Boss struct {
	x, y           float64
	w, h           float64
	health         int
	maxHealth      int
	phase          bossPhase
	invulnerable   bool
	hitFlash       int
	tick           int
	vx             float64
	entryTimer     int
	dyingTimer     int
	action         bossAction
	actionTimer    int
	patternIndex   int
	fireTimer      int
	crystalsActive bool

	justEntered  bool
	phaseChanged bool
}

func newBoss() *Boss {
	return &Boss{
		x:            ScreenWidth/2 - bossW/2,
		y:            -bossH,
		w:            bossW,
		h:            bossH,
		health:       bossMaxHealth,
		maxHealth:    bossMaxHealth,
		phase:        bossEntering,
		invulnerable: true,
		entryTimer:   bossEntryHold,
		vx:           1,
	}
}

func (b *Boss) centerX() float64 { return b.x + b.w/2 }

func (b *Boss) attacking() bool { return b.action == actionFiring }

func (b *Boss) healthRatio() float64 { return float64(b.health) / float64(b.maxHealth) }

func (b *Boss) combatPhaseForHealth() bossPhase {
	ratio := b.healthRatio()
	if ratio > 0.65 {
		return bossPhase1
	}
	if ratio > 0.30 {
		return bossPhase2
	}
	return bossPhase3
}

func (b *Boss) takeDamage(n int) {
	if b.invulnerable {
		return
	}
	b.health -= n
	b.hitFlash = hitFlashDuration
	if b.health <= 0 {
		b.health = 0
		b.startDying()
	}
}

func (b *Boss) startDying() {
	b.phase = bossDying
	b.invulnerable = true
	b.dyingTimer = bossDeathDuration
}

func (b *Boss) update(ctx *bossContext) {
	b.tick++
	if b.hitFlash > 0 {
		b.hitFlash--
	}

	switch b.phase {
	case bossEntering:
		b.doEntry()
		return
	case bossDying:
		if b.dyingTimer > 0 {
			b.dyingTimer--
		}
		if b.dyingTimer <= 0 {
			b.phase = bossDead
		}
		return
	case bossDead:
		return
	}

	b.refreshPhase()
	b.move()
	b.runActions(ctx)
}

func (b *Boss) doEntry() {
	if b.y < bossY {
		b.y += bossEntrySpeed
		if b.y >= bossY {
			b.y = bossY
		}
		return
	}
	b.entryTimer--
	if b.entryTimer <= 0 {
		b.beginCombat()
	}
}

func (b *Boss) beginCombat() {
	b.invulnerable = false
	b.phase = bossPhase1
	b.patternIndex = 0
	b.action = actionWarning
	b.actionTimer = bossWarningDuration
	b.justEntered = true
}

func (b *Boss) refreshPhase() {
	desired := b.combatPhaseForHealth()
	if desired == b.phase {
		return
	}
	b.phase = desired
	b.patternIndex = 0
	b.action = actionWarning
	b.actionTimer = bossWarningDuration
	b.phaseChanged = true
	if desired == bossPhase3 {
		b.crystalsActive = true
	}
}

func (b *Boss) speed() float64 {
	switch b.phase {
	case bossPhase2:
		return bossSpeedP2
	case bossPhase3:
		return bossSpeedP3
	default:
		return bossSpeedP1
	}
}

func (b *Boss) move() {
	b.x += b.vx * b.speed()
	minX := float64(bossMargin)
	maxX := ScreenWidth - b.w - bossMargin
	if b.x <= minX {
		b.x = minX
		b.vx = 1
	}
	if b.x >= maxX {
		b.x = maxX
		b.vx = -1
	}
}

func (b *Boss) currentPattern() bossPattern {
	pats := bossPatterns[b.phase]
	if len(pats) == 0 {
		return bossPattern{}
	}
	return pats[b.patternIndex%len(pats)]
}

func (b *Boss) runActions(ctx *bossContext) {
	b.actionTimer--
	switch b.action {
	case actionWarning:
		if b.actionTimer <= 0 {
			b.action = actionFiring
			b.actionTimer = b.currentPattern().duration
			b.fireTimer = 0
		}
	case actionFiring:
		b.fire(ctx)
		if b.actionTimer <= 0 {
			b.action = actionCooldown
			b.actionTimer = b.currentPattern().cooldown
		}
	case actionCooldown:
		if b.actionTimer <= 0 {
			b.nextPattern()
		}
	}
}

func (b *Boss) nextPattern() {
	pats := bossPatterns[b.phase]
	b.patternIndex = (b.patternIndex + 1) % len(pats)
	b.action = actionWarning
	b.actionTimer = bossWarningDuration
}

func (b *Boss) fire(ctx *bossContext) {
	b.fireTimer--
	if b.fireTimer > 0 {
		return
	}
	pat := b.currentPattern()
	b.fireTimer = pat.interval

	switch pat.kind {
	case patternAimedFire:
		b.emitAimed(ctx)
	case patternCone:
		b.emitFan(ctx, 5, 1.1)
	case patternArc:
		b.emitFan(ctx, 9, 2.0)
	case patternSummon:
		b.emitSummon(ctx)
	case patternSweep:
		b.emitSweep(ctx)
	}
}

func (b *Boss) emitBullet(ctx *bossContext, x, y, vx, vy float64) {
	*ctx.bullets = append(*ctx.bullets, newEnemyBullet(x, y, vx, vy))
}

func (b *Boss) emitAimed(ctx *bossContext) {
	vx, vy := aimVelocity(b.centerX(), b.y+b.h, ctx.playerX, ctx.playerY, enemyBulletSpeed)
	b.emitBullet(ctx, b.centerX(), b.y+b.h, vx, vy)
}

func (b *Boss) emitFan(ctx *bossContext, count int, spread float64) {
	for _, a := range fanAngles(count, spread) {
		vx := math.Sin(a) * enemyBulletSpeed
		vy := math.Cos(a) * enemyBulletSpeed
		b.emitBullet(ctx, b.centerX(), b.y+b.h, vx, vy)
	}
}

func (b *Boss) emitSummon(ctx *bossContext) {
	if len(*ctx.enemies) >= 6 {
		return
	}
	if b.tick%2 == 0 {
		*ctx.enemies = append(*ctx.enemies, newCrow(b.centerX()))
		return
	}
	*ctx.enemies = append(*ctx.enemies, newHarpy(b.centerX()))
}

// emitSweep cria uma parede horizontal com uma brecha móvel, sempre desviável.
func (b *Boss) emitSweep(ctx *bossContext) {
	gap := math.Mod(float64(b.tick)*bossSweepGapMove, ScreenWidth)
	for x := 0.0; x < ScreenWidth; x += bossSweepStep {
		if x > gap-bossSweepGapHalf && x < gap+bossSweepGapHalf {
			continue
		}
		b.emitBullet(ctx, x, b.y+b.h, 0, bossSweepSpeed)
	}
}

func (b *Boss) hitCrystal(x, y, w, h float64) bool {
	if !b.crystalsActive {
		return false
	}
	for _, s := range bossWeakSpots {
		if collides(x, y, w, h, b.x+s.dx, b.y+s.dy, s.size, s.size) {
			return true
		}
	}
	return false
}

func (b *Boss) draw(screen *ebiten.Image) {
	if b.phase == bossDead {
		return
	}
	body := color.RGBA{0x8a, 0x2a, 0x2a, 0xff}
	if b.hitFlash > 0 {
		body = color.RGBA{0xff, 0xff, 0xff, 0xff}
	}
	if b.phase == bossDying && (b.tick/4)%2 == 0 {
		body = color.RGBA{0xff, 0xf0, 0xa0, 0xff}
	}
	x, y := float32(b.x), float32(b.y)
	w, h := float32(b.w), float32(b.h)
	vector.DrawFilledRect(screen, x, y, w, h, body, false)

	wing := color.RGBA{0x5a, 0x18, 0x18, 0xff}
	flap := float32(math.Sin(float64(b.tick)*bossWingSpeed) * 4)
	vector.DrawFilledRect(screen, x-12, y+2-flap, 12, h-6, wing, false)
	vector.DrawFilledRect(screen, x+w, y+2-flap, 12, h-6, wing, false)
	vector.DrawFilledRect(screen, x+w/2-6, y+h, 12, 6, wing, false)
	// Chifres e olhos brilhantes reforçam a leitura da cabeça do dragão.
	eye := color.RGBA{0xff, 0xd0, 0x40, 0xff}
	vector.DrawFilledRect(screen, x+w/2-10, y+6, 3, 3, eye, false)
	vector.DrawFilledRect(screen, x+w/2+7, y+6, 3, 3, eye, false)

	if b.crystalsActive && b.phase != bossDying {
		crystal := color.RGBA{0x8a, 0xf0, 0xff, 0xff}
		for _, s := range bossWeakSpots {
			vector.DrawFilledRect(screen, x+float32(s.dx), y+float32(s.dy), float32(s.size), float32(s.size), crystal, false)
		}
	}
}

package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Resolução interna pensada para pixel art. A janela é escalada no main.
const (
	ScreenWidth  = 240
	ScreenHeight = 320
)

type state int

const (
	statePlaying state = iota
	statePaused
	stateBossIncoming
	stateGameOver
)

type Game struct {
	state state

	player       *Player
	bullets      []*Bullet
	enemies      []*Enemy
	enemyBullets []*EnemyBullet
	powerups     []*Powerup
	stars        []*star

	level       *Level
	timeScale   int
	resumeState state
	score       int
}

func New() *Game {
	g := &Game{}
	g.reset()
	g.initStars()
	return g
}

func (g *Game) reset() {
	g.state = statePlaying
	g.player = newPlayer()
	g.bullets = g.bullets[:0]
	g.enemies = g.enemies[:0]
	g.enemyBullets = g.enemyBullets[:0]
	g.powerups = g.powerups[:0]
	g.level = newLevel()
	g.timeScale = devTimeScale
	g.score = 0
}

func (g *Game) Update() error {
	g.handlePauseToggle()

	if g.state == statePaused {
		return nil
	}

	g.updateStars()

	if g.state == stateGameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			g.reset()
		}
		return nil
	}

	g.player.update()
	g.bullets = append(g.bullets, g.player.tryShoot()...)
	g.handleDebugSpawns()

	if g.state == statePlaying {
		g.handleTimeScaleToggle()
		for i := 0; i < g.timeScale; i++ {
			g.enemies = append(g.enemies, g.level.update()...)
		}
	}

	g.updateBullets()
	g.updateEnemies()
	g.updateEnemyBullets()
	g.updatePowerups()
	g.handleCollisions()
	g.collectPowerups()
	g.removeDead()

	if g.player.health <= 0 {
		g.state = stateGameOver
		return nil
	}
	if g.state == statePlaying && g.level.readyForBoss(len(g.enemies)) {
		g.state = stateBossIncoming
	}

	return nil
}

func (g *Game) handlePauseToggle() {
	if !inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return
	}
	if g.state == statePaused {
		g.state = g.resumeState
		return
	}
	if g.state == statePlaying || g.state == stateBossIncoming {
		g.resumeState = g.state
		g.state = statePaused
	}
}

func (g *Game) handleTimeScaleToggle() {
	if !inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		return
	}
	if g.timeScale == devFastTimeScale {
		g.timeScale = devTimeScale
		return
	}
	g.timeScale = devFastTimeScale
}

func (g *Game) updateBullets() {
	for _, b := range g.bullets {
		b.update()
		if b.offScreen() {
			b.dead = true
		}
	}
}

func (g *Game) spawnCrowFormation() {
	count := 3 + rand.Intn(2)
	const gap = crowSize + 6
	startX := rand.Float64() * (ScreenWidth - float64(count)*gap)
	for i := 0; i < count; i++ {
		g.enemies = append(g.enemies, newCrow(startX+float64(i)*gap))
	}
}

func (g *Game) spawnGargoyle() {
	fromLeft := rand.Intn(2) == 0
	targetX := 40 + rand.Float64()*(ScreenWidth-80-gargoyleSize)
	y := 30 + rand.Float64()*60
	g.enemies = append(g.enemies, newGargoyle(fromLeft, targetX, y))
}

func randX(size float64) float64 {
	return rand.Float64() * (ScreenWidth - size)
}

// handleDebugSpawns permite forçar cada tipo de inimigo para teste manual.
func (g *Game) handleDebugSpawns() {
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit1) {
		g.spawnCrowFormation()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit2) {
		g.enemies = append(g.enemies, newHarpy(randX(harpySize)))
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit3) {
		g.spawnGargoyle()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDigit4) {
		g.enemies = append(g.enemies, newWyvern(randX(wyvernSize)))
	}

	if !devMode {
		return
	}
	g.debugSpawnPowerup(ebiten.KeyZ, powerLight)
	g.debugSpawnPowerup(ebiten.KeyX, powerFire)
	g.debugSpawnPowerup(ebiten.KeyC, powerIce)
	g.debugSpawnPowerup(ebiten.KeyV, powerHeal)
	g.debugSpawnPowerup(ebiten.KeyB, powerShield)
}

func (g *Game) debugSpawnPowerup(key ebiten.Key, kind powerupType) {
	if inpututil.IsKeyJustPressed(key) {
		g.powerups = append(g.powerups, newPowerup(kind, g.player.centerX(), 24))
	}
}

func (g *Game) updateEnemies() {
	ctx := &enemyContext{
		playerX: g.player.centerX(),
		playerY: g.player.centerY(),
		bullets: &g.enemyBullets,
	}
	for _, e := range g.enemies {
		e.update(ctx)
		if e.offScreen() {
			e.dead = true
		}
	}
}

func (g *Game) updateEnemyBullets() {
	for _, b := range g.enemyBullets {
		b.update()
		if b.offScreen() {
			b.dead = true
		}
	}
}

func (g *Game) handleCollisions() {
	g.bulletEnemyCollisions()

	if g.player.canBeHit() {
		g.playerEnemyCollisions()
	}
	if g.player.canBeHit() {
		g.playerBulletCollisions()
	}
}

func (g *Game) bulletEnemyCollisions() {
	for _, b := range g.bullets {
		if b.dead {
			continue
		}
		for _, e := range g.enemies {
			if e.dead || b.alreadyHit(e) {
				continue
			}
			if !collides(b.x, b.y, b.w, b.h, e.x, e.y, e.w, e.h) {
				continue
			}
			e.takeDamage(b.damage)
			b.hitEnemies = append(b.hitEnemies, e)
			if e.dead {
				g.score += e.score
				g.spawnDrop(e)
			}
			if b.pierce > 0 {
				b.pierce--
				continue
			}
			b.dead = true
			break
		}
	}
}

func (g *Game) playerEnemyCollisions() {
	hx, hy, hw, hh := g.player.hitbox()
	for _, e := range g.enemies {
		if e.dead {
			continue
		}
		if collides(hx, hy, hw, hh, e.x, e.y, e.w, e.h) {
			e.dead = true
			g.player.hit(e.damage)
			return
		}
	}
}

func (g *Game) playerBulletCollisions() {
	hx, hy, hw, hh := g.player.hitbox()
	for _, b := range g.enemyBullets {
		if b.dead {
			continue
		}
		if collides(hx, hy, hw, hh, b.x, b.y, enemyBulletSize, enemyBulletSize) {
			b.dead = true
			g.player.hit(enemyBulletDamage)
			return
		}
	}
}

func (g *Game) spawnDrop(e *Enemy) {
	if e.hasDrop {
		g.powerups = append(g.powerups, newPowerup(e.drop, e.centerX(), e.centerY()))
		return
	}
	if rand.Float64() < dropChance {
		g.powerups = append(g.powerups, newPowerup(randomRune(), e.centerX(), e.centerY()))
	}
}

func randomRune() powerupType {
	switch rand.Intn(5) {
	case 0:
		return powerLight
	case 1:
		return powerFire
	case 2:
		return powerIce
	case 3:
		return powerHeal
	default:
		return powerShield
	}
}

func (g *Game) updatePowerups() {
	for _, p := range g.powerups {
		p.update()
		if p.offScreen() {
			p.dead = true
		}
	}
}

func (g *Game) collectPowerups() {
	for _, p := range g.powerups {
		if p.dead {
			continue
		}
		if collides(g.player.x, g.player.y, playerSize, playerSize, p.x, p.y, powerupSize, powerupSize) {
			g.player.applyPowerup(p.kind)
			p.dead = true
		}
	}
}

func (g *Game) removeDead() {
	g.bullets = filterAlive(g.bullets, func(b *Bullet) bool { return !b.dead })
	g.enemies = filterAlive(g.enemies, func(e *Enemy) bool { return !e.dead })
	g.enemyBullets = filterAlive(g.enemyBullets, func(b *EnemyBullet) bool { return !b.dead })
	g.powerups = filterAlive(g.powerups, func(p *Powerup) bool { return !p.dead })
}

// filterAlive mantém apenas os itens aprovados por keep, reaproveitando o array
// e liberando as posições restantes para que o GC recolha itens removidos.
func filterAlive[T any](items []*T, keep func(*T) bool) []*T {
	alive := items[:0]
	for _, it := range items {
		if keep(it) {
			alive = append(alive, it)
		}
	}
	for i := len(alive); i < len(items); i++ {
		items[i] = nil
	}
	return alive
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.level.background())
	g.drawStars(screen)

	for _, b := range g.bullets {
		b.draw(screen)
	}
	for _, b := range g.enemyBullets {
		b.draw(screen)
	}
	for _, p := range g.powerups {
		p.draw(screen)
	}
	for _, e := range g.enemies {
		e.draw(screen)
	}
	g.player.draw(screen)

	g.drawHUD(screen)

	if g.level.announceTimer > 0 && g.level.announce != "" {
		ebitenutil.DebugPrintAt(screen, g.level.announce, ScreenWidth/2-len(g.level.announce)*3, 70)
	}

	if g.state == statePaused {
		ebitenutil.DebugPrintAt(screen, "PAUSADO", ScreenWidth/2-24, ScreenHeight/2-4)
	}

	if g.state == stateBossIncoming {
		ebitenutil.DebugPrintAt(screen, "O CHEFE SE APROXIMA", ScreenWidth/2-57, ScreenHeight/2-4)
	}

	if g.state == stateGameOver {
		ebitenutil.DebugPrintAt(screen, "GAME OVER", ScreenWidth/2-28, ScreenHeight/2-8)
		ebitenutil.DebugPrintAt(screen, "ENTER PARA REINICIAR", ScreenWidth/2-60, ScreenHeight/2+8)
	}
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "SCORE "+itoa(g.score), 4, 4)
	ebitenutil.DebugPrintAt(screen, "VIDA "+itoa(g.player.health), 4, 16)

	weapon := weaponName(g.player.weapon) + " Nv" + itoa(g.player.weaponLevel)
	ebitenutil.DebugPrintAt(screen, weapon, 4, ScreenHeight-16)
	if g.player.hasShield() {
		ebitenutil.DebugPrintAt(screen, "ESCUDO", ScreenWidth-46, ScreenHeight-16)
	}

	// Barra de progresso discreta na borda superior.
	barWidth := float32(g.level.progress() * ScreenWidth)
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, 2, color.RGBA{0x22, 0x22, 0x33, 0xff}, false)
	vector.DrawFilledRect(screen, 0, 0, barWidth, 2, color.RGBA{0xd0, 0xc0, 0x50, 0xff}, false)

	name := g.level.sectionName()
	ebitenutil.DebugPrintAt(screen, name, ScreenWidth-len(name)*6-4, 4)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// collides testa sobreposição entre dois retângulos AABB.
func collides(ax, ay, aw, ah, bx, by, bw, bh float64) bool {
	return ax < bx+bw && ax+aw > bx && ay < by+bh && ay+ah > by
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

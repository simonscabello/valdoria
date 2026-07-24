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

	lives           int
	bombCharges     int
	bombEffectTimer int

	score      int
	highScore  int
	combo      int
	comboTimer int

	sectionDamaged bool
	lastSection    int
	formations     map[int]*formationTracker
}

type formationTracker struct {
	total   int
	killed  int
	escaped int
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

	g.lives = startingLives
	g.bombCharges = bombStartCharges
	g.bombEffectTimer = 0
	g.score = 0
	g.combo = 0
	g.comboTimer = 0
	g.sectionDamaged = false
	g.lastSection = g.level.section
	g.formations = map[int]*formationTracker{}
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
	g.handleBomb()

	if g.bombEffectTimer > 0 {
		g.bombEffectTimer--
	}
	if g.comboTimer > 0 {
		g.comboTimer--
		if g.comboTimer == 0 {
			g.combo = 0
		}
	}

	if g.state == statePlaying {
		g.handleTimeScaleToggle()
		for i := 0; i < g.timeScale; i++ {
			g.spawnFromLevel()
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
		g.loseLife()
	}
	if g.score > g.highScore {
		g.highScore = g.score
	}
	if g.state == statePlaying {
		g.checkSectionBonus()
		if g.level.readyForBoss(len(g.enemies)) {
			if !g.sectionDamaged {
				g.score += waveNoDamageBonus
			}
			g.state = stateBossIncoming
		}
	}

	return nil
}

func (g *Game) spawnFromLevel() {
	for _, e := range g.level.update() {
		if e.formationID > 0 {
			g.registerFormationMember(e.formationID)
		}
		g.enemies = append(g.enemies, e)
	}
}

func (g *Game) registerFormationMember(id int) {
	t := g.formations[id]
	if t == nil {
		t = &formationTracker{}
		g.formations[id] = t
	}
	t.total++
}

func (g *Game) loseLife() {
	g.lives--
	if g.lives <= 0 {
		g.state = stateGameOver
		return
	}
	g.respawn()
}

func (g *Game) respawn() {
	g.enemyBullets = g.enemyBullets[:0]
	g.player.respawn()
	g.combo = 0
	g.comboTimer = 0
}

func (g *Game) checkSectionBonus() {
	if g.level.section == g.lastSection {
		return
	}
	if !g.sectionDamaged {
		g.score += waveNoDamageBonus
	}
	g.sectionDamaged = false
	g.lastSection = g.level.section
}

func (g *Game) handleBomb() {
	if g.state != statePlaying || g.bombCharges <= 0 {
		return
	}
	if !inpututil.IsKeyJustPressed(ebiten.KeyX) && !inpututil.IsKeyJustPressed(ebiten.KeyControl) {
		return
	}
	g.useBomb()
}

func (g *Game) useBomb() {
	g.bombCharges--
	g.enemyBullets = g.enemyBullets[:0]
	for _, e := range g.enemies {
		if e.dead {
			continue
		}
		e.takeDamage(bombDamage)
		if e.dead {
			g.score += e.score
		}
	}
	g.player.invincible = bombInvincibility
	g.bombEffectTimer = bombEffectDuration
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
	g.debugSpawnPowerup(ebiten.KeyF, powerFire)
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
			e.escaped = true
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
				g.registerKill(e.score)
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
			if g.player.hit(e.damage) {
				g.onPlayerDamaged()
			}
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
			if g.player.hit(enemyBulletDamage) {
				g.onPlayerDamaged()
			}
			return
		}
	}
}

func (g *Game) registerKill(base int) {
	g.combo++
	g.comboTimer = comboWindow
	g.score += base * g.multiplier()
}

func (g *Game) multiplier() int {
	m := 1 + g.combo/comboStep
	if m > maxMultiplier {
		return maxMultiplier
	}
	return m
}

func (g *Game) onPlayerDamaged() {
	g.combo = 0
	g.comboTimer = 0
	g.sectionDamaged = true
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
	for _, e := range g.enemies {
		if e.dead && e.formationID > 0 {
			g.accountFormation(e)
		}
	}

	g.bullets = filterAlive(g.bullets, func(b *Bullet) bool { return !b.dead })
	g.enemies = filterAlive(g.enemies, func(e *Enemy) bool { return !e.dead })
	g.enemyBullets = filterAlive(g.enemyBullets, func(b *EnemyBullet) bool { return !b.dead })
	g.powerups = filterAlive(g.powerups, func(p *Powerup) bool { return !p.dead })
}

// accountFormation concede bônus quando toda uma formação é destruída sem
// que nenhum de seus membros escape da tela.
func (g *Game) accountFormation(e *Enemy) {
	t := g.formations[e.formationID]
	if t == nil {
		return
	}
	if e.escaped {
		t.escaped++
	} else {
		t.killed++
	}
	if t.killed+t.escaped < t.total {
		return
	}
	if t.escaped == 0 {
		g.score += formationBonus
	}
	delete(g.formations, e.formationID)
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

	if g.bombEffectTimer > 0 {
		g.drawBombEffect(screen)
	}

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

// drawBombEffect desenha o dragão ancestral atravessando a tela de baixo p/ cima.
func (g *Game) drawBombEffect(screen *ebiten.Image) {
	progress := 1 - float64(g.bombEffectTimer)/float64(bombEffectDuration)
	y := float32(ScreenHeight - progress*(ScreenHeight+60))
	body := color.RGBA{0xff, 0xd0, 0x60, 0xc0}
	vector.DrawFilledRect(screen, 0, y, ScreenWidth, 24, body, false)
	wing := color.RGBA{0xff, 0xa0, 0x30, 0xa0}
	vector.DrawFilledRect(screen, 0, y-10, ScreenWidth, 6, wing, false)
	vector.DrawFilledRect(screen, 0, y+28, ScreenWidth, 6, wing, false)
}

func (g *Game) drawHUD(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "SCORE "+itoa(g.score), 4, 4)
	ebitenutil.DebugPrintAt(screen, "HIGH "+itoa(g.highScore), 4, 16)
	if g.multiplier() > 1 {
		ebitenutil.DebugPrintAt(screen, "x"+itoa(g.multiplier()), 4, 28)
	}

	g.drawHealthBar(screen)
	ebitenutil.DebugPrintAt(screen, "VIDAS "+itoa(g.lives), 4, ScreenHeight-28)
	ebitenutil.DebugPrintAt(screen, "BOMBA "+itoa(g.bombCharges), 4, ScreenHeight-16)

	weapon := weaponName(g.player.weapon) + " Nv" + itoa(g.player.weaponLevel)
	ebitenutil.DebugPrintAt(screen, weapon, ScreenWidth-len(weapon)*6-4, ScreenHeight-16)
	if g.player.hasShield() {
		ebitenutil.DebugPrintAt(screen, "ESCUDO", ScreenWidth-46, ScreenHeight-28)
	}

	// Barra de progresso discreta na borda superior.
	barWidth := float32(g.level.progress() * ScreenWidth)
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, 2, color.RGBA{0x22, 0x22, 0x33, 0xff}, false)
	vector.DrawFilledRect(screen, 0, 0, barWidth, 2, color.RGBA{0xd0, 0xc0, 0x50, 0xff}, false)

	name := g.level.sectionName()
	ebitenutil.DebugPrintAt(screen, name, ScreenWidth-len(name)*6-4, 4)
}

// drawHealthBar mostra os pontos de vida como pequenos blocos.
func (g *Game) drawHealthBar(screen *ebiten.Image) {
	const pip = 8
	for i := 0; i < maxHealth; i++ {
		c := color.RGBA{0x40, 0x20, 0x20, 0xff}
		if i < g.player.health {
			c = color.RGBA{0xe0, 0x50, 0x50, 0xff}
		}
		vector.DrawFilledRect(screen, float32(60+i*(pip+2)), 4, pip, pip, c, false)
	}
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

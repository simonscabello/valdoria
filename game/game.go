package game

import (
	"image/color"

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

// gameMode distingue a campanha (fases + chefe) do modo de sobrevivência
// (ondas infinitas com dificuldade crescente).
type gameMode int

const (
	modeCampaign gameMode = iota
	modeSurvival
)

type state int

const (
	stateMenu state = iota
	stateControls
	stateOptions
	statePlaying
	statePaused
	stateBoss
	stateVictory
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

	clouds     []*bgItem
	hills      []*bgItem
	structures []*bgItem
	dust       []*star
	particles  []*particle
	popups     []*scorePopup

	shakeMag     float64
	shakeSetting shakeLevel
	damageFlash  int
	hitStop      int
	scene        *ebiten.Image

	level         *Level
	stages        []*stageDef
	stageIndex    int
	mode          gameMode
	survivalTimer int
	timeScale     int
	resumeState   state

	lives           int
	bombCharges     int
	bombEffectTimer int

	score        int
	highScore    int
	survivalBest int
	combo        int
	comboTimer   int

	save saveData

	corruption Corruption

	sectionDamaged bool
	lastSection    int
	formations     map[int]*formationTracker

	boss *Boss

	menuIndex int
	fade      int

	enemiesDefeated  int
	maxMult          int
	elapsedTicks     int
	victoryLifeBonus int
	victoryTime      int

	audio *AudioManager
}

const menuItemCount = 7
const optionsItemCount = 3
const pauseItemCount = 3

var gameOverOptions = []string{"Tentar novamente", "Voltar ao menu"}
var victoryOptions = []string{"Jogar novamente", "Voltar ao menu"}

type formationTracker struct {
	total   int
	killed  int
	escaped int
}

func New() *Game {
	g := &Game{}
	g.audio = getAudio()

	// Carrega recordes e preferências antes do reset (a dificuldade influencia
	// vidas/bombas iniciais).
	g.save = loadSave()
	setDifficulty(g.save.Difficulty)
	g.audio.muted = g.save.Muted
	g.audio.setMusicVolume(g.save.MusicVolume)
	g.audio.setSFXVolume(g.save.SFXVolume)

	g.reset()

	g.highScore = g.save.HighScore
	g.survivalBest = g.save.SurvivalBest
	g.shakeSetting = g.save.Shake

	g.initStars()
	g.enterMenu()
	return g
}

// persist sincroniza o estado persistível e grava o save (best-effort).
func (g *Game) persist() {
	g.save.HighScore = g.highScore
	g.save.SurvivalBest = g.survivalBest
	g.save.Difficulty = currentDifficulty
	g.save.Shake = g.shakeSetting
	g.save.Muted = g.audio.muted
	g.save.MusicVolume = g.audio.musicVol
	g.save.SFXVolume = g.audio.sfxVol
	g.save.VolumesSet = true
	g.save.write()
}

// reset cria uma sessão completamente nova, sem qualquer estado anterior.
func (g *Game) reset() {
	resetRNG()
	g.state = statePlaying
	g.player = newPlayer()
	g.bullets = g.bullets[:0]
	g.enemies = g.enemies[:0]
	g.enemyBullets = g.enemyBullets[:0]
	g.powerups = g.powerups[:0]
	g.stages = campaignStages()
	g.stageIndex = 0
	g.level = newLevelFromStage(g.stages[0])
	g.timeScale = dev.timeScale

	dp := diffParams()
	g.lives = dp.startLives
	g.bombCharges = dp.startBombs
	g.bombEffectTimer = 0
	g.score = 0
	g.combo = 0
	g.comboTimer = 0
	g.corruption.reset()
	g.sectionDamaged = false
	g.lastSection = g.level.section
	g.formations = map[int]*formationTracker{}
	g.boss = nil

	g.survivalTimer = 90
	g.enemiesDefeated = 0
	g.maxMult = 1
	g.elapsedTicks = 0
	g.victoryLifeBonus = 0
	g.victoryTime = 0

	g.particles = g.particles[:0]
	g.popups = g.popups[:0]
	g.shakeMag = 0
	g.damageFlash = 0
	g.hitStop = 0

	if dev.startBoss {
		g.startBoss()
	}
}

func (g *Game) enterMenu() {
	g.state = stateMenu
	g.menuIndex = 0
	g.fade = fadeDuration
}

func (g *Game) startNewGame() {
	g.reset()
	g.fade = fadeDuration
}

func (g *Game) startBoss() {
	g.enemies = g.enemies[:0]
	g.enemyBullets = g.enemyBullets[:0]
	// Reino caído = confronto verdadeiro. A aposta do jogador decide qual dos
	// dois finais ele vai ver.
	g.boss = newBoss(g.corruption.fallen())
	g.state = stateBoss
}

// advanceStage encadeia para a próxima fase da campanha ao concluir a atual sem
// chefe, mantendo pontuação/poder e anunciando a nova região.
func (g *Game) advanceStage() {
	g.stageIndex++
	if g.stageIndex >= len(g.stages) {
		g.startBoss() // salvaguarda: sem próxima fase, encerra no chefe
		return
	}
	g.enemyBullets = g.enemyBullets[:0]
	g.level = newLevelFromStage(g.stages[g.stageIndex])
	g.lastSection = g.level.section
	g.sectionDamaged = false
	g.formations = map[int]*formationTracker{}
	g.level.announce = "Nova regiao: " + g.stages[g.stageIndex].name
	g.level.announceTimer = announceDuration
}

func (g *Game) Update() error {
	if g.fade > 0 {
		g.fade--
	}
	g.updateStars()

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.audio.toggleMute()
		g.persist()
	}
	g.updateMusicForState()
	g.audio.update()
	g.handleDevToggles()

	switch g.state {
	case stateMenu:
		return g.updateMenu()
	case stateControls:
		g.updateControls()
		return nil
	case stateOptions:
		g.updateOptions()
		return nil
	case statePaused:
		g.updatePaused()
		return nil
	case stateGameOver:
		g.updateEndScreen(gameOverOptions)
		return nil
	case stateVictory:
		g.updateEndScreen(victoryOptions)
		return nil
	}

	g.updatePlay()
	return nil
}

// updateMusicForState escolhe a trilha conforme a tela atual. A pausa mantém a
// trilha do estado que estava em andamento; na campanha a música segue o bioma.
func (g *Game) updateMusicForState() {
	st := g.state
	if st == statePaused {
		st = g.resumeState
	}
	switch st {
	case statePlaying:
		if g.mode == modeSurvival {
			g.audio.playMusic(musicPhase)
			return
		}
		g.audio.playMusic(g.level.music())
	case stateBoss:
		g.audio.playMusic(musicBoss)
	case stateGameOver, stateVictory:
		g.audio.playMusic(musicNone)
	default:
		g.audio.playMusic(musicMenu)
	}
}

// updateMenu navega o menu inicial. Sair encerra o jogo sem erro.
func (g *Game) updateMenu() error {
	g.moveMenuSelection(menuItemCount)
	if !inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return nil
	}
	g.audio.playSFX(sfxMenu)
	switch g.menuIndex {
	case 0:
		g.mode = modeCampaign
		g.startNewGame()
	case 1:
		g.mode = modeSurvival
		g.startNewGame()
	case 2:
		g.cycleDifficulty()
	case 3:
		g.state = stateOptions
		g.menuIndex = 0
		g.fade = fadeDuration
	case 4:
		g.state = stateControls
		g.fade = fadeDuration
	case 5:
		g.cycleShake()
		g.persist()
	case 6:
		return ebiten.Termination
	}
	return nil
}

// cycleDifficulty alterna Fácil → Normal → Difícil e persiste a escolha.
func (g *Game) cycleDifficulty() {
	setDifficulty((currentDifficulty + 1) % difficultyCount)
	g.persist()
}

func (g *Game) updateControls() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.enterMenu()
	}
}

func (g *Game) updateOptions() {
	g.moveMenuSelection(optionsItemCount)
	g.handleVolumeKeys(0, 1)
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.enterMenu()
		return
	}
	if !inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return
	}
	g.audio.playSFX(sfxMenu)
	if g.menuIndex == 2 {
		g.enterMenu()
	}
}

func (g *Game) updatePaused() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = g.resumeState
		return
	}
	g.moveMenuSelection(pauseItemCount)
	g.handleVolumeKeys(1, 2)
	if !inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return
	}
	g.audio.playSFX(sfxMenu)
	if g.menuIndex == 0 {
		g.state = g.resumeState
	}
}

// handleVolumeKeys ajusta música/SFX com ←→ quando o índice atual está nas
// linhas musicIdx ou sfxIdx. Persiste a preferência ao mudar.
func (g *Game) handleVolumeKeys(musicIdx, sfxIdx int) {
	delta := 0.0
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		delta = -volumeStep
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
		delta = volumeStep
	}
	if delta == 0 {
		return
	}
	switch g.menuIndex {
	case musicIdx:
		g.audio.nudgeMusicVolume(delta)
	case sfxIdx:
		g.audio.nudgeSFXVolume(delta)
	default:
		return
	}
	g.audio.playSFX(sfxMenu)
	g.persist()
}

func (g *Game) musicVolumeLabel() string {
	return "Musica: " + itoa(volumePercent(g.audio.musicVol)) + "%"
}

func (g *Game) sfxVolumeLabel() string {
	return "Efeitos: " + itoa(volumePercent(g.audio.sfxVol)) + "%"
}

func (g *Game) updateEndScreen(options []string) {
	g.moveMenuSelection(len(options))
	if !inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return
	}
	g.audio.playSFX(sfxMenu)
	if g.menuIndex == 0 {
		g.startNewGame()
		return
	}
	g.enterMenu()
}

func (g *Game) moveMenuSelection(count int) {
	moved := false
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.menuIndex = (g.menuIndex - 1 + count) % count
		moved = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.menuIndex = (g.menuIndex + 1) % count
		moved = true
	}
	if moved {
		g.audio.playSFX(sfxMenu)
	}
}

func (g *Game) updatePlay() {
	g.handlePauseToggle()
	if g.state == statePaused {
		return
	}

	// Impact freeze: por poucos frames a lógica congela para dar peso aos golpes
	// fortes, mas os efeitos visuais seguem animando para não parecer travamento.
	if g.hitStop > 0 {
		g.hitStop--
		g.updateParticles()
		g.updatePopups()
		g.updateShake()
		if g.damageFlash > 0 {
			g.damageFlash--
		}
		return
	}

	g.elapsedTicks++
	g.player.update()
	if shots := g.player.tryShoot(); len(shots) > 0 {
		g.bullets = append(g.bullets, shots...)
		g.audio.playSFX(weaponShootSFX(g.player.weapon))
		if g.player.weapon == weaponFlame {
			// Brasas no bocal reforçam a sensação de fogo das Chamas do Dragão.
			g.spawnBurst(g.player.centerX(), g.player.y, 2, 1.4, 10, 2, false, flameColor)
		}
	}
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
			if g.mode == modeSurvival {
				g.spawnSurvival()
			} else {
				g.spawnFromLevel()
			}
		}
	}
	if g.state == stateBoss {
		g.updateBoss()
	}

	g.updateBullets()
	g.updateEnemies()
	g.updateEnemyBullets()
	g.updatePowerups()
	g.handleCollisions()
	g.collectPowerups()
	g.removeDead()

	g.corruption.update()
	g.updateParallax()
	g.updateParticles()
	g.updatePopups()
	g.updateShake()
	if g.damageFlash > 0 {
		g.damageFlash--
	}

	if g.player.health <= 0 {
		g.loseLife()
	}
	g.updateBestScore()
	if g.state == statePlaying && g.mode == modeCampaign {
		g.checkSectionBonus()
		if g.level.readyForBoss(len(g.enemies)) {
			if !g.sectionDamaged {
				g.score += waveNoDamageBonus
				g.spawnBonusPopup(ScreenWidth/2, 90, waveNoDamageBonus)
			}
			if g.level.hasBoss {
				g.startBoss()
			} else {
				g.advanceStage()
			}
		}
	}
	if g.state == stateBoss {
		g.checkBossDefeat()
	}
}

func (g *Game) updateBoss() {
	ctx := &bossContext{
		playerX: g.player.centerX(),
		playerY: g.player.centerY(),
		bullets: &g.enemyBullets,
		enemies: &g.enemies,
	}
	g.boss.update(ctx)

	if g.boss.justEntered {
		g.enemies = g.enemies[:0]
		g.enemyBullets = g.enemyBullets[:0]
		g.boss.justEntered = false
	}
	if g.boss.phaseChanged {
		g.enemyBullets = g.enemyBullets[:0]
		g.addShake(shakeBossMagnitude)
		g.hitStop = hitStopBig
		g.boss.phaseChanged = false
	}
	// Ao iniciar a morte, a arena é limpa e o jogador fica invulnerável durante
	// toda a sequência: nenhum projétil ou inimigo residual pode causar uma
	// derrota injusta no exato instante da vitória.
	if g.boss.justDied {
		g.enemies = g.enemies[:0]
		g.enemyBullets = g.enemyBullets[:0]
		if g.player.invincible < bossDeathDuration {
			g.player.invincible = bossDeathDuration
		}
		g.hitStop = hitStopBig
		g.boss.justDied = false
	}
	if g.boss.phase == bossDying {
		g.bossDeathEffects()
	}
}

// bossDeathEffects espalha explosões pelo corpo do chefe durante a morte.
func (g *Game) bossDeathEffects() {
	if g.boss.tick%6 != 0 {
		return
	}
	x := g.boss.x + fxRand.Float64()*g.boss.w
	y := g.boss.y + fxRand.Float64()*g.boss.h
	g.spawnBurst(x, y, bossDeathParticles/2, 2.6, 26, 3, true, color.RGBA{0xff, 0xc0, 0x50, 0xff})
	g.addShake(shakeBossMagnitude)
}

func (g *Game) checkBossDefeat() {
	if g.boss.phase != bossDead {
		return
	}
	g.victoryTime = g.elapsedTicks / 60
	g.victoryLifeBonus = g.lives * lifeBonus
	g.score += bossScore + g.victoryLifeBonus
	if g.score > g.highScore {
		g.highScore = g.score
	}
	g.state = stateVictory
	g.menuIndex = 0
	g.fade = fadeDuration
	g.audio.playSFX(sfxVictory)
	g.persist()
}

func (g *Game) spawnFromLevel() {
	for _, e := range g.level.update() {
		if e.formationID > 0 {
			g.registerFormationMember(e.formationID)
		}
		g.spawn(e)
	}
}

// spawn coloca um inimigo em jogo já marcado pelo estado do reino: em corrupção
// alta ele nasce mais resistente e, a partir do Cerco, corrompido.
//
// A corrupção é aplicada **no instante do spawn**, não continuamente: inimigos
// já em cena não mudam sob os pés do jogador.
func (g *Game) spawn(e *Enemy) {
	if g.corruption.mutates() {
		mutateEnemy(e)
	}
	if mul := g.corruption.healthMul(); mul > 1 {
		e.health = int(float64(e.health)*mul + 0.5)
	}
	g.enemies = append(g.enemies, e)
}

// updateBestScore mantém o recorde do modo atual atualizado em tempo real.
func (g *Game) updateBestScore() {
	if g.mode == modeSurvival {
		if g.score > g.survivalBest {
			g.survivalBest = g.score
		}
		return
	}
	if g.score > g.highScore {
		g.highScore = g.score
	}
}

// spawnSurvival gera inimigos aleatórios em ritmo e variedade crescentes ao
// longo do tempo — o modo não tem fim: dura enquanto o jogador sobreviver.
func (g *Game) spawnSurvival() {
	if g.survivalTimer > 0 {
		g.survivalTimer--
		return
	}
	minute := g.elapsedTicks / 3600
	interval := 42 - minute*6
	if interval < 16 {
		interval = 16
	}
	g.survivalTimer = interval
	g.spawnRandomEnemy(minute)
	// Chance de um segundo inimigo desde o início; a partir de 1 min, às vezes um terceiro.
	if rng.Intn(2) == 0 {
		g.spawnRandomEnemy(minute)
	}
	if minute >= 1 && rng.Intn(3) == 0 {
		g.spawnRandomEnemy(minute)
	}
}

// spawnRandomEnemy escolhe um inimigo de um repertório que cresce com o tempo.
func (g *Game) spawnRandomEnemy(minute int) {
	pool := []enemyKind{kindCrow, kindCrow, kindHarpy}
	if minute >= 1 {
		pool = append(pool, kindHarpy, kindGargoyle)
	}
	if minute >= 2 {
		pool = append(pool, kindWyvern, kindMage)
	}
	if minute >= 3 {
		pool = append(pool, kindBallista, kindWyvern, kindMage)
	}
	kind := pool[rng.Intn(len(pool))]
	g.spawn(spawnEnemy(kind, randX(enemySpawnSize(kind)), 0, rng.Intn(2) == 0))
}

// enemySpawnSize devolve a largura usada para posicionar o spawn na horizontal.
func enemySpawnSize(kind enemyKind) float64 {
	switch kind {
	case kindHarpy:
		return harpySize
	case kindGargoyle:
		return gargoyleSize
	case kindWyvern:
		return wyvernSize
	case kindBallista:
		return ballistaSize
	case kindMage:
		return mageSize
	default:
		return crowSize
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
		g.menuIndex = 0
		g.fade = fadeDuration
		g.audio.playSFX(sfxGameOver)
		g.persist()
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
		g.spawnBonusPopup(ScreenWidth/2, 90, waveNoDamageBonus)
	}
	g.sectionDamaged = false
	g.lastSection = g.level.section
}

func (g *Game) handleBomb() {
	if g.state != statePlaying && g.state != stateBoss {
		return
	}
	if g.bombCharges <= 0 {
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
			g.spawnExplosion(e.centerX(), e.centerY(), e.color)
		}
	}
	if g.boss != nil {
		g.boss.takeDamage(bombBossDamage) // dano relevante, mas não instantâneo
	}
	g.player.invincible = bombInvincibility
	g.bombEffectTimer = bombEffectDuration
	g.addShake(shakeBombMagnitude)
	g.audio.playSFX(sfxBomb)
}

func (g *Game) handlePauseToggle() {
	if !inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return
	}
	if g.state == statePaused {
		g.state = g.resumeState
		return
	}
	if g.state == statePlaying || g.state == stateBoss {
		g.resumeState = g.state
		g.menuIndex = 0
		g.state = statePaused
	}
}

func (g *Game) handleTimeScaleToggle() {
	if !dev.enabled || !inpututil.IsKeyJustPressed(ebiten.KeyTab) {
		return
	}
	if g.timeScale == dev.fastTimeScale {
		g.timeScale = dev.timeScale
		return
	}
	g.timeScale = dev.fastTimeScale
}

func (g *Game) updateBullets() {
	for _, b := range g.bullets {
		b.update()
		if b.trail {
			g.spawnTrail(b.x, b.y, b.color)
		}
		if b.offScreen() {
			b.dead = true
		}
	}
}

func (g *Game) spawnCrowFormation() {
	count := 3 + rng.Intn(2)
	const gap = crowSize + 6
	startX := rng.Float64() * (ScreenWidth - float64(count)*gap)
	for i := 0; i < count; i++ {
		g.enemies = append(g.enemies, newCrow(startX+float64(i)*gap))
	}
}

func (g *Game) spawnGargoyle() {
	fromLeft := rng.Intn(2) == 0
	targetX := 40 + rng.Float64()*(ScreenWidth-80-gargoyleSize)
	y := 30 + rng.Float64()*60
	g.enemies = append(g.enemies, newGargoyle(fromLeft, targetX, y))
}

func randX(size float64) float64 {
	return rng.Float64() * (ScreenWidth - size)
}

// handleDebugSpawns concentra as teclas do modo de desenvolvimento e só age
// quando ele está ativo, mantendo o build padrão livre dessas ferramentas.
func (g *Game) handleDebugSpawns() {
	if !dev.enabled {
		return
	}
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
	g.debugSpawnPowerup(ebiten.KeyZ, powerLight)
	g.debugSpawnPowerup(ebiten.KeyF, powerFire)
	g.debugSpawnPowerup(ebiten.KeyC, powerIce)
	g.debugSpawnPowerup(ebiten.KeyV, powerHeal)
	g.debugSpawnPowerup(ebiten.KeyB, powerShield)

	if inpututil.IsKeyJustPressed(ebiten.KeyK) && g.boss != nil {
		g.boss.takeDamage(devBossDamage)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyL) {
		g.clearEntities()
	}
}

// clearEntities remove inimigos e projéteis sem contá-los como abatidos.
func (g *Game) clearEntities() {
	g.enemies = g.enemies[:0]
	g.enemyBullets = g.enemyBullets[:0]
	g.bullets = g.bullets[:0]
}

// handleDevToggles alterna as sobreposições de diagnóstico do modo dev.
func (g *Game) handleDevToggles() {
	if !dev.enabled {
		return
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		dev.showHUD = !dev.showHUD
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		dev.showHitboxes = !dev.showHitboxes
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		dev.invincible = !dev.invincible
	}
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
			e.escapedBottom = e.crossedBottom()
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
	if g.boss != nil {
		g.bulletBossCollisions()
	}

	if g.player.canBeHit() {
		g.playerEnemyCollisions()
	}
	if g.player.canBeHit() {
		g.playerBulletCollisions()
	}
	if g.boss != nil && g.player.canBeHit() {
		g.playerBossCollision()
	}
}

func (g *Game) bulletBossCollisions() {
	boss := g.boss
	for _, b := range g.bullets {
		if b.dead {
			continue
		}
		if !collides(b.x, b.y, b.w, b.h, boss.x, boss.y, boss.w, boss.h) {
			continue
		}
		damage := b.damage
		if boss.hitCrystal(b.x, b.y, b.w, b.h) {
			damage *= crystalBonus
		}
		boss.takeDamage(damage)
		boss.applyWeaponHit(b.element)
		b.dead = true
	}
}

func (g *Game) playerBossCollision() {
	boss := g.boss
	if boss.phase == bossDying || boss.phase == bossDead {
		return
	}
	hx, hy, hw, hh := g.player.hitbox()
	if collides(hx, hy, hw, hh, boss.x, boss.y, boss.w, boss.h) {
		g.hitPlayer(bossContactDamage)
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
			e.applyWeaponHit(b.element)
			b.hitEnemies = append(b.hitEnemies, e)
			if e.dead {
				points := g.registerKill(e.score)
				g.spawnDrop(e)
				g.spawnExplosion(e.centerX(), e.centerY(), e.color)
				g.spawnScorePopup(e.centerX(), e.y, points)
				if e.kind == kindWyvern {
					g.hitStop = hitStopBig
				}
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
			g.spawnExplosion(e.centerX(), e.centerY(), e.color)
			g.hitPlayer(e.damage)
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
			g.hitPlayer(enemyBulletDamage)
			return
		}
	}
}

// onEnemyEscaped registra um inimigo que atravessou a base do reino. É o único
// lugar em que a corrupção sobe — e o coração do diferencial do jogo.
func (g *Game) onEnemyEscaped(e *Enemy) {
	if !g.corruption.add(e.kind) {
		g.audio.playSFX(sfxEscape)
		return
	}
	// A faixa mudou: o reino piorou de nível. O momento precisa ser sentido.
	g.audio.playSFX(sfxCorruption)
	g.addShake(shakeCorruptionMagnitude)
	g.spawnBurst(e.centerX(), ScreenHeight-4, 14, 2.2, 30, 3, false, g.corruption.barColor())
}

// registerKill contabiliza uma eliminação por disparo, aplica o multiplicador de
// combo e a bonificação da corrupção, e devolve os pontos concedidos.
func (g *Game) registerKill(base int) int {
	g.combo++
	g.comboTimer = comboWindow
	mult := g.multiplier()
	points := int(float64(base*mult) * g.corruption.scoreMul())
	g.score += points
	if mult > g.maxMult {
		g.maxMult = mult
	}
	return points
}

func (g *Game) multiplier() int {
	m := 1 + g.combo/comboStep
	if m > maxMultiplier {
		return maxMultiplier
	}
	return m
}

// hitPlayer centraliza a aplicação de dano ao jogador nas três origens
// (inimigo, projétil, chefe). Distingue o dano real da absorção pelo escudo,
// garantindo o feedback correto em cada caso.
func (g *Game) hitPlayer(damage int) {
	hadShield := g.player.hasShield()
	if g.player.hit(damage) {
		g.onPlayerDamaged()
		return
	}
	if hadShield {
		g.onShieldBroken()
	}
}

func (g *Game) onPlayerDamaged() {
	g.combo = 0
	g.comboTimer = 0
	g.sectionDamaged = true
	g.triggerDamageFlash()
	g.addShake(shakeHitMagnitude)
	g.audio.playSFX(sfxPlayerHit)
	if g.player != nil {
		g.spawnBurst(g.player.centerX(), g.player.centerY(), 6, 2.0, 18, 2, false, color.RGBA{0xff, 0x60, 0x60, 0xff})
	}
}

// onShieldBroken dá feedback claro de que o escudo absorveu um golpe (som,
// anel de partículas e leve vibração), sem contar como dano recebido.
func (g *Game) onShieldBroken() {
	g.addShake(shakeHitMagnitude)
	g.audio.playSFX(sfxShield)
	if g.player != nil {
		g.spawnCollectRing(g.player.centerX(), g.player.centerY(), shieldRing)
	}
}

func (g *Game) spawnDrop(e *Enemy) {
	if e.hasDrop {
		g.powerups = append(g.powerups, newPowerup(e.drop, e.centerX(), e.centerY()))
		return
	}
	if rng.Float64() < scaledDropChance(dropChance)*g.corruption.dropMul() {
		g.powerups = append(g.powerups, newPowerup(g.randomRune(), e.centerX(), e.centerY()))
	}
}

// randomRune sorteia o conteúdo de um drop. No teto de poder as runas
// elementais viram utilidade (cura/escudo): antes, um jogador no nível máximo
// via dezenas de runas caírem sem efeito nenhum pelo resto da partida.
func (g *Game) randomRune() powerupType {
	if g.player != nil && g.player.weaponLevel >= maxWeaponLevel {
		if rng.Intn(4) == 0 {
			return randomElementRune() // ainda permite trocar de elemento
		}
		if rng.Intn(2) == 0 {
			return powerHeal
		}
		return powerShield
	}
	switch rng.Intn(5) {
	case 3:
		return powerHeal
	case 4:
		return powerShield
	default:
		return randomElementRune()
	}
}

func randomElementRune() powerupType {
	switch rng.Intn(3) {
	case 0:
		return powerFire
	case 1:
		return powerIce
	default:
		return powerLight
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
			g.spawnCollectRing(p.x+powerupSize/2, p.y+powerupSize/2, p.color())
			g.audio.playSFX(sfxPickup)
		}
	}
}

func (g *Game) removeDead() {
	for _, e := range g.enemies {
		if !e.dead {
			continue
		}
		if !e.escaped {
			g.enemiesDefeated++
		}
		if e.escapedBottom {
			g.onEnemyEscaped(e)
		}
		if e.formationID > 0 {
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
		g.spawnBonusPopup(e.centerX(), e.centerY(), formationBonus)
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
	switch g.state {
	case stateMenu:
		g.drawMenu(screen)
	case stateControls:
		g.drawControls(screen)
	case stateOptions:
		g.drawOptionsScreen(screen)
	case stateGameOver:
		g.drawGameOver(screen)
	case stateVictory:
		g.drawVictory(screen)
	default:
		g.drawWorld(screen)
	}
	g.drawFade(screen)
}

// drawWorld renderiza a cena em uma imagem própria, aplica a vibração ao
// compô-la na tela e desenha o HUD por cima (sem tremer), preservando a leitura.
func (g *Game) drawWorld(screen *ebiten.Image) {
	if g.scene == nil {
		g.scene = ebiten.NewImage(ScreenWidth, ScreenHeight)
	}
	g.scene.Clear()
	g.drawScene(g.scene)

	screen.Fill(g.level.background())
	ox, oy := g.shakeOffset()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(ox, oy)
	screen.DrawImage(g.scene, op)

	if g.damageFlash > 0 {
		alpha := uint8(damageFlashAlpha * g.damageFlash / damageFlashDuration)
		vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0xff, 0x20, 0x20, alpha}, false)
	}

	// Escurecimento breve no disparo da Invocação Ancestral: dá peso ao golpe,
	// como se o dragão ancestral cobrisse o céu por um instante.
	if elapsed := bombEffectDuration - g.bombEffectTimer; g.bombEffectTimer > 0 && elapsed < bombDarkenFrames {
		alpha := uint8(bombDarkenAlpha * (bombDarkenFrames - elapsed) / bombDarkenFrames)
		vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0x10, 0x08, 0x18, alpha}, false)
	}

	g.drawOverlayHUD(screen)
}

// drawScene desenha o cenário e as entidades. Projéteis são desenhados por
// último para permanecerem visíveis acima das partículas e efeitos.
func (g *Game) drawScene(dst *ebiten.Image) {
	// O cenário perde cor conforme o reino se corrompe: o estado da partida
	// precisa ser legível sem olhar o HUD.
	theme := g.level.theme().corrupted(g.corruption.ratio())
	dst.Fill(theme.sky)
	g.drawParallax(dst, theme)

	// Escurece o cenário para que os elementos de jogo (desenhados depois)
	// saltem à frente, resolvendo a confusão entre fundo e inimigos.
	vector.DrawFilledRect(dst, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0x06, 0x05, 0x0e, sceneDimAlpha}, false)

	g.drawParticles(dst)

	for _, p := range g.powerups {
		p.draw(dst)
	}
	for _, e := range g.enemies {
		e.draw(dst)
	}
	if g.boss != nil {
		g.boss.draw(dst)
	}
	g.player.draw(dst)

	for _, b := range g.enemyBullets {
		b.draw(dst)
	}
	for _, b := range g.bullets {
		b.draw(dst)
	}

	if g.bombEffectTimer > 0 {
		g.drawBombEffect(dst)
	}
	g.drawCorruptionOverlay(dst)
	g.drawEscapeWarnings(dst)

	g.drawPopups(dst)

	if dev.showHitboxes {
		g.drawHitboxes(dst)
	}
}

// drawCorruptionOverlay tinge a cena e faz veias crescerem das bordas conforme
// o reino apodrece. É desenhado depois do cenário e antes do HUD.
func (g *Game) drawCorruptionOverlay(dst *ebiten.Image) {
	r := g.corruption.ratio()
	if r <= 0 {
		return
	}
	if tint := g.corruption.params().tint; tint.A > 0 {
		vector.DrawFilledRect(dst, 0, 0, ScreenWidth, ScreenHeight, tint, false)
	}

	// Veias magenta entrando pelas laterais: a corrupção invadindo o quadro.
	depth := float32(r * 26)
	if depth < 1 {
		return
	}
	vein := withAlpha(corruptedAccent, uint8(40+70*r))
	for i := 0; i < 5; i++ {
		y := float32(20 + i*62)
		h := float32(2)
		vector.DrawFilledRect(dst, 0, y, depth*float32(1+i%2), h, vein, false)
		vector.DrawFilledRect(dst, ScreenWidth-depth*float32(2-i%2), y+28, depth*float32(2-i%2), h, vein, false)
	}
}

// drawEscapeWarnings destaca quem está prestes a cruzar a base do reino.
//
// Sem este aviso a corrupção subiria "sozinha" aos olhos do jogador. Com ele, a
// fuga vira uma decisão consciente: dá tempo de perseguir, ou de deixar passar
// de propósito e aceitar a aposta.
func (g *Game) drawEscapeWarnings(dst *ebiten.Image) {
	warned := false
	for _, e := range g.enemies {
		if e.dead || !e.aboutToEscape() || corruptionWeight(e.kind) <= 0 {
			continue
		}
		warned = true
		if (e.animTick/4)%2 == 0 {
			vector.StrokeRect(dst,
				float32(e.x-3), float32(e.y-3), float32(e.w+6), float32(e.h+6),
				1, corruptedAccent, false)
		}
		// Seta apontando para baixo: "este vai passar".
		vector.DrawFilledRect(dst, float32(e.centerX()-1), float32(e.y+e.h+3), 2, 4, corruptedAccent, false)
	}
	if warned {
		// Filete pulsante na base: a linha que não deveria ser cruzada.
		a := uint8(80 + 60*((g.elapsedTicks/6)%2))
		vector.DrawFilledRect(dst, 0, ScreenHeight-1, ScreenWidth, 1, withAlpha(corruptedAccent, a), false)
	}
}

// drawHitboxes desenha os contornos de colisão para depuração visual.
func (g *Game) drawHitboxes(dst *ebiten.Image) {
	green := color.RGBA{0x40, 0xff, 0x40, 0xff}
	yellow := color.RGBA{0xff, 0xff, 0x40, 0xff}

	hx, hy, hw, hh := g.player.hitbox()
	vector.StrokeRect(dst, float32(hx), float32(hy), float32(hw), float32(hh), 1, green, false)
	for _, e := range g.enemies {
		vector.StrokeRect(dst, float32(e.x), float32(e.y), float32(e.w), float32(e.h), 1, green, false)
	}
	for _, b := range g.bullets {
		vector.StrokeRect(dst, float32(b.x), float32(b.y), float32(b.w), float32(b.h), 1, green, false)
	}
	for _, b := range g.enemyBullets {
		vector.StrokeRect(dst, float32(b.x), float32(b.y), enemyBulletSize, enemyBulletSize, 1, yellow, false)
	}
	if g.boss != nil {
		vector.StrokeRect(dst, float32(g.boss.x), float32(g.boss.y), float32(g.boss.w), float32(g.boss.h), 1, green, false)
	}
}

func (g *Game) drawOverlayHUD(screen *ebiten.Image) {
	g.drawHUD(screen)
	if g.boss != nil {
		g.drawBossHUD(screen)
	}
	if dev.enabled && dev.showHUD {
		g.drawDevHUD(screen)
	}

	if g.level.announceTimer > 0 && g.level.announce != "" {
		label(screen, g.level.announce, ScreenWidth/2-len(g.level.announce)*3, 70)
	}
	g.drawCorruptionAnnounce(screen)

	if g.state == statePaused {
		const pw, ph = 140, 92
		px := float32(ScreenWidth/2 - pw/2)
		py := float32(ScreenHeight/2 - ph/2)
		drawUIPanel(screen, px, py, pw, ph)
		drawTextCentered(screen, "PAUSADO", float64(ScreenHeight/2-38), uiHighlight)
		options := []string{
			"Continuar",
			g.musicVolumeLabel(),
			g.sfxVolumeLabel(),
		}
		g.drawOptions(screen, options, ScreenHeight/2-20)
		drawTextCentered(screen, "A/D volume  Esc volta", float64(ScreenHeight/2+36), uiInkDim)
	}
}

// drawDevHUD mostra FPS, contagens, tempo da fase e a semente em uso.
func (g *Game) drawDevHUD(screen *ebiten.Image) {
	lines := []string{
		"FPS " + itoa(int(ebiten.ActualFPS()+0.5)),
		"INIM " + itoa(len(g.enemies)),
		"PROJ " + itoa(len(g.bullets)+len(g.enemyBullets)),
		"PART " + itoa(len(g.particles)),
		"TEMPO " + itoa(g.elapsedTicks/60) + "s",
		"TS x" + itoa(g.timeScale),
		"SEED " + itoa(int(CurrentSeed())),
	}
	if dev.invincible {
		lines = append(lines, "INVENCIVEL")
	}
	for i, l := range lines {
		ebitenutil.DebugPrintAt(screen, l, 4, 44+i*12)
	}
}

func (g *Game) drawStarsBackground(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x0a, 0x0a, 0x1e, 0xff})
	g.drawStars(screen)
}

// Textos do menu dimensionados para caber em ScreenWidth (240) com margem.
const (
	menuTitle    = "ASAS DE VALDORIA"
	menuSubtitle = "Shoot 'em up medieval"
	menuHelp     = "Setas/WASD  Enter confirma"
)

func (g *Game) drawMenu(screen *ebiten.Image) {
	g.drawStarsBackground(screen)

	// Painel só atrás do bloco título+opções; textos longos não podem
	// ultrapassar a largura útil (causa o "texto atravessando" a moldura).
	const margin = 12
	drawUIPanel(screen, margin, 40, ScreenWidth-2*margin, 200)

	drawCentered(screen, menuTitle, 52)
	drawCentered(screen, menuSubtitle, 70)
	options := []string{
		"Iniciar jogo",
		"Sobrevivencia",
		"Dificuldade: " + difficultyName(currentDifficulty),
		"Opcoes",
		"Controles",
		"Vibracao: " + shakeLevelName(g.shakeSetting),
		"Sair",
	}
	g.drawOptions(screen, options, 96)

	if g.highScore > 0 || g.survivalBest > 0 {
		drawCentered(screen, "Recorde "+itoa(g.highScore)+"  Sobrev. "+itoa(g.survivalBest), 258)
	}
	drawCentered(screen, menuHelp, ScreenHeight-18)
	ver := "v" + Version
	drawText(screen, ver, ScreenWidth-textWidth(ver)-4, ScreenHeight-18, uiInkDim)
}

func (g *Game) drawControls(screen *ebiten.Image) {
	g.drawStarsBackground(screen)
	const margin = 12
	drawUIPanel(screen, margin, 24, ScreenWidth-2*margin, 250)
	drawCentered(screen, "CONTROLES", 36)
	lines := []string{
		"Setas/WASD: mover",
		"Shift: precisao",
		"Espaco: disparar",
		"X/Ctrl: invocacao",
		"Escape: pausar",
		"M: mudo",
		"Enter: confirmar",
	}
	for i, l := range lines {
		drawTextCentered(screen, l, float64(66+i*20), uiInk)
	}
	drawCentered(screen, "Enter/Esc volta", ScreenHeight-28)
}

func (g *Game) drawOptionsScreen(screen *ebiten.Image) {
	g.drawStarsBackground(screen)
	const margin = 12
	drawUIPanel(screen, margin, 40, ScreenWidth-2*margin, 200)
	drawCentered(screen, "OPCOES", 56)
	options := []string{
		g.musicVolumeLabel(),
		g.sfxVolumeLabel(),
		"Voltar",
	}
	g.drawOptions(screen, options, 100)
	drawCentered(screen, "A/D ajusta volume", ScreenHeight-36)
	drawCentered(screen, "Enter/Esc volta", ScreenHeight-22)
}

// reachedLabel descreve até onde o jogador avançou ao perder.
func (g *Game) reachedLabel() string {
	if g.mode == modeSurvival {
		return "Sobrevivencia - " + itoa(g.elapsedTicks/60) + "s"
	}
	if g.boss != nil {
		return "Chefe: " + g.boss.name()
	}
	if g.stageIndex < len(g.stages) {
		return g.stages[g.stageIndex].name
	}
	return g.level.sectionName()
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	g.drawStarsBackground(screen)
	const margin = 12
	drawUIPanel(screen, margin, 28, ScreenWidth-2*margin, 230)
	drawCentered(screen, "GAME OVER", 42)
	stats := []string{
		"Pontos: " + itoa(g.score),
		"Trecho: " + g.reachedLabel(),
		"Abates: " + itoa(g.enemiesDefeated),
		"Max combo: x" + itoa(g.maxMult),
	}
	stats = append(stats,
		"Corrupcao: "+itoa(g.corruption.percent())+"%",
		"Fugas: "+itoa(g.corruption.escaped))
	for i, s := range stats {
		drawTextCentered(screen, s, float64(72+i*18), uiInk)
	}
	g.drawOptions(screen, gameOverOptions, 168)
}

func (g *Game) drawVictory(screen *ebiten.Image) {
	g.drawStarsBackground(screen)
	const margin = 12
	drawUIPanel(screen, margin, 22, ScreenWidth-2*margin, 240)
	title := "FASE CONCLUIDA!"
	if g.boss != nil && g.boss.ascended {
		title = "VALDORIA CAIU"
	}
	drawCentered(screen, title, 36)
	stats := []string{
		"Pontos: " + itoa(g.score),
		"Bonus vidas: " + itoa(g.victoryLifeBonus),
		"Tempo: " + itoa(g.victoryTime) + "s",
		"Abates: " + itoa(g.enemiesDefeated),
		"Corrupcao: " + itoa(g.corruption.percent()) + "% " + g.corruption.tierName(),
	}
	for i, s := range stats {
		drawTextCentered(screen, s, float64(62+i*18), uiInk)
	}
	g.drawOptions(screen, victoryOptions, 172)
}

func (g *Game) drawOptions(screen *ebiten.Image, options []string, y int) {
	for i, opt := range options {
		col := uiInk
		text := "  " + opt
		if i == g.menuIndex {
			col = uiHighlight
			text = "> " + opt
		}
		drawTextCentered(screen, text, float64(y+i*16), col)
	}
}

func (g *Game) drawFade(screen *ebiten.Image) {
	if g.fade <= 0 {
		return
	}
	alpha := uint8(255 * g.fade / fadeDuration)
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, alpha}, false)
}

func drawCentered(screen *ebiten.Image, text string, y int) {
	drawTextCentered(screen, text, float64(y), uiInk)
}

// label desenha um texto do HUD numa placa leve, com padding confortável.
func label(screen *ebiten.Image, text string, x, y int) {
	labelBlock(screen, []string{text}, x, y)
}

// labelBlock agrupa várias linhas numa única placa (menos caixinhas soltas).
func labelBlock(screen *ebiten.Image, lines []string, x, y int) {
	if len(lines) == 0 {
		return
	}
	const padX, padY, lineH = 5, 3, 12
	maxW := 0.0
	for _, s := range lines {
		if w := textWidth(s); w > maxW {
			maxW = w
		}
	}
	w := float32(maxW) + padX*2
	h := float32(len(lines)*lineH + padY*2 - 2)
	drawUIChip(screen, float32(x-padX), float32(y-padY), w, h)
	for i, s := range lines {
		drawText(screen, s, float64(x), float64(y+i*lineH), uiInk)
	}
}

// labelRight alinha o texto de modo que sua borda direita fique em rightX.
func labelRight(screen *ebiten.Image, text string, rightX, y int) {
	label(screen, text, rightX-int(textWidth(text)), y)
}

func labelBlockRight(screen *ebiten.Image, lines []string, rightX, y int) {
	if len(lines) == 0 {
		return
	}
	maxW := 0.0
	for _, s := range lines {
		if w := textWidth(s); w > maxW {
			maxW = w
		}
	}
	labelBlock(screen, lines, rightX-int(maxW), y)
}

func (g *Game) drawBossHUD(screen *ebiten.Image) {
	b := g.boss
	if b.phase == bossDead {
		return
	}
	if b.phase == bossEntering {
		// Aviso de entrada pulsante para leitura clara do início do combate.
		if (b.tick/16)%2 == 0 {
			drawCentered(screen, "! ALERTA !", ScreenHeight/2-14)
		}
		drawTextCentered(screen, b.name(), ScreenHeight/2, uiHighlight)
		return
	}

	// Barra fina no topo absoluto, para não colidir com o placar logo abaixo.
	const barY = 2
	const margin = 6
	const barH = 5
	width := float32(ScreenWidth - 2*margin)

	fill := color.RGBA{0xc0, 0x30, 0x30, 0xff}
	switch b.phase {
	case bossPhase2:
		fill = color.RGBA{0xd0, 0x70, 0x20, 0xff}
	case bossPhase3:
		fill = color.RGBA{0xff, 0x40, 0x40, 0xff}
	}
	drawIronBar(screen, margin, barY, width, barH, fill, b.healthRatio())

	// Marcas nos limiares de troca de fase (65% e 30%), para o jogador antecipar.
	tick := color.RGBA{0x20, 0x10, 0x10, 0xff}
	for _, r := range []float64{0.65, 0.30} {
		x := float32(margin) + width*float32(r)
		vector.DrawFilledRect(screen, x, barY, 1, barH, tick, false)
	}

	if b.action == actionWarning {
		drawText(screen, "!", b.centerX()-2, b.y+b.h+4, uiHighlight)
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
	best := g.highScore
	if g.mode == modeSurvival {
		best = g.survivalBest
	} else {
		g.drawProgressBar(screen)
	}

	// Placar: uma placa só (score + recorde + combo).
	topLeft := []string{
		"SCORE " + itoa(g.score),
		"HIGH " + itoa(best),
	}
	if g.multiplier() > 1 {
		topLeft = append(topLeft, "x"+itoa(g.multiplier()))
	}
	labelBlock(screen, topLeft, 6, 10)

	name := g.level.sectionName()
	if g.mode == modeSurvival {
		name = "SOBREV " + itoa(g.elapsedTicks/60) + "s"
	}
	labelBlockRight(screen,
		[]string{name, "CORR " + itoa(g.corruption.percent()) + "%"},
		ScreenWidth-6, 10)

	g.drawCorruptionBar(screen)
	g.drawHealthBar(screen)
	labelBlock(screen, []string{
		"VIDAS " + itoa(g.lives),
		"BOMBA " + itoa(g.bombCharges),
	}, 6, ScreenHeight-28)

	bottomRight := []string{}
	if g.player.hasShield() {
		bottomRight = append(bottomRight, "ESCUDO")
	}
	bottomRight = append(bottomRight, weaponName(g.player.weapon)+" Nv"+itoa(g.player.weaponLevel))
	top := ScreenHeight - 16 - 12*(len(bottomRight)-1)
	labelBlockRight(screen, bottomRight, ScreenWidth-6, top)
	g.drawWeaponCharge(screen)
}

// drawWeaponCharge mostra o progresso rumo ao próximo nível de arma como
// marcas acima do nome da magia. Sem isso, exigir duas runas por nível tiraria
// recompensa do jogador em vez de espalhá-la pela partida.
func (g *Game) drawWeaponCharge(screen *ebiten.Image) {
	if g.player.weaponLevel >= maxWeaponLevel {
		return
	}
	const pip, gap = 5, 2
	y := float32(ScreenHeight - 24)
	width := float32(runesPerLevel*(pip+gap) - gap)
	x0 := float32(ScreenWidth-6) - width
	for i := 0; i < runesPerLevel; i++ {
		c := color.RGBA{0x2a, 0x24, 0x18, 0xc0}
		if i < g.player.weaponCharge {
			c = weaponColor(g.player.weapon)
		}
		x := x0 + float32(i*(pip+gap))
		vector.DrawFilledRect(screen, x, y, pip, 2, c, false)
	}
}

// drawProgressBar desenha a barra de avanço da fase e marca o início de cada
// trecho, para o jogador perceber quanto falta e quando muda o cenário.
func (g *Game) drawProgressBar(screen *ebiten.Image) {
	const barH = 4
	drawIronBar(screen, 0, 0, ScreenWidth, barH, color.RGBA{0xd0, 0xc0, 0x50, 0xff}, g.level.progress())

	if g.level.duration <= 0 {
		return
	}
	for i, s := range g.level.sections {
		if i == 0 {
			continue
		}
		x := float32(float64(s.startTick) / float64(g.level.duration) * ScreenWidth)
		vector.DrawFilledRect(screen, x, 0, 1, barH, color.RGBA{0x10, 0x10, 0x18, 0xff}, false)
	}
}

// drawCorruptionBar desenha o Medidor de Corrupção como uma coluna vertical na
// borda direita: o reino "enchendo" de baixo para cima.
//
// A posição é deliberada. A corrupção é a mecânica central do jogo e precisa
// estar sempre visível, mas não pode competir com o campo de tiro — uma coluna
// de 4px na borda fica no campo periférico, onde o olho registra a mudança sem
// desviar da ação.
func (g *Game) drawCorruptionBar(screen *ebiten.Image) {
	const w = 4
	const top = 44
	bottom := float32(ScreenHeight - 52)
	x := float32(ScreenWidth - w - 1)
	height := bottom - top

	vector.DrawFilledRect(screen, x, top, w, height, color.RGBA{0x0c, 0x08, 0x14, 0xc0}, false)
	vector.StrokeRect(screen, x, top, w, height, 1, withAlpha(uiChipEdge, 110), false)

	if r := g.corruption.ratio(); r > 0 {
		fillH := height * float32(r)
		c := g.corruption.barColor()
		if g.corruption.pulse > 0 { // clarão curto a cada fuga
			c = brighten(c, 70)
		}
		vector.DrawFilledRect(screen, x+1, bottom-fillH, w-2, fillH, c, false)
	}

	// Marcas das faixas, para o jogador antecipar quando o mundo vai piorar.
	tick := color.RGBA{0xf0, 0xe8, 0xcf, 0x70}
	for t := tierShadow; t < corruptionTierCount; t++ {
		y := bottom - height*float32(corruptionTable[t].min/maxCorruption)
		vector.DrawFilledRect(screen, x, y, w, 1, tick, false)
	}
}

// drawCorruptionAnnounce anuncia a mudança de faixa no centro da tela. É o
// momento em que o jogador entende que o mundo mudou por causa dele.
func (g *Game) drawCorruptionAnnounce(screen *ebiten.Image) {
	c := &g.corruption
	if c.announceTimer <= 0 || c.announce == "" {
		return
	}
	if (c.announceTimer/10)%2 == 0 {
		drawTextCentered(screen, "O REINO SE CORROMPE", ScreenHeight/2-24, uiInkDim)
	}
	drawTextCentered(screen, c.announce, ScreenHeight/2-10, c.barColor())
}

// drawHealthBar mostra os pontos de vida como pequenos blocos, no rodapé à
// esquerda (acima de VIDAS), sem colidir com o placar do topo.
func (g *Game) drawHealthBar(screen *ebiten.Image) {
	const pip = 7
	y := float32(ScreenHeight - 44)
	for i := 0; i < maxHealth; i++ {
		c := color.RGBA{0x3a, 0x18, 0x20, 0xc0}
		if i < g.player.health {
			c = color.RGBA{0xe8, 0x48, 0x48, 0xff}
		}
		x := float32(6 + i*(pip+3))
		vector.DrawFilledRect(screen, x, y, pip, pip, c, false)
		vector.StrokeRect(screen, x, y, pip, pip, 1, withAlpha(uiChipEdge, 160), false)
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

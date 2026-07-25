package game

import (
	"image/color"
	"math"
)

type formation int

const (
	formationSingle formation = iota
	formationLine
	formationV
)

// waveEvent descreve uma aparição agendada na linha do tempo da fase.
type waveEvent struct {
	startTick int
	kind      enemyKind
	count     int
	interval  int
	formation formation
	spawnX    float64
	fromLeft  bool
	hasDrop   bool
	drop      powerupType

	started       bool
	spawned       int
	nextSpawnTick int
	done          bool
}

func (ev *waveEvent) spawn() []*Enemy {
	// Jitter no âncora da onda: em formações o grupo inteiro se desloca junto.
	x := jitterSpawnX(ev.spawnX)
	switch ev.formation {
	case formationLine:
		out := make([]*Enemy, 0, lineFormationCount)
		for i := 0; i < lineFormationCount; i++ {
			out = append(out, spawnEnemy(ev.kind, x+float64(i)*formationGapX, 0, ev.fromLeft))
		}
		return out
	case formationV:
		out := make([]*Enemy, 0, vFormationCount)
		center := vFormationCount / 2
		for i := 0; i < vFormationCount; i++ {
			depth := i - center
			if depth < 0 {
				depth = -depth
			}
			px := x + float64(i-center)*formationGapX
			out = append(out, spawnEnemy(ev.kind, px, -float64(depth)*formationGapY, ev.fromLeft))
		}
		return out
	default:
		return []*Enemy{spawnEnemy(ev.kind, x, 0, ev.fromLeft)}
	}
}

func spawnEnemy(kind enemyKind, x, yOffset float64, fromLeft bool) *Enemy {
	switch kind {
	case kindHarpy:
		e := newHarpy(fairSpawnX(x, harpySize))
		e.y += yOffset
		return e
	case kindGargoyle:
		// Gárgulas entram pela lateral de propósito — sem margem "justa".
		return newGargoyle(fromLeft, x, 30)
	case kindWyvern:
		e := newWyvern(fairSpawnX(x, wyvernSize))
		e.y += yOffset
		return e
	case kindBallista:
		e := newBallista(fairSpawnX(x, ballistaSize))
		e.y += yOffset
		return e
	case kindMage:
		e := newMage(fairSpawnX(x, mageSize))
		e.y += yOffset
		return e
	default:
		e := newCrow(fairSpawnX(x, crowSize))
		e.y += yOffset
		return e
	}
}

// fairSpawnX mantém o spawn longe das bordas, para o jogador conseguir alcançar.
func fairSpawnX(x, size float64) float64 {
	x = clampSpawnX(x, size)
	min := float64(spawnFairMargin)
	max := ScreenWidth - size - float64(spawnFairMargin)
	if max < min {
		return x
	}
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func clampSpawnX(x, size float64) float64 {
	if x < 0 {
		return 0
	}
	if x > ScreenWidth-size {
		return ScreenWidth - size
	}
	return x
}

// sectionDef descreve um trecho da fase: quando começa, seu nome, o aviso
// exibido, o tema visual do bioma e a trilha de fundo.
type sectionDef struct {
	startTick int
	name      string
	warning   string
	theme     bgTheme
	music     musicID
}

// waveDef é a descrição declarativa (dados) de uma aparição agendada.
type waveDef struct {
	startTick int
	kind      enemyKind
	count     int
	interval  int
	formation formation
	spawnX    float64
	fromLeft  bool
	hasDrop   bool
	drop      powerupType
}

// stageDef é uma fase inteira descrita como dados: trechos, ondas e se termina
// em um chefe. Novas fases são acrescentadas escrevendo dados, não lógica.
type stageDef struct {
	name     string
	sections []sectionDef
	waves    []waveDef
	hasBoss  bool
}

// endBuffer devolve o instante em que a fase pode preparar o chefe/transição:
// um respiro após a última aparição agendada.
const stageEndBuffer = 400

func (d *stageDef) duration() int {
	max := 0
	for _, w := range d.waves {
		if w.startTick > max {
			max = w.startTick
		}
	}
	return max + stageEndBuffer
}

// stageBuilder monta um stageDef de forma legível, no mesmo estilo declarativo
// das ondas originais (dados, não código de controle).
type stageBuilder struct{ d stageDef }

func newStageBuilder(name string) *stageBuilder { return &stageBuilder{d: stageDef{name: name}} }

func (b *stageBuilder) section(start int, name, warning string, th bgTheme, music musicID) {
	b.d.sections = append(b.d.sections, sectionDef{
		startTick: start, name: name, warning: warning, theme: th, music: music,
	})
}

func (b *stageBuilder) wave(start int, kind enemyKind, count, interval int, f formation, x float64, fromLeft bool) {
	b.d.waves = append(b.d.waves, waveDef{startTick: start, kind: kind, count: count, interval: interval, formation: f, spawnX: x, fromLeft: fromLeft})
}

// drop marca a última onda adicionada como portadora de um power-up garantido.
func (b *stageBuilder) drop(kind powerupType) {
	w := &b.d.waves[len(b.d.waves)-1]
	w.hasDrop = true
	w.drop = kind
}

func (b *stageBuilder) def() *stageDef { return &b.d }

type Level struct {
	events          []*waveEvent
	sections        []sectionDef
	duration        int
	hasBoss         bool
	tick            int
	section         int
	announce        string
	announceTimer   int
	nextFormationID int
}

// scaleWaveCount aumenta a quantidade de aparições da onda (~densidade da campanha).
func scaleWaveCount(count int) int {
	if count <= 1 {
		return count
	}
	scaled := int(math.Round(float64(count) * campaignDensityScale))
	if scaled < count {
		return count
	}
	return scaled
}

// scaleWaveInterval encurta o intervalo entre aparições, sem acelerar abaixo de
// ~2/3 do valor original (evita rajadas impossíveis).
func scaleWaveInterval(interval int) int {
	if interval <= 0 {
		return interval
	}
	scaled := int(math.Round(float64(interval) / campaignDensityScale))
	floor := (interval * 2) / 3
	if floor < 1 {
		floor = 1
	}
	if scaled < floor {
		return floor
	}
	if scaled < 1 {
		return 1
	}
	return scaled
}

// newLevelFromStage constrói o estado de execução de uma fase a partir da sua
// descrição declarativa.
func newLevelFromStage(def *stageDef) *Level {
	l := &Level{
		sections: def.sections,
		duration: def.duration(),
		hasBoss:  def.hasBoss,
	}
	for _, wd := range def.waves {
		l.events = append(l.events, &waveEvent{
			startTick: wd.startTick, kind: wd.kind,
			count:     scaleWaveCount(wd.count),
			interval:  scaleWaveInterval(wd.interval),
			formation: wd.formation, spawnX: wd.spawnX,
			fromLeft: wd.fromLeft, hasDrop: wd.hasDrop, drop: wd.drop,
		})
	}
	if dev.startSection > 0 {
		l.skipToSection(dev.startSection)
	}
	return l
}

func newLevel() *Level { return newLevelFromStage(stage1()) }

// phase1Sections mantém a lista de trechos da fase 1 (usada por testes).
func phase1Sections() []sectionDef { return stage1().sections }

// update avança um passo da linha do tempo e devolve os inimigos criados.
func (l *Level) update() []*Enemy {
	l.tick++
	l.updateSection()

	var spawned []*Enemy
	for _, ev := range l.events {
		if ev.done || l.tick < ev.startTick {
			continue
		}
		if !ev.started {
			ev.started = true
			ev.nextSpawnTick = l.tick
		}
		if l.tick >= ev.nextSpawnTick && ev.spawned < ev.count {
			newEnemies := ev.spawn()
			if ev.hasDrop && ev.spawned == 0 && len(newEnemies) > 0 {
				newEnemies[0].hasDrop = true
				newEnemies[0].drop = ev.drop
			}
			// Grupos com mais de um inimigo contam como formação para o bônus.
			if len(newEnemies) > 1 {
				l.nextFormationID++
				for _, e := range newEnemies {
					e.formationID = l.nextFormationID
				}
			}
			spawned = append(spawned, newEnemies...)
			ev.spawned++
			ev.nextSpawnTick = l.tick + jitterSpawnInterval(ev.interval)
		}
		if ev.spawned >= ev.count {
			ev.done = true
		}
	}
	return spawned
}

func (l *Level) updateSection() {
	next := l.section + 1
	if next < len(l.sections) && l.tick >= l.sections[next].startTick {
		l.section = next
		l.announce = l.sections[next].warning
		l.announceTimer = announceDuration
	}
	if l.announceTimer > 0 {
		l.announceTimer--
	}
}

func (l *Level) skipToSection(n int) {
	if n < 0 || n >= len(l.sections) {
		return
	}
	l.section = n
	l.tick = l.sections[n].startTick
	for _, ev := range l.events {
		if ev.startTick < l.tick {
			ev.done = true
			ev.spawned = ev.count
		}
	}
}

func (l *Level) allEventsDone() bool {
	for _, ev := range l.events {
		if !ev.done {
			return false
		}
	}
	return true
}

func (l *Level) readyForBoss(enemiesRemaining int) bool {
	return l.allEventsDone() && enemiesRemaining == 0
}

func (l *Level) theme() bgTheme {
	if l.section < 0 || l.section >= len(l.sections) {
		return sectionThemes[0]
	}
	return l.sections[l.section].theme
}

func (l *Level) music() musicID {
	if l.section < 0 || l.section >= len(l.sections) {
		return musicPhase
	}
	if m := l.sections[l.section].music; m != musicNone {
		return m
	}
	return musicPhase
}

func (l *Level) background() color.RGBA {
	return l.theme().sky
}

func (l *Level) sectionName() string {
	return l.sections[l.section].name
}

func (l *Level) progress() float64 {
	if l.duration <= 0 {
		return 0
	}
	p := float64(l.tick) / float64(l.duration)
	if p > 1 {
		return 1
	}
	return p
}

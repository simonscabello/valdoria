package game

import "image/color"

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
	switch ev.formation {
	case formationLine:
		out := make([]*Enemy, 0, lineFormationCount)
		for i := 0; i < lineFormationCount; i++ {
			out = append(out, spawnEnemy(ev.kind, ev.spawnX+float64(i)*formationGapX, 0, ev.fromLeft))
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
			x := ev.spawnX + float64(i-center)*formationGapX
			out = append(out, spawnEnemy(ev.kind, x, -float64(depth)*formationGapY, ev.fromLeft))
		}
		return out
	default:
		return []*Enemy{spawnEnemy(ev.kind, ev.spawnX, 0, ev.fromLeft)}
	}
}

func spawnEnemy(kind enemyKind, x, yOffset float64, fromLeft bool) *Enemy {
	switch kind {
	case kindHarpy:
		e := newHarpy(clampSpawnX(x, harpySize))
		e.y += yOffset
		return e
	case kindGargoyle:
		return newGargoyle(fromLeft, x, 30)
	case kindWyvern:
		e := newWyvern(clampSpawnX(x, wyvernSize))
		e.y += yOffset
		return e
	default:
		e := newCrow(clampSpawnX(x, crowSize))
		e.y += yOffset
		return e
	}
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

type levelSection struct {
	startTick int
	name      string
	warning   string
	bg        color.RGBA
}

type Level struct {
	events          []*waveEvent
	sections        []levelSection
	tick            int
	section         int
	announce        string
	announceTimer   int
	nextFormationID int
}

func newLevel() *Level {
	l := &Level{sections: phase1Sections()}

	add := func(start int, kind enemyKind, count, interval int, f formation, x float64, fromLeft bool) {
		l.events = append(l.events, &waveEvent{
			startTick: start, kind: kind, count: count,
			interval: interval, formation: f, spawnX: x, fromLeft: fromLeft,
		})
	}
	// drop marca o último evento adicionado como portador de um power-up garantido.
	drop := func(kind powerupType) {
		ev := l.events[len(l.events)-1]
		ev.hasDrop = true
		ev.drop = kind
	}

	// Trecho 1 - Campos do reino: corvos espaçados, introdução acessível.
	add(120, kindCrow, 1, 0, formationLine, 60, false)
	add(600, kindCrow, 1, 0, formationV, 120, false)
	drop(powerFire)
	add(1200, kindCrow, 2, 90, formationLine, 40, false)
	add(1900, kindCrow, 1, 0, formationV, 150, false)
	add(2400, kindCrow, 2, 90, formationLine, 90, false)

	// Trecho 2 - Vila atacada: harpias com projéteis surgem entre corvos.
	add(2760, kindHarpy, 3, 70, formationSingle, 40, false)
	drop(powerIce)
	add(3200, kindCrow, 2, 80, formationLine, 120, false)
	add(3800, kindHarpy, 2, 90, formationSingle, 170, false)
	add(4300, kindCrow, 1, 0, formationV, 90, false)
	drop(powerHeal)
	add(4800, kindHarpy, 3, 80, formationSingle, 60, false)

	// Trecho 3 - Muralhas: gárgulas pelas laterais protegidas por harpias.
	add(5460, kindGargoyle, 1, 0, formationSingle, 60, true)
	drop(powerLight)
	add(5560, kindHarpy, 2, 80, formationSingle, 40, false)
	add(6100, kindGargoyle, 1, 0, formationSingle, 170, false)
	add(6600, kindHarpy, 3, 70, formationSingle, 150, false)
	add(7100, kindGargoyle, 1, 0, formationSingle, 100, true)
	drop(powerHeal)
	add(7400, kindHarpy, 2, 80, formationSingle, 90, false)

	// Trecho 4 - Aproximação do castelo: wyverns e formações combinadas.
	add(8160, kindWyvern, 1, 0, formationSingle, 60, false)
	drop(powerShield)
	add(8500, kindHarpy, 3, 70, formationSingle, 120, false)
	add(8900, kindWyvern, 1, 0, formationSingle, 160, false)
	add(9200, kindCrow, 2, 80, formationLine, 40, false)
	add(9500, kindWyvern, 1, 0, formationSingle, 100, false)

	if dev.startSection > 0 {
		l.skipToSection(dev.startSection)
	}
	return l
}

func phase1Sections() []levelSection {
	return []levelSection{
		{startTick: 0, name: "Campos do reino", warning: "", bg: color.RGBA{0x0a, 0x0a, 0x1e, 0xff}},
		{startTick: 2700, name: "Vila atacada", warning: "A vila esta sob ataque!", bg: color.RGBA{0x1e, 0x12, 0x12, 0xff}},
		{startTick: 5400, name: "Muralhas", warning: "Aproximando-se das muralhas", bg: color.RGBA{0x12, 0x14, 0x24, 0xff}},
		{startTick: 8100, name: "Aproximacao do castelo", warning: "O castelo se aproxima...", bg: color.RGBA{0x1a, 0x10, 0x22, 0xff}},
	}
}

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
			ev.nextSpawnTick = l.tick + ev.interval
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
	return themeForSection(l.section)
}

func (l *Level) background() color.RGBA {
	return l.theme().sky
}

func (l *Level) sectionName() string {
	return l.sections[l.section].name
}

func (l *Level) progress() float64 {
	p := float64(l.tick) / float64(phaseDurationTicks)
	if p > 1 {
		return 1
	}
	return p
}

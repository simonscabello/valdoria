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
	events        []*waveEvent
	sections      []levelSection
	tick          int
	section       int
	announce      string
	announceTimer int
}

func newLevel() *Level {
	l := &Level{sections: phase1Sections()}

	add := func(start int, kind enemyKind, count, interval int, f formation, x float64, fromLeft bool) {
		l.events = append(l.events, &waveEvent{
			startTick: start, kind: kind, count: count,
			interval: interval, formation: f, spawnX: x, fromLeft: fromLeft,
		})
	}

	// Trecho 1 - Campos do reino: formações de corvos, baixa dificuldade.
	add(120, kindCrow, 1, 0, formationLine, 40, false)
	add(600, kindCrow, 1, 0, formationV, 120, false)
	add(1200, kindCrow, 2, 90, formationLine, 80, false)
	add(2000, kindCrow, 1, 0, formationV, 160, false)
	add(2600, kindCrow, 2, 80, formationLine, 30, false)
	add(3200, kindCrow, 1, 0, formationV, 100, false)

	// Trecho 2 - Vila atacada: harpias em zigue-zague misturadas a corvos.
	add(3800, kindHarpy, 3, 70, formationSingle, 40, false)
	add(4200, kindCrow, 2, 80, formationLine, 120, false)
	add(4800, kindHarpy, 2, 90, formationSingle, 180, false)
	add(5400, kindCrow, 1, 0, formationV, 90, false)
	add(5800, kindHarpy, 3, 80, formationSingle, 60, false)
	add(6400, kindCrow, 2, 70, formationLine, 140, false)

	// Trecho 3 - Muralhas: gárgulas pelas laterais protegidas por harpias.
	add(7400, kindGargoyle, 1, 0, formationSingle, 60, true)
	add(7500, kindHarpy, 2, 80, formationSingle, 40, false)
	add(8200, kindGargoyle, 1, 0, formationSingle, 170, false)
	add(8300, kindHarpy, 2, 80, formationSingle, 150, false)
	add(9000, kindGargoyle, 1, 0, formationSingle, 100, true)
	add(9100, kindHarpy, 3, 70, formationSingle, 90, false)
	add(9800, kindGargoyle, 1, 0, formationSingle, 130, false)

	// Trecho 4 - Aproximação do castelo: wyverns e formações combinadas.
	add(11000, kindWyvern, 1, 0, formationSingle, 60, false)
	add(11400, kindHarpy, 3, 70, formationSingle, 120, false)
	add(11800, kindWyvern, 1, 0, formationSingle, 160, false)
	add(12200, kindCrow, 2, 80, formationLine, 40, false)
	add(12800, kindWyvern, 2, 200, formationSingle, 100, false)

	if devStartSection > 0 {
		l.skipToSection(devStartSection)
	}
	return l
}

func phase1Sections() []levelSection {
	return []levelSection{
		{startTick: 0, name: "Campos do reino", warning: "", bg: color.RGBA{0x0a, 0x0a, 0x1e, 0xff}},
		{startTick: 3600, name: "Vila atacada", warning: "A vila esta sob ataque!", bg: color.RGBA{0x1e, 0x12, 0x12, 0xff}},
		{startTick: 7200, name: "Muralhas", warning: "Aproximando-se das muralhas", bg: color.RGBA{0x12, 0x14, 0x24, 0xff}},
		{startTick: 10800, name: "Aproximacao do castelo", warning: "O castelo se aproxima...", bg: color.RGBA{0x1a, 0x10, 0x22, 0xff}},
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
			spawned = append(spawned, ev.spawn()...)
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

func (l *Level) background() color.RGBA {
	return l.sections[l.section].bg
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

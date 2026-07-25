package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// O Medidor de Corrupção.
//
// Este é o diferencial do jogo: todo shoot'em up pune quem deixa inimigos
// passarem (perde ponto, perde vida). Aqui, deixar passar **corrompe o reino** —
// e a corrupção é ao mesmo tempo a punição e a recompensa. Ela deixa os
// inimigos mais fortes e o mundo mais doente, mas multiplica pontos e drops, e
// leva a um confronto final diferente.
//
// O jogador que deixa passar por incompetência é punido. O que deixa passar de
// propósito está apostando. Essa é a decisão que nenhum concorrente oferece.
//
// A corrupção **nunca desce sozinha**: purificá-la será o prêmio dos altares
// (v0.6). Uma vez perdido, o terreno só se recupera com esforço.

// corruptionTier é uma faixa do medidor. Cada faixa muda o mundo e a economia.
type corruptionTier int

const (
	tierSteady   corruptionTier = iota // Reino Firme
	tierShadow                         // Sombra Crescente
	tierSiege                          // Cerco — os inimigos começam a mutar
	tierCollapse                       // Colapso
	tierFall                           // Queda de Valdoria — chefe verdadeiro
	corruptionTierCount
)

// corruptionParams descreve o que cada faixa faz com o mundo e com a economia.
type corruptionParams struct {
	name string
	// min é o valor a partir do qual a faixa vale.
	min float64
	// enemyHealthMul, scoreMul e dropMul são a troca central: o mundo fica mais
	// perigoso na mesma medida em que fica mais lucrativo.
	enemyHealthMul float64
	scoreMul       float64
	dropMul        float64
	// mutates liga as variantes corrompidas dos inimigos.
	mutates bool
	// tint é a cor que invade a tela nessa faixa.
	tint color.RGBA
}

var corruptionTable = [corruptionTierCount]corruptionParams{
	tierSteady: {
		name: "REINO FIRME", min: 0,
		enemyHealthMul: 1.0, scoreMul: 1.0, dropMul: 1.0,
		tint: color.RGBA{0x00, 0x00, 0x00, 0x00},
	},
	tierShadow: {
		name: "SOMBRA CRESCENTE", min: 26,
		enemyHealthMul: 1.15, scoreMul: 1.3, dropMul: 1.2,
		tint: color.RGBA{0x4a, 0x18, 0x50, 0x18},
	},
	tierSiege: {
		name: "CERCO", min: 51,
		enemyHealthMul: 1.3, scoreMul: 1.8, dropMul: 1.35,
		mutates: true,
		tint:    color.RGBA{0x6a, 0x14, 0x60, 0x30},
	},
	tierCollapse: {
		name: "COLAPSO", min: 76,
		enemyHealthMul: 1.5, scoreMul: 2.5, dropMul: 1.5,
		mutates: true,
		tint:    color.RGBA{0x8a, 0x10, 0x70, 0x48},
	},
	tierFall: {
		name: "QUEDA DE VALDORIA", min: 100,
		enemyHealthMul: 1.6, scoreMul: 3.0, dropMul: 1.5,
		mutates: true,
		tint:    color.RGBA{0xa0, 0x0c, 0x7a, 0x60},
	},
}

const (
	maxCorruption = 100.0

	// Peso de cada inimigo que atravessa a base da tela. Calibrado por medição
	// (`go run ./cmd/balance`, seção 5) para que a curva atravesse todas as
	// faixas conforme a negligência do jogador:
	//
	//	  5% de fugas -> ~15%  Reino Firme        (imprecisão normal não pune)
	//	 15% de fugas -> ~45%  Sombra Crescente
	//	 25% de fugas -> ~76%  Colapso
	//	 40% de fugas -> 100%  Queda de Valdoria  (Vharak Ascendido)
	//
	// O ponto é que chegar ao final alternativo exige **decisão**, não descuido.
	//
	// A gárgula não aparece: ela sai pela lateral por design, e sair de lado
	// nunca foi uma falha do jogador.
	corruptionCrow     = 0.6
	corruptionHarpy    = 1.2
	corruptionMage     = 2.0
	corruptionBallista = 2.0
	corruptionWyvern   = 2.8

	// Frames de aviso ao subir de faixa.
	corruptionAnnounce = 140
)

// corruptionWeight devolve quanto um inimigo corrompe o reino ao escapar.
func corruptionWeight(k enemyKind) float64 {
	switch k {
	case kindHarpy:
		return corruptionHarpy
	case kindMage:
		return corruptionMage
	case kindBallista:
		return corruptionBallista
	case kindWyvern:
		return corruptionWyvern
	case kindGargoyle:
		return 0 // sai pela lateral por design
	default:
		return corruptionCrow
	}
}

// tierFor devolve a faixa correspondente a um valor de corrupção.
func tierFor(v float64) corruptionTier {
	t := tierSteady
	for i := corruptionTier(0); i < corruptionTierCount; i++ {
		if v >= corruptionTable[i].min {
			t = i
		}
	}
	return t
}

// Corruption acumula o estado do medidor durante uma partida.
type Corruption struct {
	value float64
	tier  corruptionTier

	// escaped conta os inimigos que atravessaram a base — a estatística que o
	// jogador vê no fim e a métrica que explica o resultado dele.
	escaped int

	// announce/announceTimer exibem a mudança de faixa.
	announce      string
	announceTimer int
	// pulse dá um destaque curto na barra a cada acréscimo.
	pulse int
}

func (c *Corruption) reset() { *c = Corruption{} }

// add registra a fuga de um inimigo e devolve true se a faixa mudou.
func (c *Corruption) add(k enemyKind) bool {
	w := corruptionWeight(k)
	if w <= 0 {
		return false
	}
	c.escaped++
	c.pulse = 12
	before := c.tier
	c.value += w
	if c.value > maxCorruption {
		c.value = maxCorruption
	}
	c.tier = tierFor(c.value)
	if c.tier == before {
		return false
	}
	c.announce = corruptionTable[c.tier].name
	c.announceTimer = corruptionAnnounce
	return true
}

func (c *Corruption) update() {
	if c.announceTimer > 0 {
		c.announceTimer--
	}
	if c.pulse > 0 {
		c.pulse--
	}
}

func (c *Corruption) params() corruptionParams { return corruptionTable[c.tier] }

// ratio devolve o preenchimento da barra, de 0 a 1.
func (c *Corruption) ratio() float64 { return c.value / maxCorruption }

// percent devolve o valor inteiro exibido ao jogador.
func (c *Corruption) percent() int { return int(c.value + 0.5) }

// fallen indica o reino totalmente corrompido — desbloqueia Vharak Ascendido.
func (c *Corruption) fallen() bool { return c.value >= maxCorruption }

// mutates indica se os inimigos gerados agora nascem corrompidos.
func (c *Corruption) mutates() bool { return c.params().mutates }

// scoreMul e dropMul são o outro lado da aposta: o mundo pior paga melhor.
func (c *Corruption) scoreMul() float64 { return c.params().scoreMul }
func (c *Corruption) dropMul() float64  { return c.params().dropMul }

// healthMul escala a vida dos inimigos gerados a partir de agora.
func (c *Corruption) healthMul() float64 { return c.params().enemyHealthMul }

// tierName é o rótulo da faixa atual, para HUD e telas de fim de partida.
func (c *Corruption) tierName() string { return c.params().name }

// --- Variantes corrompidas ---
//
// A partir do Cerco (51%), os inimigos nascem corrompidos: mais resistentes,
// mais valiosos e com um traço a mais no comportamento. A ideia é que o jogador
// **veja** a consequência da sua corrupção no que desce a tela, e não só num
// número no HUD.

const (
	mutatedHealthMul = 1.4
	mutatedScoreMul  = 1.6

	// Corvo-Sombra: descarrega um tiro ao atingir esta altura. O corvo deixa de
	// ser inofensivo — é a mudança mais sentida de todas.
	shadowCrowFireY  = 110
	shadowCrowSpeed  = 1.8
	shadowHarpySpeed = 1.5 // multiplicador sobre a velocidade base
	rotWyvernFan     = 3   // projéteis por salva do wyvern corrompido
	rotWyvernSpread  = 0.6
)

var (
	corruptedBody   = color.RGBA{0x8a, 0x1e, 0x74, 0xff}
	corruptedAccent = color.RGBA{0xff, 0x64, 0xd0, 0xff}
	// mutatedTint multiplica o sprite existente: puxa o verde para baixo e
	// mantém vermelho/azul, deixando qualquer criatura com um tom doente sem
	// exigir uma segunda pixel art.
	mutatedTint = color.RGBA{0xff, 0x5a, 0xd8, 0xff}
)

// mutateEnemy converte um inimigo recém-criado na sua variante corrompida.
func mutateEnemy(e *Enemy) {
	if e.mutated {
		return
	}
	e.mutated = true
	e.health = int(float64(e.health)*mutatedHealthMul + 0.5)
	e.score = int(float64(e.score)*mutatedScoreMul + 0.5)
	e.color = corruptedBody
	e.accent = corruptedAccent

	switch e.kind {
	case kindCrow:
		e.behavior = &shadowCrowBehavior{}
	case kindHarpy:
		if b, ok := e.behavior.(*harpyBehavior); ok {
			b.speedMul = shadowHarpySpeed
			b.fireCD = harpyFireInterval / 2
		}
	case kindWyvern:
		if b, ok := e.behavior.(*wyvernBehavior); ok {
			b.fan = rotWyvernFan
		}
	}
}

// shadowCrowBehavior: o corvo corrompido desce mais rápido e dispara uma única
// vez ao cruzar a metade da tela. Deixar corvos passarem faz com que os
// próximos atirem em você — a corrupção se retroalimenta.
type shadowCrowBehavior struct {
	fired bool
}

func (b *shadowCrowBehavior) update(e *Enemy, ctx *enemyContext) {
	e.y += shadowCrowSpeed
	if b.fired || e.y < shadowCrowFireY {
		return
	}
	b.fired = true
	e.telegraph = enemyTelegraphFrames
	vx, vy := aimVelocity(e.centerX(), e.y+e.h, ctx.playerX, ctx.playerY, enemyBulletSpeed)
	*ctx.bullets = append(*ctx.bullets, newEnemyBullet(e.centerX(), e.y+e.h, vx, vy))
}

// drawCorruptedAura marca a variante corrompida com um halo pulsante, para o
// jogador distinguir de relance uma harpia de uma harpia corrompida.
func (e *Enemy) drawCorruptedAura(screen *ebiten.Image) {
	if (e.animTick/6)%2 != 0 {
		return
	}
	vector.StrokeRect(screen,
		float32(e.x-2), float32(e.y-2), float32(e.w+4), float32(e.h+4),
		1, withAlpha(corruptedAccent, 150), false)
}

// barColor devolve a cor da barra, que esquenta conforme a corrupção avança.
func (c *Corruption) barColor() color.RGBA {
	switch c.tier {
	case tierShadow:
		return color.RGBA{0x8a, 0x4a, 0xc0, 0xff}
	case tierSiege:
		return color.RGBA{0xb0, 0x30, 0xd0, 0xff}
	case tierCollapse:
		return color.RGBA{0xd8, 0x20, 0xc0, 0xff}
	case tierFall:
		return color.RGBA{0xff, 0x28, 0xa0, 0xff}
	default:
		return color.RGBA{0x5a, 0x50, 0x78, 0xff}
	}
}

package game

import "image/color"

// Parâmetros visuais centralizados para a identidade provisória e os efeitos.
// Mantê-los aqui evita números mágicos espalhados pelos desenhos.
const (
	// Animação e pose do jogador.
	playerWingSpeed = 0.28 // velocidade da batida de asas
	playerTiltMax   = 4.0  // deslocamento máximo das asas ao inclinar
	playerTiltRate  = 0.35 // suavização da inclinação lateral
	enemyWingSpeed  = 0.2  // velocidade da batida de asas dos inimigos
	bossWingSpeed   = 0.12 // batida de asas do chefe

	// Partículas.
	explosionParticles = 10
	explosionSpeed     = 2.4
	explosionLife      = 26
	collectParticles   = 10
	collectSpeed       = 1.8
	collectLife        = 20
	trailLife          = 12
	bossDeathParticles = 18
	particleGravity    = 0.04
	particleDrag       = 0.92
	maxParticles       = 400 // teto de segurança: partículas são puramente cosméticas

	// Números de pontuação flutuantes ao abater inimigos.
	popupLife = 34
	popupRise = 0.55
	maxPopups = 40

	// Impact freeze (hit-stop) em golpes fortes: poucos frames de congelamento.
	hitStopBig = 3

	// Vibração de tela (screen shake).
	shakeHitMagnitude  = 2.0
	shakeBombMagnitude = 6.0
	shakeBossMagnitude = 4.5
	// Subir de faixa de corrupção é um evento de mundo: sacode mais que um dano.
	shakeCorruptionMagnitude = 5.0
	// Mergulho: sacudida curta, só para dar arranque ao gesto.
	shakeDiveMagnitude = 2.4
	shakeDecay         = 0.86
	shakeMaxOffset     = 5.0
	shakeMinMagnitude  = 0.15

	// Flash discreto ao receber dano.
	damageFlashDuration = 12
	damageFlashAlpha    = 70

	// Escurecimento de impacto ao usar a Invocação Ancestral (bomba).
	bombDarkenFrames = 10
	bombDarkenAlpha  = 150

	// Escurecimento do cenário: empurra o parallax para trás, dando contraste e
	// leitura clara aos elementos de jogo (inimigos, chefe, projéteis).
	sceneDimAlpha = 104

	// Parallax do cenário.
	cloudLayerSpeed     = 0.25
	hillLayerSpeed      = 0.6
	structureLayerSpeed = 1.2
	dustLayerSpeed      = 0.9
	cloudCount          = 6
	hillCount           = 4
	structureCount      = 5
)

// shakeLevel controla a intensidade da vibração escolhida pelo jogador.
type shakeLevel int

const (
	shakeFull shakeLevel = iota
	shakeReduced
	shakeOff
)

func shakeLevelName(s shakeLevel) string {
	switch s {
	case shakeReduced:
		return "Reduzida"
	case shakeOff:
		return "Off"
	default:
		return "Cheia"
	}
}

func (s shakeLevel) scale() float64 {
	switch s {
	case shakeReduced:
		return 0.4
	case shakeOff:
		return 0
	default:
		return 1
	}
}

// structureStyle define a silhueta desenhada na camada de estruturas do cenário.
type structureStyle int

const (
	styleFields structureStyle = iota
	styleVillageFire
	styleWalls
	styleCastle
)

// bgTheme reúne as cores de cada camada de parallax por trecho da fase.
type bgTheme struct {
	sky       color.RGBA
	cloud     color.RGBA
	hill      color.RGBA
	structure color.RGBA
	accent    color.RGBA
	dust      color.RGBA
	style     structureStyle
}

// sectionThemes segue a ordem dos trechos da fase 1.
var sectionThemes = []bgTheme{
	{ // Campos do reino
		sky:       color.RGBA{0x0a, 0x12, 0x1e, 0xff},
		cloud:     color.RGBA{0x24, 0x2c, 0x40, 0xff},
		hill:      color.RGBA{0x12, 0x22, 0x1a, 0xff},
		structure: color.RGBA{0x1c, 0x38, 0x22, 0xff},
		accent:    color.RGBA{0x3a, 0x6a, 0x40, 0xff},
		dust:      color.RGBA{0x8a, 0xb0, 0x90, 0xff},
		style:     styleFields,
	},
	{ // Vila atacada (em chamas)
		sky:       color.RGBA{0x1e, 0x10, 0x0e, 0xff},
		cloud:     color.RGBA{0x3a, 0x22, 0x1a, 0xff},
		hill:      color.RGBA{0x22, 0x12, 0x10, 0xff},
		structure: color.RGBA{0x32, 0x1c, 0x16, 0xff},
		accent:    color.RGBA{0xff, 0x7a, 0x2a, 0xff},
		dust:      color.RGBA{0xff, 0xb0, 0x60, 0xff},
		style:     styleVillageFire,
	},
	{ // Muralhas
		sky:       color.RGBA{0x10, 0x12, 0x1e, 0xff},
		cloud:     color.RGBA{0x24, 0x28, 0x36, 0xff},
		hill:      color.RGBA{0x1a, 0x1c, 0x26, 0xff},
		structure: color.RGBA{0x30, 0x34, 0x40, 0xff},
		accent:    color.RGBA{0x58, 0x60, 0x70, 0xff},
		dust:      color.RGBA{0x9a, 0x9a, 0xb0, 0xff},
		style:     styleWalls,
	},
	{ // Aproximação do castelo
		sky:       color.RGBA{0x16, 0x0e, 0x20, 0xff},
		cloud:     color.RGBA{0x2c, 0x20, 0x3a, 0xff},
		hill:      color.RGBA{0x1c, 0x14, 0x26, 0xff},
		structure: color.RGBA{0x34, 0x28, 0x44, 0xff},
		accent:    color.RGBA{0x8a, 0x6a, 0xc0, 0xff},
		dust:      color.RGBA{0xb0, 0x9a, 0xe0, 0xff},
		style:     styleCastle,
	},
}

// Biomas adicionais das fases 2 e 3. Reutilizam os estilos de estrutura
// existentes com paletas próprias, dando variedade visual sem novo código de
// desenho.
var (
	themeForest = bgTheme{ // Bosque sombrio (fase 2)
		sky:       color.RGBA{0x08, 0x14, 0x0e, 0xff},
		cloud:     color.RGBA{0x18, 0x2a, 0x20, 0xff},
		hill:      color.RGBA{0x10, 0x20, 0x16, 0xff},
		structure: color.RGBA{0x14, 0x2a, 0x1a, 0xff},
		accent:    color.RGBA{0x3a, 0x7a, 0x44, 0xff},
		dust:      color.RGBA{0x8a, 0xc0, 0x90, 0xff},
		style:     styleFields,
	}
	themeSwamp = bgTheme{ // Pântano corrompido (fase 2)
		sky:       color.RGBA{0x12, 0x0e, 0x1c, 0xff},
		cloud:     color.RGBA{0x22, 0x1a, 0x2e, 0xff},
		hill:      color.RGBA{0x18, 0x14, 0x22, 0xff},
		structure: color.RGBA{0x24, 0x1c, 0x30, 0xff},
		accent:    color.RGBA{0x7a, 0x4a, 0xb0, 0xff},
		dust:      color.RGBA{0xa0, 0x80, 0xd0, 0xff},
		style:     styleFields,
	}
	themeCanyon = bgTheme{ // Desfiladeiro (fase 3)
		sky:       color.RGBA{0x1c, 0x10, 0x0c, 0xff},
		cloud:     color.RGBA{0x30, 0x1e, 0x14, 0xff},
		hill:      color.RGBA{0x24, 0x16, 0x10, 0xff},
		structure: color.RGBA{0x3a, 0x24, 0x18, 0xff},
		accent:    color.RGBA{0xc0, 0x6a, 0x30, 0xff},
		dust:      color.RGBA{0xe0, 0xa0, 0x60, 0xff},
		style:     styleWalls,
	}
	themeLair = bgTheme{ // Covil do dragão (fase 3, chefe)
		sky:       color.RGBA{0x20, 0x08, 0x08, 0xff},
		cloud:     color.RGBA{0x3a, 0x14, 0x12, 0xff},
		hill:      color.RGBA{0x28, 0x0e, 0x0e, 0xff},
		structure: color.RGBA{0x3c, 0x18, 0x16, 0xff},
		accent:    color.RGBA{0xff, 0x5a, 0x28, 0xff},
		dust:      color.RGBA{0xff, 0x90, 0x50, 0xff},
		style:     styleVillageFire,
	}
)

// corrupted devolve o tema degradado pela corrupção: o mundo perde saturação e
// escorrega para o violeta doente.
//
// A regra do guia de arte é preservada — o cenário nunca ganha saturação, só
// perde. Assim a leitura das entidades (que são saturadas e contornadas) fica
// *melhor* conforme o reino piora, e não pior.
func (t bgTheme) corrupted(ratio float64) bgTheme {
	if ratio <= 0 {
		return t
	}
	if ratio > 1 {
		ratio = 1
	}
	t.sky = corruptColor(t.sky, ratio)
	t.cloud = corruptColor(t.cloud, ratio)
	t.hill = corruptColor(t.hill, ratio)
	t.structure = corruptColor(t.structure, ratio)
	t.accent = corruptColor(t.accent, ratio)
	t.dust = corruptColor(t.dust, ratio)
	return t
}

// corruptRot é o alvo cromático da corrupção: violeta escuro e sem vida.
var corruptRot = color.RGBA{0x2a, 0x0e, 0x2e, 0xff}

// corruptColor dessatura em direção ao cinza e depois puxa para o violeta,
// proporcionalmente à corrupção. No máximo, o mundo perde ~55% da sua cor.
func corruptColor(c color.RGBA, ratio float64) color.RGBA {
	const maxDrain = 0.55
	k := ratio * maxDrain

	// Luminância perceptual, para dessaturar sem escurecer demais.
	lum := 0.299*float64(c.R) + 0.587*float64(c.G) + 0.114*float64(c.B)
	mix := func(v uint8, target float64) uint8 {
		return clampByte(int(float64(v)*(1-k) + target*k))
	}
	gray := color.RGBA{mix(c.R, lum), mix(c.G, lum), mix(c.B, lum), c.A}

	// Metade da degradação vai para o violeta da corrupção.
	k2 := ratio * maxDrain * 0.5
	toward := func(v, target uint8) uint8 {
		return clampByte(int(float64(v)*(1-k2) + float64(target)*k2))
	}
	return color.RGBA{
		toward(gray.R, corruptRot.R),
		toward(gray.G, corruptRot.G),
		toward(gray.B, corruptRot.B),
		c.A,
	}
}

// enemyOutline é o contorno escuro desenhado atrás dos corpos para destacá-los
// do cenário, reforçando a leitura das silhuetas.
var enemyOutline = color.RGBA{0x0a, 0x08, 0x12, 0xff}

// withAlpha devolve a cor com o canal alfa ajustado, útil para fades e rastros.
func withAlpha(c color.RGBA, a uint8) color.RGBA {
	return color.RGBA{c.R, c.G, c.B, a}
}

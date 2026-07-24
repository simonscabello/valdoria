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

	// Vibração de tela (screen shake).
	shakeHitMagnitude  = 2.0
	shakeBombMagnitude = 6.0
	shakeBossMagnitude = 4.5
	shakeDecay         = 0.86
	shakeMaxOffset     = 5.0
	shakeMinMagnitude  = 0.15

	// Flash discreto ao receber dano.
	damageFlashDuration = 12
	damageFlashAlpha    = 70

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

func themeForSection(section int) bgTheme {
	if section < 0 || section >= len(sectionThemes) {
		return sectionThemes[0]
	}
	return sectionThemes[section]
}

// withAlpha devolve a cor com o canal alfa ajustado, útil para fades e rastros.
func withAlpha(c color.RGBA, a uint8) color.RGBA {
	return color.RGBA{c.R, c.G, c.B, a}
}

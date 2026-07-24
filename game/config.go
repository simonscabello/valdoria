package game

// Configurações centrais de jogabilidade. Ajuste aqui para calibrar o controle.
const (
	initialHealth = 3
	maxHealth     = 5

	playerSize           = 16
	playerHitboxSize     = 4
	playerSpeed          = 2.5
	playerPrecisionSpeed = 1.0

	bulletInset         = 3
	muzzleFlashDuration = 4

	invincibilityDuration = 90 // frames (~1.5s a 60 TPS)
	invincibilityBlink    = 6  // frames por ciclo de piscada
)

// Configurações de inimigos e projéteis inimigos.
const (
	hitFlashDuration = 4

	enemyBulletSize   = 4
	enemyBulletSpeed  = 2.2
	enemyBulletDamage = 1

	crowSize   = 10
	crowHealth = 1
	crowScore  = 10
	crowSpeed  = 2.4
	crowDamage = 1

	harpySize         = 14
	harpyHealth       = 3
	harpyScore        = 30
	harpySpeed        = 1.1
	harpyDamage       = 1
	harpyAmplitude    = 42
	harpyFreq         = 0.05
	harpyFireInterval = 70

	gargoyleSize           = 18
	gargoyleHealth         = 6
	gargoyleScore          = 50
	gargoyleSpeed          = 1.6
	gargoyleDamage         = 1
	gargoyleAttackDuration = 180 // ~3s atacando
	gargoyleFireInterval   = 40

	wyvernSize         = 24
	wyvernHealth       = 10
	wyvernScore        = 100
	wyvernDamage       = 1
	wyvernDescend      = 0.35
	wyvernAlign        = 0.8
	wyvernFireInterval = 90
)

// Configurações da fase e do sistema de ondas.
const (
	phaseDurationTicks = 13400 // referência para a barra de progresso
	announceDuration   = 150   // frames que um aviso permanece na tela

	lineFormationCount = 4
	vFormationCount    = 5
	formationGapX      = 20
	formationGapY      = 16

	// Ajustes de desenvolvimento/teste.
	devMode          = true
	devStartSection  = 0 // inicia a fase neste trecho (0 a 3)
	devTimeScale     = 1 // passos de linha do tempo por frame
	devFastTimeScale = 6 // valor aplicado ao acelerar com Tab
)

// Configurações de armas e power-ups.
const (
	maxWeaponLevel = 3

	lightBulletSpeed     = 5.0
	lightBulletSpeedFast = 6.5
	lightBulletDamage    = 2
	lightCooldown        = 10

	flameSpread       = 0.9 // abertura do leque em radianos
	flameBulletSpeed  = 4.0
	flameBulletDamage = 1
	flameCooldown     = 12
	flameCooldownFast = 7

	iceBulletSpeed  = 3.0
	iceBulletDamage = 4
	iceCooldown     = 20
	icePierce       = 1
	icePierceMax    = 3

	powerupSize  = 10
	powerupSpeed = 1.0

	dropChance     = 0.18
	shieldDuration = 480 // ~8s a 60 TPS
)

package game

// Configurações centrais de jogabilidade. Ajuste aqui para calibrar o controle.
const (
	initialHealth = 3
	maxHealth     = 5
	startingLives = 3

	playerSize           = 16
	playerHitboxSize     = 4
	playerSpeed          = 2.5
	playerPrecisionSpeed = 1.0

	bulletInset         = 3
	muzzleFlashDuration = 4

	invincibilityDuration = 90  // frames (~1.5s a 60 TPS)
	invincibilityBlink    = 6   // frames por ciclo de piscada
	respawnInvincibility  = 150 // invencibilidade ao reaparecer
)

// Habilidade especial (Invocação Ancestral) e pontuação.
const (
	bombStartCharges   = 2
	bombDamage         = 999
	bombInvincibility  = 120
	bombEffectDuration = 45

	formationBonus    = 60
	waveNoDamageBonus = 200

	comboWindow   = 180 // frames sem eliminar até o combo decair
	comboStep     = 5   // eliminações por nível de multiplicador
	maxMultiplier = 4
)

// Configurações de inimigos e projéteis inimigos.
const (
	hitFlashDuration = 4

	// Antecedência (frames) do aviso visual antes de um inimigo atirar.
	enemyTelegraphFrames = 16

	enemyBulletSize   = 4
	enemyBulletSpeed  = 2.2
	enemyBulletDamage = 1

	crowSize   = 10
	crowHealth = 1
	crowScore  = 10
	crowSpeed  = 1.5 // mais lento: densidade alta + fuga pune, precisa de tempo de mira
	crowDamage = 1

	// Margem mínima das bordas para spawns "justos" (voadores/terrestres).
	// Evita corvos impossíveis nos extremos da tela de 240px.
	spawnFairMargin = 36

	harpySize         = 14
	harpyHealth       = 3
	harpyScore        = 30
	harpySpeed        = 0.9
	harpyDamage       = 1
	harpyAmplitude    = 42
	harpyFreq         = 0.05
	harpyFireInterval = 65

	gargoyleSize           = 18
	gargoyleHealth         = 6
	gargoyleScore          = 50
	gargoyleSpeed          = 1.5
	gargoyleDamage         = 1
	gargoyleAttackDuration = 180 // ~3s atacando
	gargoyleFireInterval   = 36

	wyvernSize         = 24
	wyvernHealth       = 8
	wyvernScore        = 100
	wyvernDamage       = 1
	wyvernDescend      = 0.28
	wyvernAlign        = 0.8
	wyvernFireInterval = 80

	// Balista corrompida: ameaça "terrestre" que desce lenta como o cenário e
	// dispara rajadas de virotes mirados. Resistente e bem telegrafada.
	ballistaSize         = 22
	ballistaHealth       = 6
	ballistaScore        = 70
	ballistaDamage       = 1
	ballistaDescend      = 0.35
	ballistaFireInterval = 100
	ballistaBurst        = 3
	ballistaBurstGap     = 9

	// Feiticeiro corrompido: para no alto, desliza de lado e solta anéis de
	// projéteis. Recompensa o jogador por pressioná-lo rapidamente.
	mageSize         = 16
	mageHealth       = 5
	mageScore        = 90
	mageDamage       = 1
	mageDescend      = 1.0
	mageStopY        = 54
	mageDrift        = 0.75
	mageFireInterval = 115
	mageRingCount    = 10
	mageBulletSpeed  = 1.9
)

// Configurações da fase e do sistema de ondas.
const (
	announceDuration = 150 // frames que um aviso permanece na tela

	lineFormationCount = 4
	vFormationCount    = 5
	formationGapX      = 20
	formationGapY      = 16

	devBossDamage = 20 // dano aplicado ao chefe pela tecla de desenvolvimento
)

// Chefe da fase 1 — Vharak, o Dragão Corrompido.
const (
	bossMaxHealth = 420 // antes 200: o confronto final precisa pesar
	bossW         = 60
	bossH         = 34
	bossY         = 26
	bossMargin    = 10

	bossEntrySpeed      = 0.6
	bossEntryHold       = 90 // frames exibindo o nome antes do combate
	bossWarningDuration = 45
	bossDeathDuration   = 120
	bossScore           = 5000
	bossContactDamage   = 1

	bossSpeedP1    = 0.65
	bossSpeedP2    = 1.35
	bossSpeedP3    = 2.0
	bossSweepSpeed = 2.2

	// Varredura: brecha larga e lenta garante que sempre haja rota de fuga.
	bossSweepStep    = 16
	bossSweepGapHalf = 28
	bossSweepGapMove = 1.0

	crystalBonus = 3 // multiplicador de dano ao acertar um cristal
)

// Fluxo de telas e bônus finais.
const (
	fadeDuration    = 30
	lifeBonus       = 500
	bombBonusPoints = 300
)

// Configurações de armas e power-ups.
// Identidades: Luz = foco/DPS; Chamas = cobertura controlada; Gelo = peso/perfuração.
const (
	maxWeaponLevel = 3

	lightBulletSpeed     = 5.5
	lightBulletSpeedFast = 7.0
	lightBulletDamage    = 3
	lightCooldown        = 8
	lightPierceMax       = 1 // Nv3 atravessa um inimigo extra

	flameSpread       = 0.52 // leque ~30°; antes 0.9 varria quase a tela toda
	flameBulletSpeed  = 3.7
	flameBulletDamage = 1
	flameCooldown     = 14
	flameCooldownFast = 11 // Nv3 um pouco mais rápido, sem virar metralhadora

	iceBulletSpeed  = 3.8
	iceBulletDamage = 5
	iceCooldown     = 14
	icePierce       = 1
	icePierceMax    = 3

	// Efeitos elementais ao acertar (identidade de cada arma).
	iceSlowDuration = 100 // frames lento
	iceSlowFactor   = 0.4 // fração dos updates de movimento/tiro

	burnDuration  = 130
	burnInterval  = 16 // 1 de dano a cada N frames
	burnDamage    = 1
	burnBonusPct  = 50 // +% de dano recebido enquanto queima

	stunDuration     = 32 // Luz: julgamento — para no lugar
	lightChainRange  = 52
	lightChainRatio  = 0.55 // fração do dano em cadeia
	lightChainStun   = 18

	powerupSize  = 10
	powerupSpeed = 1.0

	dropChance     = 0.2
	shieldDuration = 480 // ~8s a 60 TPS
)

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

	// Mergulho do Grifo — o verbo de assinatura do jogo.
	//
	// Naves não mergulham. É a mecânica que faz o jogador sentir que monta um
	// animal, e a única do gênero que avança *contra* a rolagem da tela.
	// Custa fôlego, e o instante logo depois é de vulnerabilidade aumentada:
	// sem esse custo, invulnerabilidade sob demanda apagaria o desafio.
	diveDuration     = 22  // frames de avanço
	diveSpeed        = 6.2 // px/frame para frente
	diveDamage       = 14  // dano de contato ao atravessar um inimigo
	diveRecovery     = 14  // frames de recuperação após o mergulho
	diveCameraPush   = 10.0
	diveStaminaCost  = 34
	staminaMax       = 100
	staminaRegen     = 0.55 // por frame parado/voando normalmente
	staminaRegenSlow = 0.18 // enquanto se recupera de um mergulho

	invincibilityDuration = 60  // frames (1s a 60 TPS)
	invincibilityBlink    = 6   // frames por ciclo de piscada
	respawnInvincibility  = 150 // invencibilidade ao reaparecer
)

// Habilidade especial (Invocação Ancestral) e pontuação.
const (
	// Graze: raspar um projétil inimigo sem ser atingido carrega a Invocação
	// Ancestral. Converte medo em ganância — o jogador passa a *querer* chegar
	// perto — e faz as bombas circularem em vez de ficarem guardadas até o fim.
	grazeRadius    = 13.0
	grazePerBullet = 3.0 // carga por projétil raspado
	grazeFull      = 100.0
	grazeCooldown  = 24 // frames até o mesmo projétil poder render de novo

	bombMaxCharges     = 4
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

	// Altura (px) da faixa no rodapé em que um inimigo já é avisado como
	// "prestes a escapar". Sem esse aviso a corrupção sobe sem o jogador
	// entender por quê.
	escapeWarningBand = 46

	enemyBulletSize = 4
	// Projéteis inimigos precisam ser mais rápidos que o jogador (2.5), senão
	// dá para simplesmente sair andando na frente deles.
	enemyBulletSpeed  = 3.1
	enemyBulletDamage = 1

	// Vida e cadência calibradas pela seção "AMEACA" do `make balance`, que
	// separa as duas formas de um inimigo ser difícil:
	//
	//	ameaça  = consegue disparar 1–2 vezes antes de morrer
	//	esponja = fica viva absorvendo tiros, sem gerar tensão nenhuma
	//
	// Os valores originais (1 a 8 de vida) davam alvos inofensivos: uma salva da
	// Luz Nv3 matava qualquer inimigo e a campanha inteira gerava 57 projéteis.
	// A primeira correção passou do ponto e criou esponjas — a gárgula levava
	// 2,4s para morrer. O ajuste certo veio de **cadência**, não de vida.
	//
	// Cuidado ao mexer: o tempo-para-matar depende tanto da largura do inimigo
	// quanto da vida. As lanças da Luz cobrem 18px, então um wyvern de 24px
	// leva a salva inteira e um corvo de 10px leva metade — vida igual não
	// significa tempo de vida igual.
	crowSize   = 10
	crowHealth = 2
	crowScore  = 10
	crowSpeed  = 1.5 // mais lento: densidade alta, precisa de tempo de mira
	crowDamage = 1

	// Margem mínima das bordas para spawns "justos" (voadores/terrestres).
	// Evita corvos impossíveis nos extremos da tela de 240px.
	spawnFairMargin = 36

	// Densidade da campanha: mais aparições e intervalos um pouco menores (~+45%).
	campaignDensityScale = 1.45

	harpySize         = 14
	harpyHealth       = 10
	harpyScore        = 30
	harpySpeed        = 0.9
	harpyDamage       = 1
	harpyAmplitude    = 32
	harpyFreq         = 0.038
	harpyFireInterval = 38
	// Atraso do **primeiro** disparo, bem menor que o intervalo seguinte: sem
	// isso o inimigo morre antes de a primeira ameaça existir.
	harpyFirstFire = 18

	gargoyleSize           = 18
	gargoyleHealth         = 18
	gargoyleScore          = 50
	gargoyleSpeed          = 2.3
	gargoyleDamage         = 1
	gargoyleAttackDuration = 180 // ~3s atacando
	gargoyleFireInterval   = 30

	wyvernSize         = 24
	wyvernHealth       = 34
	wyvernScore        = 100
	wyvernDamage       = 1
	wyvernDescend      = 0.28
	wyvernAlign        = 0.8
	wyvernFireInterval = 48
	wyvernFirstFire    = 24

	// Balista corrompida: ameaça "terrestre" que desce lenta como o cenário e
	// dispara rajadas de virotes mirados. Resistente e bem telegrafada.
	ballistaSize         = 22
	ballistaHealth       = 26
	ballistaScore        = 70
	ballistaDamage       = 1
	ballistaDescend      = 0.35
	ballistaFireInterval = 100
	ballistaFirstFire    = 30
	ballistaBurst        = 3
	ballistaBurstGap     = 9

	// Feiticeiro corrompido: para no alto, desliza de lado e solta anéis de
	// projéteis. Recompensa o jogador por pressioná-lo rapidamente.
	mageSize         = 16
	mageHealth       = 18
	mageScore        = 90
	mageDamage       = 1
	mageDescend      = 1.0
	mageStopY        = 54
	mageDrift        = 0.75
	mageFireInterval = 115
	mageFirstFire    = 34
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

	devBossDamage = 100 // dano aplicado ao chefe pela tecla de desenvolvimento
)

// Chefe da fase 1 — Vharak, o Dragão Corrompido.
const (
	// Vida calibrada por medição (`go run ./cmd/balance`) para um confronto de
	// 70–120s com qualquer uma das três armas. Era 420, o que dava 5s de
	// gatilho segurado — o chefe morria sem executar um único padrão.
	bossMaxHealth = 2000
	bossW         = 60
	bossH         = 34
	bossY         = 26
	bossMargin    = 10

	bossEntrySpeed      = 0.6
	bossEntryHold       = 90 // frames exibindo o nome antes do combate
	bossWarningDuration = 45
	bossDeathDuration   = 120
	bossContactDamage   = 1

	// bossScore era 5000 contra 10–100 por inimigo comum: um único evento valia
	// mais que centenas de decisões de combate, e a pontuação deixava de medir
	// habilidade. 1500 mantém o chefe como o maior prêmio sem apagar a fase.
	bossScore = 1500

	// Dano da Invocação Ancestral ao chefe: relevante (~6% da vida) sem
	// transformar duas cargas em um atalho para o final.
	bombBossDamage = 120

	bossSpeedP1    = 0.65
	bossSpeedP2    = 1.35
	bossSpeedP3    = 2.0
	bossSweepSpeed = 2.2

	// Varredura: brecha larga e lenta garante que sempre haja rota de fuga.
	bossSweepStep    = 16
	bossSweepGapHalf = 28
	bossSweepGapMove = 1.0

	crystalBonus = 3 // multiplicador de dano ao acertar um cristal

	// Vharak Ascendido: o confronto verdadeiro, liberado com o reino totalmente
	// corrompido. Mais vida e mais velocidade — mas o jogador chega nele com o
	// multiplicador de pontos no topo.
	ascendedHealthMul = 1.6
	ascendedSpeedMul  = 1.25
)

// Fluxo de telas e bônus finais.
//
// Não há bônus por bomba não usada: ele pagava 300 pontos por carga guardada e
// portanto **pagava o jogador para não usar a mecânica mais espetacular do
// jogo**. Vidas restantes seguem valendo, porque premiar quem não morre é o
// incentivo correto.
const (
	fadeDuration = 30
	lifeBonus    = 500
)

// Configurações de armas e power-ups.
//
// Cada arma é dona de um cenário e perdedora nos outros. Sem essa assimetria
// não existe escolha — existe uma resposta certa. Os números abaixo são
// verificados por `go run ./cmd/balance` e pelo TestBalanceCriteria.
//
//	Luz    -> alvo único e distante (chefe)      | tiros retos, sem perfuração
//	Chamas -> vários alvos lado a lado (formação) | leque + queimadura
//	Gelo   -> alvos empilhados na vertical (fila) | perfuração + lentidão
const (
	maxWeaponLevel = 3

	// Runas necessárias para subir um nível. Com 1 runa por nível o teto de poder
	// chegava aos 7s de uma partida de 210s e ~53 runas caíam inertes depois
	// disso. Com 3, cada runa avança algo visível (as marcas no HUD) e o teto vai
	// para ~20% da run.
	//
	// Limite estrutural conhecido: com apenas três níveis, nenhuma quantidade de
	// runas por nível faz a progressão cobrir uma partida inteira. Quem resolve
	// isso é a profundidade extra das Runas Fundidas (v0.6), não este número.
	runesPerLevel = 3

	// Luz: a maior cadência e o maior dano focado. Não perfura (perfuração é a
	// identidade do Gelo) e não aplica status (o atordoamento congelava o jogo).
	lightBulletSpeed     = 6.6
	lightBulletSpeedFast = 8.6
	lightBulletDamage    = 3
	lightCooldown        = 10

	// Chamas: contagens ímpares para o leque sempre ter um tiro central — com
	// 4 projéteis o Nv2 não acertava alvos alinhados e era *pior* que o Nv1.
	flameSpread       = 0.52 // leque ~30°
	flameSpreadWide   = 0.72 // Nv3 abre o leque: subir de nível amplia a cobertura
	flameBulletSpeed  = 3.7
	flameBulletDamage = 2
	flameCooldown     = 14
	flameCooldownFast = 11 // Nv3 um pouco mais rápido, sem virar metralhadora
	flameCountBase    = 3
	flameCountMid     = 5
	flameCountMax     = 9

	// Gelo: o tiro mais pesado e o único que perfura fundo. Cadência baixa em
	// troca de valor altíssimo contra filas e formações em profundidade.
	iceBulletSpeed  = 3.8
	iceBulletDamage = 4
	iceCooldown     = 18
	icePierce       = 1
	icePierceMax    = 3

	// Efeitos elementais ao acertar (identidade de cada arma).
	iceSlowDuration = 100 // frames lento
	iceSlowFactor   = 0.4 // fração dos updates de movimento/tiro

	burnDuration = 130
	burnInterval = 12 // 1 de dano a cada N frames
	burnDamage   = 1
	burnBonusPct = 50 // +% de dano recebido enquanto queima

	// stunDuration segue disponível para uma habilidade ativa futura (Grito do
	// Grifo). Nenhuma arma o aplica passivamente.
	stunDuration = 32

	powerupSize  = 10
	powerupSpeed = 1.0

	// Chance de drop por abate. Era 0.20: com ~460 abates caíam ~92 runas e o
	// teto de poder chegava aos 60s de uma partida de 300s.
	dropChance     = 0.07
	shieldDuration = 480 // ~8s a 60 TPS
)

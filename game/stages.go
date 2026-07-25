package game

// Campanha em dados.
//
// Regra de composição (verificada por `go run ./cmd/balance` e por
// TestBestiaryIsBalanced): **cada inimigo é apresentado sozinho, depois
// combinado, depois em massa.** É a gramática de introdução dos shmups
// clássicos, e é de graça — os seis inimigos já existem.
//
// Antes desta revisão, 63,5% da campanha eram corvos e 31,1% harpias: 94,6% de
// cinco minutos de jogo em dois inimigos, enquanto gárgula, wyvern, feiticeiro
// e balista somavam 25 aparições. O bestiário existia e não era usado.
//
// Lembre que `campaignDensityScale` multiplica as contagens maiores que 1; as
// aparições únicas (chefes de onda, ameaças pesadas) não escalam.

// campaignStages devolve a campanha na ordem de jogo. A última fase é a que
// invoca o chefe; as anteriores encadeiam para a seguinte ao serem concluídas.
func campaignStages() []*stageDef {
	stages := []*stageDef{
		stage1(),
		stage2(),
		stage3(),
	}
	stages[len(stages)-1].hasBoss = true
	return stages
}

// stage1 — O Cerco de Eldoria: apresenta corvo, harpia, gárgula e wyvern, um
// por trecho, cada um primeiro sozinho.
func stage1() *stageDef {
	b := newStageBuilder("O Cerco de Eldoria")

	b.section(0, "Campos do reino", "", sectionThemes[0], musicFields)
	b.section(1500, "Vila atacada", "A vila esta sob ataque!", sectionThemes[1], musicVillage)
	b.section(3000, "Muralhas", "Aproximando-se das muralhas", sectionThemes[2], musicWalls)
	b.section(4500, "Rumo ao castelo", "O castelo se aproxima...", sectionThemes[3], musicCastle)

	// --- Campos: só corvos. O jogador aprende a mirar e a ler formações. ---
	b.wave(60, kindCrow, 3, 45, formationSingle, 115, false)
	b.wave(300, kindCrow, 3, 40, formationSingle, 85, false)
	b.wave(330, kindCrow, 3, 40, formationSingle, 150, false)
	b.wave(620, kindCrow, 1, 0, formationV, 110, false)
	b.drop(powerFire)
	b.wave(780, kindCrow, 4, 38, formationSingle, 95, false)
	b.wave(1040, kindCrow, 1, 0, formationLine, 70, false)
	b.wave(1240, kindCrow, 3, 36, formationSingle, 160, false)

	// --- Vila: a HARPIA entra sozinha, depois acompanhada, depois em grupo. ---
	b.wave(1520, kindHarpy, 1, 0, formationSingle, 120, false)
	b.drop(powerIce)
	b.wave(1720, kindHarpy, 2, 70, formationSingle, 70, false)
	b.wave(1780, kindCrow, 3, 40, formationSingle, 170, false)
	b.wave(2000, kindHarpy, 2, 60, formationSingle, 140, false)
	b.wave(2060, kindCrow, 1, 0, formationV, 100, false)
	b.wave(2320, kindHarpy, 3, 55, formationSingle, 60, false)
	b.drop(powerHeal)
	b.wave(2440, kindCrow, 3, 36, formationSingle, 110, false)
	b.wave(2700, kindHarpy, 4, 50, formationSingle, 90, false)

	// --- Muralhas: a GÁRGULA entra sozinha e vira presença constante. ---
	b.wave(3020, kindGargoyle, 1, 0, formationSingle, 90, true)
	b.drop(powerLight)
	b.wave(3240, kindGargoyle, 1, 0, formationSingle, 150, false)
	b.wave(3300, kindHarpy, 2, 55, formationSingle, 60, false)
	b.wave(3540, kindGargoyle, 2, 120, formationSingle, 70, true)
	b.wave(3600, kindCrow, 3, 36, formationSingle, 160, false)
	b.wave(3840, kindGargoyle, 2, 110, formationSingle, 160, false)
	b.drop(powerHeal)
	b.wave(3900, kindHarpy, 3, 55, formationSingle, 100, false)
	b.wave(4120, kindGargoyle, 2, 100, formationSingle, 60, true)
	b.wave(4180, kindHarpy, 2, 50, formationSingle, 150, false)

	// --- Castelo: o WYVERN entra sozinho e passa a liderar as investidas. ---
	b.wave(4520, kindWyvern, 1, 0, formationSingle, 120, false)
	b.drop(powerShield)
	b.wave(4780, kindWyvern, 1, 0, formationSingle, 70, false)
	b.wave(4840, kindHarpy, 2, 55, formationSingle, 170, false)
	b.wave(5060, kindWyvern, 2, 150, formationSingle, 150, false)
	b.wave(5120, kindGargoyle, 1, 0, formationSingle, 60, true)
	b.wave(5180, kindCrow, 3, 36, formationSingle, 110, false)
	b.wave(5420, kindWyvern, 2, 140, formationSingle, 80, false)
	b.drop(powerHeal)
	b.wave(5480, kindHarpy, 3, 50, formationSingle, 160, false)
	b.wave(5540, kindGargoyle, 1, 0, formationSingle, 130, false)
	b.wave(5780, kindWyvern, 2, 130, formationSingle, 110, false)

	return b.def()
}

// stage2 — A Floresta Corrompida: apresenta o FEITICEIRO e passa a combinar as
// ameaças pesadas da fase anterior.
func stage2() *stageDef {
	b := newStageBuilder("A Floresta Corrompida")

	b.section(0, "Bosque sombrio", "Adentrando a floresta corrompida", themeForest, musicForest)
	b.section(1600, "Pantano corrompido", "A corrupcao se aprofunda...", themeSwamp, musicSwamp)

	// --- Bosque: o FEITICEIRO entra sozinho. Anéis completos exigem posição. ---
	b.wave(40, kindMage, 1, 0, formationSingle, 120, false)
	b.drop(powerIce)
	b.wave(280, kindHarpy, 3, 55, formationSingle, 70, false)
	b.wave(340, kindCrow, 3, 36, formationSingle, 170, false)
	b.wave(560, kindMage, 1, 0, formationSingle, 80, false)
	b.wave(620, kindHarpy, 2, 50, formationSingle, 150, false)
	b.wave(860, kindMage, 2, 200, formationSingle, 150, false)
	b.drop(powerHeal)
	b.wave(920, kindCrow, 3, 34, formationSingle, 60, false)
	b.wave(1140, kindWyvern, 1, 0, formationSingle, 110, false)
	b.wave(1380, kindMage, 2, 190, formationSingle, 90, false)
	b.wave(1440, kindCrow, 3, 34, formationSingle, 165, false)

	// --- Pântano: feiticeiro + gárgula + wyvern juntos. A fase aperta. ---
	b.wave(1620, kindMage, 2, 180, formationSingle, 60, false)
	b.drop(powerFire)
	b.wave(1700, kindGargoyle, 1, 0, formationSingle, 150, false)
	b.wave(1940, kindWyvern, 2, 150, formationSingle, 100, false)
	b.wave(2000, kindHarpy, 3, 50, formationSingle, 170, false)
	b.wave(2240, kindMage, 2, 170, formationSingle, 170, false)
	b.wave(2300, kindGargoyle, 2, 120, formationSingle, 70, true)
	b.wave(2360, kindCrow, 3, 34, formationSingle, 120, false)
	b.wave(2640, kindWyvern, 2, 140, formationSingle, 60, false)
	b.drop(powerShield)
	b.wave(2700, kindMage, 2, 160, formationSingle, 140, false)
	b.wave(2940, kindGargoyle, 2, 110, formationSingle, 110, false)
	b.wave(3000, kindHarpy, 2, 48, formationSingle, 60, false)
	b.wave(3060, kindCrow, 1, 0, formationV, 110, false)

	return b.def()
}

// stage3 — O Covil de Vharak: apresenta a BALISTA e reúne todas as ameaças
// antes do chefe.
func stage3() *stageDef {
	b := newStageBuilder("O Covil de Vharak")

	b.section(0, "Desfiladeiro", "O covil do dragao se aproxima", themeCanyon, musicCanyon)
	b.section(1900, "Covil do dragao", "Vharak aguarda...", themeLair, musicLair)

	// --- Desfiladeiro: a BALISTA entra sozinha, bem telegrafada. ---
	b.wave(50, kindBallista, 1, 0, formationSingle, 120, false)
	b.drop(powerLight)
	b.wave(300, kindBallista, 1, 0, formationSingle, 70, false)
	b.wave(360, kindHarpy, 2, 50, formationSingle, 165, false)
	b.wave(600, kindBallista, 2, 170, formationSingle, 150, false)
	b.wave(660, kindCrow, 3, 34, formationSingle, 70, false)
	b.wave(900, kindWyvern, 2, 150, formationSingle, 110, false)
	b.drop(powerHeal)
	b.wave(960, kindBallista, 2, 160, formationSingle, 60, false)
	b.wave(1200, kindMage, 2, 170, formationSingle, 160, false)
	b.wave(1260, kindGargoyle, 2, 120, formationSingle, 80, true)
	b.wave(1500, kindBallista, 2, 150, formationSingle, 110, false)
	b.wave(1560, kindHarpy, 2, 50, formationSingle, 60, false)
	b.wave(1620, kindCrow, 3, 34, formationSingle, 170, false)

	// --- Covil: tudo ao mesmo tempo, subindo até a entrada de Vharak. ---
	b.wave(1920, kindBallista, 2, 150, formationSingle, 60, false)
	b.drop(powerShield)
	b.wave(1980, kindMage, 2, 160, formationSingle, 175, false)
	b.wave(2220, kindWyvern, 2, 140, formationSingle, 110, false)
	b.wave(2280, kindGargoyle, 2, 110, formationSingle, 70, true)
	b.wave(2340, kindHarpy, 2, 48, formationSingle, 150, false)
	b.wave(2580, kindBallista, 2, 140, formationSingle, 160, false)
	b.drop(powerHeal)
	b.wave(2640, kindMage, 2, 150, formationSingle, 70, false)
	b.wave(2700, kindCrow, 3, 34, formationSingle, 120, false)
	b.wave(2940, kindWyvern, 3, 130, formationSingle, 60, false)
	b.wave(3000, kindGargoyle, 2, 110, formationSingle, 150, false)
	b.wave(3060, kindHarpy, 2, 48, formationSingle, 110, false)
	b.wave(3300, kindBallista, 2, 130, formationSingle, 90, false)
	b.drop(powerFire)
	b.wave(3360, kindMage, 2, 140, formationSingle, 160, false)
	b.wave(3420, kindCrow, 1, 0, formationV, 110, false)

	return b.def()
}

package game

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

// stage1 — O Cerco de Eldoria: introdução que sobe de corvos a wyverns.
// Densidade alta e fuga pela base pune — inimigos descem um pouco mais lentos.
func stage1() *stageDef {
	b := newStageBuilder("O Cerco de Eldoria")

	b.section(0, "Campos do reino", "", sectionThemes[0])
	b.section(1500, "Vila atacada", "A vila esta sob ataque!", sectionThemes[1])
	b.section(3000, "Muralhas", "Aproximando-se das muralhas", sectionThemes[2])
	b.section(4500, "Rumo ao castelo", "O castelo se aproxima...", sectionThemes[3])

	// Campos: abre no centro, logo vira pressão em duas pistas + formações.
	b.wave(60, kindCrow, 5, 38, formationSingle, 115, false)
	b.wave(240, kindCrow, 5, 34, formationSingle, 85, false)
	b.wave(255, kindCrow, 5, 34, formationSingle, 145, false)
	b.wave(480, kindCrow, 1, 0, formationV, 110, false)
	b.drop(powerFire)
	b.wave(520, kindCrow, 6, 32, formationSingle, 95, false)
	b.wave(540, kindCrow, 5, 34, formationSingle, 155, false)
	b.wave(780, kindCrow, 1, 0, formationLine, 70, false)
	b.wave(820, kindCrow, 5, 36, formationSingle, 120, false)
	b.wave(900, kindHarpy, 3, 50, formationSingle, 130, false)
	b.wave(1100, kindCrow, 6, 30, formationSingle, 100, false)
	b.wave(1120, kindHarpy, 2, 55, formationSingle, 70, false)

	// Vila: harpias + bandos densos.
	b.wave(1520, kindHarpy, 4, 48, formationSingle, 50, false)
	b.drop(powerIce)
	b.wave(1550, kindCrow, 7, 28, formationSingle, 170, false)
	b.wave(1750, kindCrow, 1, 0, formationV, 110, false)
	b.wave(1900, kindHarpy, 3, 50, formationSingle, 140, false)
	b.wave(1920, kindCrow, 6, 30, formationSingle, 80, false)
	b.wave(2200, kindHarpy, 4, 45, formationSingle, 70, false)
	b.drop(powerHeal)
	b.wave(2240, kindCrow, 7, 28, formationSingle, 160, false)
	b.wave(2500, kindCrow, 1, 0, formationLine, 60, false)
	b.wave(2550, kindHarpy, 3, 50, formationSingle, 120, false)

	// Muralhas: gárgulas + pressão constante.
	b.wave(3020, kindGargoyle, 1, 0, formationSingle, 70, true)
	b.drop(powerLight)
	b.wave(3020, kindHarpy, 4, 48, formationSingle, 140, false)
	b.wave(3060, kindCrow, 7, 28, formationSingle, 50, false)
	b.wave(3300, kindCrow, 1, 0, formationV, 110, false)
	b.wave(3450, kindGargoyle, 1, 0, formationSingle, 170, false)
	b.wave(3480, kindHarpy, 4, 45, formationSingle, 100, false)
	b.wave(3520, kindCrow, 6, 30, formationSingle, 175, false)
	b.wave(3800, kindGargoyle, 1, 0, formationSingle, 110, true)
	b.drop(powerHeal)
	b.wave(3840, kindHarpy, 4, 48, formationSingle, 55, false)
	b.wave(3880, kindCrow, 6, 30, formationSingle, 130, false)

	// Castelo: wyverns com apoio denso.
	b.wave(4520, kindWyvern, 1, 0, formationSingle, 60, false)
	b.drop(powerShield)
	b.wave(4520, kindHarpy, 4, 48, formationSingle, 170, false)
	b.wave(4560, kindCrow, 7, 28, formationSingle, 110, false)
	b.wave(4900, kindWyvern, 1, 0, formationSingle, 150, false)
	b.wave(4920, kindHarpy, 5, 42, formationSingle, 50, false)
	b.wave(4960, kindCrow, 1, 0, formationLine, 80, false)
	b.wave(5200, kindCrow, 6, 30, formationSingle, 140, false)
	b.wave(5350, kindWyvern, 1, 0, formationSingle, 100, false)
	b.wave(5380, kindHarpy, 4, 48, formationSingle, 130, false)
	b.wave(5420, kindCrow, 5, 32, formationSingle, 70, false)

	return b.def()
}

// stage2 — A Floresta Corrompida: feiticeiros + voadores em ritmo alto.
func stage2() *stageDef {
	b := newStageBuilder("A Floresta Corrompida")

	b.section(0, "Bosque sombrio", "Adentrando a floresta corrompida", themeForest)
	b.section(1600, "Pantano corrompido", "A corrupcao se aprofunda...", themeSwamp)

	b.wave(40, kindHarpy, 5, 42, formationSingle, 50, false)
	b.drop(powerIce)
	b.wave(60, kindCrow, 7, 28, formationSingle, 170, false)
	b.wave(280, kindCrow, 1, 0, formationV, 110, false)
	b.wave(420, kindMage, 1, 0, formationSingle, 120, false)
	b.wave(450, kindHarpy, 4, 48, formationSingle, 70, false)
	b.wave(480, kindCrow, 7, 28, formationSingle, 175, false)
	b.wave(750, kindHarpy, 4, 45, formationSingle, 50, false)
	b.drop(powerHeal)
	b.wave(780, kindCrow, 1, 0, formationLine, 70, false)
	b.wave(900, kindMage, 1, 0, formationSingle, 90, false)
	b.wave(940, kindCrow, 6, 30, formationSingle, 150, false)
	b.wave(1100, kindHarpy, 4, 48, formationSingle, 130, false)

	b.wave(1620, kindMage, 1, 0, formationSingle, 70, false)
	b.wave(1650, kindHarpy, 5, 42, formationSingle, 170, false)
	b.drop(powerFire)
	b.wave(1680, kindCrow, 7, 28, formationSingle, 110, false)
	b.wave(2000, kindWyvern, 1, 0, formationSingle, 120, false)
	b.wave(2030, kindHarpy, 4, 48, formationSingle, 50, false)
	b.wave(2060, kindCrow, 6, 30, formationSingle, 160, false)
	b.wave(2400, kindMage, 1, 0, formationSingle, 150, false)
	b.wave(2440, kindCrow, 7, 28, formationSingle, 70, false)
	b.wave(2700, kindWyvern, 1, 0, formationSingle, 70, false)
	b.drop(powerShield)
	b.wave(2740, kindWyvern, 1, 0, formationSingle, 170, false)
	b.wave(2780, kindHarpy, 5, 42, formationSingle, 110, false)
	b.wave(3000, kindCrow, 1, 0, formationV, 110, false)

	return b.def()
}

// stage3 — O Covil de Vharak: balistas e todas as ameaças antes do chefe.
func stage3() *stageDef {
	b := newStageBuilder("O Covil de Vharak")

	b.section(0, "Desfiladeiro", "O covil do dragao se aproxima", themeCanyon)
	b.section(1900, "Covil do dragao", "Vharak aguarda...", themeLair)

	b.wave(50, kindBallista, 1, 0, formationSingle, 60, false)
	b.drop(powerLight)
	b.wave(70, kindHarpy, 5, 42, formationSingle, 170, false)
	b.wave(100, kindCrow, 7, 28, formationSingle, 110, false)
	b.wave(400, kindCrow, 1, 0, formationV, 110, false)
	b.wave(550, kindGargoyle, 1, 0, formationSingle, 120, true)
	b.wave(580, kindMage, 1, 0, formationSingle, 50, false)
	b.wave(620, kindHarpy, 4, 48, formationSingle, 180, false)
	b.wave(900, kindBallista, 1, 0, formationSingle, 150, false)
	b.drop(powerHeal)
	b.wave(940, kindWyvern, 1, 0, formationSingle, 90, false)
	b.wave(980, kindCrow, 7, 28, formationSingle, 50, false)
	b.wave(1300, kindHarpy, 5, 42, formationSingle, 130, false)
	b.wave(1340, kindCrow, 6, 30, formationSingle, 170, false)

	b.wave(1920, kindBallista, 1, 0, formationSingle, 50, false)
	b.wave(1950, kindMage, 1, 0, formationSingle, 180, false)
	b.drop(powerShield)
	b.wave(1980, kindHarpy, 5, 42, formationSingle, 110, false)
	b.wave(2200, kindCrow, 1, 0, formationLine, 70, false)
	b.wave(2400, kindGargoyle, 1, 0, formationSingle, 90, false)
	b.wave(2440, kindWyvern, 1, 0, formationSingle, 50, false)
	b.wave(2480, kindCrow, 7, 28, formationSingle, 175, false)
	b.wave(2800, kindWyvern, 1, 0, formationSingle, 70, false)
	b.drop(powerHeal)
	b.wave(2840, kindWyvern, 1, 0, formationSingle, 160, false)
	b.wave(2880, kindMage, 1, 0, formationSingle, 110, false)
	b.wave(2920, kindHarpy, 5, 42, formationSingle, 50, false)
	b.wave(3100, kindCrow, 6, 30, formationSingle, 130, false)

	return b.def()
}

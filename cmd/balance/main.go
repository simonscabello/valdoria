// Comando de medição de balanceamento.
//
// Roda a lógica real do jogo sem janela e imprime um relatório com DPS por
// arma, duração do confronto com o chefe, composição do bestiário e curva de
// poder. Serve para decidir balanceamento com números, não com impressões.
//
//	go run ./cmd/balance              # dificuldade Normal
//	go run ./cmd/balance -diff facil  # outra dificuldade
//	go run ./cmd/balance -seed 42     # semente fixa
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"valdoria/game"
)

func main() {
	diffName := flag.String("diff", "normal", "dificuldade: facil, normal ou dificil")
	seed := flag.Int64("seed", 42, "semente da simulação")
	flag.Parse()

	diff, ok := game.ParseDifficulty(*diffName)
	if !ok {
		fmt.Fprintf(os.Stderr, "dificuldade desconhecida: %q (use facil, normal ou dificil)\n", *diffName)
		os.Exit(2)
	}

	r := game.MeasureBalance(diff, *seed)

	title("ASAS DE VALDORIA - RELATORIO DE BALANCEAMENTO")
	fmt.Printf("Dificuldade: %s   Semente: %d\n", r.Difficulty, *seed)

	printWeapons(r)
	printBoss(r)
	printBestiary(r)
	printProgression(r)
	printCorruption(r)
	printThreat(r)
	printPressure(r)
	printVerdict(r)
}

func title(s string) {
	fmt.Printf("\n%s\n%s\n", s, strings.Repeat("=", len(s)))
}

func section(s string) {
	fmt.Printf("\n%s\n%s\n", s, strings.Repeat("-", len(s)))
}

func printWeapons(r game.BalanceReport) {
	section("1. ARMAS - dano por segundo em cada cenario")
	fmt.Println("   Foco     = alvo unico largo e distante (chefe)   -> dominio esperado: Luz")
	fmt.Println("   Formacao = cinco alvos lado a lado               -> dominio esperado: Chamas")
	fmt.Println("   Coluna   = quatro alvos empilhados na vertical   -> dominio esperado: Gelo")
	fmt.Printf("\n%-12s %-3s %8s %9s %8s %8s %10s %6s\n",
		"ARMA", "NV", "FOCO", "FORMACAO", "COLUNA", "PICO", "MELHOR EM", "COOLD")
	for _, w := range r.Weapons {
		fmt.Printf("%-12s %-3d %8.1f %9.1f %8.1f %8.1f %10s %6d\n",
			w.Weapon, w.Level, w.FocusDPS, w.SwarmDPS, w.ColumnDPS,
			w.Peak(), w.BestScenario(), w.Cooldown)
	}

	fmt.Printf("\nDesequilibrio entre os picos no Nv3: %.2fx\n", game.PeakSpread(r.Weapons, 3))
}

func printBoss(r game.BalanceReport) {
	section("2. CHEFE - duracao real do confronto")
	fmt.Printf("%-12s %12s %12s %10s %8s\n", "ARMA", "GATILHO 100%", "REALISTA 50%", "PADROES", "FASE 3")
	for _, b := range r.Boss {
		fase3 := "nao"
		if b.ReachedPhase3 {
			fase3 = "sim"
		}
		fmt.Printf("%-12s %11.1fs %11.1fs %10d %8s\n",
			b.Weapon, b.FullUptime, b.Realistic, b.Patterns, fase3)
	}
	fmt.Println()
	fmt.Println("Vharak Ascendido (corrupcao 100%):")
	for _, b := range r.BossAscended {
		fmt.Printf("%-12s %11.1fs %11.1fs %10d %8s\n",
			b.Weapon, b.FullUptime, b.Realistic, b.Patterns, "sim")
	}
	fmt.Println("\nAlvo de design: 60-180s no cenario realista, com todos os padroes executados.")
}

func printCorruption(r game.BalanceReport) {
	section("5. MEDIDOR DE CORRUPCAO")
	fmt.Println("   Quanto o reino apodrece conforme o jogador deixa inimigos passarem.")
	fmt.Printf("\n%-8s %8s %8s %20s %8s %10s\n",
		"FUGAS", "PASSARAM", "FINAL", "FAIXA", "PONTOS", "ASCENDIDO")
	for _, c := range r.Corruption {
		asc := "nao"
		if c.AscendedBoss {
			asc = "SIM"
		}
		fmt.Printf("%7.0f%% %8d %7.0f%% %20s %7.1fx %10s\n",
			c.LeakRate*100, c.Escaped, c.Final, c.Tier, c.ScoreMul, asc)
	}
}

func printThreat(r game.BalanceReport) {
	section("6. AMEACA - perigoso ou so demorado?")
	fmt.Println("   TTK   = segundos para morrer sob fogo focado da Luz Nv3 (melhor caso)")
	fmt.Println("   TIROS = quantos projeteis ele consegue disparar antes de morrer")
	fmt.Println("   Alvo: 1-2 tiros. Zero = alvo inofensivo. 4+ = esponja.")
	fmt.Printf("\n%-14s %6s %7s %8s %9s %10s\n",
		"INIMIGO", "VIDA", "TTK", "TIROS", "TRAVESSIA", "SOB FOGO")
	for _, t := range r.Threat {
		fmt.Printf("%-14s %6d %6.2fs %8d %8.1fs %9.0f%%\n",
			t.Kind, t.Health, t.TTK, t.Shots, t.Lifetime, t.Pressure*100)
	}
}

func printPressure(r game.BalanceReport) {
	section("7. PRESSAO - o jogo e dificil de verdade?")
	fmt.Println("   Jogador simulado com poder de fogo maximo e ZERO desvio.")
	fmt.Println("   passivo = nao se move | mediano = persegue alvos, mas nunca desvia")
	fmt.Printf("\n%-10s %7s %7s %8s %10s %9s %6s %8s\n",
		"MODELO", "GOLPES", "VIDAS", "MORREU", "SOBREVIVEU", "PROJETEIS", "PICO", "FUGAS")
	for _, p := range r.Pressure {
		died := "nao"
		if p.Died {
			died = "SIM"
		}
		fmt.Printf("%-10s %7d %7d %8s %9.0fs %9d %6d %8d\n",
			p.Model, p.Hits, p.Lives, died, p.SurvivedSeconds,
			p.BulletsFired, p.PeakBullets, p.Escapes)
	}
	fmt.Println("\nAlvo: o passivo morre cedo e o mediano nao termina a campanha.")
	fmt.Println("Um jogador que desvia deve conseguir terminar — isso o playtest confirma.")
}

func printBestiary(r game.BalanceReport) {
	section("3. BESTIARIO - o que a campanha realmente gera")
	fmt.Printf("%-26s %6s %8s %8s %8s %7s\n", "FASE", "ONDAS", "INIMIGOS", "VIDA", "SEGUNDOS", "PICO")
	for _, s := range r.Stages {
		fmt.Printf("%-26s %6d %8d %8d %8.0f %7d\n",
			trunc(s.Name, 26), s.Waves, s.Enemies, s.TotalHealth, s.Seconds, s.PeakOnScreen)
	}
	t := r.Totals
	fmt.Printf("%-26s %6d %8d %8d %8.0f %7d\n", "TOTAL", t.Waves, t.Enemies, t.TotalHealth, t.Seconds, t.PeakOnScreen)

	fmt.Printf("\n%-14s %8s %9s %s\n", "INIMIGO", "QTD", "SHARE", "ALVO")
	for _, k := range game.EnemyKinds() {
		fmt.Printf("%-14s %8d %8.1f%%   %s\n",
			game.EnemyKindName(k), t.ByKind[k], t.Share(k)*100, game.TargetShareLabel(k))
	}
}

func printProgression(r game.BalanceReport) {
	p := r.Progression
	section("4. CURVA DE PODER")
	fmt.Printf("Inimigos na campanha .................. %d\n", p.TotalEnemies)
	fmt.Printf("Duracao das fases ..................... %.0fs\n", p.CampaignSeconds)
	fmt.Printf("Runas garantidas por onda ............. %d\n", p.GuaranteedDrops)
	fmt.Printf("Runas coletadas numa run completa ..... %d\n", p.TotalDrops)
	fmt.Printf("Runas do elemento ate o nivel maximo .. %d\n", p.RunesToMax)
	fmt.Printf("Poder maximo atingido no abate ........ %d de %d\n", p.KillsToMax, p.TotalEnemies)
	fmt.Printf("Poder maximo atingido aos ............. %.0fs de %.0fs (%.0f%% da run)\n",
		p.SecondsToMax, p.CampaignSeconds, p.SecondsToMax/p.CampaignSeconds*100)
	fmt.Printf("Runas inertes apos o teto ............. %d\n", p.WastedRunes)
}

func printVerdict(r game.BalanceReport) {
	section("8. VEREDITO")
	ok := true
	for _, c := range game.Checks(r) {
		mark := "OK  "
		if !c.Pass {
			mark = "FALHA"
			ok = false
		}
		fmt.Printf("[%-5s] %s\n         %s\n", mark, c.Name, c.Detail)
	}
	fmt.Println()
	if ok {
		fmt.Println("Todos os criterios de balanceamento foram atendidos.")
		return
	}
	fmt.Println("Ha criterios de balanceamento nao atendidos (ver acima).")
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "."
}

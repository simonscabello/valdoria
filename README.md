# Valdoria

Protótipo de um shoot 'em up vertical em Golang, inspirado na jogabilidade de
*Super Aleste*, com identidade própria e temática medieval fantástica. O jogador
controla um cavaleiro montado em um grifo que enfrenta criaturas aladas
corrompidas descendo pela tela.

Nesta etapa inicial o jogo é totalmente representado por formas geométricas
simples, sem imagens, sprites ou áudio.

## Tecnologias

- [Go](https://go.dev/) 1.26+
- [Ebitengine](https://ebitengine.org/) (`github.com/hajimehoshi/ebiten/v2`)

## Instalando as dependências

```bash
go mod download
```

No Linux, o Ebitengine depende de bibliotecas de sistema. Se necessário:

```bash
sudo apt install libgl1-mesa-dev xorg-dev libasound2-dev
```

## Executando o jogo

```bash
go run .
```

## Verificações

```bash
go fmt ./...
go vet ./...
go test ./...
```

## Controles

| Ação                | Teclas                         |
| ------------------- | ------------------------------ |
| Mover               | `W` `A` `S` `D` ou setas       |
| Modo de precisão    | Segurar `Shift`                |
| Atirar              | `Espaço` (contínuo ao segurar) |
| Pausar / Despausar  | `Esc`                          |
| Reiniciar           | `Enter` (na tela de Game Over) |
| Forçar inimigo (debug) | `1` corvos · `2` harpia · `3` gárgula · `4` wyvern |
| Gerar power-up (dev)   | `Z` luz · `X` fogo · `C` gelo · `V` cura · `B` escudo |
| Acelerar a fase (debug) | `Tab` |

## Funcionalidades implementadas

- Janela em resolução interna para pixel art com escala 3x.
- Game loop com `Update`, `Draw` e `Layout`.
- Jogador em oito direções com velocidade normalizada na diagonal.
- Modo de precisão (`Shift`): move mais devagar e revela o ponto real de colisão.
- Limitação do jogador à área visível.
- Disparo duplo com cooldown, tiro contínuo e efeito visual de fogo.
- Invencibilidade temporária após dano, com o jogador piscando e ignorando novas colisões.
- Pausa/despausa com `Esc` que congela toda a lógica do jogo.
- Projéteis que sobem e são removidos ao sair da tela.
- Fundo escuro com estrelas rolando verticalmente.
- Quatro tipos de inimigos com vida, pontuação, dano e padrões próprios:
  - **Corvo corrompido**: pequeno, frágil, desce reto e aparece em formações.
  - **Harpia**: vida média, zigue-zague e disparo reto para baixo.
  - **Gárgula**: resistente, entra pela lateral, ataca parada e sai pelo lado oposto.
  - **Wyvern**: grande e resistente, alinha-se ao jogador e dispara projéteis mirados (sem perseguir depois de lançados).
- Projéteis inimigos e indicação visual de dano (piscada branca) ao acertar um inimigo.
- Colisão projétil × inimigo e jogador × inimigo/projétil, sem dano repetido de inimigos já destruídos.
- **Fase 1 — O Cerco de Eldoria**: onda determinística baseada em linha do tempo, com quatro trechos (Campos, Vila, Muralhas, Aproximação do castelo), progressão de dificuldade, avisos de trecho, barra de progresso, mudança de cor de fundo por trecho e preparação da entrada do chefe ao concluir a fase.
- **Três magias, cada uma com três níveis**:
  - **Lança de Luz** (inicial): disparos retos (2 → 3 → 4 projéteis, mais velozes no nível 3).
  - **Chamas do Dragão**: leque (3 → 5 direções, maior cadência no nível 3), dano individual menor, boa cobertura.
  - **Lanças de Gelo**: lentas e fortes (1 → 2 → 3 projéteis), atravessam inimigos (mais perfuração no nível 3).
- **Power-ups** que caem de inimigos e saem da tela se não coletados: Runa da Luz/Fogo/Gelo (coletar a runa da arma atual sobe o nível, outra runa troca a arma para o nível 1, máximo 3), Cura (limitada à vida máxima) e Escudo temporário (absorve um ataque). Chance de drop configurável e alguns drops garantidos por onda.
- HUD com arma, nível e estado do escudo.
- Sistema básico de vida e pontuação.
- Tela de Game Over e reinício com `Enter`.

A lógica roda em TPS fixo do Ebitengine (60), independente da taxa de renderização.
A remoção de projéteis e inimigos reaproveita os slices e libera as referências
antigas para o coletor de lixo.

## Próximas etapas sugeridas

- Chefe da Fase 1 (a entrada já é preparada ao fim da fase).
- Sprites e animações para o cavaleiro/grifo e inimigos.
- Áudio (efeitos e trilha).
- Novas fases e mais padrões de ataque.
- Menu inicial e HUD aprimorada.

## Estrutura do projeto

```
valdoria/
├── main.go            # ponto de entrada e configuração da janela
├── go.mod
├── README.md
└── game/
    ├── game.go         # loop principal, estados, spawn, colisões e HUD
    ├── config.go       # valores centrais (jogador, combate e inimigos)
    ├── player.go       # cavaleiro/grifo (movimento, precisão, disparo, dano)
    ├── enemy.go        # tipos de inimigos e seus comportamentos
    ├── level.go        # fase 1, linha do tempo de ondas e trechos
    ├── weapon.go       # as três magias, níveis e leque
    ├── powerup.go      # runas, cura e escudo
    ├── enemybullet.go  # projéteis inimigos
    ├── bullet.go       # projéteis do jogador (velocidade, dano, perfuração)
    ├── background.go   # estrelas de fundo
    ├── game_test.go    # testes das funções puras (colisão, remoção, HUD)
    ├── enemy_test.go   # testes de inimigos (dano, morte, pontuação, mira, movimento)
    ├── level_test.go   # testes da fase (eventos, ondas, trechos, pausa, chefe)
    └── weapon_test.go  # testes de armas/power-ups (troca, nível, cura, escudo, perfuração, leque)
```

## Desenvolvimento e testes da fase

Em `game/config.go`:

- `devMode`: habilita as teclas de gerar power-ups (`Z`/`X`/`C`/`V`/`B`).
- `devStartSection` (0 a 3): inicia a fase direto em um trecho específico.
- `devTimeScale`: passos da linha do tempo por frame (padrão 1).
- `devFastTimeScale`: valor aplicado ao acelerar a fase com `Tab`.
- `dropChance`: probabilidade de um inimigo comum largar uma runa aleatória.

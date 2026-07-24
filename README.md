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
- Sistema básico de vida e pontuação.
- Tela de Game Over e reinício com `Enter`.

A lógica roda em TPS fixo do Ebitengine (60), independente da taxa de renderização.
A remoção de projéteis e inimigos reaproveita os slices e libera as referências
antigas para o coletor de lixo.

## Próximas etapas sugeridas

- Sprites e animações para o cavaleiro/grifo e inimigos.
- Áudio (efeitos e trilha).
- Tipos variados de inimigos (harpias, dragões, gárgulas) e padrões de ataque.
- Chefes e sistema de fases.
- Diferentes armas e power-ups.
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
    ├── enemybullet.go  # projéteis inimigos
    ├── bullet.go       # projéteis do jogador
    ├── background.go   # estrelas de fundo
    ├── game_test.go    # testes das funções puras (colisão, remoção, HUD)
    └── enemy_test.go   # testes de inimigos (dano, morte, pontuação, mira, movimento)
```

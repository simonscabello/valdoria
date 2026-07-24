# Asas de Valdoria

Shoot 'em up vertical em Golang, inspirado na jogabilidade de *Super Aleste*,
com identidade própria e temática medieval fantástica. O jogador controla um
cavaleiro montado em um grifo que enfrenta criaturas aladas corrompidas
descendo pela tela, até o confronto com o chefe da fase, Vharak.

Toda a arte é feita com formas geométricas simples (sem sprites externos) e o
áudio é gerado proceduralmente, podendo ser substituído por arquivos livres.

## Capturas

As capturas de tela ficam em `docs/screenshots/` (a pasta pode ser criada ao
gerar as imagens). Espaço reservado:

| Menu | Fase | Chefe |
| --- | --- | --- |
| `docs/screenshots/menu.png` | `docs/screenshots/fase.png` | `docs/screenshots/chefe.png` |

## Requisitos

- [Go](https://go.dev/) 1.26+
- [Ebitengine](https://ebitengine.org/) (`github.com/hajimehoshi/ebiten/v2`)
- **Linux**: bibliotecas de sistema para OpenGL/X11/ALSA (necessárias para rodar e compilar nativamente):

```bash
sudo apt install libgl1-mesa-dev xorg-dev libasound2-dev
```

- **Windows**: nenhuma dependência extra; o executável é autocontido.

Instale as dependências Go com:

```bash
go mod download
```

## Como executar

Modo release (padrão):

```bash
go run ./cmd/game
# ou
make run
```

Modo desenvolvimento:

```bash
VALDORIA_DEV=1 go run ./cmd/game
# ou
make dev
```

## Como compilar

```bash
# sistema atual  -> bin/valdoria
make build

# Linux (amd64)  -> bin/valdoria-linux-amd64
make build-linux

# Windows (amd64)-> bin/valdoria-windows-amd64.exe
make build-windows

# ambos
make build-all
```

Sem o Makefile, os comandos equivalentes são:

```bash
go build ./cmd/game
CGO_ENABLED=1 GOOS=linux   GOARCH=amd64 go build -o bin/valdoria-linux-amd64 ./cmd/game
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui" -o bin/valdoria-windows-amd64.exe ./cmd/game
```

A versão é embutida no build via `-ldflags "-X main.version=..."` (o Makefile usa
`git describe` automaticamente) e aparece no título da janela e no menu inicial.

Os assets de áudio são **opcionais**: o jogo os gera proceduralmente e funciona
sem nenhum arquivo. Para distribuir sons próprios, coloque os arquivos em
`assets/audio/` ao lado do executável (veja `assets/audio/README.md`).

## Verificações

```bash
gofmt -l .
go vet ./...
go test ./...   # ou: make test
```

## Controles

| Ação                | Teclas                         |
| ------------------- | ------------------------------ |
| Mover               | `W` `A` `S` `D` ou setas       |
| Modo de precisão    | Segurar `Shift`                |
| Atirar              | `Espaço` (contínuo ao segurar) |
| Invocação Ancestral (bomba) | `X` ou `Ctrl`          |
| Pausar / Despausar  | `Esc`                          |
| Silenciar áudio     | `M`                            |
| Navegar menus       | Setas / `W` `S`                |
| Confirmar           | `Enter`                        |

As teclas de desenvolvimento (`1`–`4`, `Z`/`F`/`C`/`V`/`B`, `Tab`, `F1`–`F3`,
`K`, `L`) só funcionam com o modo dev ativo — veja a seção "Modo de desenvolvimento".

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
- **Vidas e barra de energia**: barra de HP; ao zerar, perde uma vida e reaparece (limpa os projéteis inimigos, ganha invencibilidade, volta a uma área segura e perde apenas um nível de arma). Início com três vidas; sem vidas, Game Over.
- **Invocação Ancestral (bomba)**: `X`/`Ctrl` limpa os projéteis inimigos, causa dano elevado a todos os inimigos, deixa o jogador invulnerável por instantes e exibe um dragão atravessando a tela. Começa com duas cargas (mostradas no HUD) e não pode ser usada em pausa/Game Over.
- **Pontuação avançada**: pontos por tipo de inimigo, bônus por destruir uma formação completa, bônus por concluir um trecho sem sofrer dano, multiplicador por eliminações consecutivas (decai sem eliminações ou ao sofrer dano) e maior pontuação local da sessão.
- **Chefe da Fase 1 — Vharak, o Dragão Corrompido**: entrada invulnerável exibindo o nome (limpando inimigos/projéteis residuais), barra de vida própria e três fases por porcentagem de vida (100–65%, 65–30%, <30%) com padrões distintos (bolas de fogo direcionadas, cone, arcos, invocações, varredura com brecha), avisos antes de cada padrão, cristais como pontos fracos na fase final, feedback ao receber dano, limpeza de projéteis ao mudar de fase, sequência de destruição ao morrer e transição para a tela de vitória.
- **Fluxo completo de telas**, navegável por teclado, com transições em fade:
  - **Menu inicial** ("Asas de Valdoria"): Iniciar jogo, Controles e Sair.
  - **Controles**: lista todas as teclas do jogo.
  - **Game Over**: pontuação, trecho alcançado, inimigos derrotados e maior multiplicador, com opções de tentar novamente ou voltar ao menu.
  - **Vitória**: conclusão da fase, pontuação final, bônus por vidas e bombas restantes e tempo de conclusão, com opções de jogar novamente ou voltar ao menu.
- Reiniciar cria uma sessão completamente nova, sem inimigos, projéteis, power-ups ou estado anterior; a pausa não funciona no menu, o Game Over não avança a fase e a vitória interrompe os ataques. As confirmações usam borda de tecla (sem disparos múltiplos no mesmo pressionamento).
- **Identidade visual provisória e efeitos** (formas geométricas, sem assets externos):
  - Jogador como silhueta de grifo com o cavaleiro montado, asas em batida, leve inclinação ao mover para os lados, clarão ao disparar e piscada ao receber dano.
  - Inimigos com leitura distinta (corvo, harpia, gárgula, wyvern e o chefe Vharak), com asas animadas.
  - Cenário em parallax de quatro camadas (nuvens altas, colinas distantes, estruturas/árvores e partículas), que muda de tema por trecho: campos, vila em chamas, muralhas e castelo.
  - Partículas ao destruir inimigos, efeito de coleta de power-up, rastros em projéteis especiais (gelo e chamas), sequência da invocação ancestral e explosões na morte do chefe.
  - Vibração de tela discreta em ataques fortes e flash suave ao receber dano, sem prejudicar a visibilidade dos projéteis (desenhados por cima dos efeitos).
  - Opção no menu para ajustar a **Vibração** (Cheia / Reduzida / Off).
- **Áudio básico** com gerenciador próprio: músicas de menu, fase e chefe (com troca suave por fade), e efeitos de disparo, inimigo destruído, dano, coleta, invocação ancestral, vitória e Game Over. Volumes de geral, música e efeitos, tecla `M` para silenciar e proteção contra repetir o mesmo efeito no mesmo frame. Os sons são gerados proceduralmente e podem ser substituídos por arquivos livres (veja `assets/audio/`); o jogo funciona normalmente mesmo sem arquivos.

A lógica roda em TPS fixo do Ebitengine (60), independente da taxa de renderização.
A remoção de projéteis e inimigos reaproveita os slices e libera as referências
antigas para o coletor de lixo.

## Versão

A versão é definida em tempo de build (`-ldflags "-X main.version=..."`). Em
builds pelo Makefile ela vem de `git describe`; sem isso, o padrão é `dev`. A
versão aparece no título da janela ("Asas de Valdoria vX") e no menu inicial.

## Configuração

O jogo **não requer arquivo de configuração**. Os ajustes de execução são feitos
por variáveis de ambiente (`VALDORIA_DEV`, `VALDORIA_SECTION`, `VALDORIA_BOSS`,
`VALDORIA_SEED`) — veja "Modo de desenvolvimento". A opção de intensidade da
vibração fica no menu inicial. Assets de áudio opcionais são lidos de
`assets/audio/` (diretório de trabalho ou ao lado do executável).

## Tratamento de ausência de arquivos

O áudio é gerado proceduralmente; se os arquivos em `assets/audio/` não existirem
ou não puderem ser decodificados, o jogo usa os sons gerados sem falhar. Se o
dispositivo de áudio não puder ser aberto, o jogo segue sem som. Erros na
inicialização da janela são registrados no log e encerram o processo com uma
mensagem clara (prefixo `[valdoria]`).

## Licença

Código distribuído sob a licença [MIT](LICENSE).

## Créditos dos assets

Nenhum asset de terceiros é incluído: toda a arte é geométrica e o áudio é
gerado proceduralmente. Ao adicionar arquivos livres em `assets/audio/`,
registre aqui a origem e a licença de cada um.

## Limitações conhecidas

- Arte totalmente geométrica; ainda sem sprites ou animações detalhadas.
- Áudio procedural simples, sem trilha/efeitos produzidos.
- Uma única fase (Fase 1 — O Cerco de Eldoria) e um único chefe.
- Recorde (high score) mantido apenas na sessão, sem persistência em disco.
- Builds e testes cobrem apenas `amd64` (Linux e Windows).

## Próximas ideias (não implementadas)

- Sprites e animações para o cavaleiro/grifo, inimigos e chefe.
- Trilha e efeitos de áudio produzidos (substituindo os sons gerados).
- Novas fases e mais padrões de ataque.
- Persistência do recorde em arquivo.
- Suporte a gamepad e remapeamento de teclas.

## Checklist do MVP

| Item | Presente |
| --- | --- |
| Menu inicial | Sim |
| Fase completa | Sim (Fase 1, quatro trechos + chefe) |
| Quatro tipos de inimigos | Sim (corvo, harpia, gárgula, wyvern) |
| Três armas | Sim (Luz, Chamas, Gelo — três níveis cada) |
| Power-ups | Sim (runas, cura, escudo) |
| Vidas | Sim (3 vidas + respawn) |
| Pontuação | Sim (formação, sem dano, multiplicador, recorde) |
| Bomba especial | Sim (Invocação Ancestral) |
| Chefe com fases | Sim (Vharak, três fases) |
| Game Over | Sim |
| Vitória | Sim |
| Pausa | Sim |
| Áudio | Sim (música + efeitos, procedural) |
| Modo de desenvolvimento | Sim (via variáveis de ambiente) |
| Testes | Sim (66 testes automatizados) |
| Build executável | Sim (Linux e Windows) |

## Estrutura do projeto

```
valdoria/
├── cmd/
│   └── game/
│       └── main.go     # ponto de entrada, versão, logs e janela
├── go.mod
├── Makefile            # run, dev, test, build, build-linux, build-windows
├── LICENSE             # licença MIT do código
├── README.md
├── TESTING.md          # roteiro de testes manuais
├── bin/                # binários gerados (make build*)
├── assets/
│   └── audio/          # arquivos de áudio livres opcionais (veja o README de lá)
└── game/
    ├── game.go         # loop principal, telas/estados, spawn, colisões e HUD
    ├── config.go       # valores centrais (jogador, combate e inimigos)
    ├── version.go      # variável de versão exibida no título/menu
    ├── dev.go          # modo de desenvolvimento e leitura de variáveis de ambiente
    ├── rng.go          # fontes de aleatoriedade (jogabilidade semeável e efeitos)
    ├── player.go       # cavaleiro/grifo (movimento, precisão, disparo, dano)
    ├── enemy.go        # tipos de inimigos e seus comportamentos
    ├── boss.go         # chefe Vharak (fases, padrões, cristais)
    ├── level.go        # fase 1, linha do tempo de ondas e trechos
    ├── weapon.go       # as três magias, níveis e leque
    ├── powerup.go      # runas, cura e escudo
    ├── enemybullet.go  # projéteis inimigos
    ├── bullet.go       # projéteis do jogador (velocidade, dano, perfuração, rastro)
    ├── background.go   # estrelas e parallax do cenário (nuvens, colinas, estruturas)
    ├── visual.go       # parâmetros visuais centrais, temas por trecho e vibração
    ├── particle.go     # sistema de partículas (explosões, coleta, rastros)
    ├── effects.go      # vibração de tela e flash de dano
    ├── audio.go        # gerenciador de áudio (música, efeitos, volumes, mudo)
    └── *_test.go       # 66 testes (jogador, armas, inimigos, ondas, chefe, fluxo, dev, áudio)
```

## Modo de desenvolvimento

O modo de desenvolvimento **fica desligado no build padrão** e é ativado por
variáveis de ambiente (lidas em `game/dev.go` via `InitFromEnv`):

| Variável            | Efeito                                             |
| ------------------- | -------------------------------------------------- |
| `VALDORIA_DEV=1`    | Ativa o modo de desenvolvimento e o HUD de diagnóstico |
| `VALDORIA_SECTION=n`| Inicia a fase no trecho `n` (0 a 3)                |
| `VALDORIA_BOSS=1`   | Inicia direto no combate contra o chefe            |
| `VALDORIA_SEED=n`   | Fixa a semente da jogabilidade (reprodução de bugs)|

Exemplo:

```bash
VALDORIA_DEV=1 VALDORIA_SEED=42 go run ./cmd/game
```

Com o modo dev ativo, ficam disponíveis:

- **HUD de diagnóstico** (`F1` liga/desliga): FPS, número de inimigos, de
  projéteis e de partículas, tempo da fase, escala de tempo e semente em uso.
- **Hitboxes** (`F2`): contornos de colisão de jogador, inimigos, projéteis e chefe.
- **Invencibilidade** (`F3`).
- **Gerar inimigos**: `1` corvos · `2` harpia · `3` gárgula · `4` wyvern.
- **Gerar power-ups**: `Z` luz · `F` fogo · `C` gelo · `V` cura · `B` escudo.
- **Acelerar a linha do tempo** (`Tab`).
- **Causar dano ao chefe** (`K`).
- **Limpar inimigos e projéteis** (`L`).

Esses recursos não têm efeito quando `VALDORIA_DEV` não está definido.

## Reprodução por semente (seed)

Toda a aleatoriedade que afeta a jogabilidade (drops de power-up e posições de
spawn de teste) usa uma fonte semeável (`game/rng.go`); os efeitos puramente
visuais usam uma fonte separada, para não interferir na reprodução.

Para reproduzir uma sessão, fixe a semente:

```bash
VALDORIA_SEED=42 go run ./cmd/game
```

Com a semente fixa, cada reinício da sessão recomeça exatamente a mesma
sequência de sorteios. Em código/testes, use `game.SetSeed(42)`.

## Testes automatizados

São **66 testes** (`go test ./...`), cobrindo as regras independentes de renderização:

- **Jogador**: clamp na tela, escudo, respawn e redução de arma, invencibilidade dev.
- **Armas**: troca, progressão de nível, teto, ângulos do leque, perfuração do gelo.
- **Projéteis**: remoção fora da tela, ausência de dano repetido no mesmo alvo.
- **Inimigos**: dano, morte, pontuação, mira, movimento determinístico e fuga.
- **Colisões**: acerto único por projétil e um único dano ao jogador por frame.
- **Ondas**: ativação por tick, conclusão, troca de trecho, pausa na linha do tempo.
- **Pontuação**: multiplicador, bônus de formação e de trecho sem dano.
- **Power-ups**: coleta, cura, escudo, troca/upgrade de arma.
- **Chefe**: fases por vida, seleção de padrão, invulnerabilidade e vitória.
- **Mudança de estados e reinício**: menu inicial, pausa fora do jogo, Game Over sem
  avanço, sessão totalmente limpa ao reiniciar.
- **Áudio**: mudo, proteção de repetição por frame, troca de música, limite de volume.
- **Semente**: reprodutibilidade dos sorteios e determinismo do reinício.

Rodar com cobertura: `go test -cover ./game/`.

### Partes que dependem de verificação manual

Renderização e feedback visual/sonoro não são cobertos por testes automatizados:
desenho do grifo/inimigos/chefe, parallax do cenário, partículas, vibração de tela,
flash de dano, transições de fade e a reprodução efetiva do áudio no dispositivo.

# Asas de Valdoria

Shoot 'em up vertical em Golang, inspirado na jogabilidade de *Super Aleste*,
com identidade própria e temática medieval fantástica. O jogador controla um
cavaleiro montado em um grifo que enfrenta criaturas aladas corrompidas
descendo pela tela, até o confronto com o chefe da fase, Vharak.

A arte usa **pixel art no estilo arcade 16-bit** gerada por código (sem assets
externos obrigatórios), com cenário em parallax por formas geométricas; tanto os
sprites quanto o áudio podem ser substituídos por arquivos livres. O áudio também
é gerado proceduralmente.

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
- Seis tipos de inimigos com vida, pontuação, dano e padrões próprios:
  - **Corvo corrompido**: pequeno, frágil, desce reto e aparece em formações.
  - **Harpia**: vida média, zigue-zague e disparo reto para baixo.
  - **Gárgula**: resistente, entra pela lateral, ataca parada e sai pelo lado oposto.
  - **Wyvern**: grande e resistente, alinha-se ao jogador e dispara projéteis mirados (sem perseguir depois de lançados).
  - **Balista corrompida**: ameaça "terrestre" que desce lenta como o cenário e dispara rajadas de virotes mirados (bem telegrafadas). Resistente.
  - **Feiticeiro corrompido**: para no alto, desliza de lado e solta anéis completos de projéteis; frágil, mas perigoso se ignorado.
- Inimigos que atiram exibem um **telegraph** (contorno de alerta piscante) antes de disparar.
- Projéteis inimigos com cor própria (magenta/roxo com núcleo claro e contorno), deliberadamente distinta de todas as armas do jogador para leitura imediata de aliado × inimigo; indicação visual de dano (piscada branca) ao acertar um inimigo.
- **Fuga pela base**: se um inimigo atravessar a parte de baixo da tela sem ser destruído, o jogador **perde uma vida** (saídas laterais, como a gárgula, não contam). Durante a invencibilidade pós-morte a fuga ainda remove o inimigo, mas não empilha perdas.
- Colisão projétil × inimigo e jogador × inimigo/projétil, sem dano repetido de inimigos já destruídos.
- **Campanha de três fases encadeadas**, descritas de forma **declarativa (dados, não código)** em `stages.go` (`stageDef`/`waveDef`), o que torna trivial adicionar/editar fases:
  - **Fase 1 — O Cerco de Eldoria**: campos, vila em chamas, muralhas e castelo (corvos → harpias → gárgulas → wyverns).
  - **Fase 2 — A Floresta Corrompida**: bosque sombrio e pântano; introduz o **Feiticeiro**.
  - **Fase 3 — O Covil de Vharak**: desfiladeiro e covil; introduz a **Balista** e culmina no chefe.
  - Cada fase tem biomas/temas próprios, avisos de trecho, barra de progresso com marcos e encadeia para a seguinte; só a última invoca o chefe.
  - **Ritmo denso e contínuo**: ondas curtas e sobrepostas mantêm a tela sempre com ação (sem esperas mortas).
- **Modos de jogo**: **Campanha** (fases + chefe) e **Sobrevivência** (ondas infinitas com dificuldade crescente e recorde próprio).
- **Dificuldades** Fácil / Normal / Difícil (vidas e bombas iniciais, vida dos inimigos, chance de drop e **variação de spawn**: posição X e ritmo entre aparições da mesma onda — Fácil fica previsível; Difícil espalha bem), selecionáveis no menu e persistidas.
- **Três magias, cada uma com três níveis**:
  - **Lança de Luz** (inicial): disparos retos densos (2 → 3 → 4), dano alto; **atordoa** o alvo e lança um **arco dourado** no inimigo próximo (metade do dano + stun curto). Nv3 ainda perfura 1.
  - **Chamas do Dragão**: leque curto (3 → 4 → 5); aplica **queimadura** (dano ao longo do tempo) e deixa o inimigo **mais frágil** (+50% dano recebido).
  - **Lanças de Gelo**: tiros pesados com perfuração; **congela o ritmo** do inimigo (movimento e tiros ~40% mais lentos).
- **Power-ups** que caem de inimigos e saem da tela se não coletados: Runa da Luz/Fogo/Gelo (coletar a runa da arma atual sobe o nível até o máximo 3; a runa de outra arma troca o tipo **preservando o nível atual** — trocar nunca rebaixa o poder), Cura (limitada à vida máxima) e Escudo temporário (absorve um ataque). Chance de drop configurável e alguns drops garantidos por onda.
- HUD com arma, nível e estado do escudo.
- **Vidas e barra de energia**: barra de HP; ao zerar, perde uma vida e reaparece (limpa os projéteis inimigos, ganha invencibilidade, volta a uma área segura e perde apenas um nível de arma). Início com três vidas; sem vidas, Game Over.
- **Invocação Ancestral (bomba)**: `X`/`Ctrl` limpa os projéteis inimigos, causa dano elevado a todos os inimigos, deixa o jogador invulnerável por instantes e exibe um dragão atravessando a tela. Começa com duas cargas (mostradas no HUD) e não pode ser usada em pausa/Game Over.
- **Pontuação avançada**: pontos por tipo de inimigo, bônus por destruir uma formação completa, bônus por concluir um trecho sem sofrer dano, multiplicador por eliminações consecutivas (decai sem eliminações ou ao sofrer dano) e maior pontuação local da sessão.
- **Chefe da Fase 1 — Vharak, o Dragão Corrompido**: entrada invulnerável exibindo o nome (limpando inimigos/projéteis residuais), barra de vida própria e três fases por porcentagem de vida (100–65%, 65–30%, <30%) com padrões distintos (bolas de fogo direcionadas, cone, arcos, invocações, varredura com brecha), avisos antes de cada padrão, cristais como pontos fracos na fase final, feedback ao receber dano, limpeza de projéteis ao mudar de fase, sequência de destruição ao morrer e transição para a tela de vitória.
- **Fluxo completo de telas**, navegável por teclado, com transições em fade:
  - **Menu inicial** ("Asas de Valdoria"): Iniciar jogo, Sobrevivência, Dificuldade, Controles, Vibração e Sair; mostra os recordes (campanha e sobrevivência).
  - **Controles**: lista todas as teclas do jogo.
  - **Game Over**: pontuação, trecho alcançado, inimigos derrotados e maior multiplicador, com opções de tentar novamente ou voltar ao menu.
  - **Vitória**: conclusão da fase, pontuação final, bônus por vidas e bombas restantes e tempo de conclusão, com opções de jogar novamente ou voltar ao menu.
- Reiniciar cria uma sessão completamente nova, sem inimigos, projéteis, power-ups ou estado anterior; a pausa não funciona no menu, o Game Over não avança a fase e a vitória interrompe os ataques. As confirmações usam borda de tecla (sem disparos múltiplos no mesmo pressionamento).
- **Identidade visual e efeitos** (pixel art procedural + formas geométricas no cenário; sem assets externos obrigatórios):
  - **Sprites em pixel art (arcade 16-bit)** para o jogador, os seis inimigos, o chefe Vharak e as **runas/power-ups** (luz, fogo, gelo, vida, escudo), com sistema de fallback igual ao do áudio: usa `assets/sprites/<nome>.png` se existir, senão gera a arte por código. Sprites escalados à hitbox, com batida de asa nos voadores e leve flutuação nas runas. Prévia: `VALDORIA_EXPORT_SPRITES=/tmp/s.png go test ./game/ -run TestExportSpritePreviews`.
  - **Glow em projéteis**: halo emissivo nas armas do jogador (luz/chamas/gelo) e nos tiros inimigos (magenta), para nunca sumirem no fundo.
  - **Molduras de UI** (ferro/pergaminho) no HUD, menus, pausa, barra de progresso e vida do chefe.
  - **Fonte própria de interface** (bitmap embutida via `text/v2` + `basicfont`, sem asset externo): texto colorido com identidade de pergaminho, substituindo a fonte de depuração no HUD, menus e telas.
  - **Contraste fundo × jogo**: o cenário é escurecido e os elementos de jogo (jogador, inimigos, chefe) recebem contorno escuro, resolvendo a confusão entre fundo e monstros e reforçando as silhuetas.
  - Jogador como silhueta de grifo com o cavaleiro montado, asas em batida, leve inclinação ao mover para os lados, clarão ao disparar e piscada ao receber dano.
  - Inimigos com leitura distinta (corvo, harpia, gárgula, wyvern e o chefe Vharak), com asas animadas.
  - Cenário em parallax de quatro camadas (nuvens altas, colinas distantes, estruturas/árvores e partículas), que muda de tema por trecho: campos, vila em chamas, muralhas e castelo.
  - Partículas ao destruir inimigos, efeito de coleta de power-up, rastros em projéteis especiais (gelo e chamas), sequência da invocação ancestral e explosões na morte do chefe.
  - **Números de pontuação flutuantes** sobre o inimigo abatido (já com o multiplicador aplicado), reforçando o combo.
  - **Impact freeze (hit-stop)** curto ao abater inimigos grandes (wyvern) e nas trocas de fase e na morte do chefe, dando peso aos golpes fortes sem parecer travamento (os efeitos seguem animando).
  - **Feedback de escudo**: som e anel de partículas quando o escudo absorve um golpe (sem contar como dano).
  - **Telegraph de disparo**: inimigos que atiram (harpia, gárgula, wyvern) piscam um contorno de alerta antes de disparar, evitando tiros surpresa.
  - **Brasas no bocal** ao disparar as Chamas do Dragão, reforçando o peso da arma.
  - **Escurecimento de impacto** ao usar a Invocação Ancestral (o céu escurece por um instante).
  - **Popups de bônus** ("+60" por formação completa, "+200" por trecho sem dano) destacando as recompensas ocultas.
  - **HUD com contraste**: painel escuro atrás dos textos garante leitura sobre cenários claros (como a vila em chamas).
  - Vibração de tela discreta em ataques fortes e flash suave ao receber dano, sem prejudicar a visibilidade dos projéteis (desenhados por cima dos efeitos).
  - **Barra de progresso** com marcos dos quatro trechos, e **barra de vida do chefe** com moldura, cor por fase (vermelho → laranja → vermelho vivo) e marcas nos limiares de troca de fase (65% e 30%); aviso pulsante de "ALERTA" na entrada do chefe.
  - Opção no menu para ajustar a **Vibração** (Cheia / Reduzida / Off).
- **Persistência local** (`save.go`, JSON no diretório de configuração do usuário): recordes de campanha e sobrevivência e preferências (dificuldade, vibração, mudo) sobrevivem entre execuções; tolerante a falhas (o jogo nunca quebra se o save não existir).
- **Áudio básico** com gerenciador próprio: músicas de menu, fase e chefe (com troca suave por fade), e efeitos de disparo **por arma** (luz/chamas/gelo), inimigo destruído, dano, coleta, invocação ancestral, vitória, Game Over, quebra de escudo e navegação/confirmação de menu. Volumes de geral, música e efeitos, tecla `M` para silenciar e proteção contra repetir o mesmo efeito no mesmo frame. Os sons são gerados proceduralmente e podem ser substituídos por arquivos livres (veja `assets/audio/`); o jogo funciona normalmente mesmo sem arquivos.

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
vibração fica no menu inicial. Assets opcionais de áudio e sprites são lidos de
`assets/audio/` e `assets/sprites/` (diretório de trabalho ou ao lado do
executável).

## Tratamento de ausência de arquivos

O áudio e os sprites de entidades são gerados proceduralmente; se os arquivos em
`assets/audio/` ou `assets/sprites/` não existirem ou não puderem ser
decodificados, o jogo usa as versões geradas sem falhar. Se não houver
dispositivo de som disponível (por exemplo, WSL2 sem áudio), o áudio é
desativado automaticamente e o jogo segue sem som — dá para forçar isso com
`VALDORIA_NOAUDIO=1`. Erros na inicialização da janela são registrados no log e
encerram o processo com uma mensagem clara (prefixo `[valdoria]`).

## Licença

Código distribuído sob a licença [MIT](LICENSE).

## Créditos dos assets

Nenhum asset de terceiros é incluído: a arte de entidades é pixel art gerada
por código (substituível por PNGs em `assets/sprites/`), o cenário ainda usa
formas geométricas e o áudio é gerado proceduralmente. Ao adicionar arquivos
livres em `assets/audio/` ou `assets/sprites/`, registre aqui a origem e a
licença de cada um.

## Limitações conhecidas

- Cenário ainda é geométrico (parallax); UI e projéteis usam molduras/glow vetoriais, não sprites PNG.
- Áudio procedural simples, sem trilha/efeitos produzidos.
- Campanha de três fases com um **único chefe** ao final (Vharak); ainda sem chefe intermediário.
- Builds e testes cobrem apenas `amd64` (Linux e Windows).
- A nova fonte de UI e os biomas das fases 2-3 ainda dependem de verificação visual manual (sem teste de renderização).

## Próximas ideias (não implementadas)

- Mais frames de animação / dithering nos fundos / arte PNG desenhada.
- Trilha e efeitos de áudio produzidos (substituindo os sons gerados).
- Segundo chefe / generalização do sistema de chefe.
- Suporte a gamepad e remapeamento de teclas.

## Checklist do MVP

| Item | Presente |
| --- | --- |
| Menu inicial | Sim (com dificuldade, modos e recordes) |
| Campanha completa | Sim (3 fases em dados + chefe) |
| Modo Sobrevivência | Sim (ondas infinitas, recorde próprio) |
| Dificuldades | Sim (Fácil/Normal/Difícil) |
| Persistência | Sim (recordes e opções em disco) |
| Seis tipos de inimigos | Sim (corvo, harpia, gárgula, wyvern, balista, feiticeiro) |
| Três armas | Sim (Luz, Chamas, Gelo — três níveis cada, som próprio) |
| Power-ups | Sim (runas, cura, escudo) |
| Vidas | Sim (respawn; quantidade varia por dificuldade) |
| Pontuação | Sim (formação, sem dano, multiplicador, recorde) |
| Bomba especial | Sim (Invocação Ancestral) |
| Chefe com fases | Sim (Vharak, três fases) |
| Game Over | Sim |
| Vitória | Sim |
| Pausa | Sim |
| Áudio | Sim (música + efeitos, procedural) |
| Modo de desenvolvimento | Sim (via variáveis de ambiente) |
| Testes | Sim (105 testes automatizados) |
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
│   ├── audio/          # arquivos de áudio livres opcionais (veja o README de lá)
│   └── sprites/        # PNGs opcionais que substituem a pixel art (veja o README de lá)
└── game/
    ├── game.go         # loop principal, telas/estados, spawn, colisões e HUD
    ├── config.go       # valores centrais (jogador, combate e inimigos)
    ├── version.go      # variável de versão exibida no título/menu
    ├── dev.go          # modo de desenvolvimento e leitura de variáveis de ambiente
    ├── rng.go          # fontes de aleatoriedade (jogabilidade semeável e efeitos)
    ├── player.go       # cavaleiro/grifo (movimento, precisão, disparo, dano)
    ├── enemy.go        # tipos de inimigos e seus comportamentos (6 inimigos)
    ├── boss.go         # chefe Vharak (fases, padrões, cristais)
    ├── level.go        # execução da fase (linha do tempo de ondas e trechos)
    ├── stages.go       # campanha em dados: fases 1-3 (stageDef/waveDef)
    ├── difficulty.go   # presets Fácil/Normal/Difícil (multiplicadores)
    ├── save.go         # persistência de recordes e opções (JSON)
    ├── weapon.go       # as três magias, níveis, leque e som por arma
    ├── powerup.go      # runas, cura e escudo
    ├── enemybullet.go  # projéteis inimigos
    ├── bullet.go       # projéteis do jogador (velocidade, dano, perfuração, rastro)
    ├── background.go   # estrelas e parallax do cenário (nuvens, colinas, estruturas)
    ├── visual.go       # parâmetros visuais centrais, temas/biomas e vibração
    ├── particle.go     # sistema de partículas (explosões, coleta, rastros)
    ├── popup.go        # números de pontuação/bônus flutuantes
    ├── sprite.go       # sistema de sprites (PNG ou pixel art procedural, cache)
    ├── spritedata.go   # pixel art (grades ASCII + paletas) do jogador, inimigos e chefe
    ├── glow.go         # halo emissivo dos projéteis
    ├── uiframe.go      # molduras de ferro/pergaminho da UI
    ├── uitext.go       # fonte própria da interface (text/v2 + basicfont)
    ├── effects.go      # vibração de tela e flash de dano
    ├── audio.go        # gerenciador de áudio (música, efeitos, volumes, mudo)
    └── *_test.go       # 105 testes (jogador, armas, inimigos, ondas, chefe, fluxo, dev, áudio, integração, save, sobrevivência, sprites)
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

São **105 testes** (`go test ./...`), cobrindo as regras independentes de renderização:

- **Jogador**: clamp na tela, escudo, respawn e redução de arma, invencibilidade dev, normalização da diagonal.
- **Armas**: troca, progressão de nível, teto, ângulos do leque, perfuração do gelo.
- **Projéteis**: remoção fora da tela, ausência de dano repetido no mesmo alvo, cor distinta do inimigo.
- **Inimigos**: dano, morte, pontuação, mira, movimento determinístico e fuga.
- **Colisões**: acerto único por projétil e um único dano ao jogador por frame; escudo × dano real.
- **Ondas**: ativação por tick, conclusão, troca de trecho, pausa na linha do tempo.
- **Pontuação**: multiplicador, bônus de formação e de trecho sem dano, popups ao abater.
- **Power-ups**: coleta, cura, escudo, troca/upgrade de arma.
- **Chefe**: fases por vida, seleção de padrão, invulnerabilidade, vitória e limpeza segura da arena na morte.
- **Feedback/efeitos**: hit-stop congela a lógica, teto de partículas respeitado.
- **Mudança de estados e reinício**: menu inicial, pausa fora do jogo, Game Over sem
  avanço, sessão totalmente limpa ao reiniciar.
- **Integração (fumaça)**: sessão completa dirigida pela lógica até a vitória, sem panic; chefe alcançado após as ondas.
- **Áudio**: mudo, proteção de repetição por frame, troca de música, limite de volume.
- **Semente**: reprodutibilidade dos sorteios e determinismo do reinício.

Rodar com cobertura: `go test -cover ./game/`.

### Partes que dependem de verificação manual

Renderização e feedback visual/sonoro não são cobertos por testes automatizados:
desenho do grifo/inimigos/chefe, parallax do cenário, partículas, vibração de tela,
flash de dano, transições de fade e a reprodução efetiva do áudio no dispositivo.

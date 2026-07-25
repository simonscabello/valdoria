# GAME_DIRECTION.md — Asas de Valdoria

> **Documento de visão oficial do projeto.**
> Escrito na posição de Diretor do Jogo, após leitura integral do código-fonte
> (`game/`, 8.751 linhas, 120 testes passando), do `ROADMAP.md` anterior, do
> `README.md` e do `TESTING.md`.
>
> Este documento **substitui o `ROADMAP.md`** como direção de produto. O
> `ROADMAP.md` foi um bom plano de *engenharia* e cumpriu seu papel: os 10 itens
> dele foram entregues. Ele não é mais suficiente porque o problema do projeto
> deixou de ser técnico.
>
> **Nada foi implementado, alterado ou commitado.** Isto é análise, estratégia e
> planejamento.
>
> Data: 24 de julho de 2026 · Base: `main` @ `3989aa0`

---

## ✅ Estado de execução — versão 0.2 "O jogo justo" concluída

O bloco 0.2 do roadmap (seção 7) e os passos 2 e 5–12 do plano de ação
(seção 15) foram implementados. A análise das seções 0 e 1 descreve o MVP
**anterior** a essas mudanças e é mantida como registro do diagnóstico.

| Item | Antes | Depois |
| --- | --- | --- |
| Stun-lock da Lança de Luz | Congelava todo alvo sob fogo contínuo | Removido; a Luz não aplica status |
| Chefe (gatilho segurado) | **5,0 s**, 0 padrões executados | **65 s**, 43 padrões, chega à fase 3 |
| Chefe (aproveitamento realista) | 16,6 s | **61–163 s** conforme a magia |
| Desequilíbrio entre armas | **5,8×** (foco) | **1,5×** entre os picos |
| Identidade das armas | Luz vencia tudo | Luz=foco · Chamas=formação · Gelo=coluna |
| Corvo na campanha | **63,5%** | **35,2%** |
| Gárgula / Wyvern / Mago / Balista | 1,1% / 2,2% / 1,5% / 0,7% | 11,3% / 10,9% / 11,3% / 7,8% |
| Teto de poder atingido em | **3% da run** (7 s) | **33% da run** (70 s) |
| Runas inertes por partida | 53 de 92 | **2 de 29** |
| Respawn | Restaurava 5 HP (mais que os 3 iniciais) | Restaura os 3 iniciais |
| Bônus por bomba não usada | 300 pts/carga (pagava para não usar) | Removido |
| Peso do chefe no placar | 5.000 pts (50× o inimigo mais caro) | 1.500 pts |

**Novidades de sistema:** medição de balanceamento *headless*
(`go run ./cmd/balance`, também executada por `TestBestiaryIsBalanced`);
carga de runas com progresso visível no HUD (3 runas por nível, compartilhadas
entre elementos); drops que favorecem cura/escudo no teto de poder; leque das
Chamas com contagem ímpar (corrige o Nv2 que rendia **menos** que o Nv1).

**Estado técnico:** 124 testes, `gofmt`/`go vet`/`-race` limpos, 13 de 13
critérios de balanceamento atendidos.

---

## ✅ Estado de execução — versão 0.3 "A Corrupção" concluída

O diferencial da seção 3.2 está implementado e medido. O jogo agora tem uma
mecânica que nenhum dos concorrentes da seção 13 possui.

| Entrega | Estado |
| --- | --- |
| Medidor de Corrupção com as 5 faixas | ✅ `corruption.go` |
| Pesos por inimigo, só em fugas pela base | ✅ gárgula isenta (sai pela lateral por design) |
| Multiplicadores de pontos, drops e vida | ✅ até ×3,0 em pontos |
| Variantes corrompidas | ✅ corvo-sombra que atira, harpia acelerada, wyvern em leque |
| Degradação visual do mundo | ✅ dessaturação + violeta + veias nas bordas |
| Leitura no HUD | ✅ coluna na borda direita, percentual, anúncio de faixa |
| Aviso de fuga iminente | ✅ contorno, seta e filete pulsante na base |
| **Vharak Ascendido** | ✅ +60% vida, +25% velocidade, padrão extra por fase, cristais na fase 2 |
| Medição da curva na ferramenta | ✅ seção 5 do relatório + 3 critérios novos |

**Curva medida** (`make balance`), que é o que prova que a aposta funciona:

| Fugas do jogador | Corrupção final | Faixa | Pontos | Final alternativo |
| ---: | ---: | --- | ---: | --- |
| 2% | 3% | Reino Firme | ×1,0 | não |
| 5% | 22% | Reino Firme | ×1,0 | não |
| 10% | 29% | Sombra Crescente | ×1,3 | não |
| 15% | 48% | Sombra Crescente | ×1,3 | não |
| 25% | 87% | Colapso | ×2,5 | não |
| 40% | 100% | Queda de Valdoria | ×3,0 | **sim** |

O formato importa: imprecisão normal não pune, negligência real derruba o reino,
e o final alternativo exige **decisão** — 1 de 6 perfis chega lá.

**Estado técnico:** 134 testes na entrega da 0.3 (139 após a correção de
dificuldade abaixo), `gofmt`/`go vet`/`-race` limpos.

**Pendente da 0.3:** a moldura lateral decorada (o medidor ainda divide espaço
com o campo de jogo) e o áudio reativo por camadas — ambos dependem de decisões
de apresentação que cabem melhor junto da 0.5.

## ✅ Correção de dificuldade — pós-playtest

Relato do playtest: *"o jogo está muito fácil, mesmo na dificuldade normal"*.

A medição confirmou e localizou a causa, que **não** era a esperada (velocidade
de projétil ou invencibilidade), e sim a vida dos inimigos:

| Modelo simulado | Antes | Depois |
| --- | --- | --- |
| passivo (não se move) | morria aos 92s · 2.819 projéteis inimigos | morre aos 92s · 126 projéteis |
| **mediano** (mira, nunca desvia) | **terminava a campanha, 2 golpes** | **morre aos 88s, 3 vidas perdidas** |

A causa raiz: **uma salva da Lança de Luz Nv3 (12 de dano) matava qualquer
inimigo do jogo** — o wyvern tinha 8 de vida. Eles morriam antes de chegar a
disparar, e a campanha inteira gerava 57 projéteis inimigos em 210 segundos. Um
alvo que não revida não é um inimigo.

**Correções:** vida dos inimigos recalibrada (harpia 3→15, gárgula 6→30, wyvern
8→46, balista 6→32, feiticeiro 5→22; total da campanha 933→4.336); atraso do
**primeiro** disparo bem menor que o intervalo seguinte, para a ameaça existir na
chegada; projétil inimigo 2,2→3,1 (mais rápido que o jogador, que anda a 2,5);
invencibilidade 90→60 frames.

**Dificuldades agora separam de verdade.** Os presets mexiam só em vida dos
inimigos, drops, vidas e jitter — a ameaça real era idêntica nos três. Ganharam
três eixos novos (velocidade de projétil, cadência de tiro e invencibilidade):

| Preset | Sobrevivência do modelo passivo |
| --- | ---: |
| Fácil | 139s |
| Normal | 92s |
| Difícil | 67s |

### Segunda rodada: esponja não é ameaça

Relato seguinte: *"ficou bem mais difícil, gostei, mas pareceu que foi demais.
Eles demoram."* — o modo de falha oposto, e que parece idêntico numa tabela de
vida dos inimigos.

Criei a seção **AMEACA** do relatório, que separa as duas coisas medindo
tempo-para-matar e **tiros disparados antes de morrer**. Ela derrubou minha
hipótese: o wyvern, com 46 de vida, morria em 1,1s; a **gárgula (2,4s) e a
harpia (1,8s)** eram as esponjas, com um terço da vida dele.

A causa não era vida — era **largura e evasão**. A salva da Luz cobre 18px:
o wyvern de 24px leva as quatro lanças, a harpia de 14px leva metade e ainda
foge do próprio projétil com o zigue-zague. Baixar a vida da harpia de 15 para
10 não mudou nada no tempo-para-matar.

O que resolveu:

| Ajuste | Efeito |
| --- | --- |
| Lança de Luz 5,5→6,6 (Nv3 7,0→8,6) | menos tempo de voo, menos necessidade de liderar o alvo |
| Harpia: amplitude 42→32, frequência 0,05→0,038 | segue evasiva sem escapar dos próprios tiros |
| Gárgula: velocidade de entrada 1,5→2,3 | para de gastar a vida atravessando a tela |
| Cadência de tiro de todos os armados | a ameaça volta pela **frequência**, não pela absorção |

| Inimigo | TTK antes | TTK depois |
| --- | ---: | ---: |
| Harpia | 1,77s | **0,93s** |
| Gárgula | 2,38s | **1,62s** |
| Wyvern | 1,27s | **1,00s** |
| Feiticeiro | 0,87s | **0,63s** |

E a pressão **subiu** em vez de cair: os inimigos passaram a disparar 96
projéteis contra o modelo mediano (era 82), morrendo mais rápido. É a troca
correta — dificuldade por cadência, não por barra de vida.

**Lição de método registrada:** "aumentar a dificuldade" degenera naturalmente em
inflar vida, que é a forma mais preguiçosa e menos divertida de dificultar um
jogo. Sem uma métrica que separe *ameaça* de *demora*, essa degeneração não é
visível — nem para quem escreveu o código, nem para quem joga.

**Novos instrumentos:** as seções **AMEACA** e **PRESSAO** do relatório
(`make balance`). A primeira mede tempo-para-matar e tiros por inimigo; a
segunda simula a campanha com modelos de jogador que nunca desviam. Guardadas
por `TestGameDemandsDodging`, `TestEnemiesThreatenWithoutDragging` e
`TestHarderDifficultyKillsFaster`.

**Estado técnico:** 139 testes, `gofmt`/`vet`/`-race` limpos, 18 de 18 critérios.

---

**Próximo bloco:** 0.4 — o Mergulho do Grifo, o fôlego, o graze e a camada de
input com gamepad (seção 7).

> ⚠️ **O passo 14 do plano de ação continua pendente e é o mais importante do
> cronograma:** validar a Corrupção com 5 jogadores externos. O sistema está
> medido e coerente, mas medição não é diversão. Só o playtest decide se esta
> aposta sustenta o posicionamento do jogo.

---

## Sumário

| # | Seção |
| --- | --- |
| 0 | [Método e evidências](#0-método-e-evidências) |
| 1 | [Avaliação honesta do MVP](#1-avaliação-honesta-do-mvp) |
| 2 | [A visão do jogo](#2-a-visão-do-jogo) |
| 3 | [O que torna este jogo único](#3-o-que-torna-este-jogo-único) |
| 4 | [Evolução da jogabilidade](#4-evolução-da-jogabilidade) |
| 5 | [Progressão entre partidas](#5-progressão-entre-partidas) |
| 6 | [Conteúdo da primeira versão comercial](#6-conteúdo-da-primeira-versão-comercial) |
| 7 | [Roadmap 0.2 → 1.0](#7-roadmap-02--10) |
| 8 | [Game feel: 60 melhorias](#8-game-feel-60-melhorias) |
| 9 | [Direção artística](#9-direção-artística) |
| 10 | [Direção de áudio](#10-direção-de-áudio) |
| 11 | [Monetização](#11-monetização) |
| 12 | [Marketing](#12-marketing) |
| 13 | [Benchmark](#13-benchmark) |
| 14 | [O que NÃO devemos fazer](#14-o-que-não-devemos-fazer) |
| 15 | [Plano de ação: os 30 próximos passos](#15-plano-de-ação-os-30-próximos-passos) |

---

## 0. Método e evidências

Não escrevi nada aqui por impressão. Antes de opinar, medi o jogo pelo código.
Estes números aparecem repetidamente ao longo do documento e sustentam quase
todas as recomendações:

### Composição real da campanha

Extraído de [`stages.go`](game/stages.go), aplicando `campaignDensityScale = 1.45`
e expandindo formações (`formationLine` = 4, `formationV` = 5):

| Fase | Ondas | Inimigos gerados | Último spawn | Duração |
| --- | ---: | ---: | ---: | ---: |
| 1 — O Cerco de Eldoria | 40 | 233 | tick 5.578 | ~93 s |
| 2 — A Floresta Corrompida | 23 | 124 | tick 3.001 | ~50 s |
| 3 — O Covil de Vharak | 24 | 103 | tick 3.289 | ~55 s |
| **Total** | **87** | **460** | — | **~3 min 20 s** + chefe |

### Distribuição por tipo de inimigo — o dado mais importante do projeto

| Inimigo | Aparições na campanha | % |
| --- | ---: | ---: |
| Corvo | 292 | **63,5 %** |
| Harpia | 143 | **31,1 %** |
| Wyvern | 10 | 2,2 % |
| Mago | 7 | 1,5 % |
| Gárgula | 5 | 1,1 % |
| Balista | 3 | 0,6 % |

**94,6 % da campanha inteira é composta por dois inimigos.** Existem seis
inimigos implementados, testados e desenhados; quatro deles somam 25 aparições
em cinco minutos de jogo. O jogo tem um bestiário que ele não usa.

### Balanceamento das armas (nível 3, todos os projéteis acertando)

| Arma | Projéteis × dano / cooldown | DPS | Efeito |
| --- | --- | ---: | --- |
| **Lança de Luz** | 4 × 3 / 8f | **90** | Stun 32f + arco em cadeia + perfura 1 |
| Lanças de Gelo | 3 × 5 / 14f | 64 | Lentidão 60% + perfura 3 |
| Chamas do Dragão | 5 × 1 / 11f | 27 | Queimadura (~3,8 DPS) + 50% vulnerabilidade |

A Luz tem **3,3× o DPS das Chamas**, mais controle, mais alcance efetivo e uma
área de dano gratuita. Não é uma escolha entre estilos: é a resposta certa.

### O stun-lock

[`status.go:19`](game/status.go:19) aplica `stun = stunDuration` (32 frames) a
cada acerto de Luz. O cooldown da Luz é de 8 frames
([`config.go:166`](game/config.go:166)). [`enemy.go:76`](game/enemy.go:76)
retorna cedo enquanto `stun > 0`.

**Consequência: qualquer inimigo sob fogo contínuo da Lança de Luz fica
permanentemente congelado — não anda, não atira, não existe.** No chefe o stun é
metade (16f), ainda o dobro do cooldown: **Vharak também é stun-locked**. Com 420
de vida e 90 DPS, o clímax da campanha dura **4,7 segundos de gatilho segurado**,
sem que o chefe execute um único padrão.

Isto não é um detalhe de balanceamento. É o boss fight — a melhor peça de código
do jogo, com três fases, cinco padrões, cristais e telegrafia — sendo apagado por
uma interação de duas linhas.

### A regra que existe no README mas não no código

O `README.md` documenta, com destaque: *"**Fuga pela base**: se um inimigo
atravessar a parte de baixo da tela sem ser destruído, o jogador **perde uma
vida**"*.

Isso **não existe no código**. [`game.go:871-877`](game/game.go:871) apenas marca
`escaped = true` e remove a entidade. `removeDead` usa a flag só para não contar o
inimigo como abatido e para invalidar o bônus de formação. **Não há penalidade
alguma por deixar inimigos escaparem.**

Guarde esta observação. Ela é, ao mesmo tempo, o maior buraco de design do MVP e
a maior oportunidade criativa do projeto — volto a ela na seção 3.

### Curva de poder dentro de uma partida

- Chance de drop: 20% por abate ([`config.go:198`](game/config.go:198)), mais 14
  drops garantidos por onda.
- 460 abates × 0,20 ≈ **92 power-ups aleatórios + 14 garantidos**.
- Nível máximo de arma: **3** ([`config.go:161`](game/config.go:161)).

O jogador atinge o teto de poder em torno do **primeiro minuto** de uma partida de
**cinco**. Nos quatro minutos seguintes, ~85 power-ups caem na tela sem significar
absolutamente nada, exceto cura e escudo.

### Outras medições relevantes

- `respawn()` restaura a vida para `maxHealth` (5), enquanto o jogo **começa** com
  `initialHealth` (3) — [`player.go:185`](game/player.go:185). Morrer devolve mais
  vida do que nascer.
- Invencibilidade pós-dano: 90 frames (1,5 s). Pós-respawn: 150 frames (2,5 s).
- Orçamento total de erro: 3 + 5 + 5 = **13 acertos** antes do Game Over, mais
  curas. Para um shmup, é muito generoso.
- Vitória concede **300 pontos por bomba não usada** ([`config.go:155`](game/config.go:155)).
  O sistema de pontuação **paga o jogador para não usar a mecânica mais
  espetacular do jogo**.
- Não existe power-up que reponha bombas. Duas cargas, e acabou.
- Chefe vale 5.000 pontos; o inimigo mais caro vale 100. **Um único evento vale
  mais que centenas de decisões de combate.**
- Toda música é um laço de 8 notas: `genMusicFields` = 8 × 0,28 s = **2,24
  segundos em loop infinito** ([`audio.go:575`](game/audio.go:575)).
- Entrada de teclado é lida diretamente em `Player.update()`
  ([`player.go:45`](game/player.go:45)). **Não há camada de input: gamepad e
  remapeamento estão bloqueados por arquitetura.**

---

## 1. Avaliação honesta do MVP

Vou responder as perguntas diretamente, sem diplomacia. É o que serve ao projeto.

### O MVP é divertido?

**Por cerca de 90 segundos, sim. Depois, não.**

Os primeiros noventa segundos funcionam de verdade: os corvos descem, o controle
responde, o tiro sai gostoso, os números de pontuação sobem, o multiplicador
cresce, você pega uma runa e sente que ficou mais forte. Isso é um shmup
funcionando.

O que acontece depois é o problema. Aos ~60 segundos você está no nível 3 de
arma. Aos ~90 segundos você já viu tudo o que o corvo e a harpia fazem — e eles
são 94% do que resta. A partir daí o jogo não introduz nada: não há nova ameaça
significativa, não há decisão nova, não há aumento real de pressão. Você repete o
mesmo gesto por mais três minutos até um chefe que você mata em cinco segundos
segurando o botão.

O diagnóstico preciso não é "o jogo é chato". É: **o jogo esgota seu vocabulário
em 90 segundos e depois dura mais 3 minutos e meio.**

### O combate é satisfatório?

**Momento a momento, sim. Estruturalmente, não.**

O que está bom, e é mérito real:
- A hitbox de 4px revelada no modo de precisão é design de shmup correto.
- A separação cromática dos projéteis é exemplar: magenta/roxo inimigo com núcleo
  claro e contorno ([`enemybullet.go:13`](game/enemybullet.go:13)), deliberadamente
  distinto das três armas do jogador. Muitos jogos comerciais erram isso.
- Os telegraphs antes de cada disparo inimigo ([`enemy.go:99`](game/enemy.go:99))
  são justiça de design.
- Flash branco de dano, partículas, hit-stop, screen shake, popups de pontuação —
  a camada de feedback existe e é competente.

O que quebra:
- **Não há razão para escolher.** Com a Luz sendo 3,3× melhor que as Chamas e
  ainda congelando tudo, o combate não tem estratégia. Existe uma resposta certa.
- **Não há razão para mirar.** Sem penalidade por fuga, deixar passar é grátis. O
  jogo nunca me obriga a lidar com uma ameaça — só a sobreviver a ela.
- **Não há razão para arriscar.** O grande gesto de risco/recompensa de shmup
  (chegar perto, raspar, pressionar) não é recompensado por nada.

Um combate satisfatório em shmup é um diálogo: a tela pergunta, você responde.
Aqui a tela faz a mesma pergunta 460 vezes e aceita qualquer resposta.

### Existe sensação de progressão?

**Dentro da partida: existe e acaba cedo demais.** Nível 1 → 3 em ~60 s, e
depois nada por 4 minutos. A curva de poder é um degrau, não uma rampa.

**Entre partidas: não existe.** O que persiste em [`save.go`](game/save.go) é:
recorde de campanha, recorde de sobrevivência, dificuldade, vibração, mudo e
volumes. Isto é *configuração*, não progressão. Ao fechar o jogo depois de uma
vitória, nada mudou — o jogo seguinte é idêntico ao primeiro.

Um jogador que termina a campanha não tem **nenhum motivo mecânico** para abrir o
jogo de novo. Só "bater o próprio número", que é a forma mais fraca e mais nichada
de retenção que existe.

### O jogador entende rapidamente o objetivo?

**Sim — e isso é uma vitória subestimada.** Sobe, atira, os monstros descem, uma
barra enche, um dragão aparece. O gênero é autoexplicativo e o jogo não atrapalha:
os avisos de trecho ("A vila está sob ataque!") funcionam como narração
diegética, a barra de progresso com marcos comunica avanço, o HUD é legível.

As falhas de clareza são pontuais, não estruturais:
- Nunca é ensinado que **Shift revela a hitbox real** — provavelmente a informação
  mais importante do jogo, e ela está escondida numa tela de texto no menu.
- A bomba (X/Ctrl) não tem onboarding; muitos jogadores terminam a campanha sem
  usar uma única carga (e o sistema de pontos os recompensa por isso).
- As três magias não comunicam *função*. "CHAMAS Nv2" no HUD não diz "cobertura
  local com queimadura".

### Existe identidade própria?

**Existe uma boa ideia de identidade. Ela ainda não chegou à tela.**

No papel: cavaleiro montado num grifo, criaturas aladas corrompidas, três magias
elementais, invocação de um dragão ancestral, reino de Valdoria caindo, biomas
que contam o cerco (campos → vila em chamas → muralhas → castelo). Isto é bom.
É específico, é evocativo, e não é o que os outros shmups fazem.

Na tela: o grifo é um sprite de 16×16 pixels legível mas genérico; o cenário são
**retângulos** ([`background.go:122-157`](game/background.go:122) — `drawCloud`,
`drawHill`, `drawStructure` são todos `vector.DrawFilledRect`); a fonte é a
`basicfont.Face7x13` da biblioteca padrão do Go, que é uma fonte de terminal; a
música é um arpejo de 8 notas em onda quadrada repetindo a cada 2,2 segundos.

**A distância entre a identidade concebida e a identidade percebida é o maior
problema de percepção do projeto.** Um jogador que vê um GIF deste jogo hoje não
consegue nomear o que ele é. Ele vê "um shoot'em up de pixel art".

### O jogo parece um protótipo técnico ou um jogo de verdade?

**Parece um protótipo técnico excepcionalmente bem-feito.**

Vou ser específico sobre o que é sinal de qual coisa, porque isso importa para
priorizar.

Sinais de jogo de verdade (mais do que a maioria dos projetos nesta fase tem):
- Máquina de estados completa, sem vazamento entre sessões, com fade.
- Persistência funcionando, tolerante a falhas.
- Dois modos de jogo, três dificuldades.
- 120 testes automatizados, `vet` limpo, build para Linux e Windows.
- Fallback gracioso em áudio (roda sem dispositivo de som) e em sprites.
- Modo de desenvolvimento com seed fixa, HUD de diagnóstico, hitboxes, salto de
  trecho — ferramentas que times profissionais frequentemente não têm.
- Fase descrita como **dados**, não código ([`stages.go`](game/stages.go)).

Sinais de protótipo:
- Cenário 100% geométrico.
- Fonte de sistema, não fonte do jogo.
- Música de 2,2 s em loop.
- `ebitenutil.DebugPrintAt` ainda usado para os popups de pontuação
  ([`popup.go:43`](game/popup.go:43)) e para o HUD de dev — resíduo da fonte de
  debug que o resto do jogo já abandonou.
- Nenhum gamepad.
- Balanceamento que ninguém mediu (o stun-lock atravessou 120 testes intacto —
  porque nenhum teste pergunta *se o jogo é divertido*, só se ele *funciona*).

O veredito: **a engenharia está em 8/10 e o jogo está em 5/10.** Isso é raro e é
uma boa posição — é muito mais fácil colocar conteúdo e sensação sobre uma base
sólida do que consertar uma base podre embaixo de conteúdo bonito. Mas exige
admitir que os próximos meses **não são de engenharia**.

### O que mais chamou minha atenção?

Cinco coisas, em ordem de surpresa:

1. **A disciplina de fallback.** Sprites e áudio degradam para versões procedurais
   sem quebrar; o áudio se desativa sozinho em WSL2 sem dispositivo. Isso é
   pensamento de engenharia de produto, não de hobby.

2. **`stages.go` é conteúdo, não código.** A decisão de tornar as fases
   declarativas foi a decisão técnica mais valiosa já tomada neste projeto. Ela é
   o motivo pelo qual tudo que proponho na seção 7 é viável.

3. **O chefe é bem melhor que o resto do jogo.** Três fases por percentual de
   vida, cinco padrões, limpeza de projéteis na troca de fase, cristais como
   pontos fracos, varredura com brecha garantidamente desviável, morte com
   invulnerabilidade concedida ao jogador para nunca haver derrota injusta no
   instante da vitória ([`game.go:544-555`](game/game.go:544)). Isso é design de
   shmup maduro. **E está sendo desperdiçado por um stun de 16 frames.**

4. **A separação `rng` / `fxRand`** ([`rng.go`](game/rng.go)) — aleatoriedade de
   jogabilidade isolada da de efeitos, para reprodução por seed. É uma sutileza
   que só aparece em times que já sofreram para reproduzir bugs.

5. **O README descreve um jogo mais rico do que o código entrega.** A regra de
   fuga pela base é o exemplo mais claro, mas o padrão se repete: seis inimigos
   descritos em detalhe, dos quais dois aparecem; três magias apresentadas como
   escolhas, das quais uma domina. **O documento é a visão; o código é o MVP; a
   distância entre eles é exatamente o backlog deste documento.**

### O que mais prejudica a experiência?

Em ordem estrita de dano causado:

1. **O stun-lock da Lança de Luz.** Apaga o combate, apaga o chefe, apaga a
   escolha de arma. É uma linha de código.
2. **A ausência de penalidade por fuga.** Remove a tensão central do gênero: o
   jogo nunca me pune por ignorá-lo.
3. **94% da campanha ser dois inimigos.** O conteúdo existe e não é usado. É a
   correção de maior retorno sobre esforço em todo o projeto — são edições em
   `stages.go`, não código novo.
4. **A curva de poder que satura em 60 s.** 85 power-ups sem significado.
5. **Áudio de 2,2 segundos em loop.** O áudio é o que mais rapidamente diz
   "protótipo" e o que menos custa terceirizar.
6. **Cenário geométrico.** O que mais rapidamente diz "protótipo" visualmente.
7. **Zero progressão entre partidas.** O que impede a retenção.
8. **Sem gamepad.** Barreira dura para o público de shmup e para o Steam Deck.

### Notas

Escala: **5 = MVP funcional competente**. **7 = padrão de jogo comercial
aceitável**. **9 = referência do gênero.** Notas comparam com o mercado de
shmups na Steam, não com "o que dá para esperar de um projeto solo".

| Eixo | Nota | Justificativa |
| --- | :---: | --- |
| **Gameplay** | **6,0** | O núcleo é sólido: movimento em 8 direções normalizado, modo de precisão com hitbox reduzida, bomba, escudo, três arquétipos de arma. O que impede nota alta é a **ausência de decisão**: a Luz domina por 3,3× de DPS mais controle, e a falta de penalidade por fuga elimina a pressão de eliminar. Tenho verbos suficientes, mas o jogo nunca me obriga a usá-los bem. |
| **Game Feel** | **7,0** | A área mais forte do projeto e a maior surpresa da análise. Hit-stop em abates grandes e trocas de fase do chefe, screen shake com três níveis de intensidade configuráveis, flash de dano, partículas com gravidade e arrasto, popups de pontuação com o multiplicador já aplicado, rastro em projéteis elementais, brasas no bocal das Chamas, escurecimento de tela na bomba, anel de partículas na quebra do escudo. Falta: hit-stop nos abates comuns, impacto visual de projétil inimigo ao colidir, recuo do disparo, e uma câmera que participe da ação. |
| **Progressão** | **3,5** | A nota mais baixa e a mais consequente. Dentro da partida a curva satura em 60 s de 300 s (nível 3 travado, ~85 power-ups irrelevantes depois disso). Entre partidas não existe progressão: `save.go` guarda dois recordes e quatro preferências. Não há um único desbloqueio, moeda, meta ou razão mecânica para reabrir o jogo. |
| **Clareza Visual** | **6,5** | Decisões acertadas e deliberadas: escurecimento do cenário em 104/255 de alfa para empurrar o parallax para trás ([`visual.go:54`](game/visual.go:54)), contorno escuro em todas as entidades, projéteis desenhados **por último** para nunca sumirem sob partículas, glow de 1px calibrado para leitura sem ofuscar. O que segura a nota: o parallax de retângulos gera ruído de alta frequência que compete com os inimigos, e 240×320 fica apertado nos picos de densidade. |
| **Direção de Arte** | **4,5** | Os sprites de entidades em pixel art procedural são decentes e o esquema de paleta por personagem é coerente. Mas o **cenário são retângulos**, a UI é funcional e sem personalidade, e não existe um documento ou regra de estilo — cada elemento foi resolvido isoladamente. O jogo não tem uma imagem que possa ser reconhecida fora de contexto. Este é o eixo que mais separa o projeto de parecer comercial. |
| **Interface** | **6,0** | Boa engenharia de UI: molduras de ferro/pergaminho, chips com sombra no HUD, medição real de largura de texto para alinhamento à direita ([`game.go:1520`](game/game.go:1520)), barras com marcos e limiares. Mas a fonte é `basicfont.Face7x13` — uma fonte de terminal do Go, não do jogo; as opções estão fragmentadas (dificuldade e vibração no menu principal, volume numa tela separada); e os popups ainda usam a fonte de debug. |
| **Balanceamento** | **3,5** | Três defeitos estruturais mensuráveis: (a) stun-lock da Luz anula inimigos e chefe; (b) DPS de 90 / 64 / 27 entre três armas apresentadas como equivalentes; (c) a economia recompensa não usar bombas (300 pts/carga) e morrer é parcialmente premiado (respawn com 5 HP contra 3 iniciais). Some-se o chefe valendo 5.000 contra 10–100 por inimigo, e a pontuação deixa de medir habilidade. |
| **Ritmo** | **4,5** | O ritmo *interno* das ondas foi trabalhado e está bom — ondas curtas, sobrepostas, sem espaços mortos, com jitter por dificuldade. O ritmo *macro* é plano: 63% corvos e 31% harpias significam que a fase 3 sente igual à fase 1. Não há crescendo, não há mini-clímax, não há respiro estruturado, não há chefe intermediário. Cinco minutos com a intensidade de noventa segundos. |
| **Rejogabilidade** | **3,5** | Existe mais do que parece — modo Sobrevivência com escalada por minuto, três dificuldades com jitter de spawn distinto, seed reproduzível — e ainda assim quase nada convida ao retorno, porque nenhuma partida muda a seguinte. Sobrevivência é o ativo subaproveitado: é o modo com maior retenção potencial e o que menos recebeu atenção de design. |
| **Potencial Comercial** | **5,5** | A base técnica vale um lançamento comercial: roda, é estável, compila para duas plataformas, tem save e opções. O que falta é tudo o que faz alguém *comprar*: um diferencial nomeável, uma imagem reconhecível, áudio produzido e um motivo para jogar dez vezes. Hoje é um projeto que eu apresentaria com orgulho num portfólio de engenharia e não colocaria numa loja. **Com as decisões deste documento, avalio o teto realista em 7,5–8,0** — um shmup indie de nicho, bem avaliado, com vendas modestas e boa reputação. |

**Média ponderada por impacto comercial: 5,0.** Um MVP que cumpriu o que
prometeu e chegou ao fim do seu ciclo de vida. O próximo ciclo é outro jogo.

---

## 2. A visão do jogo

> **Documento interno — Direção do Jogo**
> Circulação: equipe · Versão 1.0 · Substitui qualquer definição anterior de escopo

### 2.1 Declaração de visão

**Asas de Valdoria é um shoot'em up vertical de partidas curtas sobre um
cavaleiro e seu grifo defendendo um reino que está sendo corrompido em tempo
real — onde cada inimigo que você deixa passar torna o mundo, e o chefe, mais
corrompidos.**

Uma frase. Se uma proposta de feature não puder ser justificada por essa frase,
ela não entra.

### 2.2 O que estamos construindo (e o que não estamos)

| Estamos construindo | Não estamos construindo |
| --- | --- |
| Um shmup de **partidas de 12–18 minutos** com identidade forte | Uma campanha longa de 6 horas |
| Um jogo com **decisões de build dentro da partida** | Um roguelike com árvores de talento |
| Um jogo que se **aprende em 2 minutos e se domina em 20 horas** | Um bullet-hell para especialistas |
| Um jogo **desenhado para GIF** — legível em 3 segundos de vídeo | Um jogo que precisa ser explicado |
| Um jogo com **um gancho mecânico nomeável**: a Corrupção | Um jogo com muitos sistemas e nenhum slogan |

### 2.3 Público-alvo

**Núcleo (o jogador que compra na semana de lançamento — alvo de 60% das vendas
iniciais):**

O entusiasta de shmup moderno. Joga ou jogou *Sky Force Reloaded*, *Jamestown+*,
*ZeroRanger*, *Devil Blade Reboot*, *Blue Revolver*. Entre 25 e 45 anos. Não tem
mais tempo nem paciência para decorar padrões por 40 horas, mas quer profundidade
mecânica real. **Compra por respeito ao gênero: quer sentir que quem fez sabia o
que estava fazendo.** Vai perceber a hitbox de 4px, vai perceber a separação
cromática dos projéteis, e vai perdoar arte modesta se o jogo for justo.

**Secundário (o alvo de crescimento — onde estão as unidades que dobram o
resultado):**

O jogador de roguelite de sessão curta. Joga *Vampire Survivors*, *Brotato*,
*Dead Cells*, *Hades*. **Não busca um shmup** — busca "mais uma partida". Compra
se o jogo prometer variação entre runs e progressão que persista. Este público é
5 a 10× maior que o primeiro e é a razão pela qual a Corrupção e as Relíquias
(seção 3) existem no plano.

**Terciário (não desenhar para, mas não excluir):**

Nostálgico de 16-bit; público brasileiro de indies nacionais (o jogo é em
português, com tema medieval-fantástico e nome próprio — há espaço identitário
real aqui); streamers de sessão curta.

**Anti-público — quem NÃO queremos e por quê:** o purista de bullet-hell de
contagem de padrões (quer 2.000 projéteis na tela e ranking mundial: não é o que
vamos entregar, e tentar agradá-lo destrói o secundário) e o jogador casual de
mobile (espera progressão automática e sessões de 90 segundos).

### 2.4 Plataforma

| Prioridade | Plataforma | Justificativa |
| --- | --- | --- |
| **1** | **Steam (Windows)** | Onde o público de shmup compra. `make build-windows` já produz um `.exe` autocontido, sem CGO. |
| **2** | **Steam Deck / Linux** | Ebitengine roda; o build Linux já existe. Formato retrato exige moldura decorada, mas o Deck é o público perfeito para partidas de 15 min. **Verified é meta explícita de 1.0** e exige gamepad completo. |
| 3 | itch.io | Lançamento simultâneo, custo zero, canal de demo e de comunidade. |
| Depois | Nintendo Switch | O melhor lar comercial para o gênero. **Só depois de 1.0 na Steam** e mediante viabilidade de porte do Ebitengine — é decisão de 2027, não de agora. |
| Nunca | Mobile / Web | Mobile exige repensar o controle inteiro. Web (WASM) é tentador mas divide o esforço de build e QA sem retorno de vendas. |

**Resolução e janela:** manter 240×320 interno — é a resolução correta para o
gênero e para o custo de arte. Mas 240×320 escalado em 16:9 desperdiça ~70% da
tela. **A decisão de 1.0 é adotar moldura decorada lateral** (referência: *Sky
Force*, *Ikaruga*, *Jamestown*): pergaminho/pedra nas laterais exibindo pontuação,
vidas, runa equipada e o **Medidor de Corrupção**. Isso libera o HUD de cima do
campo de jogo e transforma um problema técnico em identidade visual.

### 2.5 Duração de uma run

| Métrica | Hoje | Alvo 1.0 |
| --- | ---: | ---: |
| Campanha completa | ~5 min | **14–18 min** |
| Fase individual | 50–93 s | **2,5–3 min** |
| Boss fight | 5 s (Luz) a ~90 s | **90–150 s** |
| Sobrevivência (partida típica) | indefinido | **8–12 min** |
| Sessão típica de jogador | 1 run | **3–5 runs (~1 h)** |

**Por que 14–18 minutos é o número certo:** é curto o bastante para "mais uma
antes de dormir" (o loop de *Vampire Survivors* / *Hades*) e longo o bastante
para que uma run tenha **arco dramático** — abertura, complicação, clímax. Cinco
minutos é curto demais para uma run ter memória; trinta minutos é longo demais
para se perder e querer repetir.

### 2.6 Diferencial (o elevator pitch mecânico)

**"O reino se corrompe enquanto você joga — e você escolhe se luta contra isso ou
se aproveita disso."**

Cada inimigo que atravessa a base da tela aumenta o **Medidor de Corrupção**. A
Corrupção não é uma barra de derrota: é um **modificador global de dificuldade e
recompensa**. Corrupção alta muda o cenário, muta os inimigos, deixa Vharak mais
forte — **e multiplica pontos, drops e a qualidade das relíquias oferecidas**.

O jogador que deixa passar por incompetência é punido. O jogador que deixa passar
**de propósito** está apostando. Essa é a decisão que nenhum outro shmup oferece,
ela é 100% original em relação aos concorrentes listados na seção 13, e ela nasce
diretamente de uma mecânica que o projeto já documentou e nunca implementou.

### 2.7 A fantasia do jogador

**"Eu e minha montaria somos a última coisa entre este reino e o que está vindo
por ele."**

Três componentes que precisam estar em cada decisão de design:

1. **Você é uma dupla, não uma nave.** O grifo é um personagem: ele cansa, ele
   grita, ele reage. Quando você morre, ele cai com você. Nenhum concorrente tem
   isto — todos pilotam máquinas.
2. **Você está defendendo, não invadindo.** Os cenários são lugares que você
   conhece: os campos do reino, a vila, as muralhas. A vila queima **atrás** de
   você.
3. **Você é magia, não arsenal.** Você não coleta armas; você canaliza runas
   ancestrais. E pode **fundi-las** (seção 3).

### 2.8 A emoção principal

**Tensão heroica em crescendo.**

Não é o pânico do bullet-hell nem o poder de *Vampire Survivors*. É a emoção
específica de **estar segurando uma linha que está cedendo**. O jogo deve fazer o
jogador dizer, em voz alta, "não, não, NÃO — ah, consegui".

Consequência de design: a Corrupção precisa ser **visível e crescente**, o áudio
precisa reagir a ela, e o clímax precisa ser um confronto que o jogador
**sente que quase perdeu**.

### 2.9 Identidade em três palavras

**Nobre. Corrompido. Alado.**

Filtro para qualquer decisão futura:

| Palavra | Significa que… | Rejeitamos… |
| --- | --- | --- |
| **Nobre** | Dourado, heráldica, pergaminho, ferro, cerimônia. O jogador é um cavaleiro. | Estética sci-fi, neon, tecnologia. |
| **Corrompido** | Magenta/violeta doente, veias, cristais, deformação orgânica. O inimigo é o reino adoecendo. | Demônios genéricos, zumbis, robôs. |
| **Alado** | Tudo voa, plana, mergulha. Movimento é ar, não propulsão. | Inimigos terrestres imóveis sem justificativa; movimento mecânico. |

### 2.10 Metas mensuráveis de 1.0

Se não medirmos, é opinião. Estas são as metas contra as quais o projeto será
avaliado no lançamento:

| Meta | Alvo |
| --- | --- |
| Taxa de conclusão da campanha (Normal) | 45–60% dos jogadores |
| Runs por jogador na primeira sessão | ≥ 3 |
| Tempo mediano jogado antes do refund (2 h) | ≥ 90 min |
| Avaliações positivas na Steam | ≥ 85% |
| Jogadores que testam ≥ 2 magias diferentes | ≥ 70% (hoje seria ~0%) |
| Jogadores que usam a bomba ao menos 1× por run | ≥ 80% (hoje, provavelmente < 30%) |
| Steam Deck Verified | Sim |

---

## 3. O que torna este jogo único

Existem centenas de shoot'em ups. A pergunta correta não é "o que podemos
adicionar", é **"o que só este jogo pode fazer, dado o que ele já é?"**

Abaixo, 24 ideias. Elas não são um catálogo de desejos — cada uma nasce de algo
já presente no código ou na ficção do projeto.

### 3.1 As 24 ideias

**Grupo A — Mecânicas de sistema (mudam o que a partida É)**

| # | Ideia | Descrição | Nasce de |
| --- | --- | --- | --- |
| 1 | **Medidor de Corrupção** | Inimigos que escapam pela base enchem uma barra global. Corrupção alta muta inimigos, escurece o bioma, fortalece Vharak — e multiplica pontos, drops e raridade de relíquias. Um recurso de risco negociado. | A regra fantasma do README + `escaped` já rastreado em `game.go:874` |
| 2 | **Runas Fundidas** | Coletar a runa de um segundo elemento com o atual no nível máximo **funde** as duas: Luz+Fogo = *Julgamento Solar*; Fogo+Gelo = *Tempestade de Vapor*; Gelo+Luz = *Prisma*. 3 armas base + 3 fusões. | `gainWeapon` já preserva nível ao trocar; `weaponType` e `element` já existem em `bullet.go:23` |
| 3 | **Altar entre Trechos** | Ao passar de trecho, o scroll pausa por 4 s num altar suspenso e o jogador **voa até uma de três bênçãos**. Escolha diegética, sem menu. | `sectionDef` já marca transições; `announce` já existe |
| 4 | **Pactos da Corrupção** | Com Corrupção acima de 50%, o altar oferece uma quarta opção maldita: poder alto por um custo permanente na run ("seus tiros ferem você" / "o grifo perde uma asa"). | Ideias 1 + 3 |
| 5 | **Sistema de Rank dinâmico** | Jogar bem (combo, sem dano, sem fugas) aumenta o Rank; o jogo responde com mais inimigos e padrões mais densos. Jogar mal alivia. Auto-balanceamento invisível. | `campaignDensityScale` e `jitterSpawn*` já parametrizam densidade |
| 6 | **Reações elementais** | Congelar e depois acertar = **estilhaçar** (dano em área). Queimar + gelo = **vapor** (cobertura visual). Luz em inimigo congelado = **refração** (cadeia dobrada). | `status.go` já rastreia `slow`/`burn`/`stun` simultaneamente |
| 7 | **Rota ramificada do Reino** | Entre fases, escolher o próximo destino num mapa: rota Segura, Arriscada (mais corrupção, melhor recompensa) ou Perdida (bioma secreto). | `campaignStages()` já devolve uma lista ordenada — trivial de tornar um grafo |
| 8 | **Fôlego do Grifo** | O grifo tem estamina: gasta em mergulho e em manobras, regenera parado. Transforma movimento em recurso. | Ideia 10 |

**Grupo B — Verbos novos (mudam o que o jogador FAZ)**

| # | Ideia | Descrição | Nasce de |
| --- | --- | --- | --- |
| 9 | **Mergulho do Grifo** | Botão dedicado: o grifo mergulha para frente 0,4 s, invulnerável, causando dano por contato e **empurrando a câmera para cima** — você avança contra a rolagem. | A fantasia da montaria; `hitStop` e `shake` já existem para vender o impacto |
| 10 | **Bater de asas (repelir)** | Segurar precisão + direção gera uma rajada de vento que **desvia projéteis inimigos próximos** por meio segundo. Custa fôlego. | `enemyBullets` já é uma lista mutável acessível |
| 11 | **Grito do Grifo** | Um grito que **atordoa** todos os inimigos na tela por 1 s, sem dano. Alternativa tática à bomba. | O stun já existe e está mal alocado (hoje é passivo na Luz) |
| 12 | **Raspar (graze)** | Passar a menos de 6px de um projétil inimigo sem ser atingido carrega o medidor de bomba. Converte medo em ganância. | Hitbox de 4px já implementada e revelada na precisão |
| 13 | **Bomba com direção** | A Invocação Ancestral passa a ser mirada: o dragão atravessa na direção apontada, e não sempre de baixo para cima. | `drawBombEffect` já anima uma travessia |
| 14 | **Cavaleiro Caído** | Perder uma vida não é morte instantânea: o grifo despenca e você tem 6 segundos em queda livre, sem tiro, para reencontrá-lo. Se conseguir, mantém a arma; se falhar, perde nível. | `respawn()` já existe como momento distinto |

**Grupo C — Conteúdo com personalidade (mudam o que o jogador VÊ)**

| # | Ideia | Descrição | Nasce de |
| --- | --- | --- | --- |
| 15 | **Chefes com partes destrutíveis** | Asas, cauda e cabeça com vida própria. Destruir uma asa remove um padrão e faz o chefe voar torto. | `bossWeakSpots` já é uma lista de retângulos com bônus de dano |
| 16 | **Chefe de perseguição** | Um chefe que vem **de baixo**, invertendo a leitura da tela e forçando o jogador para cima. | `Boss` já tem entrada, fases e movimento parametrizados |
| 17 | **Cenário destrutível** | Torres, catapultas e telhados no parallax de estruturas viram alvos: destruí-los solta runas e altera a silhueta do bioma. | `structures []*bgItem` já existe e rola na tela |
| 18 | **Clima dinâmico** | Vento lateral que empurra o jogador; névoa que reduz visibilidade; tempestade que ilumina a tela em flashes. Um por bioma. | Camadas de parallax já independentes |
| 19 | **Eventos de trecho** | Escoltar uma caravana; sobreviver 30 s sem atirar; abater 20 corvos antes que fujam. Recompensa: relíquia. | `waveDef` já é dado; eventos são só um novo tipo de onda |
| 20 | **A Ordem dos Cavaleiros** | 3–4 cavaleiros jogáveis, cada um com runa inicial, passiva e montaria distintas. | `newPlayer()` centraliza toda a configuração inicial |
| 21 | **Diário de Criaturas** | Bestiário que se preenche por abates; completar uma entrada concede um bônus passivo permanente pequeno. | `enemiesDefeated` já contado; `save.go` já persiste |
| 22 | **Vharak como rival persistente** | Vharak sobrevive à primeira derrota e retorna em runs seguintes mais forte, comentando o desempenho anterior do jogador. | `save.go` pode guardar o histórico |
| 23 | **Vilas que você salva** | Vilas defendidas com sucesso reaparecem reconstruídas em runs futuras — narrativa que persiste no cenário. | Temas de bioma já por trecho |
| 24 | **Fantasma do recorde** | Um grifo espectral repete sua melhor run ao seu lado. | Seed reproduzível já implementada |

### 3.2 As 5 escolhidas

Critérios de escolha: **(a)** é nomeável numa frase de loja; **(b)** aproveita
código existente; **(c)** aumenta rejogabilidade sem virar outro gênero; **(d)**
gera imagem de marketing; **(e)** um projeto pequeno consegue executar bem.

---

#### 🥇 1. Medidor de Corrupção (#1) — *o gancho do jogo*

**O quê:** uma barra global de 0 a 100%. Cada inimigo que atravessa a base
adiciona corrupção (corvo +1, wyvern +4). Ela nunca desce sozinha — só é
purificada em altares ou por eventos.

Faixas:

| Corrupção | Efeito no mundo | Efeito na recompensa |
| --- | --- | --- |
| 0–25% *Reino Firme* | Normal | Base |
| 26–50% *Sombra Crescente* | Paleta esmaece, inimigos ganham +15% vida | Pontos ×1,3 · drops +20% |
| 51–75% *Cerco* | Inimigos **mutam** (corvo vira Corvo-Sombra e atira) | Pontos ×1,8 · relíquias raras liberadas · pactos disponíveis |
| 76–99% *Colapso* | Bioma vira versão corrompida; música distorce | Pontos ×2,5 |
| 100% *Queda de Valdoria* | **Vharak Ascendido**: chefe verdadeiro, 4 fases | Final alternativo + recompensa máxima |

**Por que é a melhor ideia do documento:**

1. **É o diferencial que não existe no gênero.** Todo shmup pune quem deixa
   passar (perde ponto, perde vida). Nenhum transforma isso numa **aposta com
   payoff crescente**. Isso é dizível em uma frase na página da Steam, é
   compreensível num GIF de 6 segundos, e não tem concorrente direto.
2. **Conserta o buraco mais grave do MVP** — a ausência de penalidade por fuga —
   sem adicionar uma punição chata. Converte um bug de design em identidade.
3. **É barato.** `escaped` já é rastreado. É um `int` no `Game`, uma barra na
   moldura lateral, um multiplicador aplicado em `registerKill`, e variantes de
   inimigo que reusam `enemyBehavior`.
4. **Dá arco dramático a cada run**, resolvendo o problema de ritmo plano: a
   intensidade agora sobe por causa das **suas** decisões, não de um cronograma.
5. **Justifica o final alternativo**, que é o gancho de rejogo mais barato que
   existe: "existe outro chefe, e você precisa querer perder para vê-lo".

**Risco:** jogadores podem farmar corrupção artificialmente. *Mitigação:* a
corrupção também aumenta a chance real de morte de forma superlinear, e não há
como reduzi-la depois de subir.

---

#### 🥈 2. Runas Fundidas (#2) — *o motor de builds*

**O quê:** com uma runa no nível 3, coletar a de outro elemento **funde**. Seis
armas finais em vez de três, cada fusão com identidade visual e sonora própria:

| Fusão | Nome | Identidade |
| --- | --- | --- |
| Luz + Fogo | **Julgamento Solar** | Raio vertical contínuo que incendeia tudo na coluna |
| Fogo + Gelo | **Tempestade de Vapor** | Nuvem persistente que causa dano por área e obscurece inimigos |
| Gelo + Luz | **Prisma** | Um projétil pesado que **refrata** em três ao acertar |

**Por quê:**

1. **Conserta a curva de poder que satura em 60 s.** Com fusões, a progressão da
   run continua até o minuto 8. Os ~85 power-ups hoje irrelevantes voltam a
   significar algo.
2. **Transforma power-up em decisão.** "Fundo agora ou seguro esperando a runa
   que quero?" — a primeira decisão interessante de recurso do jogo.
3. **Gera conversa e rejogo:** "qual fusão é a melhor?" é o tipo de pergunta que
   move comunidade e vídeos.
4. **É barato:** `element` já viaja no projétil; `status.go` já aplica efeitos por
   elemento. Uma tabela de definição de armas elimina os `switch` paralelos e
   abre espaço para as seis.

**Pré-requisito obrigatório:** rebalancear as três armas base **antes** de fundir
qualquer coisa. Fundir armas desbalanceadas produz fusões desbalanceadas.

---

#### 🥉 3. Altar entre Trechos (#3 + #4) — *o momento de respiro e escolha*

**O quê:** ao trocar de trecho, a rolagem desacelera, a música baixa, um altar de
pedra surge no centro com três luzes. Você **voa até uma delas**. Não há menu, não
há pausa — é diegético e leva 4 segundos.

Ofertas: relíquias passivas ("+20% de velocidade de mergulho"), purificação
(−15% de Corrupção), ouro, carga de bomba, ou — acima de 50% de Corrupção — um
**Pacto**: poder alto com custo permanente na run.

**Por quê:**

1. **Resolve o ritmo plano.** Dá ao jogo a estrutura de respiração
   tensão → alívio → escolha → tensão que *Hades* usa e que este projeto não tem.
2. **É o momento de identidade.** É o instante em que a run vira *sua* run — o
   que o jogador conta para outra pessoa.
3. **É a fonte de variação entre partidas** sem tocar em meta-progressão pesada.
4. **É barato:** `sectionDef` já marca as transições. Uma pausa na timeline, três
   entidades coletáveis, uma lista de relíquias em dados.

---

#### 4. Mergulho do Grifo (#9) — *o verbo de assinatura*

**O quê:** um botão dedicado. O grifo mergulha para frente por 0,4 s:
invulnerável, causando dano de contato, com a câmera empurrando para cima. Custa
fôlego e tem ~1,5 s de recarga.

**Por quê:**

1. **É a única mecânica que faz o jogador *sentir* que está montado num animal.**
   Naves não mergulham. É a diferença física entre este jogo e todos os
   concorrentes.
2. **É o melhor GIF possível do projeto** — 3 segundos, legível, com hit-stop e
   deslocamento de câmera.
3. **Adiciona profundidade de habilidade** sem adicionar complexidade de sistema:
   iniciantes ignoram, veteranos usam para atravessar padrões e maximizar dano.
4. **Cria o momento de graça:** mergulhar por dentro de uma parede de projéteis é
   exatamente a "tensão heroica" definida na seção 2.8.

**Risco:** invulnerabilidade quebra o desafio. *Mitigação:* fôlego limitado,
recarga longa, e vulnerabilidade aumentada nos 0,3 s após o mergulho.

---

#### 5. Chefes com partes destrutíveis (#15) — *o clímax que merece a build*

**O quê:** asas, cauda e cabeça com vida própria. Destruir a asa esquerda remove
o padrão de varredura e faz Vharak voar inclinado. A cauda remove as invocações.
Só depois de mutilá-lo o núcleo fica exposto.

**Por quê:**

1. **O chefe já é a melhor peça do jogo** e está sendo desperdiçado. Este é o
   investimento com melhor relação entre qualidade existente e esforço.
2. **Torna o chefe uma conversa, não uma barra.** O jogador aprende, decide onde
   atirar, e vê o resultado da decisão imediatamente no comportamento do inimigo.
3. **Torna as builds significativas no momento que mais importa:** o Gelo
   perfurante brilha contra partes alinhadas; a Luz em cadeia limpa a cauda; o
   Vapor cobre a cabeça.
4. **`bossWeakSpots` já é uma lista de retângulos com multiplicador de dano** —
   a estrutura de dados correta já está lá, faltando vida própria e consequência.

---

### 3.3 O que fica de fora, e por quê

| Ideia | Veredito | Motivo |
| --- | --- | --- |
| Rota ramificada (#7) | **Adiar para 1.1** | Excelente, mas multiplica QA de balanceamento por 3 antes de o balanceamento base estar resolvido. |
| Ordem dos Cavaleiros (#20) | **Adiar para DLC/1.2** | Melhor gancho de conteúdo pós-lançamento que existe. Fazer agora multiplica arte e balanceamento antes de haver um jogo. |
| Rank dinâmico (#5) | **Adiar** | Sobrepõe-se à Corrupção. Dois sistemas invisíveis de dificuldade competindo é confusão. |
| Cavaleiro Caído (#14) | **Considerar em 0.7** | Charmoso, mas é um minigame; risco de quebrar o ritmo. |
| Fantasma do recorde (#24) | **Considerar em 1.1** | Barato e ótimo para comunidade, mas não vende o jogo. |
| Reações elementais (#6) | **Incorporar às Fusões** | Não é sistema separado; é o comportamento natural das Runas Fundidas. |

---

## 4. Evolução da jogabilidade

Regra desta seção: **escolher poucos sistemas e executá-los bem.** Um sistema
medíocre custa o mesmo que um bom em manutenção e QA, e vale muito menos.

### 4.1 Sistemas aprovados

| Sistema | Prioridade | O que resolve |
| --- | :---: | --- |
| **Rebalanceamento das armas** | 🔴 **P0** | Stun-lock, DPS 90/64/27, ausência de escolha |
| **Medidor de Corrupção** | 🔴 **P0** | Falta de penalidade por fuga, ritmo plano, diferencial |
| **Runas Fundidas** | 🟠 **P1** | Curva de poder que satura em 60 s |
| **Altar + Relíquias** | 🟠 **P1** | Ausência de respiro e de decisão de build |
| **Mergulho do Grifo** | 🟠 **P1** | Ausência de verbo de assinatura |
| **Camada de input (ações)** | 🟠 **P1** | Bloqueio de gamepad, Steam Deck e remapeamento |
| **Partes destrutíveis + 2º chefe** | 🟡 **P2** | Clímax subaproveitado, campanha sem meio |
| **Redistribuição do bestiário** | 🔴 **P0** | 94% da campanha ser dois inimigos |
| **Graze → carga de bomba** | 🟡 **P2** | Bombas hoardadas, ausência de risco/recompensa |
| **Eventos de trecho** | 🟡 **P2** | Monotonia estrutural das ondas |
| **Clima por bioma** | 🟢 **P3** | Identidade de bioma além da paleta |

### 4.2 Detalhamento dos P0

#### Rebalanceamento das armas (o item mais urgente do projeto)

Três correções, todas pequenas:

1. **Stun não é passivo de arma.** Remover o stun automático da Luz
   ([`status.go:19`](game/status.go:19)). O stun vira recurso ativo (o Grito do
   Grifo) ou proc raro com cooldown por alvo. *Se nada mais desta seção for feito,
   faça isto.*
2. **Nivelar DPS em ~65 na janela ideal de cada arma.** Não igualar números
   brutos: dar a cada arma um contexto onde ela é a melhor.
3. **Definir identidades sem sobreposição:**

| Arma | Identidade | Melhor contra | Pior contra |
| --- | --- | --- | --- |
| **Lança de Luz** | Precisão. DPS alto **apenas** no eixo vertical estreito. Sem cadeia automática. | Alvo único, chefes | Enxames dispersos |
| **Chamas do Dragão** | Cobertura. Leque curto, queimadura forte, vulnerabilidade. | Enxames próximos, formações | Alvos distantes |
| **Lanças de Gelo** | Controle. Perfuração e lentidão pesada; DPS baixo, valor tático alto. | Colunas alinhadas, sobrevivência | Chefes solitários |

**Métrica de sucesso:** ≥ 70% dos jogadores experimentam ao menos duas magias
numa sessão; nenhuma arma acima de 45% de uso.

#### Redistribuição do bestiário (a maior vitória por esforço do projeto)

Não é código: é edição de dados em [`stages.go`](game/stages.go).

| Inimigo | Hoje | Alvo |
| --- | ---: | ---: |
| Corvo | 63,5% | **35%** |
| Harpia | 31,1% | **22%** |
| Gárgula | 1,1% | **12%** |
| Wyvern | 2,2% | **10%** |
| Mago | 1,5% | **10%** |
| Balista | 0,6% | **8%** |
| Novos (Corrompidos) | — | **3%** |

Regra de composição por fase: **cada fase apresenta um inimigo novo em isolamento
(sozinho, sem pressão), depois em combinação, depois em massa.** É a gramática de
*Super Aleste* e de todo shmup bem estruturado — e é gratuita, porque os inimigos
já existem.

### 4.3 Sistemas explicitamente rejeitados

| Sistema | Por que não |
| --- | --- |
| Árvore de habilidades | Menu antes de jogar; o gênero é sobre executar, não configurar. |
| Elementos como triângulo de fraquezas | Vira contabilidade; transforma escolha em resposta certa. |
| Crafting | Não há fantasia que sustente; adiciona tela e economia sem diversão. |
| NPCs e diálogo | Corta o ritmo; a narrativa deve ser ambiental. |
| Missões/contratos | Adia a diversão; funciona em jogo de serviço, não em shmup de 15 min. |
| Dificuldade dinâmica invisível | Já temos Corrupção, que é visível e agenciada pelo jogador. Melhor. |

---

## 5. Progressão entre partidas

### 5.1 O problema

Hoje o `save.go` persiste dois recordes e quatro preferências. Quando o jogador
vence a campanha, o jogo diz "FASE CONCLUÍDA!" e volta ao menu, exatamente igual
ao que era antes. **Não existe motivo mecânico para uma segunda partida.**

### 5.2 Os três modelos possíveis

| Modelo | Como funciona | Prós | Contras |
| --- | --- | --- | --- |
| **A. Arcade puro** | Sem progressão. Só recordes e dificuldades. *(estado atual)* | Puro; respeita o gênero; zero custo | Retenção quase nula fora do nicho hardcore; mata o público secundário |
| **B. Meta-progressão de poder** | Moeda persistente compra upgrades permanentes de vida/dano | Retenção altíssima; loop *Sky Force* comprovado | **Destrói o balanceamento**: precisa funcionar com jogador fraco e forte simultaneamente; grinding obrigatório; risco de virar pay-to-progress sem pagamento |
| **C. Meta-progressão de *variedade*** ⭐ | Moeda persistente desbloqueia **novas opções**, nunca poder bruto: novas relíquias no pool, novas runas iniciais, novos cavaleiros, novos modificadores, entradas do bestiário | Retenção real; **não inflaciona poder**; cada desbloqueio é uma razão para jogar diferente, não mais fácil; balanceamento estável | Retenção um pouco menor que B; exige conteúdo de verdade para desbloquear |

### 5.3 Modelo escolhido: C — Progressão de variedade

**Por quê:** o Modelo B é a armadilha clássica do gênero. Se upgrades permanentes
tornam o jogo mais fácil, a campanha precisa ser balanceada duas vezes — para o
jogador da primeira run e para o da vigésima — e nenhuma das duas fica boa. Pior:
o jogador que morre repetidamente aprende que a solução é *grindar*, não
*melhorar*. Isso é exatamente o oposto da promessa "aprende em 2 minutos, domina
em 20 horas".

O Modelo C mantém a curva intacta e converte tempo jogado em **leque de
possibilidades**. É o que *Hades* faz com os Espelhos e o que *Dead Cells* faz com
os planos: você não fica mais forte, você fica com mais opções.

### 5.4 Desenho concreto

**Moeda: Fragmentos de Coroa.** Caem de mini-chefes, eventos, altares e ao
concluir fases. Uma run típica rende 80–150; uma run excelente, 250. Corrupção
alta multiplica o ganho — amarrando meta-progressão ao gancho principal.

**Loja: A Cidadela** (menu principal, com identidade visual de sala do trono).

| Categoria | Itens | Custo | O que dá |
| --- | ---: | ---: | --- |
| **Relíquias** | 24 | 150–400 | Entram no pool dos altares. Não são equipadas — são **possibilidades**. |
| **Cavaleiros** | 3 além do inicial | 800–1.500 | Runa inicial, passiva e montaria distintas |
| **Runas iniciais** | 3 | 300 | Começar com Chamas ou Gelo |
| **Modificadores de run** | 8 | 200–600 | Opcionais, aumentam pontuação: "sem bombas", "corrupção inicial 30%", "vida única" |
| **Diário de Criaturas** | 14 entradas | grátis (por abates) | Lore + bônus cosmético; completar tudo desbloqueia o modo Boss Rush |
| **Cosméticos** | 12 | 100–300 | Plumagens do grifo, brasões, molduras de HUD |

**Conquistas: 20**, ancoradas em feitos, não em grind ("Vença com 90% de
Corrupção", "Derrote Vharak sem usar bombas", "Funda todas as runas numa run").

### 5.5 Vantagens e desvantagens honestas

**Vantagens:**
- Cada sessão termina com ganho tangível, mesmo em derrota.
- O balanceamento da campanha permanece estável para sempre.
- Desbloqueios são conteúdo *criado*, não números aumentados — valem mais por
  unidade de esforço.
- Conecta-se naturalmente à Corrupção e às Relíquias.

**Desvantagens (e mitigações):**
- **Retenção menor que a de upgrades de poder.** *Mitigação:* modificadores de run
  e Boss Rush dão metas de longo prazo.
- **Exige conteúdo real:** 24 relíquias e 3 cavaleiros não são triviais.
  *Mitigação:* relíquias são regras simples e baratas; cavaleiros só na 0.9+.
- **Risco de moeda insossa.** *Mitigação:* Fragmentos caem **visivelmente na
  tela** com som próprio — feedback imediato, não um número na tela de resultado.
- **Risco de paralisia por opção:** 24 relíquias no pool inicial confundem.
  *Mitigação:* pool inicial de 8; as demais entram gradualmente.

### 5.6 O que explicitamente NÃO fazer

Upgrades permanentes de dano/vida/velocidade · níveis de personagem · energia ou
tempo de espera · loot randômico com raridades · passe de temporada · qualquer
coisa que faça o jogador esperar em vez de jogar.

---

## 6. Conteúdo da primeira versão comercial

### 6.1 Premissas de produção

Antes dos números, as premissas — sem elas o cronograma é ficção:

- **Equipe:** 1 desenvolvedor (o autor), meio período (~20 h/semana).
- **Terceirização:** pixel art de cenários e trilha sonora **devem** ser
  contratadas. É a decisão de produção mais importante deste documento.
- **Unidade:** 1 "semana-dev" = 20 h de trabalho efetivo.
- **Base:** ~8.700 linhas de Go já escritas, com fases em dados e 120 testes.

### 6.2 Escopo da 1.0

| Eixo | MVP hoje | **Alvo 1.0** | Justificativa |
| --- | ---: | ---: | --- |
| **Fases** | 3 | **6** | 6 × 2,5 min = 15 min de campanha. Abaixo de 5, a run não tem arco; acima de 8, o custo de arte e QA explode. |
| **Biomas** | 6 temas / 4 estilos | **6 biomas completos** | Um por fase, cada um com paleta, estruturas próprias, clima e trilha. É o mínimo para o jogo não parecer repetitivo em vídeo. |
| **Chefes** | 1 | **4** | Vharak (final) + 2 mini-chefes (fases 2 e 4) + Vharak Ascendido (Corrupção 100%). Ritmo de shmup pede um confronto a cada ~2 fases. |
| **Inimigos** | 6 | **14** | 6 atuais + 5 novos + 3 variantes corrompidas. Regra: ~2,3 por fase, cada um introduzido isoladamente. |
| **Armas** | 3 × 3 níveis | **3 base + 3 fusões** | Seis expressões finais é variedade real sem seis balanceamentos independentes desde o zero. |
| **Power-ups** | 5 | **8** | Atuais + Fragmento de Coroa, Carga de Bomba, Purificação (−Corrupção). |
| **Relíquias** | 0 | **24** | O pool que dá variação entre runs. Cada uma é uma regra pequena — de longe o conteúdo mais barato por unidade de diversão. |
| **Personagens** | 1 | **2** (+2 em DLC) | Um segundo cavaleiro dobra a rejogabilidade percebida. Quatro é DLC. |
| **Modos** | 2 | **4** | Campanha, Sobrevivência, **Boss Rush**, **Campanha com Modificadores**. |
| **Dificuldades** | 3 | **4** | + Cavaleiro (para veteranos, com ranking). |
| **Trilhas** | 11 × 2,2 s | **10 × 2–3 min** | Produzidas. O maior salto de percepção por real gasto. |
| **Efeitos sonoros** | 11 gerados | **50 produzidos** | Inclui vozes do grifo e do chefe. |
| **Idiomas** | pt-BR | **pt-BR + inglês** | Inglês não é opcional para vender na Steam. |
| **Duração de campanha** | ~5 min | **14–18 min** | Seção 2.5. |
| **Conteúdo total até dominar** | ~30 min | **15–25 h** | 6 fases × 4 dificuldades × 2 personagens × 24 relíquias × 2 finais. |

### 6.3 Estrutura da campanha

| # | Fase | Bioma | Introduz | Clímax |
| --- | --- | --- | --- | --- |
| 1 | O Cerco de Eldoria | Campos ao amanhecer | Corvo, Harpia · tutorial diegético | Onda em formação |
| 2 | A Vila em Chamas | Vila incendiada, fumaça | Gárgula, Balista · **primeiro Altar** | 🔶 **Mini-chefe: A Catapulta Viva** |
| 3 | As Muralhas de Ferro | Fortaleza, tempestade | Mago, Cavaleiro Corrompido · clima (vento) | Corredor de balistas |
| 4 | A Floresta Corrompida | Bosque doente, névoa | Wyvern, Aracnídeo Alado · fusões liberadas | 🔶 **Mini-chefe: O Guardião Apodrecido** |
| 5 | O Desfiladeiro | Cânion vulcânico | Variantes corrompidas · cenário destrutível | Perseguição |
| 6 | O Covil de Vharak | Covil, cristais | — · último Altar | 🔴 **Vharak** (ou **Vharak Ascendido** com Corrupção 100%) |

### 6.4 Estimativa de desenvolvimento

| Bloco | Semanas-dev | Detalhe |
| --- | ---: | --- |
| Balanceamento e correções de núcleo | 4 | Armas, stun, bestiário, economia |
| Corrupção (sistema + variantes + final alternativo) | 6 | Núcleo do diferencial |
| Runas Fundidas | 4 | Tabela de armas + 3 fusões + arte/som |
| Altar + 24 relíquias | 5 | Sistema + conteúdo em dados |
| Mergulho + fôlego + graze | 3 | Verbos novos, muito game feel |
| Camada de input + gamepad + Steam Deck | 3 | Habilitador de plataforma |
| 3 fases novas (6 total) | 6 | Dados + tuning; arte em paralelo |
| 3 chefes novos + partes destrutíveis | 8 | Bloco mais caro de código |
| 8 inimigos novos | 5 | Reusa `enemyBehavior` |
| Meta-progressão + Cidadela + conquistas | 5 | Sistema + UI |
| Direção artística (integração de arte terceirizada) | 6 | Cenários, UI, molduras, efeitos |
| Integração de áudio produzido | 3 | Sistema já suporta arquivos |
| Modos extras (Boss Rush, modificadores) | 2 | Reusa tudo |
| Onboarding, i18n, telas, polimento | 4 | |
| QA, playtest externo, correções | 6 | **Não negociável** |
| Loja Steam, trailer, capturas, build | 3 | |
| **Subtotal** | **73** | |
| Contingência (25%) | 18 | Realismo de cronograma |
| **TOTAL** | **≈ 91 semanas-dev** | |

**Tradução em calendário:**

| Cenário | Prazo | Observação |
| --- | --- | --- |
| Solo, 20 h/semana | **~21 meses** | Realista mas longo — risco de esgotamento |
| Solo, 30 h/semana | **~14 meses** | ⭐ **Recomendado** |
| Solo + artista + músico contratados | **~11 meses** | Melhor relação custo/risco |
| Dupla em tempo integral | **~7 meses** | Improvável neste contexto |

**Recomendação: alvo de 12–14 meses até 1.0, com Early Access aos ~7 meses
(versão 0.6).** O EA financia parcialmente a arte e o áudio e traz dados de
balanceamento reais — que é justamente o que este projeto mais precisa e menos
tem.

**Orçamento externo estimado:** pixel art de cenários e UI R$ 8.000–15.000 ·
trilha (10 faixas) R$ 6.000–12.000 · efeitos sonoros R$ 2.000–4.000 · Steam
Direct US$ 100 · **total R$ 17.000–32.000.**

---

## 7. Roadmap 0.2 → 1.0

Cada versão tem **um objetivo dominante**. Uma versão que faz duas coisas não faz
nenhuma bem.

---

### 🔧 0.2 — "O jogo justo" · ~5 semanas

**Objetivo:** consertar o que está quebrado antes de construir sobre isso.
Nenhum conteúdo novo.

**Funcionalidades**
- Remover o stun automático da Lança de Luz
- Rebalancear as três armas para ~65 DPS em suas janelas ideais
- Redistribuir o bestiário em `stages.go` (corvo 63,5% → 35%)
- Corrigir a economia: `respawn` para `initialHealth`; remover o bônus por bomba
  não usada; reduzir o peso do `bossScore`
- Chefe: revisar vida e ritmo agora que ele não é mais stun-lockado
- **Ferramenta de balanceamento:** um modo *headless* que simula N runs com seed
  fixa e reporta DPS, tempo até nível máximo, abates e mortes

| Prioridade | Risco | Dependências |
| --- | --- | --- |
| 🔴 Crítica | 🟢 Baixo — mudanças pequenas, cobertura de testes alta | Nenhuma |

**Critério de saída:** um playtest cego escolhe armas diferentes em runs
diferentes; o chefe dura ≥ 60 s e executa todos os padrões.

---

### 💀 0.3 — "A Corrupção" · ~7 semanas

**Objetivo:** implantar o diferencial. É a versão que define se o jogo tem
identidade.

**Funcionalidades**
- Medidor de Corrupção com as 5 faixas
- Consequências: mutação de inimigos (3 variantes), degradação de paleta,
  multiplicador de pontos e drops
- Moldura lateral decorada (resolve o 3:4) exibindo o medidor
- **Vharak Ascendido**: chefe alternativo em Corrupção 100%
- Áudio reativo: a trilha distorce conforme a corrupção sobe

| Prioridade | Risco | Dependências |
| --- | --- | --- |
| 🔴 Crítica | 🟠 **Médio-alto** — é uma aposta de design não validada | 0.2 (balanceamento estável) |

**Mitigação do risco:** implementar primeiro em versão mínima (barra + multiplicador
de pontos + mais vida nos inimigos), testar com 5 jogadores externos, e só então
investir nas mutações e no chefe alternativo. **Se a Corrupção não empolgar em
teste, este é o momento de saber — não na 0.8.**

**Critério de saída:** jogadores mencionam a Corrupção espontaneamente ao
descrever o jogo; ao menos um testador tenta subi-la de propósito.

---

### ⚔️ 0.4 — "O verbo" · ~5 semanas

**Objetivo:** transformar o controle em algo que só este jogo tem.

**Funcionalidades**
- Mergulho do Grifo (dano por contato, i-frames, deslocamento de câmera)
- Fôlego do grifo (recurso do mergulho)
- Graze: raspar projéteis carrega bomba
- Camada de input por ações + **gamepad completo**
- Pacote pesado de game feel (itens 1–20 da seção 8)

| Prioridade | Risco | Dependências |
| --- | --- | --- |
| 🟠 Alta | 🟠 Médio — i-frames podem trivializar o desafio | 0.2 |

**Critério de saída:** o mergulho aparece no primeiro GIF que dá vontade de
postar; gamepad funciona sem teclado em toda a interface.

---

### 🎨 0.5 — "A cara do jogo" · ~7 semanas

**Objetivo:** matar definitivamente a aparência de protótipo. **Marco de
apresentação pública.**

**Funcionalidades**
- Cenários em pixel art de verdade (fim dos retângulos), 3 primeiros biomas
- Fonte própria do jogo (bitmap desenhada, não `basicfont`)
- Guia de direção artística aplicado (seção 9)
- HUD reconstruído nas molduras laterais
- Telas de menu, vitória e derrota com direção de arte
- 3 primeiras trilhas produzidas + 20 efeitos sonoros
- Onboarding diegético de 30 s na fase 1

| Prioridade | Risco | Dependências |
| --- | --- | --- |
| 🟠 Alta | 🔴 **Alto — dependência externa de artista e músico** | Contratação iniciada na 0.3 |

**Mitigação:** contratar na 0.3, com entregas parceladas por bioma; manter o
sistema de fallback procedural (já existe) para que atrasos de arte nunca
bloqueiem o desenvolvimento.

**Critério de saída:** um estranho que vê 10 s de vídeo consegue nomear o jogo.
Página da Steam publicada com lista de desejos aberta.

---

### 🎲 0.6 — "A run é sua" · ~6 semanas · 🚀 **Early Access**

**Objetivo:** dar variação entre partidas e abrir o jogo ao público.

**Funcionalidades**
- Altar entre trechos (escolha diegética de três)
- 16 relíquias iniciais
- Runas Fundidas (3 fusões)
- Pactos da Corrupção
- Fragmentos de Coroa (moeda, sem loja ainda)
- Build de Early Access + página + trailer

| Prioridade | Risco | Dependências |
| --- | --- | --- |
| 🟠 Alta | 🟠 Médio — EA prematuro queima reputação | 0.3, 0.4, 0.5 |

**Regra do EA:** só lançar se as 3 primeiras fases estiverem em qualidade final.
**EA com aparência de beta destrói mais valor do que qualquer atraso.**

**Critério de saída:** ≥ 3 runs por jogador na primeira sessão; feedback de EA
capturado e triado.

---

### 🐉 0.7 — "Os confrontos" · ~8 semanas

**Objetivo:** dar à campanha um meio e um clímax à altura.

**Funcionalidades**
- Generalização do `Boss` (interface + configuração em dados)
- Partes destrutíveis (asas, cauda, cabeça)
- Mini-chefe 1: A Catapulta Viva (fase 2)
- Mini-chefe 2: O Guardião Apodrecido (fase 4)
- Vharak reconstruído com partes destrutíveis
- Modo Boss Rush

| Prioridade | Risco | Dependências |
| --- | --- | --- |
| 🟠 Alta | 🟠 Médio-alto — é o bloco de código mais complexo | 0.2, feedback do EA |

---

### 🌍 0.8 — "O reino inteiro" · ~8 semanas

**Objetivo:** completar o conteúdo da campanha.

**Funcionalidades**
- Fases 4, 5 e 6 (total: 6)
- Biomas 4–6 em arte final
- 5 inimigos novos + 3 variantes corrompidas
- Clima dinâmico por bioma
- Cenário destrutível
- Eventos de trecho
- Trilhas 4–10

| Prioridade | Risco | Dependências |
| --- | --- | --- |
| 🟡 Média-alta | 🟠 Médio — balanceamento de 6 fases é trabalhoso | 0.5, 0.7 |

---

### 🏰 0.9 — "Razões para voltar" · ~6 semanas

**Objetivo:** fechar o loop de retenção.

**Funcionalidades**
- A Cidadela (loja de meta-progressão)
- 24 relíquias completas
- Segundo cavaleiro jogável
- 8 modificadores de run
- Diário de Criaturas
- 20 conquistas Steam
- Quarta dificuldade (Cavaleiro) + ranking

| Prioridade | Risco | Dependências |
| --- | --- | --- |
| 🟡 Média-alta | 🟢 Baixo — sistemas isolados e bem compreendidos | 0.6, 0.8 |

---

### 👑 1.0 — "Asas de Valdoria" · ~7 semanas

**Objetivo:** lançar um jogo terminado.

**Funcionalidades**
- Localização inglês + português
- Steam Deck Verified
- Conquistas, nuvem, cartas
- Polimento final de game feel (lista completa da seção 8)
- Balanceamento final com dados de EA
- QA em todas as plataformas
- Trailer de lançamento, capturas, GIFs, press kit
- Créditos, opções de acessibilidade (daltonismo, redução de flash, mira
  assistida)

| Prioridade | Risco | Dependências |
| --- | --- | --- |
| 🔴 Crítica | 🟡 Médio — risco é de escopo, não técnico | Tudo |

**Critério de saída:** as metas mensuráveis da seção 2.10.

---

### Mapa de dependências

```
0.2 Balanceamento ──┬── 0.3 Corrupção ──┬── 0.6 Altar/EA ──┬── 0.9 Meta
                    │                   │                  │
                    ├── 0.4 Verbos ─────┤                  │
                    │                   │                  │
                    └── 0.7 Chefes ─────┴── 0.8 Conteúdo ───┴── 1.0
                                        │
        Arte/Áudio (externo) ── 0.5 ────┘
```

**Caminho crítico:** 0.2 → 0.3 → 0.6 → 0.9 → 1.0.
**Maior risco de cronograma:** 0.5 (dependência externa). **Contratar arte e
áudio durante a 0.3, não durante a 0.5.**

---

## 8. Game feel: 60 melhorias

Nada aqui é sistema novo. São ajustes de percepção — a diferença entre "funciona"
e "é gostoso". Ordenados por relação impacto/esforço.

### Impacto e resposta (1–12)

1. **Hit-stop em todo abate**, escalado por tamanho: 1 frame para corvo, 2 para
   harpia, 4 para wyvern. Hoje só o wyvern tem ([`game.go:960`](game/game.go:960)).
2. **Hit-stop no dano ao jogador** — 6 frames. O momento mais importante do jogo é
   o único sem congelamento.
3. **Recuo do disparo:** o grifo desloca 1px para trás a cada tiro, com retorno
   elástico. Vende peso sem custar precisão.
4. **Impacto de projétil inimigo:** hoje o projétil simplesmente some
   ([`game.go:1017`](game/game.go:1017)). Precisa de faísca, som e clarão.
5. **Impacto de projétil na parede/borda** — pequena poeira ao sair da tela pelos
   lados.
6. **Escalonar o shake por evento:** hoje há três magnitudes. Precisa de sete —
   abate pequeno, abate grande, dano, escudo, bomba, troca de fase do chefe,
   morte do chefe.
7. **Shake direcional:** dano vindo da direita empurra a câmera para a esquerda.
   Muito mais informativo que shake aleatório.
8. **Squash & stretch no grifo:** achatar ao mudar de direção, esticar ao
   mergulhar. Duas linhas de `GeoM.Scale`.
9. **Squash no inimigo ao ser atingido:** 2 frames de deformação antes do flash.
10. **Anticipação na morte do inimigo:** 3 frames congelado e branco antes de
    explodir — dá tempo do olho registrar o abate.
11. **Flash do inimigo com curva:** hoje é branco liso por 4 frames. Deveria ser
    branco → cor do elemento → normal.
12. **Vibração de gamepad** por evento, com intensidade configurável (a preferência
    de "Vibração" já existe e hoje só afeta a tela).

### Projéteis e armas (13–24)

13. **Clarão de bocal por arma:** hoje é um retângulo branco genérico
    ([`player.go:207`](game/player.go:207)). Luz = estrela; Chamas = jato; Gelo =
    cristal.
14. **Rastro de projétil por elemento:** a Luz deveria deixar um risco fino
    persistente, não partículas soltas.
15. **Aceleração do projétil:** sair 20% mais devagar e acelerar em 5 frames.
    Truque clássico de leitura e de sensação.
16. **Estilhaço de gelo:** projéteis de gelo quebram em fragmentos ao acertar.
17. **Brasas persistentes:** as Chamas deixam brasas que caem e apagam.
18. **Arco elétrico desenhado:** a cadeia da Luz hoje só cria um rastro
    ([`game.go:986`](game/game.go:986)). Precisa de um raio visível ligando os
    dois pontos.
19. **Grossura do projétil por nível de arma:** nível 3 é visivelmente mais
    pesado.
20. **Som de disparo com variação de pitch** de ±5% — remove a fadiga auditiva do
    tiro contínuo.
21. **Som em camadas para o tiro:** ataque + corpo + cauda, em vez de uma nota.
22. **Aura de sobreposição de projéteis:** quando ≥ 4 projéteis do jogador se
    sobrepõem, um brilho concentrado — comunica "estou no ponto certo".
23. **Cadência com micro-swing:** variar o cooldown em ±1 frame remove a sensação
    metronômica.
24. **Fim de arma sem munição:** as fusões devem ter aviso visual e sonoro ao se
    esgotarem (se tiverem duração).

### Inimigos e leitura (25–34)

25. **Telegraph em duas etapas:** aviso suave (16f) → aviso forte (6f). Hoje é
    binário ([`enemy.go:99`](game/enemy.go:99)).
26. **Linha de mira fantasma** para tiros mirados (wyvern, balista) — 8 frames de
    linha tênue antes do disparo.
27. **Entrada com sombra:** inimigos projetam uma sombra crescente no topo antes
    de aparecer.
28. **Reação de morte por elemento:** queimado vira cinzas; congelado estilhaça;
    fulminado desintegra em faíscas.
29. **Grito de morte por tipo** — quatro vozes distintas em vez de um som único.
30. **Balanço de dor:** o inimigo recua 1px na direção do impacto.
31. **Personalidade ociosa:** a gárgula "encara" o jogador antes de atacar; o mago
    gesticula.
32. **Contorno de perigo:** inimigos que vão colidir com o jogador em < 0,5 s
    ganham um contorno vermelho pulsante.
33. **Aviso de fuga iminente:** inimigo prestes a escapar pisca em magenta —
    essencial quando a Corrupção existir.
34. **Diferenciar altura de voo:** inimigos "altos" um pouco menores e mais
    claros; "baixos" maiores e mais escuros.

### Câmera e tela (35–42)

35. **Deslocamento suave da câmera** seguindo o jogador em ~8% do offset — dá
    amplitude à tela pequena.
36. **Zoom-punch na entrada do chefe** — 1,05× por 20 frames.
37. **Zoom-out na bomba** — o oposto, para caber o espetáculo.
38. **Aceleração do parallax** em momentos de tensão (velocidade das camadas +30%
    na aproximação do chefe).
39. **Vinheta reativa:** escurecimento nas bordas com pouca vida.
40. **Aberração cromática** brevíssima ao receber dano (1–2px de deslocamento de
    canal).
41. **Barra de vida do chefe com dano atrasado:** camada branca que drena
    lentamente atrás do vermelho — comunica quanto dano acabou de sair.
42. **Congelar o parallax** durante o hit-stop (hoje ele continua rolando).

### Interface e recompensa (43–52)

43. **Popups de pontuação na fonte do jogo** — ainda usam `ebitenutil.DebugPrintAt`
    ([`popup.go:43`](game/popup.go:43)).
44. **Popup com escala e queda elástica** em vez de subida linear.
45. **Cor do popup por faixa de valor** — branco, dourado, magenta.
46. **Multiplicador com pulso** a cada incremento e **tremor** quando prestes a
    expirar.
47. **Contador de pontuação animado** (rolagem até o valor) na tela de vitória.
48. **Imã de coleta:** power-ups atraídos suavemente a partir de 40px. Nunca perder
    um item por 2 pixels.
49. **Trilha do power-up:** rastro colorido enquanto cai, para o olho seguir.
50. **Feedback de nível máximo:** coletar uma runa já no nível 3 precisa comunicar
    "convertido em pontos" ou "pronto para fundir" — hoje não acontece nada.
51. **Marca de recorde na barra de progresso:** onde você morreu da última vez.
52. **Fade de HUD em momentos limpos:** o HUD esmaece 30% durante a entrada do
    chefe.

### Áudio e ritmo (53–60)

53. **Ducking de música** ao usar a bomba ou na morte do chefe — a trilha some por
    0,5 s.
54. **Filtro passa-baixa na pausa** em vez de manter a música intacta.
55. **Camadas de trilha por intensidade:** percussão que entra conforme a
    Corrupção sobe.
56. **Batida do coração** abaixo de 2 pontos de vida.
57. **Silêncio antes do chefe:** 1,5 s de silêncio total antes da trilha do chefe.
    O recurso mais barato e mais eficaz de tensão.
58. **Estéreo posicional** para inimigos e projéteis conforme a posição em X.
59. **Voz do grifo:** grito ao mergulhar, guincho ao levar dano, ronronar ao
    coletar. **Este é o item que mais transforma a montaria em personagem.**
60. **Reverb por bioma:** o covil ecoa; os campos são secos.

**Priorização se houver tempo para apenas dez:** 1, 2, 4, 13, 25, 35, 43, 48, 57,
59.

---

## 9. Direção artística

### 9.1 Diagnóstico do estilo atual

**O que existe:** sprites em pixel art procedural gerados de grades ASCII com
paletas nomeadas ([`spritedata.go`](game/spritedata.go)) — o do jogador é 16×16 com
8 cores e é legível; cenário em parallax de quatro camadas feito exclusivamente de
`vector.DrawFilledRect`; molduras de UI em ferro/pergaminho; fonte
`basicfont.Face7x13`.

**Diagnóstico honesto:** o jogo tem **duas linguagens visuais que não conversam**.
As entidades são pixel art 16-bit; o cenário é arte vetorial abstrata. Não é um
estilo — é a ausência de uma decisão. O jogador percebe isso mesmo sem saber
nomear, e é o principal motivo de o jogo parecer inacabado.

**O que preservar:** as decisões de *legibilidade* foram excelentes e são a base
de todo o resto — escurecimento do cenário, contorno escuro nas entidades,
separação cromática dos projéteis, projéteis desenhados por último. **Nenhuma
decisão de estilo pode violar essas quatro regras.**

### 9.2 Identidade visual: "Vitral Corrompido"

**Pixel art 16-bit de arcade tardio, com estrutura de vitral: contornos escuros
firmes, campos de cor saturada, e luz que atravessa em vez de iluminar.**

Três referências deliberadas, cada uma por um motivo:

| Referência | O que tomamos | O que não tomamos |
| --- | --- | --- |
| **Jamestown** | A densidade de detalhe do cenário e a legibilidade sob caos | O barroco espacial |
| **Hollow Knight** | O uso de silhueta e a luz como narrativa | O estilo desenhado à mão |
| **Vitral gótico medieval** | Contorno de chumbo, campos de cor, luz atravessando | O realismo religioso |

**Regra de ouro:** *tudo que mata tem contorno preto e cor saturada; tudo que é
cenário é dessaturado e sem contorno.* Uma regra, aplicada sem exceção, resolve
90% dos problemas de leitura de um shmup.

### 9.3 Paleta mestra

Paleta base de 32 cores, dividida em três famílias que nunca se misturam:

**Família Nobre (jogador, UI, recompensa) — dourados e azuis frios**
```
#FFF3C4  #F0BC4C  #C08A2E  #7A5418   ← ouro (grifo, destaque, recompensa)
#B8E4FF  #46C8FF  #1E74C8  #123E6E   ← aço (armadura, UI fria)
#FFFFFF  #F0E8CF  #B8AC90  #4A4438   ← pergaminho (texto, moldura)
```

**Família Corrompida (inimigos, projéteis inimigos, corrupção) — magentas e
violetas**
```
#FFDAFF  #C828DC  #7A3AB0  #3A1A5A   ← corrupção (inimigos, tiros)
#FF5088  #B02850  #5A1428  #2A0A14   ← carne doente (variantes)
```

**Família Mundo (cenário) — dessaturada, nunca compete**
```
#8A9AA8  #5A6470  #38404A  #1E2228   ← pedra
#7A8A5A  #4A5838  #2A3220  #141A10   ← vegetação
#C08050  #8A5A34  #4A3020  #1E1408   ← terra/madeira
```

**Regra de saturação:** cenário nunca acima de 40% de saturação; entidades nunca
abaixo de 60%. Isso é mensurável e verificável automaticamente — vale um teste.

### 9.4 Iluminação

- **Fonte de luz coerente por bioma:** amanhecer nos Campos (dourada, lateral
  baixa); incêndio na Vila (laranja pulsante, de baixo); relâmpago nas Muralhas
  (branco-azulada, estroboscópica); cristais no Covil (magenta, de baixo).
- **Luz atravessa, não ilumina:** feixes volumétricos entre as camadas de parallax
  (raios de sol pela copa, luz pelas frestas da muralha).
- **Entidades emitem, cenário recebe:** projéteis, runas e cristais têm glow;
  cenário nunca.
- **A Corrupção esvazia a cor:** conforme sobe, o cenário perde saturação em até
  50% enquanto as entidades corrompidas *ganham* — leitura visual instantânea do
  estado da run, sem HUD.

### 9.5 Efeitos

| Efeito | Direção |
| --- | --- |
| Explosão de inimigo | Anel branco (2f) → partículas na cor do inimigo → cinzas que caem |
| Morte elemental | Queimado = cinzas · Congelado = estilhaços de vitral · Fulminado = faíscas |
| Coleta de runa | Anel de luz + fragmentos que sobem + linha até o HUD |
| Mergulho do grifo | Rastro dourado com penas soltas + distorção do ar |
| Bomba (Invocação) | Silhueta de dragão em contorno dourado sobre a tela escurecida |
| Corrupção subindo | Veias magenta crescendo das bordas para dentro |
| Entrada do chefe | Escurecimento total → dois olhos → revelação |
| Partes destrutíveis | Fissuras que se acumulam antes de a parte se soltar |

### 9.6 UI, HUD e menus

**A moldura lateral é a decisão estruturante.** Janela em 16:9; campo de jogo
240×320 centralizado; laterais ocupadas por molduras de pedra esculpida com
heráldica.

```
┌──────────────┬──────────────────┬──────────────┐
│  MOLDURA E   │                  │  MOLDURA D   │
│              │                  │              │
│  ⚜ VALDORIA  │                  │  RUNA        │
│              │   CAMPO DE JOGO  │  [ícone]     │
│  PONTOS      │     240×320      │  Nv 3        │
│  1 2 4 5 0   │                  │              │
│              │  (limpo: apenas  │  BOMBAS      │
│  ♥♥♥○○       │   o essencial)   │  ◆ ◆         │
│              │                  │              │
│  CORRUPÇÃO   │                  │  RELÍQUIAS   │
│  ▓▓▓▓░░░░    │                  │  ◇ ◇ ◇       │
│              │                  │              │
└──────────────┴──────────────────┴──────────────┘
```

Ganhos: o campo de jogo fica **limpo** (crítico num shmup); a tela deixa de
desperdiçar 70%; a moldura vira identidade visual reconhecível em capturas; e o
Medidor de Corrupção ganha o destaque permanente que a mecânica principal merece.

**Fonte:** encomendar uma bitmap própria de 8×10 e 12×16 (versão para títulos),
com serifas leves e acentuação completa em português. `basicfont` é um marcador de
protótipo e precisa desaparecer.

**Menus:** o menu principal é a **sala do trono de Valdoria** — não uma lista sobre
estrelas. O jogador vê o grifo pousado, o trono vazio, e as opções gravadas em
placas de bronze. A Cidadela (loja) é o mesmo espaço, com mais detalhes conforme
o progresso.

### 9.7 Chefes

Diretriz: **cada chefe precisa ser reconhecível em silhueta, num frame, em preto e
branco.** Se não for, foi mal desenhado.

- Ocupar 40–60% da largura da tela.
- Partes destrutíveis visualmente distintas **antes** de o jogador saber que são
  destrutíveis.
- Degradação visível: perder uma asa muda a pose, não só remove um retângulo.
- Cada fase muda a **cor de emissão**, não só o padrão.
- Vharak Ascendido: mesma silhueta, coberto de cristais de corrupção, luz interna
  magenta atravessando as juntas.

### 9.8 Cenários

Cada bioma precisa de: uma **silhueta característica** (o formato que aparece no
parallax), uma **cor dominante**, um **elemento animado** (fogo, chuva, folhas) e
um **detalhe narrativo** (algo que conta o que aconteceu ali).

| Bioma | Silhueta | Cor | Animação | Narrativa |
| --- | --- | --- | --- | --- |
| Campos | Trigais, cercas, moinho | Dourado-verde | Trigo ondulando | Carroças abandonadas |
| Vila | Telhados, campanário | Laranja-preto | Fogo, fagulhas | Vilarejos correndo |
| Muralhas | Ameias, torres | Cinza-azul | Chuva, bandeiras | Balistas quebradas |
| Floresta | Copas, raízes | Verde doente | Névoa, esporos | Árvores com veias magenta |
| Desfiladeiro | Colunas de rocha | Laranja-marrom | Cinzas subindo | Ossadas de grifos |
| Covil | Cristais, ossos | Magenta-preto | Pulsação dos cristais | Tesouro do reino saqueado |

---

## 10. Direção de áudio

### 10.1 Diagnóstico

O sistema de áudio é **bem construído**: gerenciador único, fade entre trilhas,
frame-gate contra repetição, volumes separados, fallback procedural, desativação
graciosa sem dispositivo. A arquitetura está pronta para receber áudio produzido —
basta colocar arquivos em `assets/audio/`.

O **conteúdo** é o problema, e é grave: onze trilhas, cada uma um laço de 8 notas
de onda quadrada. `genMusicFields` dura **2,24 segundos** e repete infinitamente
([`audio.go:575`](game/audio.go:575)). Em cinco minutos de campanha, o jogador ouve
o mesmo arpejo **130 vezes**.

**Este é o eixo com maior retorno por real investido em todo o projeto.** Trilha
produzida transforma a percepção de qualidade mais rápido que qualquer outra
coisa — e o sistema já está pronto para recebê-la.

### 10.2 Identidade sonora: "Heroísmo Sujo"

**Instrumentação medieval-orquestral tocada com a sujeira e a energia de um chip
de arcade.** Não é orquestra épica limpa (previsível, cara); não é chiptune puro
(nega a fantasia medieval). É a colisão dos dois.

Regra: **se um trecho da trilha pudesse tocar num filme de fantasia sem estranhar,
está limpo demais.**

### 10.3 Instrumentação

**Núcleo (presente em tudo):**
- **Tambores de guerra** — o pulso do jogo. Sempre.
- **Baixo sintetizado com saturação** — a ponte entre o medieval e o arcade.
- **Coro masculino em vogais**, gravado ou sampleado — a nobreza.

**Família Nobre (jogador, vitória, altar):**
- Alaúde e harpa
- Trompas naturais (sem válvula — desafinadas de propósito)
- Sinos de bronze

**Família Corrompida (inimigos, corrupção, chefes):**
- Violino em *sul ponticello* (arco sobre o cavalete — som vidrado e doente)
- Coro processado com formantes deslocados
- Metais graves com distorção
- Cristal ressonante (som de vidro sendo tensionado)

**Família Arcade (a sujeira):**
- Onda quadrada com PWM para arpejos rápidos
- Ruído para percussão
- Bend de pitch agressivo nos acentos

### 10.4 Trilhas (10 faixas)

| # | Faixa | Contexto | Duração | Caráter |
| --- | --- | --- | ---: | --- |
| 1 | **O Trono Vazio** | Menu | 2:30 | Alaúde solitário, tambor distante. Melancólico e nobre. |
| 2 | **Cavalgada** | Fase 1 | 3:00 | Heroico, andamento alto, trompas. A faixa que define o jogo. |
| 3 | **Cinzas** | Fase 2 | 3:00 | Mesma melodia da faixa 2, em modo menor e mais lenta. Narrativa por música. |
| 4 | **Ferro e Chuva** | Fase 3 | 3:00 | Percussão pesada, marcial, coro em uníssono. |
| 5 | **O Bosque que Respira** | Fase 4 | 3:00 | Ambiente inquieto, violino doente, ritmo irregular. |
| 6 | **Cânion de Ossos** | Fase 5 | 3:00 | Tribal, tambores graves, vento. |
| 7 | **Sob as Asas** | Fase 6 | 3:00 | Tensão contida antes do covil. |
| 8 | **Vharak** | Chefe | 3:30 | Coro corrompido + metais + arcade agressivo. |
| 9 | **A Queda de Valdoria** | Vharak Ascendido | 4:00 | A faixa 2 destruída: mesma melodia, irreconhecível. |
| 10 | **Altar** | Altar / Cidadela | 1:30 | Suspenso, sinos, quase silêncio. O respiro. |

**Faixas 2, 3 e 9 compartilham o mesmo tema melódico.** Um tema que degrada ao
longo da campanha e retorna arruinado no clímax — a forma mais barata e mais
eficaz de contar uma história sem uma linha de texto.

**Camadas por Corrupção:** cada faixa de fase é entregue em três camadas (base,
tensão, colapso) que entram e saem conforme o medidor sobe. Isso é uma exigência
de produção a ser passada ao compositor desde o briefing — reformar depois custa a
trilha inteira.

### 10.5 Efeitos sonoros (~50)

**Princípio: cada som em três camadas — ataque (o que chama a atenção), corpo (o
que dá peso), cauda (o que dá espaço).** É o que separa som de arcade de som de
protótipo.

| Categoria | Sons | Direção |
| --- | ---: | --- |
| **Voz do grifo** ⭐ | 6 | Grito ao mergulhar, guincho de dor, ronronar ao coletar, respiração ofegante com fôlego baixo, rugido na bomba, lamento na morte. **A prioridade número 1 dos efeitos** — é o que transforma a montaria em personagem. |
| Armas | 9 | 3 por arma (níveis 1/2/3), cada uma com timbre próprio: Luz = cristalino com ataque agudo; Chamas = ruído com corpo grave; Gelo = vidro tensionado |
| Fusões | 3 | Sons híbridos reconhecíveis como a soma dos dois |
| Impactos | 8 | Acerto em carne, pedra, metal, cristal; escudo; crítico; ponto fraco; bloqueio |
| Mortes | 6 | Por tipo de inimigo e por elemento |
| Inimigos | 8 | Telegraph (2 etapas), disparo por tipo, entrada de gárgula, canto do mago |
| Chefe | 7 | Rugido de entrada, aviso, cada padrão, quebra de parte, mudança de fase, morte |
| Interface | 6 | Navegação, confirmação, cancelamento, compra, desbloqueio, conquista |
| Recompensa | 5 | Runa, fusão, fragmento, relíquia, purificação |
| Corrupção | 4 | Inimigo escapando, faixa aumentada, mutação, colapso |

### 10.6 Ambientação

Camada de ambiente independente da trilha, por bioma: vento nos campos · fogo e
madeira estalando na vila · chuva e trovão nas muralhas · insetos e ranger de
árvores na floresta · vento em fenda no desfiladeiro · gotejamento e ressonância
de cristal no covil.

**Reverb por bioma** aplicado a todos os efeitos: seco nos campos, médio na vila,
longo no covil. É o detalhe que faz o jogador sentir que os lugares são
diferentes, mesmo sem conseguir dizer por quê.

**Silêncio como ferramenta:** 1,5 s de silêncio absoluto antes da trilha do chefe;
ducking total na bomba; corte seco na morte de Vharak. **Silêncio é o efeito mais
barato e mais poderoso do arsenal.**

---

## 11. Monetização

### 11.1 Modelo recomendado

**Venda única, premium, sem microtransações, com Early Access.**

Não há debate real aqui: shmups de nicho não sustentam free-to-play, e o público
alvo (seções 2.3) reage negativamente a qualquer monetização recorrente. O modelo
correto é o óbvio, bem executado.

### 11.2 Preço

| Fase | Preço (BRL) | Preço (USD) | Justificativa |
| --- | ---: | ---: | --- |
| **Early Access (0.6)** | R$ 24,90 | US$ 9,99 | Sinaliza "jogo de verdade" sem cobrar preço de produto terminado. Abaixo de US$ 8 sinaliza amadorismo; acima de US$ 12 em EA gera reclamação. |
| **Lançamento 1.0** | R$ 34,90 | US$ 14,99 | Alinhado ao mercado: *Sky Force Reloaded* US$ 9,99 · *Jamestown+* US$ 14,99 · *ZeroRanger* US$ 9,99 · *Devil Blade Reboot* US$ 12,99. Com 6 fases, 4 chefes, meta-progressão e 15–25 h, US$ 14,99 é defensável. |
| **Compradores de EA** | — | — | Recebem a 1.0 gratuitamente. Regra da Steam e boa prática. |

**Regionalização é obrigatória.** O preço brasileiro deve ser R$ 34,90, não a
conversão de US$ 14,99 (~R$ 80). O jogo é brasileiro, em português, e o público
nacional é um ativo estratégico.

**Descontos:** 15% na semana de lançamento; nunca mais de 50% no primeiro ano;
promoções apenas nos eventos sazonais da Steam. **Descontar cedo e fundo mata o
valor percebido de um indie de nicho.**

### 11.3 Demo

**Sim — obrigatoriamente, e desde a 0.5.**

Conteúdo: fase 1 completa + mini-chefe da fase 2. Cerca de 8 minutos, terminando
no momento de maior empolgação. Modo Sobrevivência liberado (retenção alta, custo
zero). Progresso da demo **transfere** para o jogo completo.

Por que importa tanto: shmup é um gênero que se vende jogando, não assistindo. E
o Steam Next Fest com demo é, isoladamente, o maior gerador de listas de desejos
disponível para um indie sem verba de marketing.

### 11.4 Early Access

**Sim, na 0.6, com condições rígidas.**

| Regra | Detalhe |
| --- | --- |
| Só entra em EA com **3 fases em qualidade final** | EA com cara de beta destrói reputação permanentemente |
| Duração declarada: **8–10 meses** | Prazos vagos afastam compradores |
| Atualizações a cada 6–8 semanas | Cadência previsível é o que sustenta um EA |
| Roadmap público na página | Transparência é a moeda do EA |
| **Escopo congelado ao entrar em EA** | Feedback de EA ajusta balanceamento, não adiciona sistemas |

Benefícios reais: dados de balanceamento que o projeto não tem como obter de outra
forma; receita que financia arte e áudio; comunidade formada antes do lançamento.

### 11.5 DLC e conteúdo futuro

**Nada de DLC antes de seis meses após a 1.0.** Antes disso, todo conteúdo é
atualização gratuita — é o que sustenta avaliações positivas.

| Momento | Conteúdo | Modelo |
| --- | --- | --- |
| 1.0 + 2 meses | Modo Caos (modificadores extremos), 8 relíquias, ranking | **Gratuito** |
| 1.0 + 4 meses | Bioma extra + chefe secreto | **Gratuito** |
| 1.0 + 8 meses | **DLC "A Ordem dos Cavaleiros"**: 3 cavaleiros jogáveis, 3 montarias, 12 relíquias, campanha alternativa | US$ 7,99 |
| 1.0 + 12 meses | **DLC "As Terras Além"**: 4 fases, 2 chefes, novo elemento | US$ 9,99 |

Racional: o primeiro DLC vende **personagens** (multiplica o valor do conteúdo
existente, custo relativamente baixo); o segundo vende **conteúdo** (custo alto,
mas o público já está engajado).

### 11.6 Projeção realista

Com base em desempenho típico de shmups indie na Steam com boa recepção crítica e
marketing modesto:

| Cenário | Unidades (ano 1) | Receita bruta | Líquida (pós-Steam/impostos) |
| --- | ---: | ---: | ---: |
| Pessimista | 1.500 | US$ 15.000 | ~US$ 8.000 |
| **Realista** | **5.000** | **US$ 50.000** | **~US$ 27.000** |
| Otimista | 20.000 | US$ 200.000 | ~US$ 110.000 |

**Isto não é um projeto para pagar salários.** É um projeto que, bem executado,
paga seus custos externos, deixa lucro modesto e — mais importante — constrói
reputação e um portfólio que viabilizam o projeto seguinte. Essa expectativa
precisa estar clara desde já, porque ela determina quanto se pode investir.

---

## 12. Marketing

### 12.1 Elevator pitch

**Versão de 10 segundos (a que importa):**

> *"É um shoot'em up onde você monta um grifo — e cada monstro que escapa corrompe
> o reino. Corrupção alta significa mais pontos, melhores itens e um chefe final
> completamente diferente. Você decide o quanto deixa o mundo apodrecer."*

**Versão de 5 segundos (para trailer e redes):**

> *"Defenda o reino. Ou deixe-o apodrecer — e fique mais forte."*

**Versão de uma linha (Steam, tag principal):**

> *"Shoot'em up vertical sobre um cavaleiro alado e um reino que se corrompe a cada
> erro seu."*

### 12.2 Descrição da Steam

---

# ASAS DE VALDORIA

### O reino está caindo. Cada criatura que passa por você o corrompe mais.

Você é o último cavaleiro alado de Valdoria. Montado em seu grifo, você é tudo o
que separa o reino das criaturas corrompidas que descem do céu.

Mas você não consegue detê-las todas.

**E talvez não devesse.**

---

### ⚔️ A CORRUPÇÃO É SUA ARMA E SUA CONDENAÇÃO

Cada inimigo que atravessa suas linhas enche o **Medidor de Corrupção**. Conforme
ele sobe, o reino apodrece: os inimigos se mutam, os biomas escurecem, a própria
trilha sonora se deforma.

Mas a corrupção também **paga**. Mais pontos. Melhores relíquias. Pactos sombrios
que nenhum cavaleiro honrado aceitaria.

E se você deixar o reino cair por completo, encontrará algo no covil que ninguém
deveria ver.

### 🔥 TRÊS MAGIAS. SEIS FORMAS DE DESTRUIR.

Lança de Luz. Chamas do Dragão. Lanças de Gelo. Domine uma até o limite e **funda-a**
com outra: Julgamento Solar, Tempestade de Vapor, Prisma. Nenhuma run precisa ser
igual à anterior.

### 🦅 VOCÊ NÃO PILOTA UMA NAVE. VOCÊ MONTA UM ANIMAL.

Seu grifo **mergulha** através de paredes de projéteis, bate as asas para desviar
tiros, e grita quando você o machuca. Ele cansa. Ele reage. Ele cai com você.

### 🐉 CHEFES QUE VOCÊ DESMONTA PEDAÇO POR PEDAÇO

Arranque as asas de Vharak e ele voará torto. Destrua sua cauda e ele não invocará
mais. Cada parte que você derruba muda a luta.

### ⛩️ ESCOLHAS QUE FAZEM DA RUN A SUA RUN

Entre os trechos, um altar surge no ar com três bênçãos. Você voa até uma delas.
24 relíquias, pactos amaldiçoados e 4 dificuldades garantem que nenhuma partida
se pareça com a anterior.

---

**• 6 fases · 6 biomas · 4 chefes · 14 inimigos**
**• Partidas de 15 minutos · 15–25 horas para dominar**
**• Meta-progressão que amplia opções, nunca simplifica o desafio**
**• Modos Campanha, Sobrevivência, Boss Rush e Modificadores**
**• Suporte completo a gamepad · Verificado para Steam Deck**
**• Legendas em português e inglês**

---

### 12.3 Os cinco diferenciais (ordem de comunicação)

1. **Medidor de Corrupção** — o único mecanismo genuinamente novo. Sempre primeiro.
2. **Grifo como personagem** — a fantasia que nenhum concorrente tem.
3. **Runas Fundidas** — a promessa de variedade.
4. **Chefes desmontáveis** — o espetáculo.
5. **Altares e relíquias** — o gancho do público de roguelite.

### 12.4 Capturas de tela (6 obrigatórias, nesta ordem)

| # | Conteúdo | Objetivo |
| --- | --- | --- |
| 1 | Combate denso na Vila em Chamas, moldura lateral visível, Corrupção em ~40% | Vender o jogo inteiro num frame |
| 2 | Vharak com uma asa destruída, cristais expostos, projéteis na tela | Espetáculo e escala |
| 3 | Altar com três bênçãos, o grifo voando até uma delas | Sinalizar roguelite |
| 4 | Medidor de Corrupção em 85%, cenário deformado, inimigos mutados | O diferencial, visualmente |
| 5 | Mergulho do grifo atravessando uma parede de projéteis | O verbo de assinatura |
| 6 | A Cidadela com relíquias desbloqueadas | Profundidade e retenção |

**Regra:** nenhuma captura sem ação. Nenhuma captura de menu vazio. Toda captura
com o HUD visível — ele é parte da identidade.

### 12.5 Trailer (75 segundos)

| Tempo | Conteúdo | Áudio |
| --- | --- | --- |
| 0:00–0:05 | Tela preta. Texto: *"Valdoria está caindo."* | Silêncio, depois um tambor |
| 0:05–0:15 | O grifo alçando voo sobre os campos. Primeiro combate. | Tema de "Cavalgada" entrando |
| 0:15–0:25 | Combate acelerando. Três armas em corte rápido. Fusão de runas. | Trilha subindo |
| 0:25–0:35 | **Um inimigo escapa.** Corte para o medidor subindo. O mundo escurece. | Trilha distorcendo |
| 0:35–0:45 | Inimigos mutando. Bioma corrompido. Pacto sendo aceito. | Camada de colapso |
| 0:45–0:55 | Mergulho do grifo em câmera lenta atravessando os projéteis. | **Silêncio total** |
| 0:55–1:05 | Vharak. Asa se soltando. Explosões. | Tema do chefe, máximo |
| 1:05–1:12 | Corte seco para preto. Logo. | Silêncio |
| 1:12–1:15 | *"Defenda o reino. Ou deixe-o apodrecer."* + Wishlist | Único acorde de sino |

**Erro a evitar:** trailer que mostra "features". Trailers de shmup vendem
**movimento e sensação**. Zero texto explicativo além das duas frases.

### 12.6 GIFs (a ferramenta de divulgação mais importante)

Cada GIF de 3–6 segundos, em loop perfeito, legível em miniatura:

1. **Mergulho através de uma parede de projéteis** — o melhor GIF possível
2. **Fusão de runas** com a explosão de cor
3. **Asa de Vharak se soltando** com hit-stop
4. **Medidor de Corrupção estourando** e o mundo mudando de cor
5. **Bomba (Invocação Ancestral)** com o dragão atravessando
6. **Combo alto** com popups subindo e multiplicador pulsando
7. **O grifo gritando ao levar dano** (com aberração cromática)
8. **Altar** com as três luzes e a escolha

### 12.7 Cronograma de divulgação

| Fase | Momento | Ação |
| --- | --- | --- |
| Silêncio | Até 0.4 | Nada público. Não se divulga um protótipo. |
| Semeadura | 0.5 | Página da Steam no ar. Devlog em vídeo. GIFs semanais no Bluesky/X/Reddit (r/shmups, r/indiegames). |
| Aceleração | 0.6 / EA | Demo publicada. **Steam Next Fest.** Contato com streamers de shmup e de roguelite. |
| Sustentação | 0.7–0.9 | Atualizações a cada 6–8 semanas, cada uma com GIF e nota de versão em vídeo. |
| Lançamento | 1.0 | Trailer, press kit, chaves para imprensa 2 semanas antes, participação em festival sazonal. |

**Meta de listas de desejos antes da 1.0: 10.000.** Abaixo de 7.000, adiar o
lançamento — o algoritmo da Steam pune lançamentos com pouca antecipação, e é
irrecuperável.

---

## 13. Benchmark

Não copiar. Extrair o princípio e aplicá-lo ao que este jogo é.

### Super Aleste (1992) — *a referência declarada*

**O que faz melhor que nós:** a **progressão de arma dentro da partida**. Oito
armas, cada uma com identidade absoluta e níveis que mudam o comportamento, não só
os números. O jogador experimenta e adota uma. Também: densidade de inimigos
sempre crescente, sem os platôs que temos.

**Aprendizado aplicável:** níveis de arma devem **mudar comportamento**, não
incrementar contadores. Hoje o nível 3 da Luz é "quatro projéteis em vez de dois".
Deveria ser uma arma diferente.

**O que não copiar:** a estrutura linear de dez fases sem escolha. É de 1992 e o
público mudou.

### Jamestown (2011) — *a referência de produção*

**O que faz melhor:** **coesão de identidade**. Marte colonial no século XVII é
uma premissa absurda executada com convicção total — arte, música, texto,
interface e nomes contam a mesma história. Também tem a melhor **densidade
legível** do gênero moderno: caos absoluto, leitura perfeita.

**Aprendizado aplicável:** este é o modelo a perseguir. Um projeto pequeno vence
por **convicção estética**, não por volume. Também: o sistema Vaunt (acumular e
gastar num momento de invulnerabilidade com pontos dobrados) é a melhor mecânica
de risco/recompensa do gênero — e é a inspiração direta do graze e do Mergulho.

**O que não copiar:** o foco em cooperativo local, que hoje divide esforço sem
retorno.

### Sky Force Reloaded (2016) — *a referência de retenção*

**O que faz melhor:** o **loop de retenção**, sem concorrência. Cada run rende
estrelas, cada estrela aproxima um upgrade, cada fase tem objetivos secundários
que exigem revisitar. É desenhado para "mais uma".

**Aprendizado aplicável:** cada run precisa terminar com **ganho tangível**, mesmo
em derrota (Fragmentos de Coroa). Também: objetivos secundários por fase são
retenção baratíssima.

**O que não copiar:** os upgrades permanentes de poder — a armadilha analisada na
seção 5.2. Sky Force acerta na retenção e paga o preço no balanceamento: o jogo
fica trivial no fim.

### Vampire Survivors (2022) — *a referência de acessibilidade*

**O que faz melhor:** transforma progressão numérica em **espetáculo**. A tela
enchendo de efeitos é a recompensa. E o loop de escolha (subiu de nível → escolha
1 de 4) é a mecânica de retenção mais eficiente já criada.

**Aprendizado aplicável:** o **momento de escolha frequente e curto** é o que o
Altar traz para cá. E o princípio de que **a evolução deve ser visível na tela**,
não numa barra de status — o nível 3 de uma arma precisa parecer diferente.

**O que não copiar:** a ausência de habilidade mecânica. Nosso público quer
executar bem, não só acumular.

### Hades (2020) — *a referência de estrutura de run*

**O que faz melhor:** três coisas. (a) **Bênçãos como decisão**, não como loot —
sempre 3 de N, sempre significativas. (b) **Ritmo de sala**: tensão → recompensa →
escolha → tensão, sem exceção. (c) **A derrota é conteúdo** — perder avança a
narrativa.

**Aprendizado aplicável:** o Altar é o Salão dos Deuses adaptado ao scroll
vertical. E "a derrota é conteúdo" é o princípio que justifica os Fragmentos de
Coroa e o Diário de Criaturas: perder precisa render.

**O que não copiar:** o volume de diálogo e de personagens. É um jogo de 20 pessoas
com narrativa de 300.000 palavras. Nossa narrativa é **ambiental**.

### Enter the Gungeon (2016) — *a referência de game feel*

**O que faz melhor:** o **peso de cada ação**. Cada tiro tem recuo, cada morte tem
partícula, cada porta tem som. E o **dodge roll** é o verbo que define o jogo:
uma ação de risco calculado que separa iniciantes de veteranos.

**Aprendizado aplicável:** o Mergulho do Grifo é o nosso dodge roll — e precisa da
mesma qualidade de execução (i-frames legíveis, recarga visível, momento de
vulnerabilidade após). Também: a lista da seção 8 é essencialmente "o que o
Gungeon faz que nós não fazemos".

**O que não copiar:** a dificuldade punitiva com progressão lenta. Afasta o público
secundário inteiro.

### Síntese

| De | Extraímos |
| --- | --- |
| Super Aleste | Níveis de arma que mudam comportamento |
| Jamestown | Convicção estética total; risco/recompensa (Vaunt → graze) |
| Sky Force | Toda run rende algo, mesmo perdida |
| Vampire Survivors | Escolha frequente e curta; evolução visível |
| Hades | Ritmo tensão-recompensa-escolha; derrota como conteúdo |
| Gungeon | Peso de cada ação; um verbo de assinatura arriscado |

**O que ninguém faz e nós faremos:** um recurso de dificuldade que o jogador
**escolhe deixar subir** porque é lucrativo. Essa é a coluna vertebral do
posicionamento.

---

## 14. O que NÃO devemos fazer

Esta seção existe porque a maior causa de morte de projetos indies não é falta de
ideias — é excesso delas.

### 14.1 Proibições absolutas

| ❌ Decisão | Por que destruiria o projeto |
| --- | --- |
| **Multiplayer (qualquer forma)** | Multiplica a complexidade por 4 e o QA por 10. Rede em Go com Ebitengine para um jogo de 60 TPS determinístico é um projeto inteiro. **Um shmup solo bem feito vale infinitamente mais que um shmup com netcode ruim.** |
| **Mundo aberto / hub explorável** | Nega o gênero. O valor do shmup é a densidade por segundo. Andar por um mapa entre fases é tempo morto disfarçado de conteúdo. |
| **RPG com atributos e equipamentos** | Transforma habilidade em planilha. O jogador começa a perder por ter escolhido números errados, não por ter jogado mal — a pior sensação possível no gênero. |
| **Crafting** | Não existe fantasia que o sustente. Adiciona telas, economia e tempo fora do jogo, em troca de nada. |
| **Editor de fases para jogadores** | É um segundo produto. Consome meses, exige interface complexa e beneficia 2% dos jogadores. |
| **Procedural generation das fases** | O ritmo de um shmup é composto, não sorteado. Fases geradas produzem "conteúdo infinito e igualmente medíocre". A variação vem das relíquias e da Corrupção — variar as **regras**, não o **layout**. |
| **Microtransações / passe de temporada** | O público reage com hostilidade ativa. Custaria mais em avaliações negativas do que renderia. |
| **Sistema de níveis do jogador (XP global)** | Meta-progressão de poder disfarçada. Rejeitada na seção 5.2 pelos mesmos motivos. |
| **Narrativa com diálogos e cutscenes** | Fora da nossa competência e do nosso orçamento. Corta o ritmo. A narrativa é ambiental. |
| **Port simultâneo para console** | Certificação, QA e burocracia paralelos ao desenvolvimento. Só depois de 1.0 estável na Steam. |

### 14.2 Armadilhas de execução

| ⚠️ Armadilha | Como se manifesta | Contramedida |
| --- | --- | --- |
| **Excesso de sistemas** | "E se tivesse também…" | Um sistema novo por versão. Se dois competem, o segundo vai para a lista de 1.1. |
| **Refatorar sem necessidade** | O código é bom; a tentação de deixá-lo perfeito é real | **Refatorar apenas quando um recurso planejado exigir.** A camada de input é necessária (gamepad); um ECS não é. |
| **Perfeccionismo de arte** | Redesenhar o grifo cinco vezes | Definir o guia de estilo (seção 9), aplicá-lo e seguir. Só voltar a arte antiga na 1.0. |
| **Balancear sem medir** | O stun-lock sobreviveu a 120 testes porque nenhum teste pergunta se o jogo é divertido | Construir a ferramenta de simulação headless na 0.2. **É o item de infraestrutura mais valioso do roadmap.** |
| **Adiar o gamepad** | "Depois eu faço" — e o input está espalhado por 5 arquivos | Camada de ações na 0.4. Quanto mais tarde, mais caro. |
| **Adiar a arte e o áudio externos** | Descobrir na 0.5 que o artista tem 4 meses de fila | **Contratar na 0.3.** É a decisão de produção mais importante deste documento. |
| **Escutar todo o feedback do EA** | Cada jogador quer um jogo diferente | Escopo congelado ao entrar em EA. Feedback ajusta balanceamento, não adiciona sistemas. |
| **Lançar sem listas de desejos** | 2.000 wishlists = lançamento invisível | Meta de 10.000. Abaixo de 7.000, adiar. |
| **Confundir "mais conteúdo" com "melhor jogo"** | Adicionar a fase 7 em vez de consertar as 6 existentes | Se uma fase não está boa, a fase seguinte não vai ajudar. |

### 14.3 O erro específico deste projeto

O projeto tem uma tendência clara e documentada — visível no `ROADMAP.md`, no
`README.md` e no histórico de commits: **resolver problemas de design com
engenharia.**

Foram construídos: sistema de sprites com fallback duplo, sistema de áudio
procedural com onze trilhas, seis inimigos com comportamentos distintos, três
dificuldades com jitter parametrizado, persistência, modo sobrevivência, 120
testes. Tudo bem feito.

E ainda assim: **94% da campanha são dois inimigos, uma arma domina as outras
três, e o chefe morre em cinco segundos.**

Os sistemas existem. O **conteúdo dentro deles** e o **balanceamento entre eles**
não receberam a mesma atenção. A regra para os próximos doze meses:

> **Antes de construir um sistema novo, gaste uma semana usando os que já
> existem.**

Se `stages.go` fosse reescrito para usar os seis inimigos que já existem, o jogo
melhoraria mais do que com qualquer sistema novo deste documento. Isso custa dois
dias e está na versão 0.2 por esse motivo.

---

## 15. Plano de ação: os 30 próximos passos

Se eu assumisse este projeto amanhã como Diretor do Jogo, esta é exatamente a
ordem em que eu agiria.

**Convenções:** Impacto e Dificuldade em escala 1–5 · Tempo em dias-dev de 4 h ·
Risco 🟢 baixo / 🟡 médio / 🔴 alto.

### Bloco 1 — Semana 1: entender e medir (passos 1–4)

| # | Passo | Objetivo | Imp. | Dif. | Tempo | Risco |
| --- | --- | --- | :-: | :-: | ---: | :-: |
| **1** | **Jogar 20 runs completas** nas três dificuldades, anotando cada momento de tédio, confusão ou injustiça | Substituir análise de código por experiência real | 5 | 1 | 2 d | 🟢 |
| **2** | **Construir a ferramenta de simulação headless** — roda N runs com seed fixa e reporta DPS por arma, tempo até nível máximo, abates, mortes, taxa de fuga, duração do chefe | Parar de balancear no escuro. O item de infraestrutura mais valioso do roadmap | 5 | 3 | 3 d | 🟢 |
| **3** | **Playtest com 3 pessoas que nunca jogaram**, sem instruções, gravando a tela | Descobrir o que é óbvio para o autor e opaco para todos os outros | 5 | 1 | 1 d | 🟢 |
| **4** | **Escrever o documento de uma página** com a visão da seção 2 e fixá-lo no repositório | Todo o resto do trabalho passa a ter um critério de aceitação | 4 | 1 | 0,5 d | 🟢 |

> *Motivo do bloco:* nenhuma decisão deste documento deve ser executada antes de
> ser confirmada por observação. O passo 2 é o que impede que o próximo stun-lock
> passe despercebido por mais 120 testes.

### Bloco 2 — Semanas 2–5: consertar o núcleo (passos 5–12)

| # | Passo | Objetivo | Imp. | Dif. | Tempo | Risco |
| --- | --- | --- | :-: | :-: | ---: | :-: |
| **5** | **Remover o stun automático da Lança de Luz** ([`status.go:19`](game/status.go:19)) | Devolver o combate e o chefe ao jogo. Se só uma coisa for feita, é esta | 5 | 1 | 0,5 d | 🟢 |
| **6** | **Rebalancear as três armas** para ~65 DPS em suas janelas ideais, com identidades sem sobreposição | Tornar a escolha de arma uma decisão real | 5 | 3 | 4 d | 🟡 |
| **7** | **Reescrever a distribuição de inimigos em `stages.go`** (corvo 63,5% → 35%; gárgula, wyvern, mago e balista para 8–12% cada) | Usar o bestiário que já existe. Maior ganho por esforço do projeto | 5 | 2 | 3 d | 🟡 |
| **8** | **Aplicar a gramática de introdução:** cada inimigo aparece isolado, depois combinado, depois em massa | Dar curva de aprendizado ao jogo | 4 | 2 | 2 d | 🟢 |
| **9** | **Corrigir a economia:** `respawn` para `initialHealth`; remover o bônus por bomba não usada; reduzir o peso do `bossScore` de 5.000 para ~1.500 | Fazer a pontuação medir habilidade e a bomba ser usada | 4 | 1 | 1 d | 🟢 |
| **10** | **Rebalancear o chefe** (vida, ritmo, telegrafia) agora que ele executa seus padrões | Recuperar o clímax da campanha | 4 | 2 | 2 d | 🟡 |
| **11** | **Rever a curva de drops:** reduzir a chance base e escalonar por trecho | Espalhar a progressão pela run inteira | 3 | 2 | 1 d | 🟡 |
| **12** | **Tag `v0.2` + playtest de validação** com as mesmas 3 pessoas | Confirmar que o jogo ficou mais divertido, não só diferente | 5 | 1 | 1 d | 🟢 |

> *Motivo do bloco:* é impossível avaliar se a Corrupção é uma boa ideia num jogo
> cujo combate está quebrado. Este bloco é pré-requisito de todo julgamento
> posterior.

### Bloco 3 — Semanas 6–12: implantar o diferencial (passos 13–18)

| # | Passo | Objetivo | Imp. | Dif. | Tempo | Risco |
| --- | --- | --- | :-: | :-: | ---: | :-: |
| **13** | **Corrupção mínima viável:** contador + barra + multiplicador de pontos + vida dos inimigos escalando | Testar a hipótese central com o menor investimento possível | 5 | 2 | 3 d | 🟡 |
| **14** | **Testar a Corrupção com 5 jogadores externos** e decidir explicitamente: seguir, ajustar ou abandonar | **Ponto de decisão do projeto.** Descobrir agora, não na 0.8 | 5 | 1 | 2 d | 🔴 |
| **15** | **Moldura lateral decorada** com o Medidor de Corrupção em destaque permanente | Resolver o 3:4, limpar o campo de jogo e dar palco à mecânica principal | 4 | 3 | 4 d | 🟢 |
| **16** | **Três variantes corrompidas** de inimigos existentes, ativadas acima de 50% | Tornar a Corrupção visível no que o jogador enfrenta, não só num número | 4 | 3 | 4 d | 🟡 |
| **17** | **Degradação visual e sonora** por faixa de corrupção (dessaturação do cenário, distorção da trilha) | O jogador precisa *sentir* o reino apodrecendo | 4 | 3 | 3 d | 🟢 |
| **18** | **Contratar artista de pixel art e compositor** — briefing, referências, contrato, primeira entrega agendada | Iniciar o caminho crítico externo 4 meses antes de precisar dele | 5 | 2 | 3 d | 🔴 |

> *Motivo do bloco:* o passo 14 é o momento mais importante do cronograma. Se a
> Corrupção não empolgar, todo o posicionamento deste documento precisa mudar — e
> é infinitamente mais barato descobrir isso na semana 8 do que no mês 12.
> O passo 18 é a decisão que mais afeta o prazo final e é a mais fácil de adiar
> por engano.

### Bloco 4 — Semanas 13–18: o verbo e a plataforma (passos 19–23)

| # | Passo | Objetivo | Imp. | Dif. | Tempo | Risco |
| --- | --- | --- | :-: | :-: | ---: | :-: |
| **19** | **Camada de input por ações** — extrair todo `IsKeyPressed` para um mapa `ação → tecla/botão` | Destravar gamepad, Steam Deck e remapeamento. Cada mês de adiamento encarece | 4 | 3 | 3 d | 🟢 |
| **20** | **Suporte completo a gamepad**, incluindo toda a navegação de menus | Requisito de plataforma e de público | 4 | 2 | 2 d | 🟢 |
| **21** | **Mergulho do Grifo** + fôlego (i-frames, dano de contato, deslocamento de câmera) | Dar ao jogo o verbo que o distingue de toda a concorrência | 5 | 3 | 4 d | 🟡 |
| **22** | **Graze:** raspar projéteis carrega bomba | Converter medo em ganância; fazer as bombas circularem | 3 | 2 | 2 d | 🟢 |
| **23** | **Pacote de game feel:** itens 1, 2, 4, 13, 25, 35, 43, 48 da seção 8 | Elevar a sensação de "funciona" para "é gostoso" | 4 | 2 | 4 d | 🟢 |

### Bloco 5 — Semanas 19–26: a cara e a variedade (passos 24–27)

| # | Passo | Objetivo | Imp. | Dif. | Tempo | Risco |
| --- | --- | --- | :-: | :-: | ---: | :-: |
| **24** | **Integrar a arte dos 3 primeiros biomas** (fim dos retângulos) + fonte própria + HUD reconstruído | Matar definitivamente a aparência de protótipo | 5 | 3 | 6 d | 🔴 |
| **25** | **Integrar as 3 primeiras trilhas produzidas** + 20 efeitos, incluindo a voz do grifo | O maior salto de percepção de qualidade por real investido | 5 | 2 | 3 d | 🟡 |
| **26** | **Altar entre trechos** + 16 relíquias | Dar ritmo, decisão e variação entre runs | 5 | 4 | 6 d | 🟡 |
| **27** | **Runas Fundidas** (3 fusões, com arte e som próprios) | Estender a curva de poder para a run inteira | 4 | 4 | 5 d | 🟡 |

### Bloco 6 — Semanas 27–30: preparar a exposição (passos 28–30)

| # | Passo | Objetivo | Imp. | Dif. | Tempo | Risco |
| --- | --- | --- | :-: | :-: | ---: | :-: |
| **28** | **Onboarding diegético de 30 s** na fase 1 — ensinando movimento, precisão (a hitbox!), bomba e mergulho sem uma tela de texto | O jogador precisa entender a hitbox de 4px, que hoje está escondida num menu | 4 | 2 | 3 d | 🟢 |
| **29** | **Publicar a página da Steam** com trailer, 6 capturas e descrição da seção 12 · abrir listas de desejos | Começar a acumular wishlists 8 meses antes do lançamento. **Adiar isto é o erro de marketing mais comum e mais caro** | 5 | 2 | 4 d | 🟡 |
| **30** | **Preparar a demo** (fase 1 + mini-chefe) e inscrever no próximo **Steam Next Fest** | O maior gerador de listas de desejos disponível sem verba de marketing | 5 | 2 | 3 d | 🟡 |

### Resumo do plano

| Bloco | Semanas | Entrega | Marco |
| --- | ---: | --- | --- |
| 1 | 1 | Dados reais e ferramenta de medição | — |
| 2 | 2–5 | Núcleo justo e divertido | **v0.2** |
| 3 | 6–12 | O diferencial implantado e validado | **v0.3** |
| 4 | 13–18 | Verbo próprio e plataforma destravada | **v0.4** |
| 5 | 19–26 | Cara de jogo comercial e variação entre runs | **v0.5 / v0.6** |
| 6 | 27–30 | Exposição pública e demo | **Early Access** |

**Total: 30 semanas (~7 meses) até o Early Access**, seguindo o roadmap da seção
7 até a 1.0 em ~14 meses.

### Os três passos que eu não deixaria ninguém pular

1. **Passo 5** — remover o stun-lock. Custa meia hora e devolve o jogo inteiro.
2. **Passo 14** — validar a Corrupção com jogadores externos. É o ponto de decisão
   de todo o projeto, e adiá-lo transforma um ajuste de rota em um desastre.
3. **Passo 18** — contratar arte e áudio na semana 12, não na semana 20. É a
   decisão que determina se o prazo é de 14 meses ou de 24.

---

## Palavra final

Este projeto tem algo raro: **uma base técnica melhor do que o jogo que ela
sustenta.**

Isso é uma posição privilegiada, e é preciso reconhecê-la. A grande maioria dos
projetos indies chega a este ponto com o problema inverso — uma boa ideia
apodrecendo dentro de um código que ninguém consegue mais mudar. Aqui, as fases
são dados, os inimigos são uma interface, os testes protegem contra regressão, os
assets têm fallback, o áudio degrada sem quebrar e existe uma seed reproduzível.
Cada uma dessas decisões custa esforço na hora e paga por anos. Elas já foram
tomadas, e foram tomadas certas.

O que falta não é habilidade técnica. É **convicção de design**.

O jogo hoje é seguro. Três armas equilibradas no papel, seis inimigos bem
comportados, um chefe competente, dois modos, três dificuldades. Nada nele
ofende — e nada nele é memorável. É o retrato de um projeto que foi construído
respondendo à pergunta *"o que um shoot'em up precisa ter?"* em vez de *"por que
alguém jogaria este shoot'em up e não outro?"*.

A resposta que este documento propõe é a Corrupção: um jogo onde o erro não é
apenas punido, mas **negociado**; onde deixar o mundo apodrecer é uma estratégia
válida e lucrativa; e onde o final que você vê depende de quanto você foi capaz —
ou disposto — a segurar. Não é a única resposta possível. Mas é uma resposta
específica, executável com o que já existe no repositório, e é dizível em uma
frase para um estranho.

Um jogo pequeno com uma ideia forte vence um jogo pequeno com dez ideias
razoáveis. Sempre.

Os próximos doze meses não são de engenharia. São de **escolha, corte e
convicção**: escolher a ideia, cortar tudo que a dilui, e executá-la com uma
teimosia desproporcional ao tamanho do projeto.

A base está pronta. Falta decidir que jogo ela vai sustentar.

---

*Documento vivo. Revisar a cada marco de versão. As metas da seção 2.10 são o
critério — se elas não estiverem sendo atingidas, é o plano que muda, não a
medição.*

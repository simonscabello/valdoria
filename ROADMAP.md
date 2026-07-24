# Roadmap — Asas de Valdoria

> Documento de direção de produto e técnico. Analisa o estado real do repositório
> (código em `game/`) e define os próximos ciclos para levar o projeto de **MVP
> funcional** a **primeira versão publicável**. Prático, priorizado e baseado no
> código existente — não em suposições.

Data da análise: estado após a revisão técnica (76 testes, `vet`/`race`/build limpos).

---

## Estado após execução do plano (atualização)

Os 10 itens priorizados foram implementados. Resumo do que mudou desde a análise:

- **A1** economia de armas corrigida (trocar runa preserva o nível).
- **Quick wins** de game feel (contraste do HUD, telegraphs, popups de bônus, barra de progresso, sombra da bomba, brasas).
- **B2** assinatura mágica das armas (som por arma + rastro na Luz).
- **A2** fase orientada a dados (`stageDef`/`waveDef`, `stages.go`); Fase 1 recriada em dados.
- **C1** dois inimigos novos: **Balista** (ameaça terrestre, rajada mirada) e **Feiticeiro** (anel de projéteis).
- **C3** campanha de **3 fases** encadeadas com biomas próprios, culminando em Vharak.
- **D1** **dificuldades** Fácil/Normal/Difícil (vidas, bombas, vida dos inimigos, drops).
- **A3** **persistência** de recordes e opções (`save.go`, JSON no diretório de config).
- **D2** **modo Sobrevivência** (ondas infinitas com escalada, recorde próprio).
- **B1** **fonte própria** de UI (`text/v2` + `basicfont`, colorida, sem asset externo).

Cobertura: **95 testes**, `go vet`/`-race`/build limpos. Nova dependência **direta**: `golang.org/x/image` (já era transitiva do Ebitengine).

**Pendências conhecidas:** QA visual da fonte/biomas (não há teste de renderização); 2º chefe / generalização do `Boss` (R5/C2) — Vharak segue como chefe único ao fim da campanha; tela de Opções dedicada (F1); gamepad/i18n (A4/adiado).

### Ajustes pós-playtest do usuário

Feedback jogado (dificuldade Normal) e correções aplicadas:
- **"Poucos monstros / muita espera"** → ondas das 3 fases reescritas para fluxo denso e contínuo (contagens maiores, intervalos menores, ondas sobrepostas) e fases um pouco mais curtas; sobrevivência também adensada. Guardado por `TestStagesAreDense`.
- **Textos sobrepostos** → HUD re-layoutado: alinhamento à direita agora usa a largura **medida** da fonte (`labelRight`), barra de vida movida para o rodapé (não colide mais com o placar), barra do chefe afinada no topo absoluto, rótulo do Game Over encurtado e nome de trecho mais longo renomeado.
- **Fundo se confunde com os inimigos** → cenário escurecido (overlay) para recuar o parallax e **contorno escuro** em jogador, inimigos e chefe, reforçando as silhuetas e a leitura.

---

## 1. Diagnóstico atual

### O que já está implementado (e sólido)
- **Game loop** limpo em `Update`/`Draw`/`Layout`, TPS fixo 60 (determinístico, sem dependência de FPS).
- **Máquina de estados** completa: menu, controles, jogando, pausa, chefe, vitória, game over — com transições em fade e reset de sessão sem vazamento de estado.
- **Jogador**: 8 direções normalizadas, modo de precisão com hitbox pequena revelada, invencibilidade/piscar, clamp de tela.
- **4 inimigos** com comportamentos distintos via interface (`enemyBehavior`): corvo (reto), harpia (zigue-zague + tiro), gárgula (entra/ataca/sai lateral), wyvern (alinha + tiro mirado).
- **Chefe Vharak**: 3 fases por % de vida, 5 padrões (mirado, cone, arco, invocação, varredura), cristais como pontos fracos, telegrafia de padrão, morte segura (corrigida).
- **Fase 1** como linha do tempo determinística (`newLevel()`), 4 trechos com temas de parallax próprios.
- **3 armas × 3 níveis** + **5 power-ups**; pontuação com combo/multiplicador, bônus de formação e de trecho sem dano.
- **Áudio** procedural com fallback (música/efeitos, fade, mute, frame-gate), degrada sem dispositivo.
- **Dev mode** robusto por env vars (seed fixa, pular trecho, iniciar no chefe, HUD, hitboxes, invencível).
- **Efeitos**: partículas, screen shake configurável, flash de dano, hit-stop, popups de score, projéteis inimigos com leitura própria.

### O que funciona bem
- Fundação de engenharia: separação lógica/render, RNG de jogabilidade isolado do de efeitos, testabilidade real (76 testes).
- Extensibilidade de **inimigos** (interface + construtores) e de **padrões do chefe** (mapa `bossPhase → []bossPattern`, orientado a dados).
- Feedback de combate (após a revisão): impacto, shake, flash, popups, hit-stop.

### O que ainda parece incompleto
- **Uma única fase e um único chefe** — o jogo termina em ~3–4 min e não há segundo bioma nem segundo confronto.
- **Sem progressão entre partidas**: high score só na sessão, sem persistência, sem desbloqueios.
- **Tutorial inexistente** — só a tela estática de "Controles".

### O que parece provisório
- **Todo o texto usa `ebitenutil.DebugPrintAt`** (fonte de debug do Ebitengine). É o maior sinal de "protótipo": HUD, menus, avisos e nome do chefe são renderizados com a fonte de diagnóstico.
- **Arte 100% geométrica** (`vector.DrawFilledRect`): funcional e legível, mas sem identidade forte.
- **Áudio procedural** (ondas quadradas/senoidais) — placeholder assumido no próprio código.

### Sistemas frágeis / dívida técnica
- **`newLevel()` e `phase1Sections()` são hardcoded** em Go (chamadas `add(...)`). Não há um formato de dados de fase — adicionar fases significa escrever código, não conteúdo.
- **Armas e power-ups são enums com `switch`** espalhados (`weapon.go`: `weaponName`/`weaponCooldown`/`fireWeapon`; `player.go`: `gainWeapon`/`applyPowerup`; `powerup.go`: `color`; `randomRune`). Adicionar 1 arma toca ~5 switches.
- **Chefe é monolítico e específico do Vharak** (visual, cristais, posições de fraqueza hardcoded). Um 2º chefe exige generalizar ou duplicar.
- **Estado global de pacote**: `dev`, `rng`/`fxRand`, `sharedAudio`/`sharedContext`. Funciona para uma instância; dificulta múltiplas instâncias e é acoplamento implícito.
- **Input lido direto** via `ebiten.IsKeyPressed`/`inpututil` disperso — sem camada de input (bloqueia gamepad e remapeamento).
- **Strings inline em português** — sem camada de i18n.
- **Sem persistência** (nenhum `save`/config em disco).

### Funcionalidades que existem mas ainda não entregam boa experiência
- **Troca de arma penaliza o jogador**: pegar uma runa de arma diferente da atual **reseta para o nível 1** (`gainWeapon`). O jogador aprende a **evitar** power-ups — o oposto do que um shmup quer. Este é o defeito de design mais relevante do núcleo.
- **Invocação Ancestral (bomba)** funciona, mas é um botão de pânico genérico sem leitura de "ancestralidade" — pouca personalidade para uma mecânica de assinatura.
- **Duração da fase (~2,7 min de ondas)** com só 4 inimigos tende à monotonia no meio (trechos 2–3).

### Recursos implementados "cedo demais" (antes do núcleo sólido)
- **Cenário parallax de 4 camadas com temas por bioma** e **3 trilhas de música** — bonito, mas investido antes de o *conteúdo de jogo* (fases/chefes/armas) existir. Não é para remover; é um alerta de priorização: a base de apresentação está à frente da base de conteúdo.

> **Leitura de produto:** a *engenharia* está à frente do *conteúdo* e do *design de recompensa*. O projeto tem esqueleto de qualidade e pouca carne. O próximo esforço deve ser **conteúdo + power fantasy + apresentação**, reutilizando a fundação — não mais sistemas.

---

## 2. Avaliação do núcleo do jogo

Loop: entrar → ondas → derrotar → coletar → ficar forte → sobreviver → chefe → concluir.

| Pergunta | Resposta |
| --- | --- |
| O loop está claro? | **Sim.** Legível e sem fricção de estados. |
| É divertido? | **Parcialmente.** O combate é gostoso; a recompensa (power-up) é ambígua/punitiva. |
| Sustenta vários minutos? | **No limite.** Início e chefe seguram; o miolo (trechos 2–3) esvazia. |
| Variedade suficiente? | **Não.** 4 inimigos, 1 chefe, 1 bioma jogável. |
| Progressão perceptível na partida? | **Fraca.** Subir de nível de arma ajuda, mas o reset ao trocar quebra a curva. |
| O jogador sente que fica mais forte? | **Inconsistente** — por causa do reset de arma e da perda de nível ao morrer. |
| Decisões de arma/power-up são interessantes? | **Não** hoje: viram "pegar só o que combina". Deveria ser o coração da diversão. |
| Motivo para rejogar? | **Baixo.** Só bater o próprio score da sessão (não persiste). |
| Identidade própria? | **Embrionária.** Tema medieval/grifo presente na forma, ausente na apresentação (fonte debug, formas). |

### Maior problema do núcleo
**O sistema de armas/power-ups gera decisões negativas em vez de empolgantes.** Trocar de runa rebaixa você ao nível 1; morrer tira um nível. O eixo central "coletar → ficar mais forte" — que deveria ser o vício do gênero — hoje ensina o jogador a ter medo de itens. Consertar isso é mais valioso do que qualquer conteúdo novo.

---

## 3. Avaliação como jogador experiente de shmup

| Aspecto | Nota | Observação |
| --- | --- | --- |
| Responsividade dos controles | **Bom** | Input por tick, sem latência artificial. |
| Velocidade da montaria | **Bom** | 2.5 normal / 1.0 precisão — adequado à tela 240×320. |
| Clareza da hitbox | **Bom** | Hitbox 4px revelada na precisão (marca vermelha). |
| Legibilidade dos projéteis | **Bom** (pós-revisão) | Inimigo magenta com núcleo/contorno, distinto das armas. |
| Qualidade dos padrões de ataque | **Precisa de ajustes** | Inimigos comuns têm padrões simples; variedade vem só do chefe. |
| Justiça da dificuldade | **Bom** | Varredura sempre desviável; 1 hit/frame; morte do chefe segura. |
| Quantidade de inimigos | **Precisa de ajustes** | Miolo da fase rarefeito; picos poderiam ser mais densos. |
| Sensação dos disparos | **Precisa de ajustes** | Sons de placeholder; falta "peso" nas armas fracas (chamas). |
| Feedback de impacto | **Bom** (pós-revisão) | Flash, partículas, hit-stop, popups. |
| Ritmo das ondas | **Precisa de ajustes** | Curva de intensidade plana entre trechos. |
| Frequência de power-ups | **Precisa de ajustes** | Drops garantidos por onda + 20% aleatório; ok, mas ofuscado pela regra de reset. |
| Diversidade das armas | **Bom** de conceito, **Precisa de ajustes** de balanceamento | 3 arquétipos claros (reto/leque/perfurante). |
| Uso da Invocação Ancestral | **Precisa de ajustes** | Funciona; falta personalidade e leitura de recurso. |
| Duração da fase | **Precisa de ajustes** | Longa para o volume de conteúdo (repetição). |
| Qualidade do chefe | **Bom** | Melhor peça do jogo: fases, telegrafia, pontos fracos. |
| Clareza do HUD | **Precisa de ajustes** | Informativo, mas em fonte de debug e espalhado. |
| Menus/derrota/vitória | **Precisa de ajustes** | Funcionais; apresentação de protótipo. |
| Vontade de continuar jogando | **Precisa de ajustes** | Sem gancho de progressão/rejogo. |
| Sprites/animações | **Bom (fundação)** | Entidades + flap 2 frames; glow em projéteis; molduras de UI; cenário ainda geométrico. |
| Tutorial/onboarding | **Ainda não implementado** | Só tela de controles. |
| Progressão entre partidas | **Ainda não implementado** | Sem persistência. |

---

## 4. Identidade do jogo

**Presente na forma, ausente na percepção.** Os elementos certos existem (grifo + cavaleiro, magias elementais, criaturas aladas corrompidas, dragão, biomas medievais vistos de cima), mas a apresentação (fonte de debug + retângulos) apaga a fantasia.

### Como fortalecer a identidade *percebida durante o jogo* (sem lore gigante)
- **Direção artística / silhuetas:** a leitura por silhueta já existe; formalizar um "guia de formas" (grifo dourado, corrupção em roxo/magenta, pedra para gárgulas, verde para wyvern). Manter geométrico, mas com **contornos e 2–3 tons por entidade** para dar volume.
- **Paleta por bioma:** já há temas (`sectionThemes`); reforçar contraste entre Campos (verde/azul), Vila em chamas (laranja/vermelho), Muralhas (cinza/aço), Castelo (roxo/corrupção). A cor conta a história do avanço.
- **Efeitos mágicos com "assinatura":** cada arma com uma cor/rastro inconfundível (Luz = dourado; Chamas = laranja com brasas; Gelo = ciano com estilhaços). A bomba como **dragão ancestral** já existe — dar-lhe rugido e sombra na tela.
- **Interface temática:** trocar a fonte de debug por uma **bitmap font medieval** e molduras simples (pergaminho/pedra). Impacto de identidade altíssimo por custo baixo.
- **Narrativa ambiental (diegética):** o parallax já "conta" o cerco (campos → vila em chamas → muralhas → castelo). Reforçar com 1 linha de texto por trecho (já há `warning`) e o nome do reino.
- **Nomes:** batizar inimigos e trechos com nomes do mundo (ex.: "Corvos de Morrighan", "Muralhas de Eldoria"). Barato e memorável.
- **Personalidade dos inimigos:** dar a cada um 1 traço legível (gárgula que "encara" antes de atacar; wyvern que ruge ao mirar).

> Foco: **fonte + paleta + assinatura mágica**. Esses três, percebidos em segundos, transformam a impressão de "protótipo" em "jogo".

---

## 5. Conteúdo mínimo para uma "1ª versão completa"

Meta realista para **solo/equipe pequena** (não uma campanha gigante):

| Eixo | MVP atual | Meta 1ª versão publicável |
| --- | --- | --- |
| Fases | 1 | **3** biomas encadeados (Campos → Muralhas → Castelo) |
| Duração/fase | ~2,7 min | **1,5–2 min** cada (ritmo mais denso) |
| Biomas | 4 trechos de 1 fase | **3** biomas distintos |
| Inimigos | 4 | **6–7** (2–3 novos, reusando a interface) |
| Chefes | 1 | **2** (Vharak + 1 intermediário/menor) |
| Armas | 3×3 níveis | **3×3** (manter; corrigir o reset) |
| Power-ups | 5 | **5** + 1 utilitário (ex.: "concentração"/lentidão curta) opcional |
| Personagens jogáveis | 1 | **1** (não expandir agora) |
| Dificuldades | 1 | **2–3** (Fácil/Normal/Difícil via multiplicadores) |
| Progressão entre partidas | 0 | **High score persistente** + desbloqueio simples |
| Menus | ok | + Opções (áudio/dificuldade), tutorial curto |
| Áudio | procedural | **1 trilha decente por contexto** + efeitos revisados |
| Sprites | 2 | entidades em pixel art procedural + PNG opcional; próximas: frames de asa, glow, UI |
| Tutorial | 0 | **onboarding de 20s** integrado à fase 1 |

**Não fazer agora:** múltiplos personagens/montarias, roguelike de builds, online, loja, conquistas complexas, save de progresso ramificado.

---

## 6. Progressão e rejogabilidade — o que vale a pena

| Sistema | Veredito | Motivo |
| --- | --- | --- |
| **High score persistente (local)** | **Fazer** | Custo baixo, fecha o loop de "melhorar"; base de tudo. |
| **Dificuldades (Fácil/Normal/Difícil)** | **Fazer** | Reusa valores existentes via multiplicadores; alto valor/baixo custo. |
| **Modo Sobrevivência (ondas infinitas)** | **Fazer (leve)** | Reusa spawns e chefe; grande retenção com pouco conteúdo novo. |
| **Desbloqueio simples (ex.: iniciar com arma X)** | **Considerar** | Dá meta de médio prazo; barato se ganchado ao score. |
| **Conquistas locais** | **Considerar (poucas)** | 8–12 metas; barato e motivador; evitar 50+. |
| **Segredos/chefe opcional** | **Considerar depois** | Charme, mas exige conteúdo extra. |
| Escolha de personagens/montarias | **Evitar agora** | Multiplica arte/balanceamento sem núcleo pronto. |
| Rotas alternativas / modificadores | **Evitar agora** | Combinatória de teste alta. |
| New Game Plus / progressão permanente | **Evitar agora** | Exige economia e balanceamento de longo prazo. |
| Runs roguelike com escolhas aleatórias | **Evitar agora** | Mudança de gênero; escopo enorme. |

> Combo recomendado: **score persistente + dificuldades + sobrevivência**. Três sistemas pequenos que, juntos, multiplicam as horas jogáveis sem exigir arte nova.

---

## 7. Arquitetura e preparação técnica

### Facilidade de adicionar hoje
| Adicionar… | Facilidade | Nota |
| --- | --- | --- |
| Novo inimigo | **Fácil** | Interface `enemyBehavior` + construtor + case em `spawnEnemy`. Bom. |
| Novo padrão de ataque do chefe | **Fácil** | Mapa de padrões orientado a dados. Bom. |
| Nova arma | **Médio** | Toca ~5 `switch` (nome, cooldown, fire, gainWeapon, cor). |
| Novo power-up | **Médio** | Enum + switches de `applyPowerup`/`color`/`randomRune`. |
| Novo chefe | **Difícil** | `Boss` é específico do Vharak (visual/cristais/posições). |
| Nova fase/cenário | **Difícil** | `newLevel()` é código, não dados; sem "Stage" reutilizável. |
| Salvamento | **Difícil** | Inexistente; nenhum ponto de serialização. |
| Configurações (arquivo) | **Médio** | Só env vars; falta camada de settings. |
| Gamepad | **Difícil** | Input espalhado, sem abstração de ações. |
| Localização | **Difícil** | Strings inline em pt-BR. |
| Builds multiplataforma | **Fácil** | Makefile já cobre Linux/Windows; sem assets externos. |

### Refatorações **obrigatórias** (antes de expandir conteúdo)
1. **Formato de dados de fase (`Stage`/`Wave`)** — extrair `newLevel()` para uma estrutura declarativa (slice de waves + sections + tema + chefe). Sem isso, "3 fases" vira 3× código duplicado e teste manual.
2. **Corrigir a economia de armas** (design, não motor): trocar de runa **não** deve zerar; deve manter progresso ou converter (ex.: runa diferente sobe nível da atual OU acumula até "trocar"). É pré-requisito de diversão para qualquer expansão.

### Refatorações **recomendadas** (facilitam, não bloqueiam)
3. **Camada de input (`action → tecla/botão`)** — habilita gamepad e remapeamento; centraliza o que hoje está espalhado.
4. **Tabela de definição de armas** (struct com nome/cor/cooldown/fire-fn) para eliminar os `switch` paralelos.
5. **Generalizar o chefe** (interface `bossBehavior`/config de silhueta e pontos fracos) — só quando o 2º chefe estiver no backlog imediato.
6. **Persistência mínima** (`save.go`: high score, opções) — arquivo JSON no diretório de config do usuário.

### Refatorações que **podem esperar**
7. i18n (só quando houver intenção de publicar em outro idioma).
8. Remover estado global (`dev`/rng/audio) — só se for necessário rodar múltiplas instâncias/testes paralelos.
9. Sistema de partículas mais genérico / pooling — só sob pressão de performance real (hoje o teto de 400 basta).

> **Sem reescrita.** A base é boa. As duas obrigatórias (dados de fase + economia de armas) são cirúrgicas e destravam quase todo o resto.

---

## 8. Backlog priorizado

### A. Correções e fundação
| # | Item | Problema que resolve | Valor p/ jogador | Complexidade | Risco | Dependências | Critério de conclusão |
| --- | --- | --- | --- | --- | --- | --- | --- |
| A1 | **Corrigir economia de armas** ✅ | Power-up punitivo | Alto (núcleo) | Baixa | Baixo | — | ✅ **Concluído**: trocar runa preserva o nível (nunca rebaixa); testes `TestWeaponSwitchKeepsLevel` e `TestWeaponSwitchAtLowLevelPreservesLevel` |
| A2 | **Formato de dados de fase (`Stage`)** | `newLevel()` hardcoded | Alto (destrava conteúdo) | Média | Médio | — | Fase 1 recriada 100% via dados; teste de carga da fase |
| A3 | **Persistência mínima (`save.go`)** | Score/opções não persistem | Médio | Baixa | Baixo | — | High score e opções sobrevivem ao fechar o jogo |
| A4 | **Camada de input (ações)** | Sem gamepad/remap | Médio | Média | Médio | — | Todas as entradas passam por `action`; teclado inalterado |

### B. Game feel
| # | Item | Problema | Valor | Complexidade | Risco | Dep. | Conclusão |
| --- | --- | --- | --- | --- | --- | --- | --- |
| B1 | **Bitmap font + molduras** | Fonte de debug | Alto (percepção) | Média | Baixo | — | Todo texto usa fonte própria; HUD/menus reestilizados |
| B2 | **Assinatura mágica das armas** | Armas sem "peso" | Médio | Baixa | Baixo | — | Cada arma com cor/rastro/som distintos e reconhecíveis |
| B3 | **Curva de intensidade da fase** | Miolo rarefeito | Médio | Baixa | Médio (balance) | A2 | Densidade e picos ajustados; playtest aprova ritmo |
| B4 | **Personalidade da Bomba** | Genérica | Médio | Baixa | Baixo | — | Rugido + sombra + telegrafia; leitura de recurso clara |

### C. Conteúdo
| # | Item | Valor | Complexidade | Risco | Dep. | Conclusão |
| --- | --- | --- | --- | --- | --- | --- |
| C1 | **2–3 novos inimigos** (ex.: arqueiro terrestre, mago-corvo com projétil lento, dragão-filhote) | Alto | Média | Baixo | A2 | Novos `enemyBehavior` + testes de movimento/tiro |
| C2 | **1 chefe intermediário** (menor que Vharak) | Alto | Alta | Médio | A2, R5 | Novo chefe reusando padrões; morte segura |
| C3 | **2ª e 3ª fases (biomas)** | Alto | Média | Médio (balance) | A2 | Fases em dados; temas distintos; encadeamento |
| C4 | **Ameaças terrestres** (torres/balistas na camada de estrutura) | Médio | Média | Médio | A2 | Inimigo fixo que atira de baixo; leitura clara |

### D. Progressão
| # | Item | Valor | Complexidade | Risco | Dep. | Conclusão |
| --- | --- | --- | --- | --- | --- | --- |
| D1 | **Dificuldades (Fácil/Normal/Difícil)** | Alto | Baixa | Baixo | A2 | Multiplicadores de vida/dano/spawn; selecionável no menu |
| D2 | **Modo Sobrevivência** | Alto | Média | Baixo | A2, A3 | Ondas infinitas com escalada; score persistente próprio |
| D3 | **Conquistas locais (8–12)** | Médio | Baixa | Baixo | A3 | Metas checáveis; notificação simples |

### E. Apresentação
| # | Item | Valor | Complexidade | Risco | Dep. | Conclusão |
| --- | --- | --- | --- | --- | --- | --- |
| E1 | **Trilha/efeitos revisados** | Médio | Média | Dep. de áudio | — | 1 trilha por contexto + SFX melhores (mesmo que livres) |
| E2 | **Onboarding de 20s** na fase 1 | Médio | Baixa | Baixo | B1 | Dicas diegéticas de mover/atirar/precisão/bomba |
| E3 | **Telas de menu/vitória/derrota temáticas** | Médio | Baixa | Baixo | B1 | Layout com moldura e paleta do mundo |

### F. Publicação
| # | Item | Valor | Complexidade | Risco | Dep. | Conclusão |
| --- | --- | --- | --- | --- | --- | --- |
| F1 | **Tela de Opções** (volume, dificuldade, vibração) | Médio | Baixa | Baixo | A3 | Ajustes persistidos e aplicados |
| F2 | **Ícone, capturas, página da loja** | Médio | Baixa | Baixo | B1, E3 | Assets de divulgação prontos |
| F3 | **Builds assinados + smoke test por plataforma** | Médio | Média | Médio | tudo | Build roda limpo em Win/Linux; checklist manual ok |

---

## 9. Quick wins (impacto alto, custo baixo — não são features grandes)

| Quick win | O que mudar | Por que melhora | Risco | Como validar |
| --- | --- | --- | --- | --- |
| **Runa não rebaixa** ✅ | `gainWeapon` preserva o nível ao trocar de arma (nível vira potência compartilhada) | Remove a decisão negativa central do jogo | Baixo (balance) | ✅ Feito e testado |
| **Contraste do HUD** ✅ | Painel escuro atrás dos textos do HUD (`label`) | Legibilidade sobre cenários claros (vila em chamas) | Baixo | ✅ Feito |
| **Telegraph de gárgula/wyvern/harpia** ✅ | Contorno de alerta piscante antes de atirar (`enemyTelegraphFrames`) | Justiça e leitura de ameaça | Baixo | ✅ Feito e testado |
| **Brasas nas Chamas** ✅ | Burst de brasas no bocal ao disparar Chamas | Dá "peso" à arma mais fraca | Baixo | ✅ Feito (gelo já tinha rastro) |
| **Densidade do miolo** ✅ | Ondas reescritas (mais densas, curtas e sobrepostas) após playtest do usuário | Remove momentos mortos ("poucos monstros") | Médio | ✅ Feito; guardado por `TestStagesAreDense` |
| **Sombra da Bomba** ✅ | Escurecimento breve da tela ao usar (`bombDarken*`); som grave já existente | Transforma pânico em espetáculo | Baixo | ✅ Feito |
| **Apresentação do chefe** ✅ (parcial) | Aviso "ALERTA" pulsante + nome na entrada | Momento memorável | Baixo | ✅ Base feita; faixa estilizada depende de B1 (fonte) |
| **Barra de progresso legível** ✅ | Marcos dos 4 trechos na barra do topo (`drawProgressBar`) | Sensação de avanço | Baixo | ✅ Feito |
| **Popup de bônus** ✅ | Popups "+60"/"+200" para formação/trecho sem dano | Reforça recompensas ocultas | Baixo | ✅ Feito e testado |
| **Fade de música na vitória/derrota** | Já existe base; garantir corte suave | Transições menos abruptas | Baixo | Ouvir transições |

---

## 10. Roadmap por etapas

### Etapa 1 — Consolidar o MVP
- **Objetivo:** o que existe fica sólido, justo e *divertido de verdade*.
- **Entregas:** A1 (economia de armas), quick wins de game feel (HUD, telegraphs, densidade, bomba), B2 (assinatura mágica).
- **Dependências:** nenhuma.
- **Riscos:** balanceamento (mitigar com playtests curtos e seed fixa).
- **Não fazer ainda:** novas fases, novos chefes, persistência complexa.
- **Critério p/ avançar:** um playtest de ponta a ponta é claramente divertido e recompensador; power-ups são desejados, não temidos.

### Etapa 2 — Vertical slice
- **Objetivo:** uma fatia curta com qualidade próxima da final (referência de padrão).
- **Entregas:** A2 (dados de fase) recriando a Fase 1, B1 (fonte/molduras), E2 (onboarding), E3 (telas temáticas), apresentação do chefe.
- **Dependências:** Etapa 1.
- **Riscos:** escopo de arte (limitar a geométrico + tons/contorno).
- **Não fazer ainda:** 2ª/3ª fases, modos extras.
- **Critério:** dá para mostrar a fatia a terceiros e a reação é "isso parece um jogo", não "um protótipo".

### Etapa 3 — Expansão de conteúdo
- **Objetivo:** mais mundo reusando os sistemas consolidados.
- **Entregas:** C1 (novos inimigos), C3 (2ª/3ª fases em dados), C4 (ameaças terrestres), R5 + C2 (2º chefe).
- **Dependências:** A2 e o padrão de qualidade da Etapa 2.
- **Riscos:** balanceamento entre biomas; teste manual crescente.
- **Não fazer ainda:** progressão persistente pesada.
- **Critério:** 3 fases encadeadas, 2 chefes, 6–7 inimigos, curva de dificuldade coerente.

### Etapa 4 — Rejogabilidade e progressão
- **Objetivo:** motivos para novas partidas.
- **Entregas:** A3 (persistência), D1 (dificuldades), D2 (sobrevivência), D3 (conquistas).
- **Dependências:** conteúdo da Etapa 3.
- **Riscos:** inflar escopo (limitar a esses 3–4 sistemas).
- **Não fazer ainda:** roguelike/NG+/personagens.
- **Critério:** sessão típica leva a "mais uma"; score/desbloqueios persistem.

### Etapa 5 — Polimento e publicação
- **Objetivo:** 1ª versão pública.
- **Entregas:** E1 (áudio), F1 (opções), F2 (assets de loja), F3 (builds + smoke tests), A4 (gamepad, se couber).
- **Dependências:** tudo anterior.
- **Riscos:** dependência de áudio/arte externos (planejar cedo).
- **Não fazer ainda:** DLC/roadmap pós-lançamento.
- **Critério:** build estável, opções persistidas, página pronta, checklist manual 100%.

---

## 11. Vertical slice recomendada

**Recriar a Fase 1 ("O Cerco de Eldoria") como referência de qualidade final, em pequena escala.**

- **Fase escolhida:** Fase 1, trechos 1–2 (Campos → Vila em chamas) **+ mini-clímax**.
- **Duração:** ~2 minutos até um confronto.
- **Cenário:** 2 biomas com paleta contrastante (verde/azul → laranja/vermelho), parallax e narrativa ambiental (a vila queima ao fundo).
- **Inimigos:** corvo, harpia + **1 novo** (ameaça terrestre: torre/balista que atira de baixo) para provar a extensibilidade.
- **Armas:** as 3 existentes, já com **economia corrigida** e **assinatura mágica** (cor/rastro/som).
- **Power-ups:** os 5, com feedback de coleta e regra nova (sem reset).
- **Chefe:** Vharak **ou** um mini-chefe de 1–2 fases (menor, para caber na fatia) com entrada apresentada.
- **Qualidade visual:** **bitmap font + molduras**, entidades com contorno e 2–3 tons, HUD reestilizado.
- **Áudio:** 1 trilha de fase + 1 de chefe decentes e SFX revisados das 3 armas + bomba.
- **Interface:** menu inicial temático, HUD legível em ambos os biomas, tela de vitória/derrota estilizada, onboarding de 20s.
- **Critérios de conclusão:** roda 2 min sem quedas; um novato entende os controles sem ler; power-ups dão vontade de pegar; a entrada do chefe é um momento; um espectador diz "quero jogar".

---

## 12. Riscos de escopo

| Risco | Por que é perigoso | Alternativa menor |
| --- | --- | --- |
| **Sprites/animações completas** | Dependência de arte pode travar meses | Manter geométrico + contorno/tons; animar por transformações (flap/tilt já existem) |
| **Múltiplos personagens/montarias** | Multiplica arte + balanceamento | 1 personagem; variar por **dificuldade** e **arma inicial desbloqueável** |
| **Roguelike de builds/NG+** | Muda o gênero; combinatória de teste enorme | Modo Sobrevivência com escalada fixa |
| **Sistema de fases "flexível" demais** | Editor/genérico vira produto próprio | Formato de dados **mínimo** (slice de waves), só o necessário para 3 fases |
| **2º chefe totalmente novo antes de generalizar** | Duplicação e dívida | Generalizar `Boss` primeiro (R5), depois criar reusando |
| **Trilha original produzida** | Dependência de áudio externa e cara | Usar trilhas livres bem escolhidas; procedural continua de fallback |
| **Localização precoce** | Congela UI antes da hora | Centralizar strings só quando publicar em 2º idioma |
| **Gamepad + remap completos cedo** | Camada de input não urgente | Fazer a abstração leve (A4) só quando o núcleo estiver fechado |
| **Balanceamento sem ferramentas** | Ajuste no escuro | Usar **seed fixa** + dev HUD (já existem) para reproduzir e medir |
| **Conteúdo insuficiente → repetição** | 1 fase cansa | Priorizar 3 fases curtas em vez de 1 longa |

---

## 13. Recomendação final (objetiva)

1. **Estado atual:** MVP tecnicamente sólido e bem testado, com **engenharia à frente do conteúdo e do design de recompensa**. Esqueleto de qualidade, pouca carne.
2. **Principal problema hoje:** o **sistema de armas/power-ups pune o jogador** (trocar runa reseta ao nível 1), quebrando o power fantasy — o coração do gênero. Secundado pela apresentação de protótipo (fonte de debug) e pela falta de conteúdo/rejogo.
3. **Prioridade imediata:** **corrigir a economia de armas (A1)** + quick wins de game feel. Baixo custo, impacto direto na diversão.
4. **Próxima funcionalidade:** **formato de dados de fase (A2)** — destrava fases, dificuldades e sobrevivência sem duplicar código.
5. **Evitar por enquanto:** múltiplos personagens/montarias, roguelike/NG+, online, editor de fases genérico, localização e gamepad completos.
6. **Melhor vertical slice:** Fase 1 (Campos → Vila) recriada **em dados**, com fonte própria, assinatura mágica, 1 inimigo terrestre novo e um chefe apresentado — ~2 min convincentes.
7. **Escopo realista da 1ª versão pública:** **3 fases curtas, 2 chefes, 6–7 inimigos, 3 armas (corrigidas), 5 power-ups, 2–3 dificuldades, modo sobrevivência, high score persistente, fonte/UI temática e áudio decente.** Sem múltiplos personagens nem progressão persistente pesada.
8. **Os 10 próximos itens, em ordem:**
   1. ~~**A1** — Corrigir economia de armas (runa nunca rebaixa).~~ ✅
   2. ~~**Quick wins de game feel** — HUD contraste, telegraphs, bomba com sombra, popups de bônus, barra de progresso, brasas.~~ ✅ (exceto densidade do miolo, adiada para playtest)
   3. ~~**B2** — Assinatura mágica das 3 armas (cor/rastro/som).~~ ✅
   4. ~~**A2** — Formato de dados de fase; recriar a Fase 1 sobre ele.~~ ✅
   5. ~~**B1** — Bitmap font + molduras (matar a fonte de debug).~~ ✅ (fonte `basicfont` via `text/v2`; QA visual pendente)
   6. ~~**A3** — Persistência mínima (high score + opções).~~ ✅
   7. ~~**D1** — Dificuldades (Fácil/Normal/Difícil) via multiplicadores.~~ ✅
   8. ~~**C1** — 2–3 novos inimigos (incl. 1 ameaça terrestre).~~ ✅ (Balista + Feiticeiro)
   9. ~~**C3** — 2ª e 3ª fases (biomas) em dados.~~ ✅
   10. ~~**D2** — Modo Sobrevivência (reusa spawns).~~ ✅

   > **Todos os 10 itens do plano foram implementados.** Pendências principais para os próximos ciclos: 2º chefe / generalização do `Boss` (R5/C2), densidade do miolo (playtest), tela de Opções dedicada (F1), e **QA visual** da nova fonte e biomas. Ver "Estado após execução do plano" abaixo.

> Regra de ouro para os próximos ciclos: **consertar a recompensa e a apresentação antes de adicionar conteúdo; adicionar conteúdo em dados, não em código; medir com seed fixa + dev HUD.**

# Roteiro de Testes Manuais — Asas de Valdoria

Guia para validar o MVP completo do início ao fim. Marque cada item ao concluir.

## Como executar

```bash
go run .
```

Modo de desenvolvimento (HUD, hitboxes, atalhos de teste):

```bash
VALDORIA_DEV=1 go run .
# iniciar direto num trecho (0 a 3):
VALDORIA_DEV=1 VALDORIA_SECTION=2 go run .
# iniciar direto no chefe:
VALDORIA_DEV=1 VALDORIA_BOSS=1 go run .
# fixar semente para reproduzir a partida:
VALDORIA_DEV=1 VALDORIA_SEED=42 go run .
```

## Controles

| Ação | Tecla |
| --- | --- |
| Mover | WASD / setas |
| Modo de precisão | Shift |
| Atirar | Espaço |
| Mergulho do grifo | Z |
| Invocação ancestral (bomba) | X / Ctrl |
| Pausar | Esc |
| Confirmar no menu | Enter |
| Silenciar áudio | M |

Atalhos de desenvolvimento (somente com `VALDORIA_DEV=1`): `F1` HUD, `F2` hitboxes, `F3` invencível, `Tab` acelera o tempo, `K` dano ao chefe, `L` limpa inimigos, teclas de spawn de power-ups.

---

## 0. Balanceamento (automático)

Antes do roteiro manual, rode a medição — ela verifica em segundos o que levaria
uma hora de playtest para perceber:

```bash
go run ./cmd/balance
```

- [ ] Todos os critérios aparecem como `OK`.
- [ ] Na seção AMEACA, nenhum inimigo passa de 2,0s de tempo-para-matar e todos os armados disparam ao menos uma vez.
- [ ] Na seção PRESSAO, os dois modelos de jogador **morrem** antes do fim da campanha.
- [ ] Rodando com `-diff facil` e `-diff dificil`, a sobrevivência do modelo passivo muda claramente.
- [ ] Cada cenário tem uma arma dona diferente (foco/formação/coluna).
- [ ] O confronto com Vharak dura entre 60 e 180 s com **cada uma** das magias.

---

## 0.1 Medidor de Corrupção

- [ ] Deixar um corvo passar pela base enche um pouco a coluna na borda direita e toca um som grave.
- [ ] Uma gárgula saindo pela **lateral** não muda o medidor.
- [ ] Inimigos prestes a cruzar a base ganham contorno magenta e uma seta; a linha do rodapé pulsa.
- [ ] Ao cruzar 26%, aparece "O REINO SE CORROMPE / SOMBRA CRESCENTE" e a tela treme.
- [ ] Acima de 51% os inimigos nascem com halo magenta; os corvos passam a atirar uma vez ao descer.
- [ ] O cenário vai perdendo cor e ganhando veias magenta nas bordas conforme o medidor sobe.
- [ ] Com o medidor em 100%, o chefe se apresenta como "VHARAK ASCENDIDO" e a vitória diz "VALDORIA CAIU".
- [ ] Game Over e Vitória mostram a corrupção final e o número de fugas.
- [ ] Reiniciar zera o medidor.

---

## 1. Menu inicial

- [ ] O título "Asas de Valdoria" e a descrição aparecem, em **fonte própria** (não a fonte de debug).
- [ ] Navegação circula entre as opções (Iniciar, Sobrevivência, Dificuldade, Controles, Vibração, Sair); a opção selecionada fica destacada em dourado.
- [ ] "Dificuldade" alterna Fácil/Normal/Difícil e é lembrada ao reabrir o jogo.
- [ ] Os recordes de campanha e sobrevivência aparecem no menu e persistem após fechar/reabrir.
- [ ] "Iniciar" começa a partida com transição de fade.
- [ ] "Controles" abre a tela de controles e retorna ao menu com Esc/Enter.
- [ ] A opção de vibração alterna entre Completa/Reduzida/Desligada.
- [ ] "Sair" encerra o jogo sem erro no terminal.
- [ ] `M` silencia e reativa o áudio (o indicador some/volta).
- [ ] Esc no menu **não** dispara pausa.
- [ ] Navegar e confirmar no menu emitem um bipe curto de feedback.

## 2. Fase (jogabilidade)

- [ ] Movimento 8 direções; diagonal não é mais rápida que reta.
- [ ] Movimento é limitado às bordas da tela.
- [ ] Shift ativa precisão (mais lento) e mostra a marca de colisão menor.
- [ ] Espaço dispara de forma contínua, com intervalo consistente e dois projéteis.
- [ ] Coletar uma runa acende uma marca de carga no HUD; a cada três marcas o nível sobe e as marcas zeram.
- [ ] Trocar de elemento preserva nível e carga (nunca é uma perda).
- [ ] No nível 3 as marcas de carga somem e passam a cair mais curas/escudos.
- [ ] Nenhuma magia congela inimigos no lugar (a Luz causa só dano).
- [ ] `Z` lança o grifo na direção apontada, com rastro dourado, e a câmera acompanha o arranque.
- [ ] Durante o mergulho o grifo atravessa inimigos ferindo-os e **não** toma dano.
- [ ] A barra de fôlego (dourada, rodapé) esvazia ao mergulhar e não deixa emendar dois seguidos.
- [ ] Passar raspando um tiro inimigo acende um anel ciano, toca um tique e enche a barra de graze.
- [ ] Ao encher a barra de graze, aparece "INVOCACAO!" e o contador de bombas sobe (teto de 4).
- [ ] Com um gamepad conectado: analógico move, botão sul atira, leste mergulha, gatilhos precisão/bomba, Start pausa.
- [ ] Trecho 1 (Campos) é acessível: apenas corvos, bem espaçados.
- [ ] Aparece a Runa de Fogo logo no início e troca a arma ao coletar.
- [ ] Coletar a mesma runa aumenta o nível (até 3); runa diferente troca a arma **mantendo o nível** (trocar nunca rebaixa o poder — dá vontade de pegar qualquer runa).
- [ ] Cura recupera HP até o máximo; Escudo absorve um único acerto.
- [ ] O HUD mostra arma atual, nível, vidas, cargas de bomba e escudo.
- [ ] A dificuldade sobe por trecho: Vila (harpias com tiro), Muralhas (gárgulas), Castelo (wyverns).
- [ ] Avisos de trecho aparecem ("A vila esta sob ataque!", etc.) e a barra de progresso avança.
- [ ] Colisões: projétil destrói inimigo; encostar em inimigo causa dano.
- [ ] Se um inimigo atravessar a base da tela sem ser abatido, ele some sem render pontos e invalida o bônus da formação (não há perda de vida — a penalidade por fuga chega na v0.3).
- [ ] Projéteis inimigos (magenta/roxo) são claramente distinguíveis dos disparos do jogador, inclusive das chamas laranjas.
- [ ] Inimigos que atiram (harpia, gárgula, wyvern) piscam um contorno de alerta **antes** de disparar.
- [ ] O ritmo é contínuo: quase sempre há inimigos na tela, sem longas esperas.
- [ ] Os inimigos e o jogador se **destacam do cenário** (fundo escurecido + contorno escuro nas silhuetas).
- [ ] Nenhum texto do HUD se sobrepõe: placar (topo-esq.), trecho (topo-dir.), vida/vidas/bomba (rodapé-esq.), arma/escudo (rodapé-dir.); barra do chefe fina no topo não invade o placar.
- [ ] O HUD é legível mesmo sobre o cenário claro da vila em chamas (painel escuro atrás dos textos).
- [ ] A barra de progresso do topo mostra marcos dos quatro trechos.
- [ ] Bônus de formação ("+60") e de trecho sem dano ("+200") aparecem como popups.
- [ ] Disparar as Chamas gera brasas no bocal; usar a bomba escurece a tela por um instante.
- [ ] Ao abater um inimigo aparece o número de pontos ganhos subindo (já com o multiplicador).
- [ ] Abater um wyvern (inimigo grande) produz um breve "impacto" (hit-stop) com as partículas seguindo animadas.
- [ ] O escudo, ao absorver um golpe, emite som e anel de partículas (e o combo **não** zera).
- [ ] Ao tomar dano há invencibilidade curta com piscar; sem dano múltiplo no mesmo instante.
- [ ] Pontuação sobe; multiplicador de combo aumenta em sequência e zera ao tomar dano.
- [ ] Bônus por formação completa e por trecho sem dano são somados.
- [ ] Ao zerar o HP perde uma vida e reaparece com invencibilidade, **vida inicial (3, não o máximo)**, arma reduzida em 1 nível, carga zerada e projéteis inimigos limpos.
- [ ] Com 0 vidas vai para Game Over.
- [ ] Invocação ancestral (X/Ctrl): limpa projéteis inimigos, causa dano em área, deixa o jogador invulnerável e consome uma carga (HUD atualiza).
- [ ] Bomba não funciona no menu, pausa, Game Over ou vitória.

## 3. Pausa

- [ ] Esc pausa e retoma; a lógica congela e o tempo/linha do tempo não avança.
- [ ] "PAUSADO" é exibido.
- [ ] Fechar a janela continua funcionando durante a pausa.

## 3b. Campanha, dificuldades e sobrevivência

- [ ] A campanha encadeia três fases (O Cerco de Eldoria → A Floresta Corrompida → O Covil de Vharak) com biomas distintos e aviso "Nova regiao" na transição.
- [ ] A **Balista** (terrestre) aparece na fase 3 e dispara rajadas miradas telegrafadas.
- [ ] O **Feiticeiro** aparece na fase 2 e solta anéis completos de projéteis.
- [ ] Só a última fase invoca o chefe; as anteriores encadeiam direto.
- [ ] Cada dificuldade muda vidas/bombas iniciais, a resistência dos inimigos e a previsibilidade dos spawns (Fácil fixo; Normal/Difícil com X e ritmo variando).
- [ ] **Sobrevivência**: ondas infinitas com dificuldade crescente; ao morrer, o recorde de sobrevivência é salvo; nunca entra em chefe.

## 4. Chefe (Vharak)

- [ ] Após as ondas e sem inimigos restantes, o chefe entra pelo topo.
- [ ] Na entrada aparece o aviso pulsante "ALERTA" e o nome do chefe; ele é invulnerável e inimigos/projéteis residuais somem.
- [ ] Barra de vida do chefe é exibida com moldura, cor por fase e marcas nos limiares (65% e 30%).
- [ ] Fase 1 (100%–65%): bolas de fogo miradas e cone; movimento lento; padrões avisam antes de atirar.
- [ ] Fase 2 (65%–30%): movimento mais rápido, leque em arco e invocação de corvos/harpias.
- [ ] Fase 3 (<30%): ataques mais rápidos, cristais como pontos fracos e varredura com brecha.
- [ ] A varredura **sempre tem uma passagem** possível de desviar (sem dano inevitável).
- [ ] Ao mudar de fase os projéteis na tela são limpos.
- [ ] Acertar os cristais causa dano extra.
- [ ] O chefe para de atacar ao morrer e há sequência de explosões.
- [ ] Durante a morte do chefe a arena é limpa e o jogador fica invulnerável — não há como perder no instante da vitória.
- [ ] A derrota do chefe concede pontuação alta e leva à tela de Vitória.

## 5. Reinício e reset

- [ ] Game Over → "Tentar novamente" inicia partida totalmente nova.
- [ ] Game Over → "Voltar ao menu" retorna ao menu.
- [ ] Vitória → "Jogar novamente" inicia partida nova; "Voltar ao menu" retorna ao menu.
- [ ] Nova partida pelo menu começa limpa.
- [ ] Não restam inimigos, projéteis, power-ups ou chefe de partidas anteriores.
- [ ] Pontuação, vidas, arma (Lança de Luz nível 1), cargas de bomba e tempo são reiniciados.
- [ ] Teclas de confirmação não disparam múltiplas ações em um único toque.

## 6. Balanceamento (verificação geral)

- [ ] O início é acessível e ensina os controles naturalmente.
- [ ] A dificuldade cresce de forma gradual entre os trechos.
- [ ] Todos os power-ups aparecem ao menos uma vez durante a fase.
- [ ] É possível chegar ao chefe após algumas tentativas.
- [ ] O chefe é desafiador, porém vencível.
- [ ] Nenhum dano é inevitável em condições normais de jogo.
- [ ] A partida completa dura aproximadamente 3 a 5 minutos.

## 7. Áudio

- [ ] Trilhas mudam entre menu, fase e chefe.
- [ ] Efeitos tocam para tiro, inimigo destruído, dano, power-up, invocação, vitória, game over, quebra de escudo e menu.
- [ ] `M` silencia/reativa a qualquer momento.
- [ ] O jogo funciona normalmente mesmo sem arquivos em `assets/audio/`.

## Testes automatizados

```bash
go test ./...
```

Os testes cobrem regras de jogo sem renderização (jogador, armas, projéteis, inimigos, colisões, ondas, pontuação, power-ups, chefe, troca de estados e reinício de sessão). Use `VALDORIA_SEED` para reproduzir uma partida específica.

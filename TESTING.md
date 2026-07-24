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
| Invocação ancestral (bomba) | X / Ctrl |
| Pausar | Esc |
| Confirmar no menu | Enter |
| Silenciar áudio | M |

Atalhos de desenvolvimento (somente com `VALDORIA_DEV=1`): `F1` HUD, `F2` hitboxes, `F3` invencível, `Tab` acelera o tempo, `K` dano ao chefe, `L` limpa inimigos, teclas de spawn de power-ups.

---

## 1. Menu inicial

- [ ] O título "Asas de Valdoria" e a descrição aparecem.
- [ ] Navegação para cima/baixo circula entre as opções (Iniciar, Controles, Vibração, Sair).
- [ ] "Iniciar" começa a partida com transição de fade.
- [ ] "Controles" abre a tela de controles e retorna ao menu com Esc/Enter.
- [ ] A opção de vibração alterna entre Completa/Reduzida/Desligada.
- [ ] "Sair" encerra o jogo sem erro no terminal.
- [ ] `M` silencia e reativa o áudio (o indicador some/volta).
- [ ] Esc no menu **não** dispara pausa.

## 2. Fase (jogabilidade)

- [ ] Movimento 8 direções; diagonal não é mais rápida que reta.
- [ ] Movimento é limitado às bordas da tela.
- [ ] Shift ativa precisão (mais lento) e mostra a marca de colisão menor.
- [ ] Espaço dispara de forma contínua, com intervalo consistente e dois projéteis.
- [ ] Trecho 1 (Campos) é acessível: apenas corvos, bem espaçados.
- [ ] Aparece a Runa de Fogo logo no início e troca a arma ao coletar.
- [ ] Coletar a mesma runa aumenta o nível (até 3); runa diferente reinicia no nível 1.
- [ ] Cura recupera HP até o máximo; Escudo absorve um único acerto.
- [ ] O HUD mostra arma atual, nível, vidas, cargas de bomba e escudo.
- [ ] A dificuldade sobe por trecho: Vila (harpias com tiro), Muralhas (gárgulas), Castelo (wyverns).
- [ ] Avisos de trecho aparecem ("A vila esta sob ataque!", etc.) e a barra de progresso avança.
- [ ] Colisões: projétil destrói inimigo; encostar em inimigo causa dano.
- [ ] Ao tomar dano há invencibilidade curta com piscar; sem dano múltiplo no mesmo instante.
- [ ] Pontuação sobe; multiplicador de combo aumenta em sequência e zera ao tomar dano.
- [ ] Bônus por formação completa e por trecho sem dano são somados.
- [ ] Ao zerar o HP perde uma vida e reaparece com invencibilidade, arma reduzida em 1 nível e projéteis inimigos limpos.
- [ ] Com 0 vidas vai para Game Over.
- [ ] Invocação ancestral (X/Ctrl): limpa projéteis inimigos, causa dano em área, deixa o jogador invulnerável e consome uma carga (HUD atualiza).
- [ ] Bomba não funciona no menu, pausa, Game Over ou vitória.

## 3. Pausa

- [ ] Esc pausa e retoma; a lógica congela e o tempo/linha do tempo não avança.
- [ ] "PAUSADO" é exibido.
- [ ] Fechar a janela continua funcionando durante a pausa.

## 4. Chefe (Vharak)

- [ ] Após as ondas e sem inimigos restantes, o chefe entra pelo topo.
- [ ] Durante a entrada o chefe é invulnerável, o nome aparece e inimigos/projéteis residuais somem.
- [ ] Barra de vida do chefe é exibida.
- [ ] Fase 1 (100%–65%): bolas de fogo miradas e cone; movimento lento; padrões avisam antes de atirar.
- [ ] Fase 2 (65%–30%): movimento mais rápido, leque em arco e invocação de corvos/harpias.
- [ ] Fase 3 (<30%): ataques mais rápidos, cristais como pontos fracos e varredura com brecha.
- [ ] A varredura **sempre tem uma passagem** possível de desviar (sem dano inevitável).
- [ ] Ao mudar de fase os projéteis na tela são limpos.
- [ ] Acertar os cristais causa dano extra.
- [ ] O chefe para de atacar ao morrer e há sequência de explosões.
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
- [ ] Efeitos tocam para tiro, inimigo destruído, dano, power-up, invocação, vitória e game over.
- [ ] `M` silencia/reativa a qualquer momento.
- [ ] O jogo funciona normalmente mesmo sem arquivos em `assets/audio/`.

## Testes automatizados

```bash
go test ./...
```

Os testes cobrem regras de jogo sem renderização (jogador, armas, projéteis, inimigos, colisões, ondas, pontuação, power-ups, chefe, troca de estados e reinício de sessão). Use `VALDORIA_SEED` para reproduzir uma partida específica.

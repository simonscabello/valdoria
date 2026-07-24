# Áudio de Asas de Valdoria

Coloque aqui arquivos de áudio **livres de direitos autorais** para substituir os
sons gerados proceduralmente pelo jogo. Se um arquivo existir, ele é usado
automaticamente; caso contrário, o jogo gera um som simples no lugar. O jogo
funciona normalmente mesmo sem nenhum arquivo.

## Formatos aceitos

- `.ogg` (Vorbis) — preferido
- `.wav` (PCM)

A busca é feita nesta ordem (`.ogg` antes de `.wav`). Qualquer taxa de amostragem
é aceita: o áudio é reamostrado para 44100 Hz automaticamente.

## Nomes esperados (dentro de `assets/audio/`)

### Músicas (em laço contínuo)

| Arquivo          | Uso                         |
| ---------------- | --------------------------- |
| `music_menu`     | Menu inicial                |
| `music_phase`    | Sobrevivência / fallback    |
| `music_boss`     | Combate contra o chefe      |
| `music_fields`   | Campos do reino (fase 1)    |
| `music_village`  | Vila atacada (fase 1)       |
| `music_walls`    | Muralhas (fase 1)           |
| `music_castle`   | Rumo ao castelo (fase 1)    |
| `music_forest`   | Bosque sombrio (fase 2)     |
| `music_swamp`    | Pântano corrompido (fase 2) |
| `music_canyon`   | Desfiladeiro (fase 3)       |
| `music_lair`     | Covil do dragão (fase 3)    |

### Efeitos sonoros

| Arquivo        | Uso                          |
| -------------- | ---------------------------- |
| `shoot`        | Disparo (Lança de Luz)       |
| `shoot_flame`  | Disparo (Chamas do Dragão)   |
| `shoot_ice`    | Disparo (Lanças de Gelo)     |
| `enemy_down`   | Inimigo destruído            |
| `player_hit`   | Dano no jogador              |
| `pickup`       | Coleta de power-up           |
| `bomb`         | Invocação ancestral (bomba)  |
| `victory`      | Vitória                      |
| `game_over`    | Game Over                    |
| `shield_break` | Escudo absorve um golpe      |
| `menu`         | Navegação/confirmação de menu |

Exemplo: `assets/audio/shoot.ogg` ou `assets/audio/music_menu.wav`.

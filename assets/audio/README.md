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

| Arquivo         | Uso                    |
| --------------- | ---------------------- |
| `music_menu`    | Menu inicial           |
| `music_phase`   | Durante a fase         |
| `music_boss`    | Combate contra o chefe |

### Efeitos sonoros

| Arquivo        | Uso                          |
| -------------- | ---------------------------- |
| `shoot`        | Disparo do jogador           |
| `enemy_down`   | Inimigo destruído            |
| `player_hit`   | Dano no jogador              |
| `pickup`       | Coleta de power-up           |
| `bomb`         | Invocação ancestral (bomba)  |
| `victory`      | Vitória                      |
| `game_over`    | Game Over                    |
| `shield_break` | Escudo absorve um golpe      |
| `menu`         | Navegação/confirmação de menu |

Exemplo: `assets/audio/shoot.ogg` ou `assets/audio/music_menu.wav`.

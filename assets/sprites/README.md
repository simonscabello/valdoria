# Sprites de Asas de Valdoria

Coloque aqui arquivos **PNG livres de direitos** para substituir a pixel art
gerada proceduralmente pelo jogo. Se um arquivo existir, ele é usado
automaticamente; caso contrário, o jogo desenha o sprite gerado por código. O
jogo funciona normalmente mesmo sem nenhum arquivo.

## Como funciona

- O jogo procura `assets/sprites/<nome>.png` (no diretório de trabalho e ao lado
  do executável).
- O PNG deve ter **fundo transparente** e ser **pixel art na resolução nativa**
  (pequeno): o jogo escala o sprite para cobrir a hitbox da entidade, então
  desenhe no tamanho aproximado indicado abaixo, não ampliado.
- O sprite é desenhado **centrado** na entidade; partes que passam da hitbox
  (asas, por exemplo) são bem-vindas e ficam por fora da caixa de colisão.

## Nomes esperados (jogador, inimigos e chefe)

| Arquivo         | Entidade                         | Tamanho aprox. (px) |
| --------------- | -------------------------------- | ------------------- |
| `player.png`    | Cavaleiro no grifo               | 16×16               |
| `crow.png`      | Corvo corrompido                 | 12×10               |
| `harpy.png`     | Harpia (asas amplas)             | 16×14               |
| `gargoyle.png`  | Gárgula de pedra                 | 16×16               |
| `wyvern.png`    | Wyvern (grande)                  | 20×18               |
| `ballista.png`  | Balista corrompida (terrestre)   | 16×12               |
| `mage.png`      | Feiticeiro corrompido            | 14×16               |
| `boss.png`      | Vharak (dragão chefe)            | 32×16               |
| `power_light.png`  | Runa da Luz                   | 12×12               |
| `power_fire.png`   | Runa do Fogo                  | 12×12               |
| `power_ice.png`    | Runa do Gelo                  | 12×12               |
| `power_heal.png`   | Runa de Vida                  | 12×12               |
| `power_shield.png` | Escudo                        | 12×12               |

Frames opcionais de batida de asa (mesmo tamanho do base): `player_flap.png`,
`crow_flap.png`, `harpy_flap.png`, `wyvern_flap.png`, `boss_flap.png`.

Exemplo: `assets/sprites/wyvern.png`.

## Estilo

Direção atual: **arcade 16-bit vibrante** — cores saturadas, alto contraste,
contorno escuro e 3–4 tons por elemento, com um toque sombrio de corrupção
(roxos/magentas, brasas). Mantenha silhuetas legíveis em 1 frame.

## Prévia da arte procedural

Para inspecionar a pixel art gerada por código (útil ao criar substitutos):

```bash
VALDORIA_EXPORT_SPRITES=/tmp/sprites.png go test ./game/ -run TestExportSpritePreviews
```

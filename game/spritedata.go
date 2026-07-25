package game

import "image/color"

// spriteDefs guarda a pixel art procedural (primeira passada) no estilo arcade
// 16-bit vibrante. Cada sprite é uma grade ASCII: cada caractere vira uma cor da
// paleta; caracteres fora da paleta (incluindo '.') ficam transparentes.
//
// Colunas ausentes à direita são transparentes, então pequenas variações de
// largura entre linhas não deslocam o desenho — o que importa é a coluna de cada
// pixel a partir da esquerda. Para substituir por arte desenhada, basta colocar
// um PNG de mesmo nome em assets/sprites/.
var spriteDefs = map[string]spriteSource{
	"player":       {rows: playerRows, pal: playerPal},
	"player_flap":  {rows: playerFlapRows, pal: playerPal},
	"crow":         {rows: crowRows, pal: crowPal},
	"crow_flap":    {rows: crowFlapRows, pal: crowPal},
	"harpy":        {rows: harpyRows, pal: harpyPal},
	"harpy_flap":   {rows: harpyFlapRows, pal: harpyPal},
	"gargoyle":     {rows: gargoyleRows, pal: gargoylePal},
	"wyvern":       {rows: wyvernRows, pal: wyvernPal},
	"wyvern_flap":  {rows: wyvernFlapRows, pal: wyvernPal},
	"ballista":     {rows: ballistaRows, pal: ballistaPal},
	"mage":         {rows: mageRows, pal: magePal},
	"boss":         {rows: bossRows, pal: bossPal},
	"boss_flap":    {rows: bossFlapRows, pal: bossPal},
	"power_light":  {rows: powerLightRows, pal: powerLightPal},
	"power_fire":   {rows: powerFireRows, pal: powerFirePal},
	"power_ice":    {rows: powerIceRows, pal: powerIcePal},
	"power_heal":   {rows: powerHealRows, pal: powerHealPal},
	"power_shield": {rows: powerShieldRows, pal: powerShieldPal},
}

// --- Jogador: cavaleiro (armadura azul) montado no grifo (dourado) ---
var playerPal = map[rune]color.RGBA{
	'o': {0x18, 0x10, 0x0a, 0xff}, // contorno
	'p': {0xff, 0x44, 0x55, 0xff}, // plumo do elmo
	'b': {0x46, 0xc8, 0xff, 0xff}, // armadura
	'c': {0x1e, 0x74, 0xc8, 0xff}, // sombra da armadura
	'm': {0xf0, 0xbc, 0x4c, 0xff}, // corpo dourado
	'd': {0xb0, 0x82, 0x2e, 0xff}, // sombra do corpo
	'l': {0xff, 0xe8, 0x92, 0xff}, // brilho
	'a': {0x7a, 0x4c, 0x22, 0xff}, // asas
}
var playerRows = []string{
	".......ll.......",
	"......ommo......",
	"......obbo......",
	".....obccbo.....",
	"..oo.ommmmo.oo..",
	".oaaodmmmmdoaao.",
	"oaaaadmmmmdaaaao",
	"oaaaadmllmdaaaao",
	".oaaadmmmmdaaao.",
	"..oaodmmmmdoao..",
	"....ommmmmmo....",
	".....ommmmo.....",
	".....odmmdo.....",
	"......ommo......",
	"......oddo......",
	".......oo.......",
}

// Asas erguidas (batida).
var playerFlapRows = []string{
	".......ll.......",
	"..oo..ommo..oo..",
	".oaa..obbo..aao.",
	"oaaaaobccboaaaao",
	".oaa.ommmmo.aao.",
	"..aaodmmmmdoaa..",
	"...aadmmmmdaa...",
	"....admllmda....",
	"....odmmmmdo....",
	"....ommmmmmo....",
	".....ommmmo.....",
	".....ommmmo.....",
	".....odmmdo.....",
	"......ommo......",
	"......oddo......",
	".......oo.......",
}

// --- Corvo corrompido: pequeno, roxo vibrante, olhos rosados ---
var crowPal = map[rune]color.RGBA{
	'o': {0x14, 0x0a, 0x1e, 0xff},
	'm': {0x9a, 0x54, 0xd0, 0xff},
	'd': {0x5e, 0x2e, 0x92, 0xff},
	'e': {0xff, 0x50, 0x88, 0xff},
	'a': {0xff, 0xc8, 0x44, 0xff},
}
var crowRows = []string{
	"...o....o...",
	"..om....mo..",
	"..omd..dmo..",
	"...ommmmo...",
	"..omemmemo..",
	"..omdmmdmo..",
	"...ommmmo...",
	"....ommo....",
	"....oaao....",
	".....oo.....",
}

var crowFlapRows = []string{
	"o...........o",
	"om.........mo",
	".omd.....dmo.",
	"..ommmmmmmo..",
	"..omemmemo...",
	"..omdmmdmo...",
	"...ommmmo....",
	"....ommo.....",
	"....oaao.....",
	".....oo......",
}

// --- Harpia: corpo turquesa, grandes asas claras, cabelo dourado ---
var harpyPal = map[rune]color.RGBA{
	'o': {0x0e, 0x22, 0x1e, 0xff},
	'm': {0x3a, 0xd8, 0xb8, 0xff},
	'd': {0x1e, 0x94, 0x7c, 0xff},
	'e': {0xff, 0x66, 0x66, 0xff},
	'a': {0xff, 0xd8, 0x54, 0xff},
	'w': {0xd8, 0xe8, 0xf4, 0xff},
	'W': {0xa4, 0xc0, 0xdc, 0xff},
}
var harpyRows = []string{
	".......aa.......",
	"......oaao......",
	"w.....ommo.....w",
	"ww....ommo....ww",
	"www..omeemo..www",
	"wwww.ommmmo.wwww",
	"WWwwwommmmowwwWW",
	"WWWwwommmmowwWWW",
	".WWwwommmmowwWW.",
	"..wwwommmmowww..",
	"....dommmmod....",
	".....ommmmo.....",
	".....odmmdo.....",
	"......ommo......",
}

var harpyFlapRows = []string{
	"w......aa......w",
	"ww....oaao....ww",
	"www...ommo...www",
	"wwww..ommo..wwww",
	"WWwwwomeemowwwWW",
	"WWWwwommmmowwWWW",
	".WWwwommmmowwWW.",
	"..wwwommmmowww..",
	"...wwommmmoww...",
	"....dommmmod....",
	".....ommmmo.....",
	".....ommmmo.....",
	".....odmmdo.....",
	"......ommo......",
}

// --- Gárgula: pedra robusta com chifres e olhos flamejantes ---
var gargoylePal = map[rune]color.RGBA{
	'o': {0x18, 0x1a, 0x1e, 0xff},
	'm': {0x88, 0x8e, 0x9a, 0xff},
	'd': {0x56, 0x5c, 0x68, 0xff},
	'l': {0xc2, 0xc8, 0xd2, 0xff},
	'e': {0xff, 0x7a, 0x2a, 0xff},
	'w': {0x4a, 0x50, 0x5a, 0xff},
}
var gargoyleRows = []string{
	"...o........o...",
	"...o........o...",
	"...ommmmmmmmo...",
	"..oommmmmmmmoo..",
	"..ommllmmllmmo..",
	"w..ommemmemmo..w",
	"ww.ommmmmmmmo.ww",
	"w..ommddmmddo..w",
	"...ommmmmmmmo...",
	"..oommmmmmmmoo..",
	"..ommdmmmmdmmo..",
	"..oommmmmmmmoo..",
	"...ommo..ommo...",
	"...oo......oo...",
}

// --- Wyvern: dragão verde grande, asas membranosas amplas ---
var wyvernPal = map[rune]color.RGBA{
	'o': {0x10, 0x24, 0x0e, 0xff},
	'm': {0x54, 0xa8, 0x3e, 0xff},
	'd': {0x2e, 0x6e, 0x22, 0xff},
	'l': {0x86, 0xd0, 0x54, 0xff},
	'e': {0xff, 0xd0, 0x40, 0xff},
	'w': {0x3a, 0x7e, 0x2c, 0xff},
	'b': {0xd0, 0x50, 0x40, 0xff},
}
var wyvernRows = []string{
	"......oooo......",
	".....ommmmo.....",
	".....omeemo.....",
	"....odmllmdo....",
	"wwwo.ommmmo.owww",
	"owwwoommmmoowwwo",
	".owwwommmmowwwo.",
	"..owwommmmowwo..",
	"...o.ommmmo.o...",
	".....ommmmo.....",
	".....odmmdo.....",
	".....ommmmo.....",
	"......obbo......",
	"......ommo......",
	".......oo......",
	".......o.......",
}

var wyvernFlapRows = []string{
	"ww....oooo....ww",
	"oww..ommmmo..wwo",
	".oww.omeemo.wwo.",
	"..owodmllmdowo..",
	"...woommmmoow...",
	"....oommmmoo....",
	".....ommmmo.....",
	".....ommmmo.....",
	".....ommmmo.....",
	".....ommmmo.....",
	".....odmmdo.....",
	".....ommmmo.....",
	"......obbo......",
	"......ommo......",
	".......oo......",
	".......o.......",
}

// --- Balista corrompida: máquina de cerco de madeira com virote metálico ---
var ballistaPal = map[rune]color.RGBA{
	'o': {0x1e, 0x14, 0x0a, 0xff},
	'm': {0x84, 0x60, 0x38, 0xff},
	'd': {0x52, 0x3a, 0x20, 0xff},
	'l': {0xb0, 0x88, 0x54, 0xff},
	's': {0xb8, 0xbe, 0xc8, 0xff},
}
var ballistaRows = []string{
	".......ss.......",
	"o......ss......o",
	"moo....ss....oom",
	".mmoooossoooomm.",
	"..mmmmmssmmmmm..",
	"..omdllsslldmo..",
	"..ommmmssmmmmo..",
	"...ommmssmmmo...",
	"....ommssmmo....",
	"...ommo..ommo...",
	"...oo......oo...",
}

// --- Feiticeiro corrompido: manto roxo, capuz e orbe mágico brilhante ---
var magePal = map[rune]color.RGBA{
	'o': {0x14, 0x0a, 0x20, 0xff},
	'm': {0x6a, 0x30, 0xa0, 0xff},
	'd': {0x40, 0x18, 0x6e, 0xff},
	'l': {0x9a, 0x60, 0xd8, 0xff},
	'h': {0x2a, 0x10, 0x44, 0xff},
	'e': {0xff, 0xe0, 0x60, 0xff},
	'g': {0xc0, 0xf0, 0xff, 0xff},
}
var mageRows = []string{
	"....oooo....",
	"...ohhhho...",
	"..ohheehho..",
	"..ohhhhhho..",
	".ommllllmmo.",
	".ommmmmmmmo.",
	"ommdllllmmmo",
	"ommmmmmmmmmo",
	"ommdllllmmmo",
	".ommmmmmmmo.",
	".ommmmmmmmo.",
	"..ommmmmmo..",
	"...oggggo...",
	"...gggggg...",
	"...oggggo...",
	"....oggo....",
}

// --- Vharak: dragão chefe, carmesim com asas amplas e olhos de brasa ---
var bossPal = map[rune]color.RGBA{
	'o': {0x1a, 0x08, 0x0a, 0xff}, // contorno
	'm': {0xc8, 0x2e, 0x3a, 0xff}, // corpo
	'd': {0x7a, 0x18, 0x22, 0xff}, // sombra
	'l': {0xf0, 0x6a, 0x58, 0xff}, // brilho
	'e': {0xff, 0xd0, 0x40, 0xff}, // olhos
	'h': {0xff, 0x88, 0x30, 0xff}, // chifres / brasas
	'w': {0x8a, 0x1e, 0x28, 0xff}, // asas
	'W': {0x5a, 0x12, 0x18, 0xff}, // sombra das asas
	'b': {0xff, 0x50, 0x20, 0xff}, // fogo na boca
}
var bossRows = []string{
	"....h......h....",
	"...oh......ho...",
	"...ohmmmmmmho...",
	"..oommeeemmoo...",
	"Wwoo.omllmo.oowW",
	"WWww.ommmmo.wwWW",
	"wWWWwommmmowwWWw",
	"wwWwwommmmowwWww",
	".wwwwomlllmo.www",
	"..ww.ommmmo.ww..",
	"...o.oddddo.o...",
	".....ommmmo.....",
	".....odbbdo.....",
	"......obbo......",
	"......oddo......",
	".......oo.......",
}

var bossFlapRows = []string{
	"W...h......h...W",
	"Ww.oh......ho.wW",
	"wWwohmmmmmmhowWw",
	"wwwoommeeemmooww",
	".ww.oomllmoo.ww.",
	"..ww.ommmmo.ww..",
	"...wwommmmoww...",
	"....wommmmow....",
	".....omlllmo....",
	".....ommmmo.....",
	".....oddddo.....",
	".....ommmmo.....",
	".....odbbdo.....",
	"......obbo......",
	"......oddo......",
	".......oo.......",
}

// --- Runas / power-ups (12×12, silhuetas distintas) ---

// Runa da Luz: estrela/lança dourada.
var powerLightPal = map[rune]color.RGBA{
	'o': {0x2a, 0x20, 0x0a, 0xff},
	'm': {0xff, 0xf0, 0x8c, 0xff},
	'd': {0xd0, 0xb0, 0x40, 0xff},
	'l': {0xff, 0xff, 0xd0, 0xff},
}
var powerLightRows = []string{
	".....oo.....",
	"....ommo....",
	"...olmmlo...",
	"..oommmmoo..",
	".olmmmmmmlo.",
	"oommmlmmmmoo",
	".olmmmmmmlo.",
	"..oommmmoo..",
	"...olmmlo...",
	"....ommo....",
	".....oo.....",
	"............",
}

// Runa do Fogo: chama laranja.
var powerFirePal = map[rune]color.RGBA{
	'o': {0x2a, 0x10, 0x08, 0xff},
	'm': {0xff, 0x7a, 0x2a, 0xff},
	'd': {0xc0, 0x40, 0x18, 0xff},
	'l': {0xff, 0xd0, 0x60, 0xff},
	'y': {0xff, 0xf0, 0xa0, 0xff},
}
var powerFireRows = []string{
	"......o.....",
	".....omo....",
	"....omlmo...",
	"...omylmo...",
	"..ommyymmo..",
	".ommlmllmmo.",
	".ommddddmmo.",
	"..ommmmmmo..",
	"...ommmmo...",
	"....ommo....",
	".....oo.....",
	"............",
}

// Runa do Gelo: cristal ciano.
var powerIcePal = map[rune]color.RGBA{
	'o': {0x0a, 0x18, 0x28, 0xff},
	'm': {0x6a, 0xd0, 0xff, 0xff},
	'd': {0x2a, 0x80, 0xc0, 0xff},
	'l': {0xd0, 0xf4, 0xff, 0xff},
}
var powerIceRows = []string{
	".....oo.....",
	"....omlo....",
	"...ommmlo...",
	"..omlmlmmo..",
	".ommmmmmmo..",
	"omldmmmdlmo.",
	".ommmmmmmo..",
	"..omlmlmmo..",
	"...ommmlo...",
	"....omlo....",
	".....oo.....",
	"............",
}

// Runa de Vida: coração/poção verde.
var powerHealPal = map[rune]color.RGBA{
	'o': {0x0a, 0x20, 0x12, 0xff},
	'm': {0x4a, 0xe0, 0x6a, 0xff},
	'd': {0x28, 0x90, 0x40, 0xff},
	'l': {0xb0, 0xff, 0xc0, 0xff},
}
var powerHealRows = []string{
	"...oo..oo...",
	"..ommoommo..",
	".omlmmmmlmo.",
	"ommmmmmmmmo.",
	"ommmllmmmmo.",
	".ommmmmmmo..",
	"..ommmmmmo..",
	"...ommmmo...",
	"....ommo....",
	".....oo.....",
	"............",
	"............",
}

// Escudo: broquel prateado com brilho.
var powerShieldPal = map[rune]color.RGBA{
	'o': {0x14, 0x14, 0x1c, 0xff},
	'm': {0xc8, 0xd0, 0xe0, 0xff},
	'd': {0x70, 0x78, 0x90, 0xff},
	'l': {0xf0, 0xf4, 0xff, 0xff},
	'b': {0x50, 0xa0, 0xe0, 0xff},
}
var powerShieldRows = []string{
	"...oooooo...",
	"..ommmmmmo..",
	".omllbbbllmo",
	"ommmbbbbmmmo",
	"ommbbbbbbmmo",
	"ommbbbbbbmmo",
	".ommbbbbmmo.",
	"..ommmmmmo..",
	"...omddmo...",
	"....ommo....",
	".....oo.....",
	"............",
}

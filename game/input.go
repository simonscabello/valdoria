package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Camada de entrada por ações.
//
// O jogo lia teclas direto de `ebiten.IsKeyPressed` espalhado por player.go e
// game.go. Isso funcionava para teclado e bloqueava tudo o mais: sem gamepad,
// sem Steam Deck, sem remapeamento — e sem um lugar onde acrescentar um botão
// novo (o Mergulho) sem tocar em cinco arquivos.
//
// Aqui cada **ação** do jogo tem uma lista de teclas e uma lista de botões
// padrão de gamepad. O resto do código pergunta pela ação, nunca pela tecla.
//
// As teclas de desenvolvimento continuam lidas direto em dev.go: elas não fazem
// parte do jogo e não devem poluir o mapa de ações.

type inputAction int

const (
	actLeft inputAction = iota
	actRight
	actUp
	actDown
	actFire
	actPrecision
	actBomb
	actDive
	actPause
	actConfirm
	actMute
	actCount
)

// binding reúne as formas de acionar uma ação. Um gamepad padrão cobre todas as
// ações do jogo; o teclado mantém os atalhos históricos.
type binding struct {
	keys    []ebiten.Key
	buttons []ebiten.StandardGamepadButton
}

var bindings = [actCount]binding{
	actLeft: {
		keys:    []ebiten.Key{ebiten.KeyA, ebiten.KeyArrowLeft},
		buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonLeftLeft},
	},
	actRight: {
		keys:    []ebiten.Key{ebiten.KeyD, ebiten.KeyArrowRight},
		buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonLeftRight},
	},
	actUp: {
		keys:    []ebiten.Key{ebiten.KeyW, ebiten.KeyArrowUp},
		buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonLeftTop},
	},
	actDown: {
		keys:    []ebiten.Key{ebiten.KeyS, ebiten.KeyArrowDown},
		buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonLeftBottom},
	},
	actFire: {
		keys:    []ebiten.Key{ebiten.KeySpace},
		buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonRightBottom},
	},
	actPrecision: {
		keys:    []ebiten.Key{ebiten.KeyShift},
		buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonFrontBottomLeft},
	},
	actBomb: {
		keys:    []ebiten.Key{ebiten.KeyX, ebiten.KeyControl},
		buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonFrontBottomRight},
	},
	actDive: {
		keys:    []ebiten.Key{ebiten.KeyZ, ebiten.KeyShiftRight},
		buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonRightRight},
	},
	actPause: {
		keys:    []ebiten.Key{ebiten.KeyEscape},
		buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonCenterRight},
	},
	actConfirm: {
		keys:    []ebiten.Key{ebiten.KeyEnter},
		buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonRightBottom},
	},
	actMute: {
		keys:    []ebiten.Key{ebiten.KeyM},
		buttons: []ebiten.StandardGamepadButton{ebiten.StandardGamepadButtonCenterLeft},
	},
}

// axisDeadzone ignora o repouso impreciso do analógico. Abaixo disso o eixo é
// tratado como parado — sem isso o grifo escorrega sozinho.
const axisDeadzone = 0.35

// gamepadIDs é reaproveitado a cada frame para não alocar no game loop.
var gamepadIDs []ebiten.GamepadID

// activeGamepad devolve o primeiro gamepad com layout padrão disponível.
// Sem layout padrão não dá para mapear botões de forma portátil, então o
// controle é ignorado em vez de gerar entradas erradas.
func activeGamepad() (ebiten.GamepadID, bool) {
	gamepadIDs = ebiten.AppendGamepadIDs(gamepadIDs[:0])
	for _, id := range gamepadIDs {
		if ebiten.IsStandardGamepadLayoutAvailable(id) {
			return id, true
		}
	}
	return 0, false
}

// actionPressed diz se a ação está sendo mantida, por teclado ou gamepad.
func actionPressed(a inputAction) bool {
	if a < 0 || a >= actCount {
		return false
	}
	b := bindings[a]
	for _, k := range b.keys {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}
	id, ok := activeGamepad()
	if !ok {
		return false
	}
	for _, btn := range b.buttons {
		if ebiten.IsStandardGamepadButtonPressed(id, btn) {
			return true
		}
	}
	return analogPressed(id, a)
}

// actionJustPressed dispara uma única vez na borda de pressionamento — é o que
// evita menus avançando várias linhas com um toque.
func actionJustPressed(a inputAction) bool {
	if a < 0 || a >= actCount {
		return false
	}
	b := bindings[a]
	for _, k := range b.keys {
		if inpututil.IsKeyJustPressed(k) {
			return true
		}
	}
	id, ok := activeGamepad()
	if !ok {
		return false
	}
	for _, btn := range b.buttons {
		if inpututil.IsStandardGamepadButtonJustPressed(id, btn) {
			return true
		}
	}
	return analogJustPressed(id, a)
}

// analogState guarda a leitura do analógico do frame anterior, para derivar
// bordas de pressionamento a partir de um eixo contínuo (navegação de menu).
var analogState [actCount]bool

func analogPressed(id ebiten.GamepadID, a inputAction) bool {
	x := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
	y := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)
	switch a {
	case actLeft:
		return x < -axisDeadzone
	case actRight:
		return x > axisDeadzone
	case actUp:
		return y < -axisDeadzone
	case actDown:
		return y > axisDeadzone
	}
	return false
}

func analogJustPressed(id ebiten.GamepadID, a inputAction) bool {
	now := analogPressed(id, a)
	was := analogState[a]
	analogState[a] = now
	return now && !was
}

// resetInputState limpa o estado derivado do analógico. Usado nos testes para
// não vazar leitura entre casos.
func resetInputState() {
	for i := range analogState {
		analogState[i] = false
	}
}

// moveVector devolve a direção do movimento já normalizada, combinando teclado
// e analógico. O analógico tem prioridade quando está fora da zona morta, para
// permitir movimento em ângulo — mas o teclado nunca deixa de funcionar.
func moveVector() (float64, float64) {
	if id, ok := activeGamepad(); ok {
		x := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
		y := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)
		if x*x+y*y > axisDeadzone*axisDeadzone {
			return clampUnit(x), clampUnit(y)
		}
	}
	var dx, dy float64
	if actionPressed(actLeft) {
		dx--
	}
	if actionPressed(actRight) {
		dx++
	}
	if actionPressed(actUp) {
		dy--
	}
	if actionPressed(actDown) {
		dy++
	}
	return normalizeDiagonal(dx, dy)
}

func clampUnit(v float64) float64 {
	if v > 1 {
		return 1
	}
	if v < -1 {
		return -1
	}
	return v
}

// GamepadConnected informa se há um controle utilizável, para a interface poder
// mostrar os rótulos certos.
func GamepadConnected() bool {
	_, ok := activeGamepad()
	return ok
}

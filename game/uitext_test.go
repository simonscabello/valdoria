package game

import "testing"

// Garante que textos centrais do menu/HUD cabem na largura útil da tela
// (240px). Foi exatamente isso que quebrou o menu com a moldura.
func TestKeyUIStringsFitScreen(t *testing.T) {
	const pad = 24
	max := float64(ScreenWidth - pad)
	cases := []string{
		menuTitle,
		menuSubtitle,
		menuHelp,
		"Recorde 999999  Sobrev. 999999",
		"> Dificuldade: Dificil",
		"VHARAK, O DRAGAO CORROMPIDO",
		"FASE CONCLUIDA!",
		"X/Ctrl: invocacao",
	}
	for _, s := range cases {
		if w := textWidth(s); w > max {
			t.Errorf("%q tem largura %.0f > max %.0f (ScreenWidth=%d)", s, w, max, ScreenWidth)
		}
	}
}

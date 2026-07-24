package game

import "testing"

func TestSpriteDefsRenderNonEmpty(t *testing.T) {
	for name, def := range spriteDefs {
		img := def.render()
		b := img.Bounds()
		if b.Dx() == 0 || b.Dy() == 0 {
			t.Errorf("sprite %q renderizou vazio (%dx%d)", name, b.Dx(), b.Dy())
			continue
		}
		opaque := 0
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				if _, _, _, a := img.At(x, y).RGBA(); a > 0 {
					opaque++
				}
			}
		}
		if opaque == 0 {
			t.Errorf("sprite %q não tem nenhum pixel visível", name)
		}
	}
}

func TestAllEntitySpritesExist(t *testing.T) {
	names := []string{
		"player", "player_flap",
		"crow", "crow_flap",
		"harpy", "harpy_flap",
		"gargoyle",
		"wyvern", "wyvern_flap",
		"ballista", "mage",
		"boss", "boss_flap",
		"power_light", "power_fire", "power_ice", "power_heal", "power_shield",
	}
	for _, n := range names {
		if _, ok := spriteDefs[n]; !ok {
			t.Errorf("faltou o sprite procedural %q", n)
		}
	}
}

func TestWingFrameNameAlternates(t *testing.T) {
	if got := wingFrameName("player", 0); got != "player" {
		t.Fatalf("sin>=0 deveria ser player, veio %q", got)
	}
	if got := wingFrameName("player", 3.5); got != "player_flap" {
		t.Fatalf("sin<0 deveria ser player_flap, veio %q", got)
	}
	if got := wingFrameName("gargoyle", 3.5); got != "gargoyle" {
		t.Fatalf("sem flap deveria permanecer base, veio %q", got)
	}
}

func TestWhitenImageKeepsAlpha(t *testing.T) {
	src := spriteDefs["crow"].render()
	white := whitenImage(src)
	b := src.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, sa := src.At(x, y).RGBA()
			wr, wg, wb, wa := white.At(x, y).RGBA()
			if (sa > 0) != (wa > 0) {
				t.Fatalf("alfa divergente em (%d,%d)", x, y)
			}
			if wa > 0 && (wr>>8 != 0xff || wg>>8 != 0xff || wb>>8 != 0xff) {
				t.Fatalf("pixel visível deveria ser branco em (%d,%d)", x, y)
			}
		}
	}
}

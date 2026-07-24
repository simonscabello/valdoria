package game

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"
)

// TestExportSpritePreviews renderiza os sprites procedurais ampliados em um único
// PNG, para inspeção visual da pixel art. Só roda quando
// VALDORIA_EXPORT_SPRITES aponta para o caminho de saída do PNG.
func TestExportSpritePreviews(t *testing.T) {
	out := os.Getenv("VALDORIA_EXPORT_SPRITES")
	if out == "" {
		t.Skip("defina VALDORIA_EXPORT_SPRITES=/caminho.png para exportar")
	}

	const scale = 9
	const pad = 8
	bg := color.RGBA{0x14, 0x12, 0x20, 0xff} // fundo escuro como o do jogo

	// Alvo em pixels de jogo (cobre a hitbox com folga, como no jogo real).
	fit := func(name string, r *image.RGBA) *image.RGBA {
		base := strings.TrimSuffix(name, "_flap")
		tw, th := r.Bounds().Dx(), r.Bounds().Dy()
		var hw, hh float64
		switch base {
		case "crow":
			hw, hh = crowSize, crowSize
		case "harpy":
			hw, hh = harpySize, harpySize
		case "gargoyle":
			hw, hh = gargoyleSize, gargoyleSize
		case "wyvern":
			hw, hh = wyvernSize, wyvernSize
		case "ballista":
			hw, hh = ballistaSize, ballistaSize*0.7
		case "mage":
			hw, hh = mageSize, mageSize
		case "boss":
			hw, hh = bossW, bossH
		case "power_light", "power_fire", "power_ice", "power_heal", "power_shield":
			hw, hh = powerupSize*1.25, powerupSize*1.25
		default:
			return r // player: escala 1
		}
		s := 1.1 * mathMax(hw/float64(tw), hh/float64(th))
		return nearestScale(r, s)
	}

	bases := []string{
		"player", "crow", "harpy", "gargoyle", "wyvern", "ballista", "mage", "boss",
		"power_light", "power_fire", "power_ice", "power_heal", "power_shield",
	}
	names := make([]string, 0, len(bases)*2)
	for _, n := range bases {
		names = append(names, n)
		if _, ok := spriteDefs[n+"_flap"]; ok {
			names = append(names, n+"_flap")
		}
	}

	totalW, maxH := pad, 0
	imgs := map[string]*image.RGBA{}
	for _, n := range names {
		r := fit(n, spriteDefs[n].render())
		imgs[n] = r
		totalW += r.Bounds().Dx()*scale + pad
		if h := r.Bounds().Dy() * scale; h > maxH {
			maxH = h
		}
	}
	sheet := image.NewRGBA(image.Rect(0, 0, totalW, maxH+2*pad))
	for y := 0; y < sheet.Bounds().Dy(); y++ {
		for x := 0; x < sheet.Bounds().Dx(); x++ {
			sheet.SetRGBA(x, y, bg)
		}
	}

	cursor := pad
	for _, n := range names {
		r := imgs[n]
		w, h := r.Bounds().Dx(), r.Bounds().Dy()
		oy := pad + (maxH - h*scale)
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				c := r.RGBAAt(x, y)
				if c.A == 0 {
					continue
				}
				for sy := 0; sy < scale; sy++ {
					for sx := 0; sx < scale; sx++ {
						sheet.SetRGBA(cursor+x*scale+sx, oy+y*scale+sy, c)
					}
				}
			}
		}
		cursor += w*scale + pad
	}

	f, err := os.Create(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, sheet); err != nil {
		t.Fatal(err)
	}
	t.Logf("sprites exportados para %s (ordem: %v)", out, names)
}

func mathMax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// nearestScale amplia/reduz por vizinho mais próximo (mantém a pixel art nítida).
func nearestScale(src *image.RGBA, s float64) *image.RGBA {
	w := int(float64(src.Bounds().Dx()) * s)
	h := int(float64(src.Bounds().Dy()) * s)
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		sy := int(float64(y) / s)
		for x := 0; x < w; x++ {
			sx := int(float64(x) / s)
			out.SetRGBA(x, y, src.RGBAAt(sx, sy))
		}
	}
	return out
}

package game

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
)

// Sistema de sprites com fallback, na mesma filosofia do áudio: se existir um
// arquivo assets/sprites/<nome>.png ele é usado; caso contrário, o sprite é
// gerado proceduralmente a partir de uma grade ASCII (pixel art em código).
//
// As imagens do Ebiten só podem ser criadas com o contexto gráfico ativo, então
// o cache é preenchido de forma preguiçosa no primeiro uso (sempre em Draw).
const spriteDir = "assets/sprites"

// spriteSource descreve um sprite como linhas de texto + uma paleta que mapeia
// cada caractere para uma cor (caracteres fora da paleta ficam transparentes).
type spriteSource struct {
	rows []string
	pal  map[rune]color.RGBA
}

// render converte a grade ASCII em uma imagem RGBA (puro, sem contexto gráfico
// — portanto testável).
func (s spriteSource) render() *image.RGBA {
	h := len(s.rows)
	w := 0
	for _, r := range s.rows {
		if len(r) > w {
			w = len(r)
		}
	}
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y, row := range s.rows {
		for x, ch := range row {
			if c, ok := s.pal[ch]; ok {
				img.SetRGBA(x, y, c)
			}
		}
	}
	return img
}

var (
	spriteCache      = map[string]*ebiten.Image{}
	spriteFlashCache = map[string]*ebiten.Image{}
)

// sprite devolve a imagem pronta do Ebiten para o nome dado, construindo-a no
// primeiro uso. Só deve ser chamada em Draw. Devolve nil se não houver PNG nem
// gerador (o chamador então cai no desenho geométrico de reserva).
func sprite(name string) *ebiten.Image {
	if img, ok := spriteCache[name]; ok {
		return img
	}
	src := spriteImage(name)
	var img *ebiten.Image
	if src != nil {
		img = ebiten.NewImageFromImage(src)
	}
	spriteCache[name] = img
	return img
}

// spriteImage devolve a fonte de pixels (PNG livre ou grade procedural).
func spriteImage(name string) image.Image {
	if img := loadSpritePNG(name); img != nil {
		return img
	}
	if def, ok := spriteDefs[name]; ok {
		return def.render()
	}
	return nil
}

// spriteFlash devolve a silhueta branca do sprite (para o flash de dano),
// preservando o formato e a transparência.
func spriteFlash(name string) *ebiten.Image {
	if img, ok := spriteFlashCache[name]; ok {
		return img
	}
	src := spriteImage(name)
	var img *ebiten.Image
	if src != nil {
		img = ebiten.NewImageFromImage(whitenImage(src))
	}
	spriteFlashCache[name] = img
	return img
}

// whitenImage mantém o canal alfa e pinta de branco todo pixel visível.
func whitenImage(src image.Image) *image.RGBA {
	b := src.Bounds()
	out := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := src.At(x, y).RGBA()
			if a > 0 {
				out.SetRGBA(x, y, color.RGBA{0xff, 0xff, 0xff, uint8(a >> 8)})
			}
		}
	}
	return out
}

func loadSpritePNG(name string) image.Image {
	for _, dir := range spriteDirs() {
		f, err := os.Open(filepath.Join(dir, name+".png"))
		if err != nil {
			continue
		}
		img, err := png.Decode(f)
		f.Close()
		if err == nil {
			return img
		}
	}
	return nil
}

func spriteDirs() []string {
	dirs := []string{spriteDir}
	if exe, err := os.Executable(); err == nil {
		dirs = append(dirs, filepath.Join(filepath.Dir(exe), spriteDir))
	}
	return dirs
}

// spriteBounds devolve as dimensões do sprite (Draw-only). ok=false se não houver.
func spriteBounds(name string) (w, h int, ok bool) {
	img := sprite(name)
	if img == nil {
		return 0, 0, false
	}
	b := img.Bounds()
	return b.Dx(), b.Dy(), true
}

// wingFrameName escolhe o frame de asa (base ou base+"_flap") pela fase.
// Se não existir o frame alternativo, devolve o nome base.
func wingFrameName(base string, phase float64) string {
	if math.Sin(phase) >= 0 {
		return base
	}
	flap := base + "_flap"
	if _, ok := spriteDefs[flap]; ok {
		return flap
	}
	if loadSpritePNG(flap) != nil {
		return flap
	}
	return base
}

// drawSprite desenha o sprite centrado em (cx, cy), com escala, espelhamento
// horizontal e rotação opcionais. Se flash for verdadeiro, usa a silhueta branca.
func drawSprite(dst *ebiten.Image, name string, cx, cy, scale float64, flipX bool, rot float64, flash bool) bool {
	return drawSpriteTinted(dst, name, cx, cy, scale, flipX, rot, flash, nil)
}

// drawSpriteTinted é a versão com tingimento: multiplica o sprite por uma cor,
// usada para as variantes corrompidas reaproveitarem a mesma pixel art com uma
// identidade cromática própria (sem custo de arte nova).
func drawSpriteTinted(dst *ebiten.Image, name string, cx, cy, scale float64, flipX bool, rot float64, flash bool, tint *color.RGBA) bool {
	img := sprite(name)
	if img == nil {
		return false
	}
	if flash {
		if f := spriteFlash(name); f != nil {
			img = f
		}
	}
	w := float64(img.Bounds().Dx())
	h := float64(img.Bounds().Dy())
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-w/2, -h/2)
	sx := scale
	if flipX {
		sx = -scale
	}
	op.GeoM.Scale(sx, scale)
	if rot != 0 {
		op.GeoM.Rotate(rot)
	}
	op.GeoM.Translate(cx, cy)
	if tint != nil && !flash {
		op.ColorScale.Scale(
			float32(tint.R)/0xff, float32(tint.G)/0xff, float32(tint.B)/0xff, 1)
	}
	dst.DrawImage(img, op)
	return true
}

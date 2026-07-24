package game

import (
	"image/color"
	"testing"
)

func TestBrightenClamps(t *testing.T) {
	c := brighten(color.RGBA{0xf0, 0x10, 0xff, 0xff}, 90)
	if c.R != 0xff || c.G != 0x6a || c.B != 0xff {
		t.Fatalf("brighten inesperado: %#v", c)
	}
}

func TestClampByte(t *testing.T) {
	if clampByte(-1) != 0 || clampByte(300) != 255 || clampByte(40) != 40 {
		t.Fatal("clampByte fora do esperado")
	}
}

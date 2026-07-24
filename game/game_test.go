package game

import "testing"

func TestCollides(t *testing.T) {
	cases := []struct {
		name string
		ax, ay, aw, ah,
		bx, by, bw, bh float64
		want bool
	}{
		{"sobrepostos", 0, 0, 10, 10, 5, 5, 10, 10, true},
		{"encostados sem sobrepor", 0, 0, 10, 10, 10, 0, 10, 10, false},
		{"separados", 0, 0, 5, 5, 100, 100, 5, 5, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := collides(c.ax, c.ay, c.aw, c.ah, c.bx, c.by, c.bw, c.bh)
			if got != c.want {
				t.Errorf("collides = %v, quero %v", got, c.want)
			}
		})
	}
}

func TestFilterAliveRemovesAndNilsTail(t *testing.T) {
	items := []*Bullet{{}, {dead: true}, {}}
	items[1].dead = true

	filtered := filterAlive(items, func(b *Bullet) bool { return !b.dead })

	if len(filtered) != 2 {
		t.Fatalf("len = %d, quero 2", len(filtered))
	}
	if items[len(filtered)] != nil {
		t.Error("posição remanescente deveria ser nil para liberar o GC")
	}
}

func TestItoa(t *testing.T) {
	cases := map[int]string{0: "0", 7: "7", 42: "42", -13: "-13"}
	for in, want := range cases {
		if got := itoa(in); got != want {
			t.Errorf("itoa(%d) = %q, quero %q", in, got, want)
		}
	}
}

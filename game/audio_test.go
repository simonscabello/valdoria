package game

import (
	"os"
	"testing"
)

// TestMain desliga o áudio para os testes não abrirem dispositivo de som.
func TestMain(m *testing.M) {
	audioEnabled = false
	os.Exit(m.Run())
}

func TestFrameGatePreventsRepeat(t *testing.T) {
	a := newAudioManager()

	a.playSFX(sfxShoot)
	if !a.frameGate[sfxShoot] {
		t.Fatal("o primeiro disparo deveria marcar o gate do frame")
	}
	a.update()
	if a.frameGate[sfxShoot] {
		t.Error("o gate deveria ser limpo a cada frame")
	}
}

func TestMuteToggle(t *testing.T) {
	a := newAudioManager()
	if a.muted {
		t.Fatal("o áudio não deveria iniciar mudo")
	}
	a.toggleMute()
	if !a.muted {
		t.Error("a tecla de mudo deveria silenciar")
	}
	if a.effectiveMusicVolume() != 0 {
		t.Error("mudo deveria zerar o volume de música")
	}
}

func TestMusicTransitionSwitchesTrack(t *testing.T) {
	a := newAudioManager()

	a.playMusic(musicMenu)
	a.update()
	if a.current != musicMenu {
		t.Fatalf("deveria estar tocando o menu, current %d", a.current)
	}

	a.playMusic(musicBoss)
	for i := 0; i < 40 && a.current != musicBoss; i++ {
		a.update()
	}
	if a.current != musicBoss {
		t.Errorf("a troca deveria chegar à música do chefe, current %d", a.current)
	}
}

func TestVolumeClamped(t *testing.T) {
	a := newAudioManager()
	a.setMasterVolume(2)
	a.setSFXVolume(-1)
	if a.master != 1 {
		t.Errorf("volume geral deveria ser limitado a 1, foi %v", a.master)
	}
	if a.sfxVol != 0 {
		t.Errorf("volume de efeitos deveria ser limitado a 0, foi %v", a.sfxVol)
	}
}

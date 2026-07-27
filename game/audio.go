package game

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

// Diretório e formatos aceitos para trilhas/efeitos livres opcionais.
// Se um arquivo existir, ele substitui o som gerado proceduralmente.
const (
	audioDir        = "assets/audio"
	audioSampleRate = 44100
	musicFadeStep   = 0.05 // ~20 frames para uma troca suave de música

	defaultMasterVolume = 0.8
	defaultMusicVolume  = 0.5
	defaultSFXVolume    = 0.7
)

// audioEnabled permite desligar o áudio (usado nos testes para não abrir
// dispositivo de som). sharedContext e sharedAudio garantem um único contexto.
var (
	audioEnabled  = true
	sharedContext *audio.Context
	sharedAudio   *AudioManager
)

type soundID int

const (
	sfxShoot soundID = iota
	sfxShootFlame
	sfxShootIce
	sfxEnemyDown
	sfxPlayerHit
	sfxPickup
	sfxBomb
	sfxVictory
	sfxGameOver
	sfxShield
	sfxMenu
	sfxEscape
	sfxCorruption
	sfxDive
	sfxGraze
	sfxCount
)

type musicID int

const (
	musicNone musicID = iota
	musicMenu
	musicPhase // sobrevivência / fallback
	musicBoss
	musicFields
	musicVillage
	musicWalls
	musicCastle
	musicForest
	musicSwamp
	musicCanyon
	musicLair
)

const volumeStep = 0.1

// AudioManager centraliza reprodução de música e efeitos, volumes e mudo.
type AudioManager struct {
	ctx        *audio.Context
	sampleRate int

	master   float64
	musicVol float64
	sfxVol   float64
	muted    bool

	sfx       map[soundID]*audio.Player
	frameGate [sfxCount]bool // evita repetir o mesmo som várias vezes no frame

	music     map[musicID]*audio.Player
	current   musicID
	target    musicID
	musicGain float64
}

func getAudio() *AudioManager {
	if sharedAudio == nil {
		sharedAudio = newAudioManager()
	}
	return sharedAudio
}

func newAudioManager() *AudioManager {
	a := &AudioManager{
		sampleRate: audioSampleRate,
		master:     defaultMasterVolume,
		musicVol:   defaultMusicVolume,
		sfxVol:     defaultSFXVolume,
		sfx:        map[soundID]*audio.Player{},
		music:      map[musicID]*audio.Player{},
		current:    musicNone,
		target:     musicNone,
	}
	if !audioEnabled {
		return a
	}
	if !audioAvailable() {
		log.Println("audio desativado: nenhum dispositivo de som disponivel")
		return a
	}
	if sharedContext == nil {
		sharedContext = audio.NewContext(audioSampleRate)
	}
	a.ctx = sharedContext
	a.buildLibrary()
	return a
}

// audioAvailable evita abrir o dispositivo de som quando ele não existe (por
// exemplo, WSL2 ou servidores sem áudio), o que quebraria a inicialização.
// No Linux o oto usa ALSA, então exigimos um dispositivo PCM de reprodução
// real (arquivos pcm*p em /dev/snd). Pode ser forçado com VALDORIA_NOAUDIO=1.
func audioAvailable() bool {
	if v := os.Getenv("VALDORIA_NOAUDIO"); v == "1" || v == "true" {
		return false
	}
	if runtime.GOOS != "linux" {
		return true
	}
	entries, err := os.ReadDir("/dev/snd")
	if err != nil {
		return false
	}
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, "pcm") && strings.HasSuffix(name, "p") {
			return true
		}
	}
	return false
}

// buildLibrary prepara todos os efeitos e trilhas, preferindo arquivos livres
// em assets/audio quando presentes e caindo para sons gerados caso contrário.
func (a *AudioManager) buildLibrary() {
	a.sfx[sfxShoot] = a.loadSFX("shoot", genShoot)
	a.sfx[sfxShootFlame] = a.loadSFX("shoot_flame", genShootFlame)
	a.sfx[sfxShootIce] = a.loadSFX("shoot_ice", genShootIce)
	a.sfx[sfxEnemyDown] = a.loadSFX("enemy_down", genEnemyDown)
	a.sfx[sfxPlayerHit] = a.loadSFX("player_hit", genPlayerHit)
	a.sfx[sfxPickup] = a.loadSFX("pickup", genPickup)
	a.sfx[sfxBomb] = a.loadSFX("bomb", genBomb)
	a.sfx[sfxVictory] = a.loadSFX("victory", genVictory)
	a.sfx[sfxGameOver] = a.loadSFX("game_over", genGameOver)
	a.sfx[sfxShield] = a.loadSFX("shield_break", genShieldBreak)
	a.sfx[sfxMenu] = a.loadSFX("menu", genMenuBlip)
	a.sfx[sfxEscape] = a.loadSFX("escape", genEscape)
	a.sfx[sfxCorruption] = a.loadSFX("corruption", genCorruption)
	a.sfx[sfxDive] = a.loadSFX("dive", genDive)
	a.sfx[sfxGraze] = a.loadSFX("graze", genGraze)

	a.music[musicMenu] = a.loadMusic("music_menu", genMusicMenu)
	a.music[musicPhase] = a.loadMusic("music_phase", genMusicPhase)
	a.music[musicBoss] = a.loadMusic("music_boss", genMusicBoss)
	a.music[musicFields] = a.loadMusic("music_fields", genMusicFields)
	a.music[musicVillage] = a.loadMusic("music_village", genMusicVillage)
	a.music[musicWalls] = a.loadMusic("music_walls", genMusicWalls)
	a.music[musicCastle] = a.loadMusic("music_castle", genMusicCastle)
	a.music[musicForest] = a.loadMusic("music_forest", genMusicForest)
	a.music[musicSwamp] = a.loadMusic("music_swamp", genMusicSwamp)
	a.music[musicCanyon] = a.loadMusic("music_canyon", genMusicCanyon)
	a.music[musicLair] = a.loadMusic("music_lair", genMusicLair)
}

func (a *AudioManager) update() {
	if a == nil {
		return
	}
	if a.current != a.target {
		a.musicGain -= musicFadeStep
		if a.musicGain <= 0 {
			a.musicGain = 0
			a.pauseMusic(a.current)
			a.current = a.target
			a.startMusic(a.current)
		}
	} else if a.musicGain < 1 {
		a.musicGain += musicFadeStep
		if a.musicGain > 1 {
			a.musicGain = 1
		}
	}
	a.applyMusicVolume()
	for i := range a.frameGate {
		a.frameGate[i] = false
	}
}

func (a *AudioManager) playMusic(id musicID) {
	if a == nil {
		return
	}
	a.target = id
}

func (a *AudioManager) startMusic(id musicID) {
	if id == musicNone {
		return
	}
	if p := a.music[id]; p != nil {
		_ = p.Rewind()
		p.Play()
	}
}

func (a *AudioManager) pauseMusic(id musicID) {
	if id == musicNone {
		return
	}
	if p := a.music[id]; p != nil {
		p.Pause()
	}
}

func (a *AudioManager) applyMusicVolume() {
	if a.current == musicNone {
		return
	}
	if p := a.music[a.current]; p != nil {
		p.SetVolume(a.effectiveMusicVolume())
	}
}

func (a *AudioManager) effectiveMusicVolume() float64 {
	if a.muted {
		return 0
	}
	return clamp01(a.musicGain * a.musicVol * a.master)
}

func (a *AudioManager) playSFX(id soundID) {
	if a == nil || a.muted {
		return
	}
	if id < 0 || int(id) >= len(a.frameGate) {
		return
	}
	if a.frameGate[id] {
		return
	}
	a.frameGate[id] = true
	p := a.sfx[id]
	if p == nil {
		return
	}
	p.SetVolume(clamp01(a.sfxVol * a.master))
	_ = p.Rewind()
	p.Play()
}

func (a *AudioManager) toggleMute() {
	if a == nil {
		return
	}
	a.muted = !a.muted
}

func (a *AudioManager) setMasterVolume(v float64) { a.master = clamp01(v) }
func (a *AudioManager) setMusicVolume(v float64)  { a.musicVol = clamp01(v) }
func (a *AudioManager) setSFXVolume(v float64)    { a.sfxVol = clamp01(v) }

func (a *AudioManager) nudgeMusicVolume(delta float64) {
	if a == nil {
		return
	}
	a.musicVol = clamp01(snapVolume(a.musicVol + delta))
}

func (a *AudioManager) nudgeSFXVolume(delta float64) {
	if a == nil {
		return
	}
	a.sfxVol = clamp01(snapVolume(a.sfxVol + delta))
}

// snapVolume alinha o volume aos degraus de volumeStep (ex.: 0.5, 0.6…).
func snapVolume(v float64) float64 {
	return math.Round(v/volumeStep) * volumeStep
}

func volumePercent(v float64) int {
	return int(math.Round(clamp01(v) * 100))
}

// loadSFX devolve um player pronto: arquivo livre, se existir, ou som gerado.
func (a *AudioManager) loadSFX(name string, gen func(sampleRate int) []float64) *audio.Player {
	if a.ctx == nil {
		return nil
	}
	if b, ok := a.loadFileBytes(name); ok {
		return a.ctx.NewPlayerFromBytes(b)
	}
	return a.ctx.NewPlayerFromBytes(pcmFromSamples(gen(a.sampleRate)))
}

func (a *AudioManager) loadMusic(name string, gen func(sampleRate int) []float64) *audio.Player {
	if a.ctx == nil {
		return nil
	}
	if stream, length, ok := a.loadFileStream(name); ok {
		loop := audio.NewInfiniteLoop(stream, length)
		if p, err := a.ctx.NewPlayer(loop); err == nil {
			return p
		}
	}
	pcm := pcmFromSamples(gen(a.sampleRate))
	loop := audio.NewInfiniteLoop(bytes.NewReader(pcm), int64(len(pcm)))
	p, err := a.ctx.NewPlayer(loop)
	if err != nil {
		return nil
	}
	return p
}

// decoded expõe a interface comum dos streams de wav/vorbis.
type decoded interface {
	io.ReadSeeker
	Length() int64
}

// assetDirs lista os diretórios onde procurar assets opcionais: o diretório de
// trabalho e o diretório do executável (útil para pacotes distribuídos).
func assetDirs() []string {
	dirs := []string{audioDir}
	if exe, err := os.Executable(); err == nil {
		dirs = append(dirs, filepath.Join(filepath.Dir(exe), audioDir))
	}
	return dirs
}

func (a *AudioManager) decodeFile(name string) (decoded, bool) {
	for _, dir := range assetDirs() {
		for _, ext := range []string{"ogg", "wav"} {
			path := filepath.Join(dir, name+"."+ext)
			raw, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			var (
				stream decoded
				decErr error
			)
			switch ext {
			case "wav":
				stream, decErr = wav.DecodeWithSampleRate(a.sampleRate, bytes.NewReader(raw))
			case "ogg":
				stream, decErr = vorbis.DecodeWithSampleRate(a.sampleRate, bytes.NewReader(raw))
			}
			if decErr != nil {
				continue
			}
			return stream, true
		}
	}
	return nil, false
}

func (a *AudioManager) loadFileStream(name string) (io.ReadSeeker, int64, bool) {
	stream, ok := a.decodeFile(name)
	if !ok {
		return nil, 0, false
	}
	return stream, stream.Length(), true
}

func (a *AudioManager) loadFileBytes(name string) ([]byte, bool) {
	stream, ok := a.decodeFile(name)
	if !ok {
		return nil, false
	}
	b, err := io.ReadAll(stream)
	if err != nil {
		return nil, false
	}
	return b, true
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// --- Geração procedural de áudio (PCM 16 bits, estéreo, sem cabeçalho) ---

// pcmFromSamples converte amostras em [-1,1] para PCM estéreo de 16 bits.
func pcmFromSamples(samples []float64) []byte {
	buf := make([]byte, len(samples)*4)
	for i, s := range samples {
		if s > 1 {
			s = 1
		}
		if s < -1 {
			s = -1
		}
		v := uint16(int16(s * 32767))
		binary.LittleEndian.PutUint16(buf[i*4:], v)
		binary.LittleEndian.PutUint16(buf[i*4+2:], v)
	}
	return buf
}

type waveform int

const (
	waveSine waveform = iota
	waveSquare
	waveNoise
)

func envelope(i, n int, attack, release float64) float64 {
	t := float64(i) / float64(n)
	if t < attack {
		return t / attack
	}
	if t > 1-release {
		return (1 - t) / release
	}
	return 1
}

// synth gera uma nota com varredura linear de frequência e envelope simples.
func synth(sampleRate int, f0, f1, dur, amp float64, wave waveform, attack, release float64) []float64 {
	n := int(dur * float64(sampleRate))
	out := make([]float64, n)
	phase := 0.0
	for i := 0; i < n; i++ {
		t := float64(i) / float64(n)
		f := f0 + (f1-f0)*t
		phase += 2 * math.Pi * f / float64(sampleRate)
		var s float64
		switch wave {
		case waveSquare:
			if math.Sin(phase) >= 0 {
				s = 1
			} else {
				s = -1
			}
		case waveNoise:
			s = rand.Float64()*2 - 1
		default:
			s = math.Sin(phase)
		}
		out[i] = s * amp * envelope(i, n, attack, release)
	}
	return out
}

func mix(a, b []float64) []float64 {
	if len(b) > len(a) {
		a, b = b, a
	}
	out := make([]float64, len(a))
	copy(out, a)
	for i := range b {
		out[i] += b[i]
	}
	return out
}

func genShoot(sr int) []float64 {
	return synth(sr, 900, 480, 0.08, 0.35, waveSquare, 0.02, 0.6)
}

// genShootFlame: "whoosh" mais grave e ruidoso para as Chamas do Dragão.
func genShootFlame(sr int) []float64 {
	body := synth(sr, 520, 240, 0.10, 0.30, waveSquare, 0.02, 0.7)
	noise := synth(sr, 0, 0, 0.08, 0.12, waveNoise, 0.01, 0.7)
	return mix(body, noise)
}

// genShootIce: tom cristalino agudo e curto para as Lanças de Gelo.
func genShootIce(sr int) []float64 {
	return synth(sr, 1400, 1000, 0.07, 0.28, waveSine, 0.01, 0.6)
}

func genEnemyDown(sr int) []float64 {
	noise := synth(sr, 0, 0, 0.14, 0.25, waveNoise, 0.01, 0.7)
	low := synth(sr, 260, 90, 0.14, 0.3, waveSquare, 0.01, 0.6)
	return mix(noise, low)
}

func genPlayerHit(sr int) []float64 {
	return synth(sr, 220, 80, 0.22, 0.4, waveSquare, 0.01, 0.5)
}

func genPickup(sr int) []float64 {
	a := synth(sr, 660, 660, 0.07, 0.3, waveSine, 0.05, 0.3)
	b := synth(sr, 990, 990, 0.10, 0.3, waveSine, 0.05, 0.4)
	return append(a, b...)
}

func genBomb(sr int) []float64 {
	low := synth(sr, 300, 50, 0.5, 0.45, waveSquare, 0.02, 0.5)
	noise := synth(sr, 0, 0, 0.5, 0.2, waveNoise, 0.02, 0.6)
	return mix(low, noise)
}

func genVictory(sr int) []float64 {
	out := synth(sr, 523, 523, 0.14, 0.32, waveSine, 0.05, 0.3)
	out = append(out, synth(sr, 659, 659, 0.14, 0.32, waveSine, 0.05, 0.3)...)
	out = append(out, synth(sr, 784, 784, 0.22, 0.32, waveSine, 0.05, 0.3)...)
	return out
}

func genGameOver(sr int) []float64 {
	out := synth(sr, 392, 392, 0.2, 0.3, waveSquare, 0.03, 0.4)
	out = append(out, synth(sr, 330, 330, 0.2, 0.3, waveSquare, 0.03, 0.4)...)
	out = append(out, synth(sr, 262, 262, 0.3, 0.3, waveSquare, 0.03, 0.5)...)
	return out
}

// genShieldBreak: "clink" curto e brilhante (ruído + tom agudo) para o momento
// em que o escudo absorve um golpe, distinto do som de dano.
func genShieldBreak(sr int) []float64 {
	tone := synth(sr, 1200, 700, 0.10, 0.28, waveSine, 0.01, 0.6)
	noise := synth(sr, 0, 0, 0.06, 0.12, waveNoise, 0.01, 0.7)
	return mix(tone, noise)
}

// genMenuBlip: bipe curtíssimo para navegação/confirmação de menu.
func genMenuBlip(sr int) []float64 {
	return synth(sr, 660, 660, 0.045, 0.22, waveSine, 0.05, 0.5)
}

// genEscape: descida grave e curta — o som de algo passando por você. É de
// propósito desconfortável: o jogador precisa registrar cada fuga.
func genEscape(sr int) []float64 {
	tone := synth(sr, 300, 120, 0.16, 0.22, waveSine, 0.02, 0.7)
	noise := synth(sr, 0, 0, 0.10, 0.08, waveNoise, 0.02, 0.8)
	return mix(tone, noise)
}

// genCorruption: acorde grave e sujo ao subir de faixa — o reino piorando.
func genCorruption(sr int) []float64 {
	low := synth(sr, 160, 70, 0.55, 0.34, waveSquare, 0.02, 0.5)
	grind := synth(sr, 0, 0, 0.45, 0.16, waveNoise, 0.05, 0.6)
	wail := synth(sr, 420, 300, 0.50, 0.16, waveSine, 0.15, 0.5)
	return mix(mix(low, grind), wail)
}

// genDive: grito curto e ascendente do grifo ao se lançar.
func genDive(sr int) []float64 {
	cry := synth(sr, 620, 980, 0.16, 0.26, waveSine, 0.02, 0.55)
	air := synth(sr, 0, 0, 0.14, 0.10, waveNoise, 0.05, 0.7)
	return mix(cry, air)
}

// genGraze: tique agudo e curtíssimo ao raspar um projétil. Precisa ser
// discreto — acontece muitas vezes seguidas.
func genGraze(sr int) []float64 {
	return synth(sr, 1750, 1500, 0.035, 0.14, waveSine, 0.02, 0.6)
}

// sequence encadeia notas (freq, duração) formando um laço curto de música.
func sequence(sr int, wave waveform, amp float64, notes [][2]float64) []float64 {
	var out []float64
	for _, note := range notes {
		out = append(out, synth(sr, note[0], note[0], note[1], amp, wave, 0.05, 0.2)...)
	}
	return out
}

func genMusicMenu(sr int) []float64 {
	return sequence(sr, waveSine, 0.18, [][2]float64{
		{440, 0.4}, {523, 0.4}, {659, 0.4}, {523, 0.4},
		{494, 0.4}, {587, 0.4}, {494, 0.4}, {392, 0.4},
	})
}

func genMusicPhase(sr int) []float64 {
	return sequence(sr, waveSquare, 0.15, [][2]float64{
		{330, 0.25}, {392, 0.25}, {494, 0.25}, {392, 0.25},
		{440, 0.25}, {523, 0.25}, {440, 0.25}, {349, 0.25},
	})
}

func genMusicBoss(sr int) []float64 {
	return sequence(sr, waveSquare, 0.18, [][2]float64{
		{196, 0.22}, {233, 0.22}, {196, 0.22}, {174, 0.22},
		{220, 0.22}, {262, 0.22}, {220, 0.22}, {165, 0.22},
	})
}

func genMusicFields(sr int) []float64 {
	return sequence(sr, waveSine, 0.16, [][2]float64{
		{392, 0.28}, {440, 0.28}, {523, 0.28}, {587, 0.28},
		{523, 0.28}, {440, 0.28}, {392, 0.28}, {349, 0.28},
	})
}

func genMusicVillage(sr int) []float64 {
	return sequence(sr, waveSquare, 0.14, [][2]float64{
		{311, 0.20}, {370, 0.20}, {311, 0.20}, {277, 0.20},
		{349, 0.20}, {415, 0.20}, {349, 0.20}, {277, 0.20},
	})
}

func genMusicWalls(sr int) []float64 {
	return sequence(sr, waveSquare, 0.15, [][2]float64{
		{220, 0.30}, {262, 0.15}, {220, 0.30}, {196, 0.15},
		{247, 0.30}, {294, 0.15}, {247, 0.30}, {185, 0.15},
	})
}

func genMusicCastle(sr int) []float64 {
	return sequence(sr, waveSine, 0.17, [][2]float64{
		{175, 0.35}, {208, 0.35}, {233, 0.35}, {262, 0.35},
		{233, 0.35}, {208, 0.35}, {175, 0.35}, {156, 0.35},
	})
}

func genMusicForest(sr int) []float64 {
	return sequence(sr, waveSine, 0.14, [][2]float64{
		{262, 0.32}, {294, 0.32}, {349, 0.32}, {330, 0.32},
		{294, 0.32}, {247, 0.32}, {262, 0.32}, {220, 0.32},
	})
}

func genMusicSwamp(sr int) []float64 {
	return sequence(sr, waveSquare, 0.13, [][2]float64{
		{185, 0.36}, {208, 0.18}, {196, 0.36}, {233, 0.18},
		{175, 0.36}, {220, 0.18}, {165, 0.36}, {196, 0.18},
	})
}

func genMusicCanyon(sr int) []float64 {
	return sequence(sr, waveSine, 0.15, [][2]float64{
		{147, 0.40}, {165, 0.20}, {185, 0.40}, {196, 0.20},
		{175, 0.40}, {156, 0.20}, {147, 0.40}, {131, 0.20},
	})
}

func genMusicLair(sr int) []float64 {
	return sequence(sr, waveSquare, 0.17, [][2]float64{
		{110, 0.24}, {131, 0.24}, {147, 0.24}, {131, 0.24},
		{123, 0.24}, {147, 0.24}, {165, 0.24}, {110, 0.24},
	})
}

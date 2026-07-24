package game

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// saveEnabled permite desligar a persistência em disco (usado nos testes para
// não criar/alterar o arquivo de save real do usuário).
var saveEnabled = true

// saveData é o estado persistido entre execuções: recordes e preferências.
type saveData struct {
	HighScore    int             `json:"high_score"`
	SurvivalBest int             `json:"survival_best"`
	Difficulty   difficultyLevel `json:"difficulty"`
	Shake        shakeLevel      `json:"shake"`
	Muted        bool            `json:"muted"`
}

func defaultSave() saveData {
	return saveData{Difficulty: diffNormal, Shake: shakeFull}
}

// savePath devolve o caminho do arquivo de save no diretório de configuração do
// usuário. Retorna false quando o sistema não expõe esse diretório.
func savePath() (string, bool) {
	dir, err := os.UserConfigDir()
	if err != nil || dir == "" {
		return "", false
	}
	return filepath.Join(dir, "valdoria", "save.json"), true
}

// loadSave lê o save do disco, caindo para os padrões se ele não existir ou
// estiver corrompido — o jogo nunca falha por causa disso.
func loadSave() saveData {
	s := defaultSave()
	if !saveEnabled {
		return s
	}
	path, ok := savePath()
	if !ok {
		return s
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return s
	}
	if json.Unmarshal(b, &s) != nil {
		return defaultSave()
	}
	s.sanitize()
	return s
}

func (s *saveData) sanitize() {
	if s.Difficulty < 0 || s.Difficulty >= difficultyCount {
		s.Difficulty = diffNormal
	}
	if s.Shake < shakeFull || s.Shake > shakeOff {
		s.Shake = shakeFull
	}
	if s.HighScore < 0 {
		s.HighScore = 0
	}
	if s.SurvivalBest < 0 {
		s.SurvivalBest = 0
	}
}

// write grava o save em disco de forma tolerante a falhas (best-effort).
func (s saveData) write() {
	if !saveEnabled {
		return
	}
	path, ok := savePath()
	if !ok {
		return
	}
	if os.MkdirAll(filepath.Dir(path), 0o755) != nil {
		return
	}
	if b, err := json.Marshal(s); err == nil {
		_ = os.WriteFile(path, b, 0o644)
	}
}

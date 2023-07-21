// Анализатор для проекта Gowatcher. Использованы как стандартные,
// так и несколько установленных извне, таких как ineffassign и errcheck,
// а также был добавлен свой анализатор exitchecker.
package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gordonklaus/ineffassign/pkg/ineffassign"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"

	"github.com/nmramorov/gowatcher/cmd/staticlint/exitchecker"
)

// Config — имя файла конфигурации.
// Path - путь к файлу конфигурации в проекте.
const (
	Config = `config.json`
	Path   = `/cmd/staticlint`
)

// ConfigData описывает структуру файла конфигурации.
type ConfigData struct {
	Staticcheck []string
}

func main() {
	abspath, err := filepath.Abs(Config)
	if err != nil {
		panic(err)
	}
	dir := filepath.Join(filepath.Dir(abspath), Path)
	data, err := os.ReadFile(filepath.Join(dir, Config))
	if err != nil {
		panic(err)
	}
	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	mychecks := []*analysis.Analyzer{
		printf.Analyzer,
		shadow.Analyzer,
		structtag.Analyzer,
		shift.Analyzer,
		errcheck.Analyzer,
		ineffassign.Analyzer,
		exitchecker.OsExitAnalyzer,
	}

	checks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		checks[v] = true
	}
	// добавляем анализаторы из staticcheck, которые указаны в файле конфигурации.
	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	multichecker.Main(
		mychecks...,
	)
}

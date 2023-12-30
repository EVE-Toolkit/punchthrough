package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx   context.Context
	Chain Chain
}

type Chain struct {
	Systems []System `json:"systems"`
}

type System struct {
	Name        string   `json:"name"`
	Comments    []string `json:"comments"`
	Sigs        []Sig    `json:"sigs"`
	ConnectsTo  []string `json:"connects_to"`
	ConnectedTo string   `json:"connected_to"`
}

type Sig struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetActiveChain() Chain {
	return a.Chain
}

func (a *App) NewChain(system System) Chain {
	chain := Chain{
		Systems: []System{
			system,
		},
	}

	a.Chain = chain

	return chain
}

func (a *App) ConnectSystem(newSystem string, oldSystem string) {
	newSystemIndex := 0
	oldSystemIndex := 0

	for i, system := range a.Chain.Systems {
		if system.Name == newSystem {
			newSystemIndex = i
		}

		if system.Name == oldSystem {
			oldSystemIndex = i
		}
	}

	a.Chain.Systems[newSystemIndex].ConnectedTo = oldSystem
	a.Chain.Systems[oldSystemIndex].ConnectsTo = append(a.Chain.Systems[oldSystemIndex].ConnectsTo, newSystem)
}

func (a *App) AddSystem(system System) {
	a.Chain.Systems = append(a.Chain.Systems, system)
}

func (a *App) CreateComment(systemName, text string) {
	fmt.Println(text)

	systemIndex := 0

	for i, system := range a.Chain.Systems {
		if system.Name == systemName {
			systemIndex = i
		}
	}

	a.Chain.Systems[systemIndex].Comments = append(a.Chain.Systems[systemIndex].Comments, text)

	fmt.Println(a.Chain.Systems[systemIndex].Comments)
}

func (a *App) CreateSigs(text string, systemName string) {
	systemIndex := 0

	for i, system := range a.Chain.Systems {
		if system.Name == systemName {
			systemIndex = i
		}
	}

	lines := strings.Split(text, "\n")

	for i, line := range lines {
		if line == "" || line == "\n" {
			lines = append(lines[:i], lines[i+1:]...)

			continue
		}

		fields := strings.Split(line, "\t")

		for j, field := range fields {
			if field == "" || field == " " {
				fields = append(fields[:j], fields[j+1:]...)

				continue
			}
		}

		sig := Sig{}

		if len(fields) < 2 {
			continue
		}

		if len(fields) == 6 && len(strings.Trim(fields[2], " ")) > 0 && len(strings.Trim(fields[3], " ")) > 0 {
			sig.ID = fields[0]
			sig.Type = fields[2]
			sig.Name = fields[3]
		} else if len(strings.Trim(fields[2], " ")) == 0 || len(strings.Trim(fields[3], " ")) == 0 {
			sig.ID = fields[0]
			sig.Type = "Unknown"
			sig.Name = "Unknown"
		}

		exists := false

		for _, s := range a.Chain.Systems[systemIndex].Sigs {
			if s.ID == sig.ID {
				exists = true
			}
		}

		if !exists {
			a.Chain.Systems[systemIndex].Sigs = append(a.Chain.Systems[systemIndex].Sigs, sig)
		}
	}
}

func (a *App) ImportChain() (Chain, error) {
	fileName, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Punchthrough JSON",
				Pattern:     "*.json",
			},
		},
	})

	if err != nil {
		return Chain{}, err
	}

	if fileName == "" {
		return Chain{}, errors.New("It appears the dialog was cancelled. Please try again.")
	}

	file, err := os.Open(fileName)

	if err != nil {
		return Chain{}, err
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)

	if err != nil {
		return Chain{}, err
	}

	chain := Chain{}

	err = json.Unmarshal(bytes, &chain)

	if err != nil {
		return Chain{}, err
	}

	if len(chain.Systems) == 0 {
		return Chain{}, errors.New("This chain is empty. Please start a new chain.")
	}

	a.Chain = chain

	return chain, nil
}

func (a *App) SaveChain() {

}

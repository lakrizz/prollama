package ollama

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"slices"
	"strings"

	"github.com/jonathanhecl/gollama"

	"github.com/lakrizz/prollama/config"
	"github.com/lakrizz/prollama/pkg/slicex"
)

type Ollama struct {
	userPrompt        string
	systemPrompt      string
	gollama           *gollama.Gollama
	ctx               context.Context
	cfg               *config.Config
	autoselectedModel bool
}

type OllamaOption func(*Ollama) error

func New(opts ...OllamaOption) (*Ollama, error) {
	o := &Ollama{
		userPrompt: `
	Please review the following Git Diff Hunk. Identify issues such as syntax errors, logical flaws, unhandled errors, code smells, and testing gaps. Suggest improvements, refactoring opportunities, or missing tests where applicable. Inform about State of the Art implementations, Best Practices and Design Patterns. Combine multiple findings for the same line into a single comment. 

Return a string that contains all errors in a summarized manner, highlight the origins of your comments. Use markdown if you want to.

Ignore metadata lines that indicate information for the Diffpatch (e.g., lines that contain four numbers).

This is the hunk: %v`,
		gollama: gollama.New(""),
	}

	for _, opt := range opts {
		opt(o)
	}

	// extract config from context
	if o.ctx != nil {
		cfg, ok := config.FromContext(o.ctx)
		if !ok {
			return nil, fmt.Errorf("could not get config from context")
		}

		if cfg != nil {
			o.cfg = cfg
		}
	}

	// some fallback mechanisms for properties
	if o.cfg.Model == "" || !o.ValidModel() {
		err := o.AutoSelectModel()
		if err != nil {
			return nil, fmt.Errorf("could not select any model")
		}
	}

	// check if everything actually works and is valid
	if err := o.Validate(); err != nil {
		return nil, err
	}

	return o, nil
}

func (o *Ollama) Validate() error {
	// add check if ollama cmd is available
	_, err := exec.LookPath("ollama")
	if err != nil {
		return fmt.Errorf("could not find ollama path: %w", err)
	}

	return nil
}

func (o *Ollama) GetModels() ([]string, error) {
	models, err := o.gollama.ListModels(o.ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list models: %w", err)
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("no models found")
	}

	return slicex.ConvertSlice(models, func(e gollama.ModelInfo) string { return e.Model }), nil
}

func (o *Ollama) ValidModel() bool {
	models, err := o.GetModels()
	if err != nil {
		return false
	}

	return slices.Contains(models, o.cfg.Model)
}

func (o *Ollama) AutoSelectModel() error {
	if o.cfg.Debug {
		slog.Debug("given model not found, auto selecting model", "given_model", o.cfg.Model)
	}

	prioList := []string{
		"qwen2.5-coder",
		"starcoder2",
		"deepseek-coder-v2",
		"deepseek-coder",
		"starcoder",
		"sqlcoder",
		"dolphincoder",
		"yi-coder",
		"opencoder",
		"deepseek-v2.5",
		"codellama",
		"dolphin-mixtral",
		"codegemma",
		"codestral",
		"codegeex4",
		"codeqwen",
		"mistral-large",
		"codeup",
		"codebooga",
		"tulu3",
		"coder",
		"code",
	}

	models, err := o.GetModels()
	if err != nil {
		return fmt.Errorf("could not get ollama models: %w", err)
	}

	if len(models) == 0 {
		return fmt.Errorf("no models in ollama found")
	}

	for _, prio := range prioList {
		for _, m := range models {
			if strings.Contains(m, prio) {
				if o.cfg.Debug {
					slog.Debug("autoselected model", "model", m)
				}

				o.autoselectedModel = true
				o.cfg.Model = m
				return nil
			}
		}
	}

	// still nothing? select the first available model
	slog.Debug("autoselected model", "model", models[0])
	o.autoselectedModel = true
	o.cfg.Model = models[0]
	return nil
}

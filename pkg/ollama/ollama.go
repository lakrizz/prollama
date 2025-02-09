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
	Please review the following unified Diff file. Identify issues such as syntax errors, logical flaws, unhandled errors, code smells, and testing gaps. Suggest improvements, refactoring opportunities, or missing tests where applicable. Inform about State of the Art implementations, Best Practices and Design Patterns. Combine multiple findings for the same line into a single comment. 

Return all feedback as an array of JSON objects, where each object contains the fields:  
- 'line': The line number in the changed file, as indicated by the patch metadata.  
- 'body': A detailed explanation of the issue and actionable suggestions for improvement.  
- 'affected_line': a copy of the line this comment belongs to. Include all (this also applies to repeated instances) control characters and the leading '+' or '-'

If no issues are found, return an empty array ('[]'). Always return valid JSON. 

Rules for your response:
- ONLY return valid JSON data in a single line
- Sanitize and clean the resulting JSON string before returning it
- Remove all tabs and newlines from the JSON string before returning it

Ignore metadata lines that indicate information for the Diffpatch (e.g., lines that contain four numbers).
Assume all brackets, quotes, parantheses are closed at some point, so do not mark a missing closing or an unclosed pair as an error. 

This is the patch: %v`,
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

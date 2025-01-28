package ollama

import "context"

func WithUserPrompt(prompt string) OllamaOption {
	return func(o *Ollama) error {
		o.userPrompt = prompt
		return nil
	}
}

func WithSystemPrompt(prompt string) OllamaOption {
	return func(o *Ollama) error {
		o.systemPrompt = prompt
		return nil
	}
}

func WithContext(ctx context.Context) OllamaOption {
	return func(o *Ollama) error {
		o.ctx = ctx
		return nil
	}
}

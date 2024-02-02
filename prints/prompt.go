package prints

import (
	"bufio"
	"os"
	"strings"
)

type Prompt struct {
	question string

	defaultAnswer string
	secret        bool // TODO: support secret
}

type PromptOption func(*Prompt)

func WithDefaultAnswer(defaultAnswer string) PromptOption {
	return func(p *Prompt) {
		p.defaultAnswer = defaultAnswer
	}
}

func WithSecret() PromptOption {
	return func(p *Prompt) {
		p.secret = true
	}
}

func NewPrompt(question string, opts ...PromptOption) *Prompt {
	p := &Prompt{
		question: question,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Prompt) Ask() (string, error) {
	if _, err := Infof("%s ", p.question); err != nil {
		return "", err
	}
	if p.defaultAnswer != "" {
		if _, err := Alertf("[%s] ", p.defaultAnswer); err != nil {
			return "", err
		}
	}
	if _, err := Linef("\n> "); err != nil {
		return "", err
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return p.defaultAnswer, nil
	}
	return input, nil
}

func Ask(question string, defaults ...string) (string, error) {
	var opts []PromptOption
	if len(defaults) > 0 {
		opts = append(opts, WithDefaultAnswer(defaults[0]))
	}

	return NewPrompt(question, opts...).Ask()
}

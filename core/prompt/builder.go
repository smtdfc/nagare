package prompt

import (
	"bytes"
	"text/template"
)

type PromptTemplate struct {
	template *template.Template
}

func NewPromptTemplate(name string, content string) (PromptTemplate, error) {
	parsedTemplate, err := template.New(name).Parse(content)
	if err != nil {
		return PromptTemplate{}, err
	}

	return PromptTemplate{template: parsedTemplate}, nil
}

func MustNewPromptTemplate(name string, content string) PromptTemplate {
	promptTemplate, err := NewPromptTemplate(name, content)
	if err != nil {
		panic(err)
	}

	return promptTemplate
}

func (p PromptTemplate) Build(params any) (string, error) {
	var output bytes.Buffer
	if err := p.template.Execute(&output, params); err != nil {
		return "", err
	}

	return output.String(), nil
}

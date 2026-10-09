package prompt

import (
	"bytes"
	"text/template"
)

type Template struct {
	template *template.Template
}

func NewPromptTemplate(name string, content string) (Template, error) {
	parsedTemplate, err := template.New(name).Parse(content)
	if err != nil {
		return Template{}, err
	}

	return Template{template: parsedTemplate}, nil
}

func MustNewPromptTemplate(name string, content string) Template {
	promptTemplate, err := NewPromptTemplate(name, content)
	if err != nil {
		panic(err)
	}

	return promptTemplate
}

func (p Template) Build(params any) (string, error) {
	var output bytes.Buffer
	if err := p.template.Execute(&output, params); err != nil {
		return "", err
	}

	return output.String(), nil
}

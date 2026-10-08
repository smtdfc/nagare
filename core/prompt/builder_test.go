package prompt

import (
	"testing"
)

func TestNewPromptTemplate(t *testing.T) {
	// Verify valid template creation
	tmpl, err := NewPromptTemplate("greeting", "Hello, {{.Name}}!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	res, err := tmpl.Build(map[string]string{"Name": "Nagare"})
	if err != nil {
		t.Fatalf("unexpected error on build: %v", err)
	}
	if res != "Hello, Nagare!" {
		t.Errorf("expected 'Hello, Nagare!', got '%s'", res)
	}

	// Verify invalid template syntax
	_, err = NewPromptTemplate("invalid", "Hello, {{.Name")
	if err == nil {
		t.Errorf("expected error for invalid template syntax")
	}
}

func TestMustNewPromptTemplate(t *testing.T) {
	// Verify MustNewPromptTemplate does not panic on valid syntax
	tmpl := MustNewPromptTemplate("test", "Value: {{.Val}}")
	res, err := tmpl.Build(map[string]int{"Val": 42})
	if err != nil || res != "Value: 42" {
		t.Errorf("expected 'Value: 42', got '%s'", res)
	}

	// Verify panic on invalid syntax
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic on invalid template")
		}
	}()
	MustNewPromptTemplate("panic", "{{.Unclosed")
}

func TestPromptTemplate_Build(t *testing.T) {
	// Verify Build with struct params
	type contextData struct {
		Role string
		User string
	}
	tmpl := MustNewPromptTemplate("role_template", "Role: {{.Role}}, User: {{.User}}")
	res, err := tmpl.Build(contextData{Role: "assistant", User: "john"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "Role: assistant, User: john" {
		t.Errorf("expected 'Role: assistant, User: john', got '%s'", res)
	}
}

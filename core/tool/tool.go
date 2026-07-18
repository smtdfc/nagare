package tool

type ToolDeclaration[I any, O any] struct {
	Name        string
	Args        I
	Description string
	Callback    func(args I) (O, error)
}

func Declare[I any, O any](name string, description string, cb func(args I) (O, error)) *ToolDeclaration[I, O] {
	return &ToolDeclaration[I, O]{
		Name:        name,
		Description: description,
		Callback:    cb,
	}
}

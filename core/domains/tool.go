package domains

type Tool interface {
	GetName() string
	GetArgs() string
	GetDescription() string
	Execute(Context) (any, error)
}

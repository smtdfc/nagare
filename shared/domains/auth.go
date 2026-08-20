package shared_domains

type AuthPayload interface {
	GetID() string
	GetRole() string
	GetKind() string
	GetScopes() []string
}

package dto

type JwtAuthPayload struct {
	ID     string
	Role   string
	Kind   string
	Scopes []string
}

// GetID implements [AuthPayload].
func (j *JwtAuthPayload) GetID() string {
	return j.ID
}

// GetKind implements [AuthPayload].
func (j *JwtAuthPayload) GetKind() string {
	return j.Kind
}

// GetRole implements [AuthPayload].
func (j *JwtAuthPayload) GetRole() string {
	return j.Role
}

// GetScopes implements [AuthPayload].
func (j *JwtAuthPayload) GetScopes() []string {
	return j.Scopes
}

//go:build !goverter

package session

// @Injectable
func NewSessionMapper() ISessionMapper {
	return &ISessionMapperImpl{}
}

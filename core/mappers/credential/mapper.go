//go:build !goverter

package credential

// @Injectable
func NewCredentialMapper() ICredentialMapper {
	return &ICredentialMapperImpl{}
}

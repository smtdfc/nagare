//go:build !goverter

package credential

// @Injectable
func NewCredentialDtoMapper() ICredentialDtoMapper {
	return &ICredentialDtoMapperImpl{}
}

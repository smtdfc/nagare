//go:build !goverter

package model_provider

// @Injectable
func NewModelProviderDtoMapper() IModelProviderDtoMapper {
	return &IModelProviderDtoMapperImpl{}
}

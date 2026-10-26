//go:build !goverter

package model_provider

// @Injectable
func NewModelProviderMapper() IModelProviderMapper {
	return &IModelProviderMapperImpl{}
}

//go:build !goverter

package plugin

// @Injectable
func NewPluginDtoMapper() IPluginDtoMapper {
	return &IPluginDtoMapperImpl{}
}

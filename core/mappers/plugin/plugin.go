//go:build !goverter

package plugin

// @Injectable
func NewPluginMapper() IPluginMapper {
	return &IPluginMapperImpl{}
}

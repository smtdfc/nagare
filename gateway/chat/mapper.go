//go:build !goverter

package chat

// @Injectable
func NewChatMapper() ISessionDtoMapper {
	return &ISessionDtoMapperImpl{}
}

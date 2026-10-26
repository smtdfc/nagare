//go:build !goverter

package task

// @Injectable
func NewTaskMapper() ITaskMapper {
	return &ITaskMapperImpl{}
}

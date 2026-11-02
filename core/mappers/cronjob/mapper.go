//go:build !goverter

package cronjob

// @Injectable
func NewCronJobMapper() ICronJobMapper {
	return &ICronJobMapperImpl{}
}

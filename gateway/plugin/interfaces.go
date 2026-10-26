package plugin

import (
	"github.com/google/uuid"
	"github.com/smtdfc/nagare/core/plugin"
	"github.com/smtdfc/nagare/dtos/rest"
)

func uuidToString(id uuid.UUID) string {
	return id.String()
}

func featureToString(f plugin.Feature) string {
	return f.ToString()
}

// goverter:converter
// goverter:output:file ./converter_generated.go
// goverter:ignoreMissing
// goverter:extend uuidToString featureToString
type IPluginDtoMapper interface {
	ToPluginDTO(domain *plugin.Plugin) *rest.Plugin
	ToPluginDTOs(domains []*plugin.Plugin) []*rest.Plugin
}

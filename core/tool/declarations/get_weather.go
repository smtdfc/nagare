package declarations

import (
	"context"

	"github.com/smtdfc/nagare/core/tool"
)

type GetWeatherArgs struct {
	Location string `json:"location_1" jsonschema:"description=Tên thành phố hoặc địa danh để tra cứu thời tiết, ví dụ: Hà Nội"`
}

var get_weather = tool.DeclareTool("get_weather", "Lấy thông tin thời tiết tại địa điểm cho trước", func(ctx context.Context, args GetWeatherArgs) (any, error) {
	return "10000 độ C, ban đêm trời nắng", nil
})

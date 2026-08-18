package ws

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/smtdfc/nagare/shared/dto"
)

func GetPayload[T any](msg *dto.WsMessage[any]) (*T, error) {
	var result T
	err := mapstructure.Decode(msg.Payload, &result)
	return &result, err
}

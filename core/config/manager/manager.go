package manager

import (
	"context"

	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/mappers"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
)

const GENERAL_CONFIG_SCOPE_NAME = "nagare.config.general"

type ConfigManager struct {
	kvRepo   *repositories.KVRepository
	kvMapper *mappers.KVMapper
}

func (c *ConfigManager) GetGeneralConfig(ctx context.Context) (*config.GeneralConfig, error) {
	var conf config.GeneralConfig

	kvs, err := c.kvRepo.GetByScope(ctx, GENERAL_CONFIG_SCOPE_NAME)
	if err != nil {
		return nil, custom_errors.ErrGetGeneralConfigFailed
	}

	kvMap := c.kvMapper.ToDomains(kvs)
	conf.CurrentModel = kvMap["CurrentModel"]
	conf.CurrentProvider = kvMap["CurrentProvider"]

	return &conf, nil
}

func (c *ConfigManager) SetGeneralConfig(ctx context.Context, conf *config.GeneralConfig) error {
	var kvMap = map[string]string{
		"CurrentModel":    conf.CurrentModel,
		"CurrentProvider": conf.CurrentProvider,
	}
	err := c.kvRepo.Upsert(ctx, c.kvMapper.ToEntities(kvMap, GENERAL_CONFIG_SCOPE_NAME))
	if err != nil {
		return custom_errors.ErrSetGeneralConfigFailed
	}

	return nil
}

// @Injectable
func NewConfigManager(kvRepo *repositories.KVRepository, kvMapper *mappers.KVMapper) *ConfigManager {
	return &ConfigManager{
		kvRepo:   kvRepo,
		kvMapper: kvMapper,
	}
}

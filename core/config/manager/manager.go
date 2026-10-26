package manager

import (
	"context"

	"github.com/smtdfc/nagare/core/config"
	"github.com/smtdfc/nagare/core/custom_errors"
	"github.com/smtdfc/nagare/core/logger"
	kv_mapper "github.com/smtdfc/nagare/core/mappers/kv"
	"github.com/smtdfc/nagare/core/persistence/database/repositories"
)

const GeneralConfigScopeName = "nagare.config.general"

type ConfigManager struct {
	kvRepo   *repositories.KVRepository
	kvMapper kv_mapper.IKVMapper
	logger   *logger.BaseLogger
}

func (c *ConfigManager) GetGeneralConfig(ctx context.Context) (*config.GeneralConfig, error) {
	var conf config.GeneralConfig

	kvs, err := c.kvRepo.FindByScope(ctx, GeneralConfigScopeName)
	if err != nil {
		return nil, custom_errors.ErrGetGeneralConfigFailed
	}

	kvMap := c.kvMapper.ToDomains(kvs)
	conf.DefaultModelProvider = kvMap["DefaultModelProvider"]
	if conf.DefaultModelProvider == "" {
		conf.DefaultModelProvider = kvMap["DefaultLLMProvider"]
	}
	conf.DefaultModel = kvMap["DefaultModel"]
	if conf.DefaultModel == "" {
		conf.DefaultModel = kvMap["DefaultLLMModel"]
	}

	return &conf, nil
}

func (c *ConfigManager) SetGeneralConfig(ctx context.Context, conf *config.GeneralConfig) error {
	var kvMap = map[string]string{
		"DefaultModelProvider": conf.DefaultModelProvider,
		"DefaultModel":         conf.DefaultModel,
	}
	err := c.kvRepo.Upsert(ctx, c.kvMapper.ToEntities(kvMap, GeneralConfigScopeName))
	if err != nil {
		return custom_errors.ErrSetGeneralConfigFailed
	}

	return nil
}

// @Injectable
func NewConfigManager(
	kvRepo *repositories.KVRepository,
	kvMapper kv_mapper.IKVMapper,
	logger *logger.BaseLogger,
) *ConfigManager {
	return &ConfigManager{
		kvRepo:   kvRepo,
		kvMapper: kvMapper,
		logger:   logger,
	}
}

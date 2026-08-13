package mappers

import (
	"github.com/smtdfc/nagare/core/domains"
	"github.com/smtdfc/nagare/core/persistence/database/models"
)

type KVMapper struct {
}

func (k *KVMapper) ToGeneralConfig(models []models.KV) *domains.GeneralConfig {
	config := &domains.GeneralConfig{}
	for _, model := range models {
		switch model.Key {
		case "current_model":
			config.CurrentModel = model.Value
		case "current_provider":
			config.CurrentProvider = model.Value
		}
	}
	return config
}

func (k *KVMapper) FromGeneralConfig(config *domains.GeneralConfig, target string) []models.KV {
	ms := make([]models.KV, 0, 2)
	if config.CurrentModel != "" {
		ms = append(ms, models.KV{Target: target, Key: "current_model", Value: config.CurrentModel})
	}
	if config.CurrentProvider != "" {
		ms = append(ms, models.KV{Target: target, Key: "current_provider", Value: config.CurrentProvider})
	}
	return ms
}

func (k *KVMapper) ToMap(models []models.KV) map[string]string {
	m := make(map[string]string)
	for _, model := range models {
		m[model.Key] = model.Value
	}
	return m
}

func (k *KVMapper) ToModel(m map[string]string, target string) []models.KV {
	var ms []models.KV
	for key, value := range m {
		ms = append(ms, models.KV{Target: target, Key: key, Value: value})
	}
	return ms
}

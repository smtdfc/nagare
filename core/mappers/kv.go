package mappers

import "github.com/smtdfc/nagare/core/persistence/database/entities"

type KVMapper struct{}

func (k *KVMapper) ToDomains(entities []*entities.KV) map[string]string {
	result := make(map[string]string)
	for _, kv := range entities {
		result[kv.Key] = kv.Value
	}

	return result
}

func (k *KVMapper) ToEntities(domains map[string]string, scope string) []*entities.KV {
	result := make([]*entities.KV, 0, len(domains))
	for key, value := range domains {
		result = append(result, &entities.KV{
			Key:   key,
			Value: value,
			Scope: scope,
		})
	}

	return result
}

// @Injectable
func NewKVMapper() *KVMapper {
	return &KVMapper{}
}

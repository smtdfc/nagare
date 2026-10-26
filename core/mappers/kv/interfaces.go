package kv

import "github.com/smtdfc/nagare/core/persistence/database/entities"

type IKVMapper interface {
	ToEntities(domains map[string]string, scope string) []*entities.KV
	ToDomains(entities []*entities.KV) map[string]string
}

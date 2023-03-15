package cache

import "github.com/minhmannh2001/sms/entity"

type ServerCache interface {
	Set(key string, value *entity.Server)
	Get(key string) *entity.Server
	ASet(key string, value []entity.Server)
	AGet(key string) []entity.Server
}

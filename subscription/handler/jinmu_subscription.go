package handler

import (
	db "github.com/jinmukeji/jiujiantang-services/subscription/mysqldb"
)

// SubscriptionService 订阅服务
type SubscriptionService struct {
	datastore                db.Datastore
	activationCodeEntryptKey string
}

// NewSubscriptionService 构建SubscriptionService
func NewSubscriptionService(datastore db.Datastore, entrypt string) *SubscriptionService {
	j := &SubscriptionService{
		datastore:                datastore,
		activationCodeEntryptKey: entrypt,
	}
	return j
}

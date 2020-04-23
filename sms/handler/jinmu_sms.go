package handler

import (
	db "github.com/jinmukeji/jiujiantang-services/sms/mysqldb"
	sms "github.com/jinmukeji/jiujiantang-services/sms/sms_client"
)

// SMSGateway 短信网关
type SMSGateway struct {
	datastore       db.Datastore
	aliyunSMSClient *sms.AliyunSMSClient
}

// NewSMSGateway 构建SMSGateway
func NewSMSGateway(datastore db.Datastore, aliyunClient *sms.AliyunSMSClient) *SMSGateway {
	j := &SMSGateway{
		datastore:       datastore,
		aliyunSMSClient: aliyunClient,
	}
	return j
}

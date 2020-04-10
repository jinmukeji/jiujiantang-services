package handler

import (
	db "github.com/jinmukeji/gf-api2/sms/mysqldb"
	sms "github.com/jinmukeji/gf-api2/sms/sms_client"
)

// SMSGateway 短信网关
type SMSGateway struct {
	datastore           db.Datastore
	aliyunSMSClient     *sms.AliyunSMSClient
	tencentYunSmsClient *sms.TencentYunSMSClient
}

// NewSMSGateway 构建SMSGateway
func NewSMSGateway(datastore db.Datastore, aliyunClient *sms.AliyunSMSClient, tencentYunClient *sms.TencentYunSMSClient) *SMSGateway {
	j := &SMSGateway{
		datastore:           datastore,
		aliyunSMSClient:     aliyunClient,
		tencentYunSmsClient: tencentYunClient,
	}
	return j
}

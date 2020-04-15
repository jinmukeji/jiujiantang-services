package handler

import (
	"context"
	"fmt"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/core/v1"
)

// JinmuLMakePayment 支付
func (j *JinmuHealth) JinmuLMakePayment(ctx context.Context, req *proto.JinmuLMakePaymentRequest, resp *proto.JinmuLMakePaymentResponse) error {
	err := j.datastore.UpdateRecordHasPaid(ctx, req.RecordId)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to update record %d to the status of has paid: %s", req.RecordId, err.Error()))
	}
	return nil
}

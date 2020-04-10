package handler

import (
	"context"
	"fmt"

	mysql "github.com/jinmukeji/jiujiantang-services/service/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
)

// GetAnalysisReportTransactionNumber 获取流水号
func (j *JinmuHealth) GetAnalysisReportTransactionNumber(ctx context.Context, req *proto.GetAnalysisReportTransactionNumberRequest, resp *proto.GetAnalysisReportTransactionNumberResponse) error {
	transactionNumber, err := j.GetTransactionNumber(ctx)
	if err != nil {
		return err
	}
	resp.TransactionNumber = transactionNumber
	return nil
}

// GetTransactionNumber 获取TransactionNumber
func (j *JinmuHealth) GetTransactionNumber(ctx context.Context) (int32, error) {
	isExist, _ := j.datastore.IsExistTransactionNumberByCurrentDate(ctx)
	var t *mysql.TransactionNumber
	var err error

	if !isExist {
		t, err = j.datastore.CreateTransactionNumber(ctx)
		if err != nil {
			return 0, NewError(ErrDatabase, fmt.Errorf("failed to create transaction number: %s", err.Error()))
		}
	} else {
		t, err = j.datastore.FindTransactionNumberByCurrentDate(ctx)
		if err != nil {
			return 0, NewError(ErrDatabase, fmt.Errorf("failed to find transaction number by current date: %s", err.Error()))
		}
	}
	return int32(t.TransactionNumber), nil
}

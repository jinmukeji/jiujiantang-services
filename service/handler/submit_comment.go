package handler

import (
	"context"
	"errors"
	"time"

	"fmt"

	valid "github.com/asaskevich/govalidator"
	"github.com/jinmukeji/jiujiantang-services/service/auth"
	"github.com/jinmukeji/jiujiantang-services/service/mysqldb"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
)

// SubmitFeedback 账号意见反馈
func (j *JinmuHealth) SubmitFeedback(ctx context.Context, req *proto.SubmitFeedbackRequest, resp *proto.SubmitFeedbackResponse) error {
	if ok, err := validateSubmitFeedbackRequest(req); !ok || err != nil {
		return NewError(ErrNULLCommentOrContactWay, errors.New("content or contactway of feedback is null"))
	}
	userID, _ := auth.UserIDFromContext(ctx)
	now := time.Now()
	if err := j.datastore.CreateFeedback(ctx, &mysqldb.Feedback{
		UserID:     userID,
		ContactWay: req.ContactWay,
		Content:    req.Content,
		IsValid:    mysqldb.DbValidValue,
		CreatedAt:  now.UTC(),
		UpdatedAt:  now.UTC(),
	}); err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create feeback: %s", err.Error()))
	}
	return nil
}

// validateSubmitFeedbackRequest 验证意见反馈提交内容不能为空
func validateSubmitFeedbackRequest(req *proto.SubmitFeedbackRequest) (bool, error) {
	if valid.IsNull(req.Content) {
		return false, nil
	}
	if valid.IsNull(req.ContactWay) {
		return false, nil
	}
	return true, nil
}

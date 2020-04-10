package handler

import (
	"context"
	"fmt"

	auth "github.com/jinmukeji/gf-api2/service/auth"
	"github.com/jinmukeji/gf-api2/service/mysqldb"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
)

const (
	_ = iota
	isSportOrDrunk
	cold
	menstrualCycle
	ovipositPeriod
	lactation
	pregnancy
	cmStatusA
	cmStatusB
	cmStatusC
	cmStatusD
	cmStatusE
	cmStatusF
)

var statusDescription = map[int]string{
	isSportOrDrunk: "运动或饮酒",
	cold:           "感冒或病毒感染期",
	menstrualCycle: "生理周期",
	ovipositPeriod: "排卵期",
	lactation:      "哺乳期",
	pregnancy:      "怀孕",
	cmStatusA:      "口苦口黏，皮肤瘙痒，大便不成形，头重身痛",
	cmStatusB:      "急躁易怒，头晕胀痛",
	cmStatusC:      "口苦听力下降女性带下异味小便黄短",
	cmStatusD:      "口中异味反酸便秘喉咙干痒牙龈出血",
	cmStatusE:      "胃部冷痛，得温缓解",
	cmStatusF:      "失眠多梦健忘眩晕",
}

// SubmitMeasurementStatus 提交测量时身体状态
func (j *JinmuHealth) SubmitMeasurementStatus(ctx context.Context, req *corepb.SubmitMeasurementStatusRequest, repl *corepb.SubmitMeasurementStatusResponse) error {
	userID, _ := auth.UserIDFromContext(ctx)
	r, _ := j.datastore.FindRecordByID(ctx, int(req.RecordId))
	o, _ := j.datastore.FindOrganizationByUserID(ctx, r.UserID)
	isOwner, _ := j.datastore.CheckOrganizationOwner(ctx, int(userID), o.OrganizationID)
	if !isOwner {
		return NewError(ErrUserNotOrganizationOwner, fmt.Errorf("user %d is not the owner of organization %d", userID, o.OrganizationID))
	}
	record := &mysqldb.Record{
		RecordID:       int(req.RecordId),
		IsSportOrDrunk: int(req.IsSportOrDrunk),
		Cold:           int(req.Cold),
		MenstrualCycle: int(req.MenstrualCycle),
		OvipositPeriod: int(req.OvipositPeriod),
		Lactation:      int(req.Lactation),
		Pregnancy:      int(req.Pregnancy),
		StatusA:        int(req.CmStatusA),
		StatusB:        int(req.CmStatusB),
		StatusC:        int(req.CmStatusC),
		StatusD:        int(req.CmStatusD),
		StatusE:        int(req.CmStatusE),
		StatusF:        int(req.CmStatusF),
	}
	if err := j.datastore.UpdateRecordStatus(ctx, record); err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to update status of record %d: %s", req.RecordId, err.Error()))
	}
	statusDescriptions := make([]string, 0)
	if req.IsSportOrDrunk == corepb.Status_STATUS_SELECTED_STATUS {
		statusDescriptions = append(statusDescriptions, statusDescription[isSportOrDrunk])
	}
	if req.Cold == corepb.Status_STATUS_SELECTED_STATUS {
		statusDescriptions = append(statusDescriptions, statusDescription[cold])
	}
	if req.MenstrualCycle == corepb.Status_STATUS_SELECTED_STATUS {
		statusDescriptions = append(statusDescriptions, statusDescription[menstrualCycle])
	}
	if req.OvipositPeriod == corepb.Status_STATUS_SELECTED_STATUS {
		statusDescriptions = append(statusDescriptions, statusDescription[ovipositPeriod])
	}
	if req.Lactation == corepb.Status_STATUS_SELECTED_STATUS {
		statusDescriptions = append(statusDescriptions, statusDescription[lactation])
	}
	if req.Pregnancy == corepb.Status_STATUS_SELECTED_STATUS {
		statusDescriptions = append(statusDescriptions, statusDescription[pregnancy])
	}
	if req.CmStatusA == corepb.Status_STATUS_SELECTED_STATUS {
		statusDescriptions = append(statusDescriptions, statusDescription[cmStatusA])
	}
	if req.CmStatusB == corepb.Status_STATUS_SELECTED_STATUS {
		statusDescriptions = append(statusDescriptions, statusDescription[cmStatusB])
	}
	if req.CmStatusC == corepb.Status_STATUS_SELECTED_STATUS {
		statusDescriptions = append(statusDescriptions, statusDescription[cmStatusC])
	}
	if req.CmStatusD == corepb.Status_STATUS_SELECTED_STATUS {
		statusDescriptions = append(statusDescriptions, statusDescription[cmStatusD])
	}
	if req.CmStatusE == corepb.Status_STATUS_SELECTED_STATUS {
		statusDescriptions = append(statusDescriptions, statusDescription[cmStatusE])
	}
	if req.CmStatusF == corepb.Status_STATUS_SELECTED_STATUS {
		statusDescriptions = append(statusDescriptions, statusDescription[cmStatusF])
	}
	repl.StatusDescriptions = statusDescriptions
	return nil
}

// SubmitRecordAnswers 提交答案
func (j *JinmuHealth) SubmitRecordAnswers(ctx context.Context, req *corepb.SubmitRecordAnswersRequest, resp *corepb.SubmitRecordAnswersResponse) error {
	err := j.datastore.UpdateRecordAnswers(ctx, req.RecordId, req.Answers)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to update record %d answers: %s", req.RecordId, err.Error()))
	}
	return nil
}

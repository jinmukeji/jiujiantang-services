package handler

import (
	"context"
	"errors"

	"fmt"

	"github.com/jinmukeji/go-pkg/age"

	"github.com/golang/protobuf/ptypes"
	"github.com/jinmukeji/gf-api2/service/auth"
	"github.com/jinmukeji/gf-api2/service/mysqldb"
	corepb "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
	generalpb "github.com/jinmukeji/proto/gen/micro/idl/ptypes/v2"
	jinmuidpb "github.com/jinmukeji/proto/gen/micro/idl/jinmuid/v1"
)

// GetUserProfile 查看用户个人档案
func (j *JinmuHealth) GetUserProfile(ctx context.Context, req *corepb.GetUserProfileRequest, resp *corepb.GetUserProfileResponse) error {
	reqGetUserAndProfileInformation := new(jinmuidpb.GetUserAndProfileInformationRequest)
	reqGetUserAndProfileInformation.UserId = req.UserId
	reqGetUserAndProfileInformation.IsSkipVerifyToken = true
	respGetUserAndProfileInformation, errGetUserAndProfileInformation := j.jinmuidSvc.GetUserAndProfileInformation(ctx, reqGetUserAndProfileInformation)
	if errGetUserAndProfileInformation != nil {
		return errGetUserAndProfileInformation
	}
	birthday, _ := ptypes.Timestamp(respGetUserAndProfileInformation.Profile.BirthdayTime)
	resp.User = &corepb.User{
		UserId:             req.UserId,
		Username:           respGetUserAndProfileInformation.SigninUsername,
		RegisterType:       respGetUserAndProfileInformation.RegisterType,
		IsProfileCompleted: respGetUserAndProfileInformation.HasSetUserProfile,
		IsRemovable:        respGetUserAndProfileInformation.IsRemovable,
		Profile: &corepb.UserProfile{
			Nickname:        respGetUserAndProfileInformation.Profile.Nickname,
			Gender:          respGetUserAndProfileInformation.Profile.Gender,
			Age:             int32(age.Age(birthday)),
			Height:          respGetUserAndProfileInformation.Profile.Height,
			Weight:          respGetUserAndProfileInformation.Profile.Weight,
			BirthdayTime:    respGetUserAndProfileInformation.Profile.BirthdayTime,
			Remark:          respGetUserAndProfileInformation.Remark,
			UserDefinedCode: respGetUserAndProfileInformation.UserDefinedCode,
			NicknameInitial: respGetUserAndProfileInformation.Profile.NicknameInitial,
		},
	}
	return nil
}

// GetUserByRecordID 通过RecordID 获取 User
func (j *JinmuHealth) GetUserByRecordID(ctx context.Context, req *corepb.GetUserByRecordIDRequest, resp *corepb.GetUserByRecordIDResponse) error {
	isExsit, err := j.datastore.ExistRecordByRecordID(ctx, req.RecordId)
	if err != nil || !isExsit {
		return NewError(ErrGetUserFailure, fmt.Errorf("failed to check record existence by recordID %d, database error: %s", req.RecordId, err.Error()))
	}
	userID, errGetUserIDByRecordID := j.datastore.GetUserIDByRecordID(ctx, req.RecordId)
	if errGetUserIDByRecordID != nil {
		return NewError(ErrNotFoundUser, fmt.Errorf("failed to get userID by recordID %d: %s", req.RecordId, errGetUserIDByRecordID.Error()))
	}
	ownerID, ok := auth.UserIDFromContext(ctx)
	if !ok {
		return NewError(ErrInvalidUser, errors.New("failed to get userID from context"))
	}
	userOrganization, _ := j.datastore.FindOrganizationByUserID(ctx, int(userID))
	ownerOrganization, _ := j.datastore.FindOrganizationByUserID(ctx, int(ownerID))
	if ownerOrganization.OrganizationID != userOrganization.OrganizationID {
		return NewError(ErrGetUserFailure, errors.New("organization of user and the organization of owner are inconsistent"))
	}
	u, err := j.datastore.GetUserByRecordID(ctx, req.RecordId)
	if err != nil {
		return NewError(ErrGetUserFailure, fmt.Errorf("failed to get user by recordID %d: %s", req.RecordId, err.Error()))
	}
	birthday, _ := ptypes.TimestampProto(u.Birthday)
	gender, errMapDBGenderToProto := mapDBGenderToProto(u.Gender)
	if errMapDBGenderToProto != nil {
		return NewError(ErrInvalidGender, errMapDBGenderToProto)
	}

	resp.User = &corepb.User{
		UserId: int32(u.UserID),
		Profile: &corepb.UserProfile{
			Nickname:     u.Nickname,
			BirthdayTime: birthday,
			Gender:       gender,
			Weight:       int32(u.Weight),
			Height:       int32(u.Height),
		},
	}
	return nil
}

// mapDBGenderToProto 将数据库使用的 gender 映射为 proto 格式
func mapDBGenderToProto(gender mysqldb.Gender) (generalpb.Gender, error) {
	switch gender {
	case mysqldb.GenderFemale:
		return generalpb.Gender_GENDER_FEMALE, nil
	case mysqldb.GenderMale:
		return generalpb.Gender_GENDER_MALE, nil
	case mysqldb.GenderInvalid:
		return generalpb.Gender_GENDER_INVALID, fmt.Errorf("invalid database gender %s", gender)
	}
	return generalpb.Gender_GENDER_INVALID, fmt.Errorf("invalid database gender %s", gender)
}

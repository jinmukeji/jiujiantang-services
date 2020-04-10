package handler

import (
	"context"
	"fmt"
	"time"

	"github.com/jinmukeji/jiujiantang-services/service/auth"

	"github.com/jinmukeji/jiujiantang-services/service/mysqldb"

	"github.com/golang/protobuf/ptypes"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
)

// GetWxmpTempQrCodeUrl 得到微信临时二维码的URL
func (j *JinmuHealth) GetWxmpTempQrCodeUrl(ctx context.Context, req *proto.GetWxmpTempQrCodeUrlRequest, resp *proto.GetWxmpTempQrCodeUrlResponse) error {
	now := time.Now()
	qrcode := &mysqldb.QRCode{
		CreatedAt: now,
		UpdatedAt: now,
	}
	qrcode, errCreateQRCode := j.datastore.CreateQRCode(ctx, qrcode)
	if errCreateQRCode != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create qr code: %s", errCreateQRCode.Error()))
	}

	wxQrCode, err := j.wechat.GetTempQrCodeUrl(qrcode.SceneID)
	if err != nil {
		return NewError(ErrGetTempQrCodeURLFaliure, err)
	}

	account, _ := auth.AccountFromContext(ctx)
	machineUUID, _ := auth.MachineUUIDFromContext(ctx)
	updateQRCode := &mysqldb.QRCode{
		SceneID:     qrcode.SceneID,
		RawURL:      wxQrCode.RawURL,
		Account:     account,
		MachineUUID: machineUUID,
		Ticket:      wxQrCode.Ticket,
		ExpiredAt:   wxQrCode.ExpiredAt,
		OriginID:    wxQrCode.OriID,
		UpdatedAt:   now,
	}
	errUpdateQRCode := j.datastore.UpdateQRCode(ctx, updateQRCode)
	if errUpdateQRCode != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to update qr code: %s", errUpdateQRCode.Error()))
	}
	expiredAt, _ := ptypes.TimestampProto(wxQrCode.ExpiredAt)
	resp.ExpiredTime = expiredAt
	resp.ImageUrl = wxQrCode.ImageURL
	resp.RawUrl = wxQrCode.RawURL
	resp.SceneId = wxQrCode.SceneID
	return nil
}

// ScanQRCode 扫码二维码
func (j *JinmuHealth) ScanQRCode(ctx context.Context, req *proto.ScanQRCodeRequest, resp *proto.ScanQRCodeResponse) error {
	createdAt, _ := ptypes.Timestamp(req.CreatedTime)
	record := &mysqldb.ScannedQRCodeRecord{
		SceneID:   req.SceneId,
		CreatedAt: createdAt,
		UpdatedAt: time.Now().UTC(),
	}
	err := j.datastore.CreateScannedQRCodeRecord(ctx, record)
	if err != nil {
		return NewError(ErrDatabase, fmt.Errorf("failed to create scanned qr record record: %s", err.Error()))
	}
	return nil
}

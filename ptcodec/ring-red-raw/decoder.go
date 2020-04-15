package RingRedredraw

import (
	"fmt"

	"github.com/jinmukeji/plat-pkg/v2/micro/errors/errmsg"
)

const (
	// PointSize 摄像头每四个字节表示一个点
	PointSize int = 4
)

// RingRedRawDecoder 摄像头解码器
type RingRedRawDecoder struct {
}

func NewRingRedRawDecoder() *RingRedRawDecoder {
	return &RingRedRawDecoder{}
}

// Decode 按照 PointSize 解码
func (d *RingRedRawDecoder) Decode(data []byte) ([]uint32, error) {
	if len(data)%PointSize != 0 {
		return nil, fmt.Errorf(errmsg.InvalidValue("data", len(data)))
	}
	length := len(data) / PointSize
	payload := make([]uint32, length)
	for i := 0; i < len(data); i += PointSize {
		val := int32(data[i])<<24 + int32(data[i+1])<<16 + int32(data[i+2])<<8 + int32(data[i+3])
		payload[i/PointSize] = uint32(val)
	}
	return payload, nil
}

// DecodeBytes 按照字节解码
func (d *RingRedRawDecoder) DecodeBytes(data []byte) ([]byte, error) {
	return data, nil
}

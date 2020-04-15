package ptcodec

import (
	"strings"

	ringraw "github.com/jinmukeji/jiujiantang-services/ptcodec/ring-red-raw"
)

const (
	CodecRingRedRaw = "ring-red-raw"
)

type Decoder interface {
	// 按照 PointSize 解码
	Decode([]byte) ([]uint32, error)
	// 按照字节解码
	DecodeBytes([]byte) ([]byte, error)
}

type Encoder interface {
	Encode([]uint32) ([]byte, error)
}

func NewDecoder(codec string) Decoder {
	switch strings.ToLower(codec) {
	case CodecRingRedRaw:
		return ringraw.NewRingRedRawDecoder()
	default:
		return nil
	}
}

func NewEncoder(codec string) Encoder {
	switch strings.ToLower(codec) {
	case CodecRingRedRaw:
		return ringraw.NewRingRedRawEncoder()
	default:
		return nil
	}
}

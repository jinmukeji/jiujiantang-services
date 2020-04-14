package RingRedredraw

import (
	"bytes"
)

type RingRedRawEncoder struct {
}

func NewRingRedRawEncoder() *RingRedRawEncoder {
	return &RingRedRawEncoder{}
}

// Decode 按照 PointSize 解码
func (en *RingRedRawEncoder) Encode(data []uint32) ([]byte, error) {
	buf := make([]byte, 0, len(data)*PointSize)
	w := bytes.NewBuffer(buf)

	for _, d := range data {
		b1 := byte((d & 0xFF000000) >> 24)
		w.WriteByte(b1)

		b2 := byte((d & 0xFF0000) >> 16)
		w.WriteByte(b2)

		b3 := byte((d & 0xFF00) >> 8)
		w.WriteByte(b3)

		b4 := byte(d & 0xFF)
		w.WriteByte(b4)
	}

	return w.Bytes(), nil
}

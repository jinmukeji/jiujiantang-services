package handler

import (
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/ptypes/v2"
)

// ConvertToStandingCValues 转化站姿中的C0-C7
func ConvertToStandingCValues(gender proto.Gender, c0, c1, c2, c3, c4, c5, c6, c7 int) (sc0, sc1, sc2, sc3, sc4, sc5, sc6, sc7 int) {
	if c2 >= 7 && c2 <= 10 {
		c2 = c2 - 3
	} else if c2 >= 4 && c2 <= 6 {
		c2 = c2 - 2
	} else if c2 >= 2 && c2 <= 3 {
		c2 = c2 - 1
	}
	c4 = c4 + 1
	c0 = c0 + 1
	if gender == proto.Gender_GENDER_FEMALE {
		c3 = c3 - 1
	}
	return IntValBoundedBy10(c0), IntValBoundedBy10(c1), IntValBoundedBy10(c2), IntValBoundedBy10(c3), IntValBoundedBy10(c4), IntValBoundedBy10(c5), IntValBoundedBy10(c6), IntValBoundedBy10(c7)
}

// ConvertToStandingCInt32Values 转化站姿的C0-C7值
func ConvertToStandingCInt32Values(gender proto.Gender, c0, c1, c2, c3, c4, c5, c6, c7 int) (sc0, sc1, sc2, sc3, sc4, sc5, sc6, sc7 int32) {
	rsc0, rsc1, rsc2, rsc3, rsc4, rsc5, rsc6, rsc7 := ConvertToStandingCValues(gender, c0, c1, c2, c3, c4, c5, c6, c7)
	return int32(rsc0), int32(rsc1), int32(rsc2), int32(rsc3), int32(rsc4), int32(rsc5), int32(rsc6), int32(rsc7)
}

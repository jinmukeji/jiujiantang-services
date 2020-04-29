package handler

import (
	"context"
	"fmt"
	"strconv"

	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/user/v1"
)

// SecureQuestions 密保问题列表
var SecureQuestions = map[string]string{
	"1":  "你少年时代最好的朋友叫什么名字？",
	"2":  "你的第一个宠物叫什么名字？",
	"3":  "你学会做的第一道菜是什么？",
	"4":  "你第一次去电影院看的是哪一部电影？",
	"5":  "你第一次坐飞机是去哪里？",
	"6":  "你上小学时最喜欢的老师姓什么？",
	"7":  "你的父母是在哪里认识的？",
	"8":  "你的第一个上司叫什么名字？",
	"9":  "你从小长大的那条街叫什么？",
	"10": "你去过的第一个游乐场是哪一个？",
	"11": "你购买的第一张专辑是什么？",
	"12": "你最喜欢哪个球队？",
	"13": "你的理想工作是什么？",
	"14": "你小时候最喜欢哪一本书？",
	"15": "你童年时的绰号是什么？",
	"16": "你拥有的第一辆车是什么型号？",
	"17": "你在学生时代最喜欢的电影明星是谁？",
	"18": "你最喜欢哪个乐队或歌手？",
}

// UserGetSecureQuestionList 获取密保问题列表
func (j *JinmuIDService) UserGetSecureQuestionList(ctx context.Context, req *proto.UserGetSecureQuestionListRequest, resp *proto.UserGetSecureQuestionListResponse) error {
	// 判断UserID是否存在
	exist, errExistUserByUserID := j.datastore.ExistUserByUserID(ctx, req.UserId)
	if !exist || errExistUserByUserID != nil {
		return NewError(ErrInvalidUser, fmt.Errorf("userId %d doesn't exist", req.UserId))
	}
	secureQuestions := getSecureQuestionsByLanguage()
	protoSecureQuestions := make([]*proto.SecureQuestionKeyAndQuestion, len(secureQuestions))
	for idx, item := range secureQuestions {
		idxInt, errParseInt := strconv.Atoi(idx)
		if errParseInt != nil {
			log.Errorf("invalid parameter value of %s for index", idx)
			return fmt.Errorf("invalid parameter value of %s for index when getting secure question list: %s", idx, errParseInt.Error())
		}
		protoSecureQuestions[idxInt-1] = &proto.SecureQuestionKeyAndQuestion{
			Key:      idx,
			Question: item,
		}
	}
	resp.SecureQuestions = protoSecureQuestions
	return nil
}

func getSecureQuestionsByLanguage() map[string]string {

	secureQuestions := map[string]string{}
	for idx, item := range SecureQuestions {
		secureQuestions[idx] = item
	}
	return secureQuestions
}

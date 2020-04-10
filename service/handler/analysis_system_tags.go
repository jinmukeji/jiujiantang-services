package handler

import (
	"context"

	"github.com/jinmukeji/ae-v1/core"
	"github.com/jinmukeji/gf-api2/service/auth"
	proto "github.com/jinmukeji/proto/gen/micro/idl/jm/core/v1"
)

const (
	// dengyunClient 登云的client
	dengyunClient = "dengyun-10001"
	// moaiClient 摩爱的client
	moaiClient = "moai-10001"
	// hongjingtangClient 弘经堂的client
	hongjingtangClient = "hongjingtang-10001"
)

// GetAnalysisSystemTags 得到分析的tags
func (j *JinmuHealth) GetAnalysisSystemTags(ctx context.Context, req *proto.GetAnalysisSystemTagsRequest, resp *proto.GetAnalysisSystemTagsResponse) error {
	client, _ := clientFromContext(ctx)
	tags := make([]string, 0)
	switch client.ClientID {
	case dengyunClient:
		tags = append(tags, core.EnabledCCSystemTag)
	case moaiClient:
		tags = append(tags, core.EnabledCDSystemTag)
		tags = append(tags, core.EnabledCCSystemTag)
		tags = append(tags, core.EnabledSDSystemTag)
	case hongjingtangClient:
		tags = append(tags, core.EnabledCDSystemTag)
		tags = append(tags, core.EnabledCCSystemTag)
		tags = append(tags, core.EnabledSDSystemTag)
	case kangmeiClient:
		tags = append(tags, core.EnabledCDSystemTag)
		tags = append(tags, core.EnabledCCSystemTag)
		tags = append(tags, core.EnabledSDSystemTag)
	default:
		tags = append(tags, core.EnabledCDSystemTag)
		tags = append(tags, core.EnabledCCSystemTag)
		tags = append(tags, core.EnabledSDSystemTag)
		tags = append(tags, core.EnabledFactorSystemTag)
	}
	accessTokenType, _ := auth.AccessTokenTypeFromContext(ctx)
	// 女性部分内容是根据一体机客户进行开关的
	if accessTokenType == AccessTokenTypeLValue || accessTokenType == AccessTokenTypeWeChatValue {
		tags = append(tags, core.EnabledBreastSystemTag)
		tags = append(tags, core.EnabledEmotionalSystemTag)
		tags = append(tags, core.EnabledFacialSystemTag)
		tags = append(tags, core.EnabledGynecologicalSystemTag)
		tags = append(tags, core.EnabledHormoneSystemTag)
		tags = append(tags, core.EnabledLymphaticSystemTag)
		tags = append(tags, core.EnabledMammarySystemTag)
		tags = append(tags, core.EnabledMenstruationSystemTag)
		tags = append(tags, core.EnabledReproductiveSystemTag)
		tags = append(tags, core.EnabledUterineSystemTag)
	}
	resp.Tags = tags
	return nil
}

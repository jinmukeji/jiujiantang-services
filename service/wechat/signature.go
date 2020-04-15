package wechat

import (
	"sort"
	"strings"

	"github.com/jinmukeji/go-pkg/v2/crypto/hash"
)

// CheckWxSignature 验证微信接入的 Signature
func (u *Wxmp) CheckWxSignature(signature, timestamp, nonce string) bool {
	// 验证规范参考:
	//	https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421135319

	strs := []string{u.Options.WxToken, timestamp, nonce}
	sort.Strings(strs)
	v := strings.Join(strs, "")
	sv := hash.HexString(hash.SHA1String(v))
	return sv == signature
}

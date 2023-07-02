package utils

import (
	"github.com/xiaoyiEdu/Open-IM-Server/pkg/common/constant"
	sdkws "github.com/xiaoyiEdu/Open-IM-Server/pkg/proto/sdkws"
	"google.golang.org/protobuf/proto"
)

func GetContent(msg *sdkws.MsgData) string {
	if msg.ContentType >= constant.NotificationBegin && msg.ContentType <= constant.NotificationEnd {
		var tips sdkws.TipsComm
		_ = proto.Unmarshal(msg.Content, &tips)
		content := tips.JsonDetail
		return content
	} else {
		return string(msg.Content)
	}
}

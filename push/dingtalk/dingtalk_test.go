package dingtalk

import (
	"github.com/hq2005001/modules/push"
	"testing"
)

func TestNewDingTalk(t *testing.T) {
	NewDingTalk("da9e249dda5470a35669f3ed586d7f7bb7e7614f204bd868fc8c2e42c79e171a").SetMessageType("markdown").
		SetKeyword("监控报警").Push(push.Message{
		Body: "### 监控报警\n> 时代的脚步滚滚向前",
	})
}

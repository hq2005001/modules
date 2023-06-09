package wechat

import "testing"

func TestWechat_Verify(t *testing.T) {
	t.Log(New().Verify("60388835443409927", nil, nil, 0))
}

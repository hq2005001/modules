package sign

import (
	"testing"
)

func TestCustom(t *testing.T) {
	t.Log(NewCustomSigner().SetKeyName("sign").SetKey("u0vNaqciovqpjiPy").Sign(map[string]interface{}{
		"address":        "",
		"goods_id":       "611f647951d6e91ce3832baa",
		"hmac":           "mGSu|p1}^}RpyUq/",
		"password":       "5f4cb602d22072342341afe03697c60f",
		"x-device-id":    "e1b4a03d84965df4",
		"x-device-model": "Mozilla/5.0 (Linux; Android 7.1.1; OS105 Build/NGI77B; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/62.0.3202.84 Mobile Safari/537.36",
		"x-platform":     "android",
		"x-timestamp":    1629448098,
		"sign":           "AD6F3F766918B150478CAA7C4FB0460F",
	}))
}

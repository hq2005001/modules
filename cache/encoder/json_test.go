package encoder

import "testing"

type A struct {
	Name string `json:"name"`
}

func TestJSONEncoder_Encode(t *testing.T) {
	var a = A{
		Name: "111",
	}
	b, _ := NewJSONEncoder[A]().Encode(a)
	t.Log(b)
}

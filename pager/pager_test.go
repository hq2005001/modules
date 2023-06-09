package pager

import (
	"encoding/json"
	"testing"
)

func TestFilter_Predicate(t *testing.T) {
	var or struct {
		Or string `json:"or"`
	}

	t.Log(json.Unmarshal([]byte(`{"or": "x"}`), &or))
	t.Log(or)
}

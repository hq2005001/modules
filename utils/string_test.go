package utils

import "testing"

func TestSnake(t *testing.T) {
	t.Log(Snake("GoodsName"))
}

func TestReverseSnake(t *testing.T) {
	t.Log(ReverseSnake("goods_name"))
}

func TestWord(t *testing.T) {
	t.Log(Word("show-me-the_word"))
}

func TestRandomString(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(RandomString(8))
	}
}

package kv

import (
	"github.com/hq2005001/modules/cache/cacher"
	"github.com/hq2005001/modules/cache/encoder"
	"testing"
	"time"
)

type Test struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type TestCache[T any] struct {
	cacher.Cacher[T]
}

func (t TestCache[T]) Key() string {
	return "test_key"
}

func New[T any](encoder encoder.Encoder[T]) cacher.Cacher[T] {
	redisCache := NewRedisCache[T](encoder, time.Second*30)
	return TestCache[T]{
		Cacher: redisCache,
	}
}

func TestRedisCache_Get(t *testing.T) {
	test := Test{
		Name:     "111",
		Password: "222",
	}
	var item = New(encoder.New[Test]("json"))
	if err := item.Set(item.Key(), test); err != nil {
		t.Error(err)
		return
	}
	rs, err := item.Get(item.Key())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(item.Key())
	t.Log(rs.Name, rs.Password)
}

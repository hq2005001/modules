package set

//type Test struct {
//	Name     string `json:"name"`
//	Password string `json:"password"`
//}
//
//type TestCache struct {
//	cacher.SetCacher
//}
//
//func (t TestCache) Key() string {
//	return "test_key"
//}
//
//func New() cacher.SetCacher {
//	redisCache := NewRedisCache(time.Second * 30)
//	return TestCache{
//		SetCacher: redisCache,
//	}
//}
//
//func TestRedisCache_Get(t *testing.T) {
//	redis.InitRedis()
//	test := Test{
//		Name:     "111",
//		Password: "222",
//	}
//	item := bridge.NewJSONBridge(New())
//	if err := item.Set(item.Cache.Key(), test); err != nil {
//		t.Error(err)
//		return
//	}
//	rs, err := item.Get(item.Cache.Key(), &Test{})
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	t.Log(item.Cache.Key())
//	t.Log(rs.(*Test).Name, rs.(*Test).Password)
//}

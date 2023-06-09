package cache

import "github.com/hq2005001/modules/cache/cacher"

type Cachable[T any] func() (T, error)

func Cache[T any](cacher cacher.Cacher[T], cachable Cachable[T]) (T, error) {
	rs, err := cacher.Get(cacher.Key())
	if err != nil {
		rs, err = cachable()
		if err = cacher.Set(cacher.Key(), rs); err != nil {
			return rs, err
		}
	}
	return rs, err
}

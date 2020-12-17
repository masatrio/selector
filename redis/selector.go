package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v7"
	"gitlab.com/ruangguru/selector/selector"
)

type Selector struct {
	inner redis.UniversalClient
	calc  selector.Calculate
}

func Create(client redis.UniversalClient, calc selector.Calculate) selector.Selector {

	return &Selector{
		inner: client,
		calc:  calc,
	}
}

func (s *Selector) Select(
	ctx context.Context,
	namespace string,
	kind string,
	selection selector.Selection,
) (string, bool, error) {

	val, _, exists, err := s.Lookup(ctx, namespace, kind, selection)

	// val ini sekarang hostname:port tapi kedepannya pake service id consul

	if err != nil {
		return "", false, err
	}

	if !exists {
		return val, exists, err
	}

	_, err = s.calc.Increment(ctx, val, kind)

	if err != nil {
		return "", false, err
	}

	return val, true, nil
}

// add member pakai zadd
// remove member pakai zref

// zadd setelah berhasil daftarkan ke consul
// zrev saat service nya down

func (s *Selector) Lookup(
	ctx context.Context,
	namespace string,
	kind string,
	selection selector.Selection,
) (string, float64, bool, error) {

	var data []redis.Z
	var err error

	key := fmt.Sprintf("%s:%s", namespace, kind)

	if selection == selector.Min {
		data, err = s.inner.ZRangeWithScores(key, 0, 1).Result()
	}

	if selection == selector.Max {
		data, err = s.inner.ZRevRangeWithScores(key, 0, 1).Result()
	}

	if err != nil {
		return "", 0.0, false, err
	}

	if len(data) != 1 {
		return "", 0.0, false, nil
	}

	val := data[0]

	return val.Member.(string), val.Score, true, nil
}

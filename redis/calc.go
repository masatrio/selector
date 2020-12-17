package redis

import (
	"context"

	"github.com/go-redis/redis/v7"
)

type Calculate struct {
	inner  redis.UniversalClient
	weight map[string]float64
}

func (c *Calculate) Increment(ctx context.Context, name, kind string) (float64, error) {
	return c.inner.ZIncrBy(kind, c.weight[kind], name).Result()

}
func (c *Calculate) Decrement(ctx context.Context, name, kind string) (float64, error) {
	return c.inner.ZIncrBy(kind, -c.weight[kind], name).Result()
}

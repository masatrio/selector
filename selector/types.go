package selector

import (
	"context"
)

type Selection int

const (
	Min Selection = 0
	Max Selection = 1
)

type Calculate interface {
	Increment(
		ctx context.Context,
		name,
		kind string,
	) (float64, error)

	Decrement(
		ctx context.Context,
		name,
		kind string,
	) (float64, error)
}

type Selector interface {
	Select(
		ctx context.Context,
		namespace string,
		kind string,
		s Selection,
	) (string, bool, error)

	Lookup(
		ctx context.Context,
		namespace string,
		kind string,
		s Selection,
	) (string, float64, bool, error)
	// Keys(ctx context.Context, name string, kind string) ([]string, error)
}

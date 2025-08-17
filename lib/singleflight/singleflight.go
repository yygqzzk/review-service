package singleflight

import (
	"context"

	"golang.org/x/sync/singleflight"
)

var g singleflight.Group

func Do(ctx context.Context, fn func(context.Context) (any, error), key string) (any, error, bool) {
	data, err, shared := g.Do(key, func() (any, error) {
		data, err := fn(ctx)
		if err != nil {
			return nil, err
		}
		return data, nil
	})
	if err != nil {
		return nil, err, shared
	}
	return data, nil, shared
}

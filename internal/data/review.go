package data

import (
	"context"

	"github.com/yygqzzk/review-service/internal/biz"
	"github.com/yygqzzk/review-service/internal/data/model"

	"github.com/go-kratos/kratos/v2/log"
)

type reviewRepo struct {
	data *Data
	log  *log.Helper
}

// NewReviewRepo .
func NewReviewRepo(data *Data, logger log.Logger) biz.ReviewRepo {
	return &reviewRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *reviewRepo) SaveReview(ctx context.Context, g *model.ReviewInfo) error {
	return r.data.dbClient.ReviewInfo.WithContext(ctx).Create(g)
}

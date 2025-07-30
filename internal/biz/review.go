package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/yygqzzk/review-service/internal/data/model"
)

type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) error
}

type ReviewUsecase struct {
	repo ReviewRepo
	log  *log.Helper
}

func NewReviewUsecase(repo ReviewRepo, logger log.Logger) *ReviewUsecase {
	return &ReviewUsecase{repo: repo, log: log.NewHelper(logger)}
}

// 实现主要业务逻辑
// 创建评价
func (uc *ReviewUsecase) CreateReview(ctx context.Context, r *model.ReviewInfo) error {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReview: %v \n", r)
	// 1. 数据校验
	// 2. 生成review ID
	// 3. 查询订单和商品快照信息
	// 4. 拼装数据
	// 5. 保存数据
	return uc.repo.SaveReview(ctx, r)
}

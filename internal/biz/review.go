package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/yygqzzk/review-service/api/review/v1"
	"github.com/yygqzzk/review-service/internal/data/model"
	"github.com/yygqzzk/review-service/pkg/snowflake"
)

type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) error
	GetReviewByOrderID(context.Context, int64) ([]*model.ReviewInfo, error)
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
	// 1.1 参数校验 （不应该在biz层实现）
	// 1.2 业务校验
	review, err := uc.repo.GetReviewByOrderID(ctx, r.OrderID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("[biz] GetReviewByOrderID: %v \n", err)
		return pb.ErrorReviewInternalError("数据库查询失败")
	}
	if len(review) > 0 {
		uc.log.WithContext(ctx).Errorf("[biz] 订单: %d 已存在评价", r.OrderID)
		return pb.ErrorReviewAlreadyExists("订单已存在评价")
	}
	// 2. 生成review ID (雪花算法生成唯一ID)
	r.ReviewID = snowflake.GenID()
	// 3. 查询订单和商品快照信息
	// 4. 拼装数据
	// 5. 保存数据
	if err := uc.repo.SaveReview(ctx, r); err != nil {
		uc.log.WithContext(ctx).Errorf("[biz] SaveReview: %v \n", err)
		return pb.ErrorReviewInternalError("创建评价失败")
	}
	return nil
}

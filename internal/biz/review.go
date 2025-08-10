package biz

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/yygqzzk/review-service/api/review/v1"
	"github.com/yygqzzk/review-service/pkg/snowflake"
	"gorm.io/gorm"
)

type ReviewRepo interface {
	SaveReview(ctx context.Context, review *ReviewEntity) error
	GetReviewByOrderID(ctx context.Context, orderId int64) (*ReviewEntity, error)
	SaveReply(ctx context.Context, reply *ReplyEntity) error

	GetReviewById(ctx context.Context, reviewId int64) (*ReviewEntity, error)
	UpdateReviewReplyStatus(ctx context.Context, reviewId int64, status int32) error
	GetAppealByReviewID(ctx context.Context, reviewId int64) (*AppealEntity, error)
	SaveAppeal(ctx context.Context, appeal *AppealEntity) error
	UpdateReviewAuditStatus(ctx context.Context, audit *AuditReviewEntity) error
	UpdateAppealAuditStatus(ctx context.Context, audit *AuditAppealEntity) error
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
func (uc *ReviewUsecase) CreateReview(ctx context.Context, r *ReviewEntity) error {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReview: %v \n", r)
	// 1. 数据校验
	// 1.1 参数校验 （不应该在biz层实现）
	// 1.2 业务校验
	review, err := uc.repo.GetReviewByOrderID(ctx, r.OrderID)
	// 其他查询错误
	if err != nil {
		uc.log.WithContext(ctx).Errorf("[biz] GetReviewByOrderID: %v \n", err)
		return pb.ErrorReviewInternalError("数据库查询失败")
	}
	if review != nil {
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

func (uc *ReviewUsecase) CreateReply(ctx context.Context, r *ReplyEntity) error {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReply: %v \n", r)

	// 业务校验
	// 数据合法性校验 (已回复的评价不允许商家再次回复)
	// 同时进行 水平越权校验 (商家不能回复其他商家的评价)
	review, err := uc.repo.GetReviewById(ctx, r.ReviewID)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("[biz] GetReplyByReviewId: %v \n", err)
		return pb.ErrorReviewInternalError("查询回复失败")
	}
	if review.HasReply == REVIEW_HAS_REPLY {
		uc.log.WithContext(ctx).Errorf("[biz] 评价: %d 已存在回复", r.ReviewID)
	}
	if r.StoreID != review.StoreID {
		uc.log.WithContext(ctx).Errorf("[biz] 商家: %d 不能回复其他商家的评价", r.StoreID)
		return pb.ErrorReviewInternalError("水平越权")
	}
	// 水平越权校验(商家不能回复其他商家的评价)
	// review, err := uc.repo.GetReviewById(ctx, r.ReviewID)
	// if err != nil {
	// 	uc.log.WithContext(ctx).Errorf("[biz] GetReviewById: %v \n", err)
	// 	return pb.ErrorReviewInternalError("查询评价失败")
	// }
	// if review.UserID != r.StoreID {
	// 	uc.log.WithContext(ctx).Errorf("[biz] 商家: %d 不能回复其他商家的评价", r.StoreID)
	// }
	// 事务操作
	// 同时更新评价表及回复表
	r.ReplyID = snowflake.GenID()
	if err := uc.repo.SaveReply(ctx, r); err != nil {
		uc.log.WithContext(ctx).Errorf("[biz] SaveReply: %v \n", err)
		return pb.ErrorReviewInternalError("创建回复失败")
	}

	return nil
}

func (uc *ReviewUsecase) SaveAppeal(ctx context.Context, a *AppealEntity) error {
	uc.log.WithContext(ctx).Debugf("[biz] SaveAppeal: %v \n", a)

	// 查询是否存在对应评价
	review, err := uc.repo.GetReviewById(ctx, a.ReviewID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		uc.log.WithContext(ctx).Errorf("[biz] GetReviewById: %v \n", err)
		return pb.ErrorReviewInternalError("查询评价失败")
	}
	if review == nil {
		uc.log.WithContext(ctx).Errorf("[biz] 评价: %d 不存在", a.ReviewID)
		return pb.ErrorReviewInternalError("评价不存在")
	}
	if review.StoreID != a.StoreID {
		uc.log.WithContext(ctx).Errorf("[biz] 商家: %d 不能申诉其他商家的评价", a.StoreID)
		return pb.ErrorReviewInternalError("评价不属于商家")
	}

	// 查询是否已有申诉
	appeal, err := uc.repo.GetAppealByReviewID(ctx, a.ReviewID)
	// 其他查询错误
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 申诉已存在
	if appeal != nil {
		// 水平越权校验(商家不能申诉其他商家的评价)
		if appeal.StoreID != a.StoreID {
			uc.log.WithContext(ctx).Errorf("[biz] 商家: %d 不能申诉其他商家的评价", a.StoreID)
			return pb.ErrorReviewInternalError("水平越权")
		}
		// 申诉状态已被审核过
		if appeal.Status > 10 {
			return errors.New("该评价已有申诉记录")
		}
		// 共用申诉ID
		a.AppealID = appeal.AppealID
	} else {
		// 生成新的申诉ID
		a.AppealID = snowflake.GenID()
	}

	// 若申诉不存在，则插入申诉数据；若存在，则更新申诉
	if err := uc.repo.SaveAppeal(ctx, a); err != nil {
		uc.log.WithContext(ctx).Errorf("[biz] SaveAppeal: %v \n", err)
		return pb.ErrorReviewInternalError("申诉失败")
	}
	return nil
}

func (uc *ReviewUsecase) AuditReview(ctx context.Context, a *AuditReviewEntity) error {
	uc.log.WithContext(ctx).Debugf("[biz] AuditReview: %v \n", a)

	uc.repo.UpdateReviewAuditStatus(ctx, a)

	return nil
}

func (uc *ReviewUsecase) AuditAppeal(ctx context.Context, a *AuditAppealEntity) error {
	uc.log.WithContext(ctx).Debugf("[biz] AuditAppeal: %v \n", a)

	uc.repo.UpdateAppealAuditStatus(ctx, a)

	return nil
}

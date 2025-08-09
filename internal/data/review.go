package data

import (
	"context"

	"github.com/yygqzzk/review-service/internal/biz"
	"github.com/yygqzzk/review-service/internal/data/model"
	"github.com/yygqzzk/review-service/internal/data/query"
	"gorm.io/gorm/clause"

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

func (r *reviewRepo) SaveReview(ctx context.Context, g *biz.ReviewEntity) error {
	review := &model.ReviewInfo{
		ReviewID:     g.ReviewID,
		UserID:       g.UserID,
		OrderID:      g.OrderID,
		StoreID:      g.StoreID,
		Score:        g.Score,
		ServiceScore: g.ServiceScore,
		ExpressScore: g.ExpressScore,
		Content:      g.Content,
		PicInfo:      g.PicInfo,
		VideoInfo:    g.VideoInfo,
		Anonymous:    g.Anonymous,
		HasReply:     g.HasReply,
	}
	return r.data.dbClient.ReviewInfo.WithContext(ctx).Create(review)
}

func (r *reviewRepo) GetReviewByOrderID(ctx context.Context, orderID int64) (*biz.ReviewEntity, error) {
	review, err := r.data.dbClient.ReviewInfo.WithContext(ctx).Where(r.data.dbClient.ReviewInfo.OrderID.Eq(orderID)).First()
	if err != nil {
		r.log.WithContext(ctx).Errorf("[data] GetReviewByOrderID: %v \n", err)
		return nil, err
	}
	reviewEntity := &biz.ReviewEntity{
		ReviewID:     review.ReviewID,
		UserID:       review.UserID,
		OrderID:      review.OrderID,
		StoreID:      review.StoreID,
		Score:        review.Score,
		ServiceScore: review.ServiceScore,
		ExpressScore: review.ExpressScore,
		Content:      review.Content,
		PicInfo:      review.PicInfo,
		VideoInfo:    review.VideoInfo,
		Anonymous:    review.Anonymous,
		HasReply:     review.HasReply,
	}
	return reviewEntity, nil
}

func (r *reviewRepo) SaveReply(ctx context.Context, g *biz.ReplyEntity) error {
	reply := &model.ReviewReplyInfo{
		ReplyID:   g.ReplyID,
		ReviewID:  g.ReviewID,
		StoreID:   g.StoreID,
		Content:   g.Content,
		PicInfo:   g.PicInfo,
		VideoInfo: g.VideoInfo,
	}

	return r.data.dbClient.Transaction(func(client *query.Query) error {
		err := client.ReviewReplyInfo.WithContext(ctx).Create(reply)
		if err != nil {
			return err
		}
		_, err = client.ReviewInfo.WithContext(ctx).Where(client.ReviewInfo.ReviewID.Eq(g.ReviewID)).Update(client.ReviewInfo.HasReply, biz.REVIEW_HAS_REPLY)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *reviewRepo) GetReplyByReviewIdWithStoreId(ctx context.Context, reviewId int64, storeId int64) (*biz.ReplyEntity, error) {
	reply, err := r.data.dbClient.ReviewReplyInfo.WithContext(ctx).Where(r.data.dbClient.ReviewInfo.ReviewID.Eq(reviewId), r.data.dbClient.ReviewInfo.StoreID.Eq(storeId)).First()
	if err != nil {
		return nil, err
	}
	reviewEntity := &biz.ReplyEntity{
		ReplyID:   reply.ReplyID,
		ReviewID:  reply.ReviewID,
		StoreID:   reply.StoreID,
		Content:   reply.Content,
		PicInfo:   reply.PicInfo,
		VideoInfo: reply.VideoInfo,
	}
	return reviewEntity, nil

}

func (r *reviewRepo) GetReviewById(ctx context.Context, reviewId int64) (*biz.ReviewEntity, error) {
	review, err := r.data.dbClient.ReviewInfo.WithContext(ctx).Where(r.data.dbClient.ReviewInfo.ReviewID.Eq(reviewId)).First()
	if err != nil {
		return nil, err
	}
	reviewEntity := &biz.ReviewEntity{
		ReviewID:     review.ReviewID,
		UserID:       review.UserID,
		OrderID:      review.OrderID,
		StoreID:      review.StoreID,
		Score:        review.Score,
		ServiceScore: review.ServiceScore,
		ExpressScore: review.ExpressScore,
		Content:      review.Content,
		PicInfo:      review.PicInfo,
		VideoInfo:    review.VideoInfo,
		Anonymous:    review.Anonymous,
		HasReply:     review.HasReply,
	}
	r.log.WithContext(ctx).Debugf("[data] GetReviewById: %v \n", reviewEntity)
	return reviewEntity, nil
}

func (r *reviewRepo) UpdateReviewReplyStatus(ctx context.Context, reviewId int64, status int32) error {
	_, err := r.data.dbClient.ReviewInfo.WithContext(ctx).Where(r.data.dbClient.ReviewInfo.ReviewID.Eq(reviewId)).Update(r.data.dbClient.ReviewInfo.HasReply, status)

	return err
}

func (r *reviewRepo) GetAppealByReviewID(ctx context.Context, reviewId int64) (*biz.AppealEntity, error) {
	appeal, err := r.data.dbClient.ReviewAppealInfo.WithContext(ctx).Where(r.data.dbClient.ReviewAppealInfo.ReviewID.Eq(reviewId)).First()
	if err != nil {
		return nil, err
	}
	appealEntity := &biz.AppealEntity{
		AppealID:  appeal.AppealID,
		ReviewID:  appeal.ReviewID,
		StoreID:   appeal.StoreID,
		Reason:    appeal.Reason,
		Content:   appeal.Content,
		PicInfo:   appeal.PicInfo,
		VideoInfo: appeal.VideoInfo,
		Status:    appeal.Status,
	}
	return appealEntity, nil
}

func (r *reviewRepo) SaveAppeal(ctx context.Context, a *biz.AppealEntity) (err error) {
	appeal := &model.ReviewAppealInfo{
		AppealID:  a.AppealID,
		ReviewID:  a.ReviewID,
		StoreID:   a.StoreID,
		Reason:    a.Reason,
		Content:   a.Content,
		PicInfo:   a.PicInfo,
		VideoInfo: a.VideoInfo,
		Status:    10,
	}
	err = r.data.dbClient.ReviewAppealInfo.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "review_id"}}, // 唯一索引校验
			DoUpdates: clause.Assignments(map[string]interface{}{ // 若有冲突则执行更新操作
				"status":     10,
				"content":    appeal.Content,
				"reason":     appeal.Reason,
				"pic_info":   appeal.PicInfo,
				"video_info": appeal.VideoInfo,
			}),
		},
	).Create(appeal)

	r.log.WithContext(ctx).Debugf("[data] SaveAppeal: %v \n", appeal)
	return err
}

func (r *reviewRepo) UpdateReviewAuditStatus(ctx context.Context, a *biz.AuditReviewEntity) error {

	ret, err := r.data.dbClient.ReviewInfo.WithContext(ctx).Where(r.data.dbClient.ReviewInfo.ReviewID.Eq(a.ReviewID)).Updates(
		map[string]interface{}{
			"status":     a.Status,
			"op_user":    a.OpUser,
			"op_reason":  a.OpReason,
			"op_remarks": a.OpRemarks,
		},
	)
	if err != nil {
		r.log.WithContext(ctx).Errorf("[data] UpdateReviewAuditStatus: %v \n", err)
		return err
	}
	r.log.WithContext(ctx).Debugf("[data] UpdateReviewAuditStatus: %v \n", ret)
	return nil
}

func (r *reviewRepo) UpdateAppealAuditStatus(ctx context.Context, a *biz.AuditAppealEntity) error {
	return r.data.dbClient.Transaction(func(client *query.Query) error {
		ret, err := client.ReviewAppealInfo.WithContext(ctx).Where(
			client.ReviewAppealInfo.AppealID.Eq(a.AppealID),
			client.ReviewAppealInfo.ReviewID.Eq(a.ReviewID),
		).Updates(
			map[string]interface{}{
				"status":     a.Status,
				"op_user":    a.OpUser,
				"op_remarks": a.OpRemarks,
			},
		)
		if err != nil {
			r.log.WithContext(ctx).Errorf("[data] UpdateAppealAuditStatus: %v \n", err)
			return err
		}

		if a.Status == 20 {
			_, err = client.ReviewInfo.WithContext(ctx).Where(client.ReviewInfo.ReviewID.Eq(a.ReviewID)).Update(client.ReviewInfo.Status, 40)
			if err != nil {
				r.log.WithContext(ctx).Errorf("[data] UpdateAppealAuditStatus: %v \n", err)
				return err
			}
		}

		r.log.WithContext(ctx).Debugf("[data] UpdateAppealAuditStatus: %v \n", ret)
		return nil
	})
}

package data

import (
	"context"

	"github.com/yygqzzk/review-service/internal/biz"
	"github.com/yygqzzk/review-service/internal/data/model"
	"github.com/yygqzzk/review-service/internal/data/query"

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

func (r *reviewRepo) SaveReplyWithTransaction(ctx context.Context, g *biz.ReplyEntity) error {
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
	return reviewEntity, nil
}

func (r *reviewRepo) UpdateReviewReplyStatus(ctx context.Context, reviewId int64, status int32) error {
	_, err := r.data.dbClient.ReviewInfo.WithContext(ctx).Where(r.data.dbClient.ReviewInfo.ReviewID.Eq(reviewId)).Update(r.data.dbClient.ReviewInfo.HasReply, status)

	return err
}

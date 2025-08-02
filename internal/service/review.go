package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	pb "github.com/yygqzzk/review-service/api/review/v1"
	"github.com/yygqzzk/review-service/internal/biz"

	"github.com/spf13/cast"
)

type ReviewService struct {
	pb.UnimplementedReviewServer
	uc  *biz.ReviewUsecase
	log *log.Helper
}

func NewReviewService(uc *biz.ReviewUsecase, logger log.Logger) *ReviewService {
	return &ReviewService{uc: uc, log: log.NewHelper(logger)}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewReq) (*pb.CreateReviewRsp, error) {
	s.log.WithContext(ctx).Debugf("[service] CreateReview: %v \n", req)
	review := &biz.ReviewEntity{
		UserID:       req.UserId,
		OrderID:      req.OrderId,
		Score:        req.Score,
		ServiceScore: req.ServiceScore,
		ExpressScore: req.ExpressScore,
		Content:      req.Content,
		PicInfo:      req.PicInfo,
		VideoInfo:    req.VideoInfo,
		Anonymous:    cast.ToInt32(req.Anonymous),
	}
	// 调用biz层创建评价
	if err := s.uc.CreateReview(ctx, review); err != nil {
		return nil, err
	}

	// 拼装返回值
	return &pb.CreateReviewRsp{
		ReviewId: review.ReviewID,
	}, nil
}

func (s *ReviewService) ReplyReview(ctx context.Context, req *pb.ReplyReviewReq) (*pb.ReplyReviewRsp, error) {
	s.log.WithContext(ctx).Debugf("[service] ReplyReview: %v \n", req)

	reply := &biz.ReplyEntity{
		ReviewID:  req.ReviewId,
		StoreID:   req.StoreId,
		Content:   req.Content,
		PicInfo:   req.PicInfo,
		VideoInfo: req.VideoInfo,
	}
	if err := s.uc.CreateReply(ctx, reply); err != nil {
		return nil, err
	}

	return &pb.ReplyReviewRsp{
		ReplyId: reply.ReplyID,
	}, nil

}

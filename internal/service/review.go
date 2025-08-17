package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/samber/lo"
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

// 创建评价接口
func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewReq) (*pb.CreateReviewRsp, error) {
	s.log.WithContext(ctx).Debugf("[service] CreateReview: %v \n", req)
	review := &biz.ReviewEntity{
		UserID:       req.UserID,
		OrderID:      req.OrderID,
		StoreID:      req.StoreID,
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
		ReviewID: review.ReviewID,
	}, nil
}

// B端回复评价接口
func (s *ReviewService) ReplyReview(ctx context.Context, req *pb.ReplyReviewReq) (*pb.ReplyReviewRsp, error) {
	s.log.WithContext(ctx).Debugf("[service] ReplyReview: %v \n", req)

	reply := &biz.ReplyEntity{
		ReviewID:  req.ReviewID,
		StoreID:   req.StoreID,
		Content:   req.Content,
		PicInfo:   req.PicInfo,
		VideoInfo: req.VideoInfo,
	}
	if err := s.uc.CreateReply(ctx, reply); err != nil {
		return nil, err
	}

	return &pb.ReplyReviewRsp{
		ReplyID: reply.ReplyID,
	}, nil

}

// C端获取评价详情
func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewReq) (*pb.GetReviewRsp, error) {
	return nil, nil
}

// C端查看userID下所有评价
func (s *ReviewService) ListReviewByUserID(context.Context, *pb.ListReviewByUserIDReq) (*pb.ListReviewByUserIDRsp, error) {
	return nil, nil
}

// B端申诉评价
func (s *ReviewService) AppealReview(ctx context.Context, req *pb.AppealReviewReq) (*pb.AppealReviewRsp, error) {
	s.log.WithContext(ctx).Debugf("[service] AppealReview: %v \n", req)
	appeal := &biz.AppealEntity{
		ReviewID:  req.ReviewID,
		StoreID:   req.StoreID,
		Reason:    req.Reason,
		Content:   req.Content,
		PicInfo:   req.PicInfo,
		VideoInfo: req.VideoInfo,
	}
	if err := s.uc.SaveAppeal(ctx, appeal); err != nil {
		return nil, err
	}

	return &pb.AppealReviewRsp{
		AppealID: appeal.AppealID,
	}, nil
}

// O端审核评价
func (s *ReviewService) AuditReview(ctx context.Context, req *pb.AuditReviewReq) (*pb.AuditReviewRsp, error) {
	s.log.WithContext(ctx).Debugf("[service] AuditReview: %v \n", req)

	audit := &biz.AuditReviewEntity{
		ReviewID:  req.ReviewID,
		Status:    req.Status,
		OpUser:    req.OpUser,
		OpReason:  req.OpReason,
		OpRemarks: *req.OpRemarks,
	}

	if err := s.uc.AuditReview(ctx, audit); err != nil {
		return nil, err
	}

	return &pb.AuditReviewRsp{}, nil
}

// O端评价申诉审核
func (s *ReviewService) AuditAppeal(ctx context.Context, req *pb.AuditAppealReq) (*pb.AuditAppealRsp, error) {

	s.log.WithContext(ctx).Debugf("[service] AuditAppeal: %v \n", req)

	appeal := &biz.AuditAppealEntity{
		AppealID:  req.AppealID,
		ReviewID:  req.ReviewID,
		Status:    req.Status,
		OpUser:    req.OpUser,
		OpRemarks: *req.OpRemarks,
	}

	if err := s.uc.AuditAppeal(ctx, appeal); err != nil {
		return nil, err
	}

	return &pb.AuditAppealRsp{}, nil
}

func (s *ReviewService) ListReviewByStoreID(ctx context.Context, req *pb.ListReviewByStoreIDReq) (*pb.ListReviewByStoreIDRsp, error) {
	s.log.WithContext(ctx).Debugf("[service] ListReviewByStoreID: %v \n", req)
	reviewList, err := s.uc.ListReviewByStoreID(ctx, req.StoreID, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	list := lo.Map(reviewList, func(item *biz.ReviewEntity, _ int) *pb.ReviewInfo {
		return &pb.ReviewInfo{
			ReviewID:     item.ReviewID,
			UserID:       item.UserID,
			OrderID:      item.OrderID,
			Score:        item.Score,
			ServiceScore: item.ServiceScore,
			ExpressScore: item.ExpressScore,
			Content:      item.Content,
			PicInfo:      item.PicInfo,
			VideoInfo:    item.VideoInfo,
		}
	})

	return &pb.ListReviewByStoreIDRsp{
		List: list,
	}, nil
}

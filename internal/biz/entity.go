package biz

const (
	REVIEW_HAS_REPLY = 1 // 已回复
	REVIEW_NO_REPLY  = 0 // 未回复
)

type ReplyEntity struct {
	ReplyID   int64  // 回复ID
	ReviewID  int64  // 评价ID
	StoreID   int64  // 商家ID
	Content   string // 回复内容
	PicInfo   string // 回复图片
	VideoInfo string // 回复视频
}

type ReviewEntity struct {
	ReviewID     int64  // 评价ID
	UserID       int64  // 用户ID
	OrderID      int64  // 订单ID
	StoreID      int64  // 商家ID
	Score        int32  // 评分
	ServiceScore int32  // 服务评分
	ExpressScore int32  // 物流评分
	Content      string // 评价内容
	PicInfo      string // 评价图片
	VideoInfo    string // 评价视频
	Anonymous    int32  // 是否匿名
	HasReply     int32  // 是否已回复
}

type AppealEntity struct {
	AppealID  int64  // 申诉ID
	ReviewID  int64  // 评价ID
	StoreID   int64  // 商家ID
	Reason    string // 申诉原因
	Content   string // 申诉内容
	PicInfo   string // 申诉图片
	VideoInfo string // 申诉视频
	Status    int32  // 申诉状态
}

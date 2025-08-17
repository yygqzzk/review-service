package biz

const (
	REVIEW_HAS_REPLY = 1 // 已回复
	REVIEW_NO_REPLY  = 0 // 未回复
)

type ReplyEntity struct {
	ReplyID   int64  `json:"reply_id,string"`  // 回复ID
	ReviewID  int64  `json:"review_id,string"` // 评价ID
	StoreID   int64  `json:"store_id,string"`  // 商家ID
	Content   string // 回复内容
	PicInfo   string `json:"pic_info"`   // 回复图片
	VideoInfo string `json:"video_info"` // 回复视频
}

type ReviewEntity struct {
	ReviewID     int64  `json:"review_id,string"`     // 评价ID
	UserID       int64  `json:"user_id,string"`       // 用户ID
	OrderID      int64  `json:"order_id,string"`      // 订单ID
	StoreID      int64  `json:"store_id,string"`      // 商家ID
	Score        int32  `json:"score,string"`         // 评分
	ServiceScore int32  `json:"service_score,string"` // 服务评分
	ExpressScore int32  `json:"express_score,string"` // 物流评分
	Content      string `json:"content"`              // 评价内容
	PicInfo      string `json:"pic_info"`             // 评价图片
	VideoInfo    string `json:"video_info"`           // 评价视频
	Anonymous    int32  `json:"anonymous,string"`     // 是否匿名
	HasReply     int32  `json:"has_reply,string"`     // 是否已回复
}

type AppealEntity struct {
	AppealID  int64  `json:"appeal_id,string"`  // 申诉ID
	ReviewID  int64  `json:"review_id,string"`  // 评价ID
	StoreID   int64  `json:"store_id,string"`   // 商家ID
	Reason    string `json:"reason,string"`     // 申诉原因
	Content   string `json:"content,string"`    // 申诉内容
	PicInfo   string `json:"pic_info,string"`   // 申诉图片
	VideoInfo string `json:"video_info,string"` // 申诉视频
	Status    int32  `json:"status,string"`     // 申诉状态
}

type AuditReviewEntity struct {
	ReviewID  int64  `json:"review_id,string"` // 评价ID
	Status    int32  `json:"status,string"`    // 审核状态
	OpUser    string `json:"op_user"`          // 操作人
	OpReason  string `json:"op_reason"`        // 操作原因
	OpRemarks string `json:"op_remarks"`       // 操作备注
}

type AuditAppealEntity struct {
	AppealID  int64  `json:"appeal_id,string"` // 申诉ID
	ReviewID  int64  `json:"review_id,string"` // 评价ID
	Status    int32  `json:"status,string"`    // 审核状态
	OpUser    string `json:"op_user"`          // 操作人
	OpRemarks string `json:"op_remarks"`       // 操作备注
}

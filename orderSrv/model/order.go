package model

import "time"

type Cart struct {
	BaseID
	UserID  int32 `json:"user_id" gorm:"column:user_id;index:idx_user_id;comment:'用户id'"`
	GoodsID int32 `json:"goods_id" gorm:"column:goods_id;index:idx_goods_id;comment:'商品id'"`
	// TODO 的部分
	GoodsImg string `json:"goods_img" gorm:"column:goods_img;type:varchar(255);comment:'商品的图片'"`
	// TODO 的部分
	Nums    int32 `json:"nums" gorm:"column:nums;comment:'数量'"`
	Checked bool  `json:"checked" gorm:"column:checked;comment:'是否选中'"`
	BaseModel
}

// 看了之后自己设计吧。不要太依赖视频

// Order 订单表
type Order struct {
	BaseID
	UserID       int32      `json:"user_id" gorm:"column:user_id;comment:'用户id'"`
	OrderSN      string     `json:"order_sn" gorm:"column:order_sn;type:varchar(64);comment:'系统创建时的订单号'"`
	PayType      string     `json:"pay_type" gorm:"column:pay_type;type:varchar(64);comment:'alipay wechat'"`
	Status       string     `json:"status" gorm:"column:status;type:varchar(64);comment:'订单的支付状态：PAYING(待支付),TRADE_SUCCESS(支付成功),TRADE_CLOSED(超时关闭),WAITING_BUYER_PAY(交易创建),TRADE_FINISHED(交易完成)'"`
	TradeNo      string     `json:"trade_no" gorm:"column:trade_no;type:varchar(64);comment:'交易订单号'"`
	SubjectTitle string     `json:"subject_title" gorm:"column:subject_title;type:varchar(64);comment:'商品标题'"`
	OrderPrice   float32    `json:"amount_price" gorm:"column:amount_price;comment:'订单金额'"`
	FinalPrice   float32    `json:"final_price" gorm:"column:final_price;comment:'实际支付金额'"`
	PayTime      *time.Time `json:"pay_time" gorm:"column:pay_time;comment:'支付时间'"`

	Address         string `json:"address" gorm:"column:address;type:varchar(64);comment:'地址'"`
	RecipientName   string `json:"signer_name" gorm:"column:signer_name;type:varchar(64);comment:'收件人名字'"`
	RecipientMobile string `json:"signer_mobile" gorm:"column:signer_mobile;type:varchar(64);comment:'收件人电话'"`
	Message         string `json:"message" gorm:"column:message;type:varchar(255);comment:'留言'"`
	Snapshot        string `json:"snapshot" gorm:"column:snapshot;type:text;comment:'支付快照。支付时所包含的信息'"`
	BaseModel
}

// OrderGoods 订单-商品详细记录表
type OrderGoods struct {
	BaseID
	OrderId    int32      `json:"order_id" gorm:"column:order_id;comment:'订单id'"`
	OrderSN    string     `json:"order_sn" gorm:"column:order_sn;type:varchar(64);comment:'订单号'"`
	TradeNo    string     `json:"trade_no" gorm:"column:trade_no;type:varchar(64);comment:'第三方订单号'"`
	GoodsId    int32      `json:"goods_id" gorm:"column:goods_id;comment:'商品id'"`
	GoodsPrice float32    `json:"goods_price" gorm:"column:goods_price;comment:'商品价格'"`
	PayPrice   float32    `json:"pay_price" gorm:"column:pay_price;comment:'实际支付金额'"`
	PayTime    *time.Time `json:"pay_time" gorm:"column:pay_time;comment:'支付时间'"`
	GoodsName  string     `json:"goods_name" gorm:"column:goods_name;type:varchar(255);comment:'商品名称'"`
	Nums       int32      `json:"nums" gorm:"column:nums;comment:'数量'"`
	GoodsImg   string     `json:"goods_img" gorm:"column:goods_img;type:varchar(255);comment:'商品图片'"`
	BaseModel
}

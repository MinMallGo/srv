package model

import "time"

type Cart struct {
	BaseID
	UserID  int32 `json:"userId" gorm:"column:userId;index:idx_user_id;comment:'用户id'"`
	GoodsID int32 `json:"goodsId" gorm:"column:goodsId;index:idx_goods_id;comment:'商品id'"`
	// TODO 的部分
	GoodsImg string `json:"goodsImg" gorm:"column:goodsImg;type:varchar(255);comment:'商品的图片'"`
	// TODO 的部分
	Nums    int32 `json:"nums" gorm:"column:nums;comment:'数量'"`
	Checked bool  `json:"checked" gorm:"column:checked;comment:'是否选中'"`
	BaseModel
}

// 看了之后自己设计吧。不要太依赖视频

// Order 订单表
type Order struct {
	UserID       int32     `json:"userId" gorm:"column:userId;comment:'用户id'"`
	OrderSN      string    `json:"orderSN" gorm:"column:orderSN;type:varchar(64);comment:'系统创建时的订单号'"`
	PayType      string    `json:"payType" gorm:"column:payType;type:varchar(64);comment:'alipay wechat'"`
	Status       string    `json:"status" gorm:"column:status;type:varchar(64);comment:'订单的支付状态：PAYING(待支付),TRADE_SUCCESS(支付成功),TRADE_CLOSED(超时关闭),WAITING_BUYER_PAY(交易创建),TRADE_FINISHED(交易完成)'"`
	TradeNo      string    `json:"tradeNo" gorm:"column:tradeNo;type:varchar(64);comment:'交易订单号'"`
	SubjectTitle string    `json:"subjectTitle" gorm:"column:subjectTitle;type:varchar(64);comment:'商品标题'"`
	OrderPrice   float32   `json:"amountPrice" gorm:"column:amountPrice;comment:'订单金额'"`
	FinalPrice   float32   `json:"finalPrice" gorm:"column:finalPrice;comment:'实际支付金额'"`
	PayTime      time.Time `json:"payTime" gorm:"column:payTime;comment:'支付时间'"`

	Address         string `json:"address" gorm:"column:address;type:varchar(64);comment:'地址'"`
	RecipientName   string `json:"signerName" gorm:"column:signerName;type:varchar(64);comment:'收件人名字'"`
	RecipientMobile string `json:"signerMobile" gorm:"column:signerMobile;type:varchar(64);comment:'收件人电话'"`
	Message         string `json:"message" gorm:"column:message;type:varchar(255);comment:'留言'"`
	Snapshot        string `json:"snapshot" gorm:"column:snapshot;type:text;comment:'支付快照。支付时所包含的信息'"`
}

// OrderGoods 订单-商品详细记录表
type OrderGoods struct {
	BaseID
	OrderId    int32     `json:"orderId" gorm:"column:orderId;comment:'订单id'"`
	OrderSN    string    `json:"orderSN" gorm:"column:orderSN;type:varchar(64);comment:'订单号'"`
	TradeNo    string    `json:"tradeNo" gorm:"column:tradeNo;type:varchar(64);comment:'第三方订单号'"`
	GoodsId    int32     `json:"goodsId" gorm:"column:goodsId;comment:'商品id'"`
	GoodsPrice float32   `json:"goodsPrice" gorm:"column:goodsPrice;comment:'商品价格'"`
	PayPrice   float32   `json:"payPrice" gorm:"column:payPrice;comment:'实际支付金额'"`
	PayTime    time.Time `json:"payTime" gorm:"column:payTime;comment:'支付时间'"`
	GoodsName  string    `json:"goodsName" gorm:"column:goodsName;type:varchar(255);comment:'商品名称'"`
	Nums       int32     `json:"nums" gorm:"column:nums;comment:'数量'"`
	GoodsImg   string    `json:"goodsImg" gorm:"column:goodsImg;type:varchar(255);comment:'商品图片'"`
	BaseModel
}

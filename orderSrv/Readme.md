## 订单服务

### 定义表结构

1. 购物车表
2. 订单表
3. 订单详细表

### proto文件定义

1. 购物车
    1. 获取用户购物车详细
    2. 添加商品到购物车
    3. 移除商品
    4. 选中商品
    5. 更新商品数量
2. 订单
    1. 创建订单
    2. 订单列表 // user_id | 分页
    3. 获取订单详情：order_sn | id

### 跨服务调用
1. 居然是把其他服务的proto文件放过来，然后在global里面连接服务
2. TODO 跨服务调用需要弄分布式事务。这个还没讲到

### TODO完成订单部分，以及联调
完成。但是没有前端代码没法联调

### 关于分布式事务一致性的问题。
1. 教程里面用的rocketmq的事务消息来保证一致性，这个要搞清楚原理手写
2. 看看其他的分布式事务是什么玩的，比如go-zero,奎托斯

### gorm 关于嵌套成json
相当于自定义了一个类型，对于用这个类型的字段而言，需要填充Detail的内容，然后它会被序列化成json存在数据库
```go 
package main

type Detail struct {
	ID int 
	Name string
}

type AType []Detail

func (g *AType) Scan(value interface{}) error {
   return json.Unmarshal(value.([]byte), g)
}

func (g AType) Value() (driver.Value, error) {
   return json.Marshal(g)
}

```

### 明天继续完成分布式事务一致性的问题


### 关于链路追踪的问题，要么是从span里面拿
### 如果是从消息队列里面获取，就拿出来数据，构造成类似的api层调用的context就好了，把所有需要的内容全部写到消息头里面就好了
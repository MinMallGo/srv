# 微服务拆分的核心

## protobuf 和grpc通信

### protobuf
1. 定义接口和消息
2. 可以通过protoc工具生成需要的语言的代码
3. 通过服务可以实现跨语言之间的调用

### grpc
1. 自定义的消息编解码格式，加快传输速度
2. 通过HTTP2加快传输

## 注册中心和服务发现

## 配置中心

## 分布式事务一致性
### 抽象思路图
```text 
[订单服务]
   |
   | BeginTransaction -> RocketMQ 半消息(order_created)
   |
   | 本地事务: 创建订单 + 清购物车
   | COMMIT事务消息(order_created)
   |--> Send 延迟消息(order_cancel, 30min)
   |
   |<-- [order_error] ← 来自库存、优惠券服务失败时发送
   ↓
  标记订单作废

[库存服务]
   |
   |<-- [order_created] ← 消费事务消息
   |--> 扣减库存（幂等）
   |--> fail 发送 order_error

[优惠券服务 / 积分服务]
   |
   |<-- 同上
   |--> 使用资源
   |--> fail 发送 order_error

[支付服务]
   |
   |<-- [order_created]
   |--> 成功支付后，发送 order_paid
   |
   |<-- [order_cancel]
   |--> 如果未支付，做库存归还

[消息体字段：orderSN + 商品快照 + 用户ID + 下单时间 + 状态变更]
```

## 高可用
### 服务雪崩以及解决方案（超时机制）
1. 因为某个服务不可以用导致的后续服务都不可用的问题
2. ![img.png](img.png)
3. 通过超时机制来控制，避免长时间等待导致不可用

### 超时，重试，幂等
1. 可能因为网络原因导致失败，可以重试一下
2. 幂等即：相同的请求，只能处理一次。比如同一个订单的支付，同一个用户的创建等


## 链路追踪
1. openTelemetry + jaeger
2. [用法](https://github.com/ucdo/everydayNormalGo2/tree/main/Ladon/tracer)
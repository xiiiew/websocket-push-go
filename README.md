# websocket-push-go

一套基于golang的websocket订阅、推送服务。客户端使用websocket连接至服务，并可订阅任意频道。业务端可使用HTTP和RPC接口推送消息至一个或者所有频道。

## websocket消息约定

* 订阅频道

```json
{"action":"sub","ch":"ch_name"}
```

* 取消订阅频道

```json
{"action":"unsub","ch":"ch_name"}
```

* ping

```json
{"ping":1605751068}
```

* pong

```json
{"pong":1605751068}
```

## HTTP推送

HTTP推送接口面向业务端。业务端将消息封装到Request Body中，调用API接口，推送至相应WS连接。

#### API

* 推送到频道

`POST http://localhost:2181/push/ch/ch_name`

* 推送到所有连接

`POST http://localhost:2181/push/all`

## RPC推送

使用RPC接口推送对比使用HTTP接口推送有较大的性能优势，建议其他业务端使用此类接口推送消息。该项目使用gRPC，`proto`文件位于`websocket-push-go/protos/`中，其他语言可通过相应工具生成对应的代理类。

#### API

###### Golang版

* 推送到频道

```go
PushCh(ctx context.Context, in *PushChRequest, opts ...grpc.CallOption) (*PushChReply, error)
```

* 推送到所有连接

```go
PushAll(ctx context.Context, in *PushAllRequest, opts ...grpc.CallOption) (*PushAllReply, error)
```

## 待新增功能

* 服务状态监控及websocket连接统计
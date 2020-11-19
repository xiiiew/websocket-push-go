package handles

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/xiiiew/websocket-push-go/common"
	pb "github.com/xiiiew/websocket-push-go/protos"
	"github.com/xiiiew/websocket-push-go/response"
	"github.com/xiiiew/websocket-push-go/server"
)

type WsPushServer struct {
}

// 推送到频道
func (s *WsPushServer) PushCh(ctx context.Context, in *pb.PushChRequest) (*pb.PushChReply, error) {
	ch := in.Ch
	message := common.WsMessageBuilder(websocket.TextMessage, in.Message)
	if ch == "" {
		return &pb.PushChReply{
			Message: response.HttpErrorResponseBuilder("channel cannot be empty"),
		}, nil
	}
	if message == nil {
		return &pb.PushChReply{
			Message: response.HttpErrorResponseBuilder("message cannot be empty"),
		}, nil
	}

	server.GetBucketInstance().PushCh(ch, message)
	return &pb.PushChReply{
		Message: response.HttpSuccessResponseBuilder(nil),
	}, nil
}

// 推送给所有用户
func (s *WsPushServer) PushAll(ctx context.Context, in *pb.PushAllRequest) (*pb.PushAllReply, error) {
	message := common.WsMessageBuilder(websocket.TextMessage, in.Message)
	if message == nil {
		return &pb.PushAllReply{
			Message: response.HttpErrorResponseBuilder("message cannot be empty"),
		}, nil
	}

	server.GetBucketInstance().PushAll(message)
	return &pb.PushAllReply{
		Message: response.HttpSuccessResponseBuilder(nil),
	}, nil
}

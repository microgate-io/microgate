package queue

import (
	"context"

	"github.com/emicklei/xconnect"
	apiqueue "github.com/microgate-io/microgate-lib-go/v1/queue"
	mlog "github.com/microgate-io/microgate/v1/log"
)

type QueueingServiceImpl struct {
	apiqueue.UnimplementedQueueingServiceServer
}

func NewQueueingServiceImpl(config xconnect.Document) *QueueingServiceImpl {
	return new(QueueingServiceImpl)
}

func (s *QueueingServiceImpl) Publish(ctx context.Context, req *apiqueue.PublishRequest) (*apiqueue.PublishResponse, error) {
	mlog.Debugw(ctx, "Publish", "req", req)
	return nil, nil
}
func (s *QueueingServiceImpl) Subscribe(ctx context.Context, req *apiqueue.SubscribeRequest) (*apiqueue.SubscribeResponse, error) {
	mlog.Debugw(ctx, "Subscribe", "req", req)
	return nil, nil
}

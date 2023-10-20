package listener

import (
	"context"

	"github.com/outofoffice3/common/logger"
	s3deh "github.com/outofoffice3/neorestore/internal/interface/s3-delete-event-handler"
	"github.com/outofoffice3/neorestore/pkg/types"
)

var (
	sos  logger.Logger
	s3dh s3deh.S3DeleteEventHandler
)

func Init() {
	sos = logger.NewConsoleLogger(logger.LogLevelDebug)
	s3dh = s3deh.NewS3DeleteEventHandler(context.Background())
	sos.Infof("init listener")
}

func Handle(ctx context.Context, event types.ListenerEvent) (types.ListenerResponse, error) {
	result, err := s3dh.HandleEvent(ctx, event)
	if err != nil {
		return types.ListenerResponse{}, err
	}
	// type assert response to listener response
	listenerResponse, ok := result.(types.ListenerResponse)
	if !ok {
		return types.ListenerResponse{}, err
	}
	sos.Debugf("Listener Response [%+v]", listenerResponse)

	return listenerResponse, nil
}

package rdispatch

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/outofoffice3/common/logger"
)

var (
	sos logger.Logger
)

type RestoreJob struct {
	jobId  string
	Input  *s3.DeleteObjectsInput
	client *s3.Client
}

func Init() {
	sos = logger.NewConsoleLogger(logger.LogLevelDebug)
	sos.Infof("init restore job success")
}

func NewRestoreJob(client *s3.Client) *RestoreJob {
	return &RestoreJob{
		jobId:  "",
		client: client,
	}
}

func (rj *RestoreJob) Do() {
	result, err := rj.client.DeleteObjects(context.Background(), rj.Input)
	if err != nil {
		panic(err)
	}
	if result.Deleted != nil {
		for _, deleted := range result.Deleted {
			sos.Debugf("deleted : [%+v]", deleted)
		}
	}
	if result.Errors != nil {
		for _, error := range result.Errors {
			sos.Debugf("error : [%+v]", error)
		}
	}

}

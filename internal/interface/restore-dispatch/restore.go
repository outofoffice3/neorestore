package rdispatch

import (
	"github.com/YSZhuoyang/go-dispatcher/dispatcher"
)

type RestoreDispatcher struct {
	dispatcher dispatcher.Dispatcher
}

func NewRestoreDispatcher() *RestoreDispatcher {
	dispatcher, err := dispatcher.NewDispatcher(4)
	if err != nil {
		panic(err)
	}
	return &RestoreDispatcher{
		dispatcher: dispatcher,
	}
}

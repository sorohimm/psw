package infrastructure

import (
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"net/http"
	"psw_server/internal/controllers"
)

type IInjector interface {
	InjectBalanceController() controllers.QueueController
}

var env *environment

type environment struct {
	logger *zap.SugaredLogger
	client *http.Client
	nc     *nats.Conn
	js     nats.JetStreamContext
	sub    *nats.Subscription
}

func (e *environment) InjectBalanceController() controllers.QueueController {
	return controllers.QueueController{
		Log:      e.logger,
		NatsConn: e.nc,
		Js:       e.js,
		Sub:      e.sub,
	}
}

func Injector(log *zap.SugaredLogger, nc *nats.Conn, js nats.JetStreamContext, sub *nats.Subscription) (IInjector, error) {
	env = &environment{
		logger: log,
		client: http.DefaultClient,
		nc:     nc,
		js:     js,
		sub:    sub,
	}

	return env, nil
}

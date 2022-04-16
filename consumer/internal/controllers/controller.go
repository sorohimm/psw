package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"psw_server/internal/models"
)

type QueueController struct {
	Log      *zap.SugaredLogger
	NatsConn *nats.Conn
	Js       nats.JetStreamContext
	Sub      *nats.Subscription
}

func (c *QueueController) GetAllQuotes(ctx *gin.Context) {
	msgs, err := c.Sub.Fetch(100) // fetch all available messages (msg limit 100)
	switch err {
	case errors.Cause(nats.ErrTimeout):
		c.Log.Warn("timeout")
		ctx.JSON(http.StatusNotFound, fmt.Sprintf("error: stream is empty"))
		return
	default:
		if err != nil {
			c.Log.Warn("fetch error ", err.Error())
			ctx.JSON(http.StatusInternalServerError, "internal server error")
			return
		}
	}

	err = c.Js.PurgeStream("QUOTES") // purge stream
	if err != nil {
		c.Log.Warn("purge error")
		ctx.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	var res []models.QuotesMsg

	if msgs != nil {
		for _, msg := range msgs {
			var msgData models.QuotesMsg
			json.Unmarshal(msg.Data, &msgData)
			res = append(res, msgData)
		}
	} else {
		ctx.JSON(http.StatusNotFound, "error: no messages")
		return
	}

	ctx.JSON(http.StatusOK, res)
}

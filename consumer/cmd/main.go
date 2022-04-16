package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"os"
	"psw_server/internal/infrastructure"
)

var (
	log *zap.SugaredLogger
)

const (
	subSubjectName = "QUOTES.publish"
)

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("error loading logger: %s", err)
		os.Exit(1)
		return
	}

	log = logger.Sugar()
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, _ := nc.JetStream()
	sub, err := js.PullSubscribe(subSubjectName, "MONITOR")

	injector, err := infrastructure.Injector(log, nc, js, sub)
	if err != nil {
		log.Fatal("main :: inject failing")
	}

	balanceController := injector.InjectBalanceController()

	router := gin.Default()

	v1 := router.Group("/quotes/v1")
	{
		v1.GET("/all", balanceController.GetAllQuotes)
	}

	err = router.Run(":8081")
	if err != nil {
		log.Fatal("main: router deployment error")
	}
}

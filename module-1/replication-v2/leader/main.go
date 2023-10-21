package main

import (
  "os"
  "strings"
  "strconv"

  "github.com/gin-gonic/gin"
  "github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v2/leader/internal/leader"
  "go.uber.org/zap"
)

type MessageRequest struct {
  Message string
  W int
}

func main() {
  logger, _ := zap.NewDevelopment()
  r := gin.Default()

  w, err :=  strconv.Atoi(os.Getenv("W"))
  if err != nil {
    logger.Fatal("Can't convert value of env var W into int")
    return
  }
  replicas := strings.Split(os.Getenv("REPLICAS"), ","),

  l := leader.NewLeader(
    replicas,
    logger,
  )

  r.GET("logs", func (c *gin.Context) {
    c.JSON(200, l.GetLogs())
  })

  r.POST("message", func (c *gin.Context) {
    m := MessageRequest{}
    c.Bind(&m)

    l.AddMessage(m.Message, m.W)
    c.Status(200)
  })

  r.Run()
}

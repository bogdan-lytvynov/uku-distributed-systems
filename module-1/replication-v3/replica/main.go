package main

import (
  "os"
  "strconv"
  "time"
  "net"
  "net/rpc"
  "net/http"
  "github.com/gin-gonic/gin"
  "github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v1/proto"
  "go.uber.org/zap"
)

type Replica struct{
  logger *zap.Logger
  logs []string
  pendingMessages []*proto.ReplicateMessage
  delay int
}

func NewReplica (delay int, logger *zap.Logger) *Replica {
  return &Replica{
    logger: logger,
    delay: delay,
  }
}

func (r *Replica) maybeDelay() {
  if r.delay != 0 {
    r.logger.Info("Start delay")
    time.Sleep(time.Duration(r.delay) * time.Second)
    r.logger.Info("End delay")
  }
}

//insert sort into pending messages queue
func (r *Replica) addPendingMessage(m *proto.ReplicateMessage) {

}

// check if there are pending messages we can process
func (r *Replica) reviewPendingMessages() {

}

func (r *Replica) Replicate(m *proto.ReplicateMessage, reply *proto.ReplicateMessageReply) error {
  r.maybeDelay()

  r.logger.Info("Add message to replica log",
    zap.Int("order", m.Order),
    zap.String("message", m.Message),
  )
  //if we received next message in order - add it to the logs
  if  m.Order == len(r.logs) {
    r.logs = append(r.logs, m.Message)
    r.reviewPendingMessages()
  } else if m.Order > len(r.logs) {
    r.addPendingMessage(m)
  }

  reply.Ack = true
  return nil
}

func (r *Replica) GetLogs() []string {
  return r.logs
}

func startRpcServer(r *Replica, logger *zap.Logger) {
  //start RPC server
  rpc.Register(r)
  rpc.HandleHTTP()
  logger.Info("Start RPC server")
  l, e := net.Listen("tcp", ":3001")
  if e != nil {
    logger.Error("Failed to net.Listen", zap.Error(e))
  }
  go http.Serve(l, nil)
}

func startHttpServer(r *Replica) {
  e := gin.Default()

  e.GET("logs", func (c *gin.Context) {
    c.JSON(200, r.GetLogs())
  })

  e.Run()
}

func parseDelay(delayStr string) int {
  if delayStr == "" {
    return 0
  }
  d, err := strconv.Atoi(delayStr)
  if err != nil {
    return 0
  }
  return d
}

func main() {
  logger, _ := zap.NewDevelopment()
  logger.Info("Start replica")
  r := NewReplica(parseDelay(os.Getenv("DELAY")), logger)
  startRpcServer(r, logger)
  startHttpServer(r)
}

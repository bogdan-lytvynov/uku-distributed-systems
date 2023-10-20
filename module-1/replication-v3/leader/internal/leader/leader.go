package leader

import (
  "errors"
  "time"
  "sync"
  "net/rpc"
  "go.uber.org/zap"
  "github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v2/proto"
)

var replicaClientTimeout = errors.New("ReplicaClientTimeout")
var defaultReplicaClientTimeout = time.Duration(5) // 5 Seconds

type Leader struct {
  order int
  mx sync.Mutex
  logs []string
  logger *zap.Logger
  replicaClients []ReplicaClient
  minAcks int
}

func NewLeader(replicas []string, minAcks int, logger *zap.Logger) Leader {
  replicaClients := []ReplicaClient{}
  for _, r := range replicas {
    client := NewReplicaClient(r, defaultReplicaClientTimeout, logger)
    replicaClients = append(replicaClients, client)
  }

  return Leader{
    logger: logger,
    logs: []string{},
    replicaClients: replicaClients,
    minAcks: minAcks,
  }
}

type ReplicaClient struct {
  address string
  timeout time.Duration
  logger *zap.Logger
}

func NewReplicaClient(address string, timeout time.Duration, logger *zap.Logger) ReplicaClient {
  return ReplicaClient{
    address,
    timeout,
    logger,
  }
}

func (rc *ReplicaClient) Replicate(order int, message string) error {
  client ,err := rpc.DialHTTP("tcp", rc.address) 
  defer client.Close()

  if err != nil {
    rc.logger.Fatal("Failed to rpc.DialHTTP", zap.Error(err))
    return err
  }
  args := &proto.ReplicateMessage{
    Order: order,
    Message: message,
  }
  reply := &proto.ReplicateMessageReply{}
  replicateCall := client.Go("Replica.Replicate", args, reply, nil)

  select {
  case <- replicateCall.Done:
    if replicateCall.Error != nil {
      return replicateCall.Error
    }
    return nil
  case <-time.After(rc.timeout * time.Second):
    rc.logger.Info("Replica didn't respond")
    return replicaClientTimeout
  }
}

func (l *Leader) GetLogs() []string {
  return l.logs
}

func (l *Leader) AddMessage(m string) error {
  l.mx.Lock()
  l.order++
  l.mx.Unlock()
  l.replicateToAll(l.order, m)
  l.logs[l.order] = m 
  return nil
}

func (l *Leader) replicateToAll(order int, m string) error {
  l.logger.Info("Start replication from leader to all the replicas", zap.String("message", m))
  c := make(chan bool, len(l.replicaClients))

  for _, rc := range l.replicaClients {
    go func(rc ReplicaClient) {
      // repeat until delivered
      // TODO: make it smart backoof exponential repeat
      for {
        err := rc.Replicate(order, m) 
        if err == nil {
          c <- true
          return
        }
      }
    }(rc)
  }

  ackCount := 0
  for ackCount < l.minAcks {
    <- c
    ackCount++
    l.logger.Debug("Replica acked", zap.Int("ack count", ackCount), zap.Int("minAck", l.minAcks))
  }

  l.logger.Info("End replication from leader to all replicas", zap.String("message", m))
  return nil
}

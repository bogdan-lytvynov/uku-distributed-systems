package leader

import (
  "time"
  "sync"
  "net/rpc"
  "go.uber.org/zap"
  "github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v2/proto"
)
var defaultReplicaClientTimeout time.Duration = 5

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

func (rc *ReplicaClient) replicate(order int, message string) error {
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
  }
}

type pendingMessage struct {
  order int
  message string
}

type Leader struct {
  lastOrder int
  orderMX sync.Mutex

  logs []string
  logsMX sync.Mutex

  pending []pendingMessage
  logger *zap.Logger
  replicaClients []ReplicaClient
}

func NewLeader(replicas []string, logger *zap.Logger) Leader {
  replicaClients := []ReplicaClient{}
  for _, r := range replicas {
    client := NewReplicaClient(r, defaultReplicaClientTimeout, logger)
    replicaClients = append(replicaClients, client)
  }
  return Leader{
    lastOrder: 0,
    logs: make([]string, 100),
    replicaClients: replicaClients,
    logger: logger,
  }
}

func (l *Leader) GetLogs() []string {
  return l.logs
}

func (l *Leader) AddMessage(m string, w int) error {
  var order int

  l.orderMX.Lock()
  order = l.lastOrder
  l.lastOrder++
  l.orderMX.Unlock()

  l.replicate(order, m, w)
  l.processReplicatedMessage(order, m)

  return nil
}

func (l *Leader) processReplicatedMessage(order int, m string) {
  l.logsMX.Lock()
  defer l.logsMX.Unlock()

  pendingMessage := pendingMessage{
    order,
    m,
  }

  // append pending message using shift sort
  for i, p := range l.pending {
    if order < p.order {
      before := l.pending[0:i]
      after := l.pending[i:len(l.pending)]
      l.pending = append(append(before, pendingMessage), after...)
      break
    }
  }

  // add pending messages to the log if they come in expected order
  lastExpectedIndex := len(l.logs)
  for _, p := range l.pending {
    if p.order == lastExpectedIndex { //expectly the next message in order
      l.logs = append(l.logs, p.message)
    } else if p.order <  lastExpectedIndex{ // message duplicate but might have new value
      l.logs[p.order] = p.message 
    } else {
      break
    }
  }

}

func (l *Leader) replicate(order int, m string, w int) error {
  l.logger.Info("Start replication from leader to all the replicas", zap.String("message", m))
  c := make(chan bool, len(l.replicaClients))

  for _, r := range l.replicaClients {
    go func (r ReplicaClient) {
      _ = r.replicate(order, m)
      c <- true
    }(r)
  }

  replicaACKs := 0
  for replicaACKs < (w - 1) {
    <- c
    replicaACKs++
    l.logger.Debug("Replica ACKed", zap.Int("ACK count", replicaACKs), zap.Int("W", w))
  }

  l.logger.Info("End replication from leader to all replicas", zap.String("message", m))
  return nil
}

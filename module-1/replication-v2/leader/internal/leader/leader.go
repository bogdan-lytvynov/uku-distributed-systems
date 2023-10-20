package leader

import (
  "sync"
  "net/rpc"
  "go.uber.org/zap"
  "github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v1/proto"
)

type Leader struct {
  order int
  mx sync.Mutext
  logs []string
  logger *zap.Logger
  replicas []string
  minAcks int
}

func NewLeader(replicas []string, minAcks int, logger *zap.Logger) Leader {
  return Leader{
    logger: logger,
    logs: []string{},
    replicas: replicas,
    minAcks: minAcks,
  }
}

func (l *Leader) GetLogs() []string {
  return l.logs
}

func (l *Leader) AddMessage(m string) error {
  mx.Lock()
  l.order++
  mx.Unlock()
  l.replicateToAll(order, m)
  l.logs[order] = m 
  return nil
}

func (l *Leader) replicateToOneReplica(args ReplicateArgs, replicaAdress string, c chan bool) {
  client ,err := rpc.DialHTTP("tcp", replicaAdress) 
  if err != nil {
    logger.Fatal("Failed to rpc.DialHTTP", zap.Error(err))
    return
  }
  reply := &proto.ReplicateReply{}
  replicateCall := client.Go("Replica.Replicate", args, reply, nil)
  <- replicateCall.Done

  c <- reply.Ack

  return nil
}

func (l *Leader) replicateToAll(order int, m string) error {
  l.logger.Info("Start replication from leader to all the replicas", zap.String("message", m))
  c := make(chan bool, len(l.replicas))

  for _, r := range l.replicas {
    args := &proto.ReplicateArgs{
      Order: order,
      Message: m,
    }
    go l.replicateToOneReplica(args, r, c)
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

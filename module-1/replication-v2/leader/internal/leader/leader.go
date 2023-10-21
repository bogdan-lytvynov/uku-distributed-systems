package leader

import (
  "time"
  "sync"
  "net/rpc"
  "go.uber.org/zap"
  "github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v2/log"
  "github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v2/proto"
)
var defaultReplicaClientTimeout time.Duration = 5

type Leader struct {
  log log.Log
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
    log: log.NewLog(),
    replicaClients: replicaClients,
    logger: logger,
  }
}

func (l *Leader) GetLog() []string {
  return l.log.GetLog()
}

func (l *Leader) AddMessage(m string, w int) error {
  index := l.log.nextIndex()

  l.replicate(index, m, w)
  l.processReplicatedMessage(index, m)

  return nil
}

func (l *Leader) replicate(index int, m string, w int) error {
  l.logger.Info("Start replication from leader to all the replicas",
    zap.String("message", m),
  )
  c := make(chan bool, len(l.replicaClients))

  for _, r := range l.replicaClients {
    go func (r ReplicaClient) {
      _ = r.replicate(index, m)
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

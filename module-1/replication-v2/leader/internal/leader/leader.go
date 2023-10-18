package leader

import (
  "net/rpc"
  "sync"
  "go.uber.org/zap"
  "github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v1/proto"
)

type Leader struct {
  logs []string
  logger *zap.Logger
  replicas []string
}

func NewLeader(replicas []string, logger *zap.Logger) Leader {
  return Leader{
    logger: logger,
    logs: []string{},
    replicas: replicas,
  }
}

func (l *Leader) GetLogs() []string {
  return l.logs
}

func (l *Leader) AddMessage(m string) error {
  l.replicateToAll(m)
  l.logs = append(l.logs, m)
  l.logger.Info("Add message to the log", zap.String("message", m))
  return nil
}

func (l *Leader) replicateToOneReplica(m string, replicaAdress string, wg *sync.WaitGroup) error {
  defer wg.Done()
  client ,err := rpc.DialHTTP("tcp", replicaAdress) 
  if err != nil {
    //fmt.Println("dialing:", err)
  }
  args := &proto.ReplicateArgs{
    Message: m,
  }
  reply := &proto.ReplicateReply{}
  replicateCall := client.Go("ReplicaRPC.Replicate", args, reply, nil)
  <- replicateCall.Done

  return nil
}

func (l *Leader) replicateToAll(m string) error {
  l.logger.Info("Start replication", zap.String("message", m))
  var wg sync.WaitGroup
  for _, r := range l.replicas {
    wg.Add(1)
    go l.replicateToOneReplica(m, r, &wg)
  }

  wg.Wait()
  l.logger.Info("Finished replication", zap.String("message", m))

  return nil
}

package leader
import (
  "net/rpc"
  "go.uber.org/zap"
)

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

func (rc *ReplicaClient) replicate(index int, message string) error {
  client ,err := rpc.DialHTTP("tcp", rc.address) 
  defer client.Close()

  if err != nil {
    rc.logger.Fatal("Failed to rpc.DialHTTP", zap.Error(err))
    return err
  }
  args := &proto.ReplicateMessage{
    Index: index,
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

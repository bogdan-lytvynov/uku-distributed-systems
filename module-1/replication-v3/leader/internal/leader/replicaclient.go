package leader
import (
  "math"
  "math/rand"
  "time"
  "errors"
  "net/rpc"
  "go.uber.org/zap"
  "github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v2/proto"
)
var replicaClientTimeout = errors.New("ReplicaClientTimeout")

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

func expBackoffWithJitter(n int) int {
  jitter := rand.Intn(5)
  // we now that n is int and 2^n will be int all the time
  expBackoff := int(math.Pow(2, float64(n)))

  return expBackoff + jitter
}

func (rc *ReplicaClient) replicate(index int, message string) error {
  i:=0
  // it is wrong to replicate in the loop but there is only way to achive convergance of the system to desired state
  for {
    err := rc.replicateAttempt(index, message)
    if err == nil {
      return nil
    } else {
      backoff := expBackoffWithJitter(i)
      rc.logger.Error("failed to replicate message to replica, backoff and try again")
      time.Sleep(time.Duration(backoff) * time.Second)
    }
  }
}

func (rc *ReplicaClient) replicateAttempt(index int, message string) error {
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
  case <-time.After(rc.timeout * time.Second):
    rc.logger.Info("Replica didn't respond")
    return replicaClientTimeout
  }
}

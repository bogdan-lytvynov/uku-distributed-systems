package proto

type ReplicateArgs struct {
  Order int
  Message string
}

type ReplicateReply struct {
  Ack bool
}

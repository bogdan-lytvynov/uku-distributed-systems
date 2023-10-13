package proto

type ReplicateArgs struct {
  Message string
}

type ReplicateReply struct {
  Ack bool
}

package proto

type ReplicateMessage struct {
  Order int
  Message string
}

type ReplicateMessageReply struct {
  Ack bool
}

package proto

type ReplicateMessage struct {
  Index int
  Message string
}

type ReplicateMessageReply struct {
  Ack bool
}

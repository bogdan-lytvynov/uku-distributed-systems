package main

import (
  "fmt"
  "net"
  "net/http"
  "net/rpc"
  "encoding/json"
  "github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v1/proto"

)
const HTTP_PORT = 3000
const RPC_PORT = 3001

var logs = []string{}

func getLogs(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":   
    w.Header().Set("Content-Type", "application/json")
    m, err := json.Marshal(logs)
    if err != nil {
      fmt.Println("Marshal error", err)
    }
    w.Write(m)

  default:
    fmt.Fprintf(w, "Sorry, only GET method is supported.")
  }
}

type ReplicaRPC struct{
}

func (r *ReplicaRPC) Replicate(args *proto.ReplicateArgs, reply *proto.ReplicateReply) error {
  logs = append(logs, args.Message)
  reply.Ack = true
  return nil
}

func startHttpServer() {
  http.HandleFunc("/logs", getLogs)

  fmt.Println("Start http server")
  err := http.ListenAndServe(fmt.Sprintf(":%d", HTTP_PORT), nil)

  if err != nil {
    fmt.Println("Error happened on server start", err)
  }
}

func startRpcServer() {
  //start RPC server
  r := new(ReplicaRPC)
  rpc.Register(r)
  rpc.HandleHTTP()
  l, e := net.Listen("tcp", fmt.Sprintf(":%d", RPC_PORT))
  if e != nil {
    fmt.Println("listen error:", e)
  }
  fmt.Println("Start rpc server")
  go http.Serve(l, nil)
}

func main() {
  startRpcServer()
  startHttpServer()
}

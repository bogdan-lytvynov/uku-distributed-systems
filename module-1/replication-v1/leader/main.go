package main

import (
  "os"
  "fmt"
  "net/rpc"
  "net/http"
  "encoding/json"
  "strings"
  "github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v1/proto"
  "sync"
)
const HTTP_PORT = 4000

var logs = []string{}

type MessageRequest struct {
  Message string
}

// return leader logs
func getLogs(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":   
    w.Header().Set("Content-Type", "application/json")
    m, err := json.Marshal(logs)
    if err != nil {
      fmt.Println("Marshal error: %w", err)
    }
    w.Write(m)

  default:
    fmt.Fprintf(w, "Sorry, only GET method is supported.")
  }
}

//hanlde new message
func postMessage(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "POST":   
    m := MessageRequest{}
    err := json.NewDecoder(r.Body).Decode(&m)
    if err != nil {
      fmt.Println("Can't decode", err)
      fmt.Fprintf(w, "Sorry, failed to decode request body")
      return
    }
    replicateOnAll(m.Message)
    logs = append(logs, m.Message)
    w.WriteHeader(http.StatusOK)

  default:
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprintf(w, "Sorry, only GET method is supported.")
  }
}

func replicateOnOneReplica(m string, replicaAdress string, wg *sync.WaitGroup) error {
  defer wg.Done()
  client ,err := rpc.DialHTTP("tcp", replicaAdress) 
  if err != nil {
    fmt.Println("dialing:", err)
  }
  args := &proto.ReplicateArgs{
    Message: m,
  }
  reply := &proto.ReplicateReply{}
  replicateCall := client.Go("ReplicaRPC.Replicate", args, reply, nil)
  <- replicateCall.Done

  return nil
}

func replicateOnAll(m string) error {
  fmt.Println("Replicate message:", m)
  replicas := strings.Split(os.Getenv("REPLICAS"), ",")

  var wg sync.WaitGroup
  for _, r := range replicas {
    wg.Add(1)
    go replicateOnOneReplica(m, r, &wg)
  }

  wg.Wait()
  fmt.Println("Finished replication of message", m)

  return nil
}

func main() {
  http.HandleFunc("/logs", getLogs)
  http.HandleFunc("/message", postMessage)

  err := http.ListenAndServe(fmt.Sprintf(":%d", HTTP_PORT), nil)

  if err != nil {
    fmt.Println("Error happened on server start %w", err)
  }
}

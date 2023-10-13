package main

import (
  "os"
  "fmt"
  "net/http"
  "net/rpc"
  "encoding/json"
  "github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v1/proto"
)
var logs = []string{}

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

type Replica struct{}

func (r Replica) Replicate(args *proto.ReplicateArgs, reply *proto.ReplicateReply) {
}

func main() {
  //start http server
  port := os.Getenv("PORT")
  http.HandleFunc("/logs", getLogs)

  host := fmt.Sprintf(":%s", port)
  fmt.Printf("Starting server: %s \n", host)
  err := http.ListenAndServe(host, nil)


  if err != nil {
    fmt.Println("Error happened on server start %w", err)
  }

  //start RPC

	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

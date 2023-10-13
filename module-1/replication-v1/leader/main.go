package main

import (
  "os"
  "fmt"
  "net/http"
  "encoding/json"
  "strings"
  //"github.com/bogdan-lytvynov/uku-distributed-systems/module-1/replication-v1/proto"
)

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
    replicateMessage(m.Message)

  default:
    fmt.Fprintf(w, "Sorry, only GET method is supported.")
  }
}

func replicateMessage(m string) error {
  replicas := strings.Split(os.Getenv("REPLICAS"), ",")

  fmt.Println(replicas)
  return nil
}

func main() {
  port := os.Getenv("PORT")
  http.HandleFunc("/logs", getLogs)
  http.HandleFunc("/message", postMessage)

  host := fmt.Sprintf(":%s", port)
  fmt.Printf("Starting server: %s \n", host)
  err := http.ListenAndServe(host, nil)

  if err != nil {
    fmt.Println("Error happened on server start %w", err)
  }

}

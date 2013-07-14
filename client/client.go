package main

import (
  "bufio"
  "fmt"
  "net"
  "log"
)

func main() {
  conn, err := net.Dial("tcp", ":2000")
  if err != nil {
    log.Fatal(err)
  }
  fmt.Fprintf(conn, "Hello, server\n")
  status, err := bufio.NewReader(conn).ReadString('\n')
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(status)
}

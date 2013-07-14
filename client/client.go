package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}

	counter := 0

	for {
		fmt.Printf("SEND Message %d\n", counter)
		fmt.Fprintf(conn, "Message %d\n", counter)
		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(status)

		counter++
		time.Sleep(2 * 1e9) // 2 seconds
	}
}

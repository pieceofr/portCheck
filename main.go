package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"
)

/* ./nodePortCheck -host=118.163.120.180 -port=2136 */
var peerPortReachable = false

const retryDelay = time.Duration(200 * time.Millisecond)
const retryTimes = 3
const checkInterMs = 1000
const dialTimeout = 2 * time.Second

func main() {
	var host = flag.String("host", "127.0.0.1", "enter the host ip for examine the node")
	var port = flag.String("port", "2136", "enter the host port for examine the node")
	flag.Parse()
	log.Println("Detecting host port = ", *host, ":", *port)

	go CheckPortReachableRoutine(*host, *port)

	for {
		fmt.Println("Reach Peer : ", peerPortReachable)
		time.Sleep(5 * time.Second)
	}

}

//CheckPortReachableRoutine is a Connection Check Routine
func CheckPortReachableRoutine(host, port string) {
	status := make(chan bool)
	for {
		go func(updateStatus chan<- bool) {
			connected := true
			for retry := 0; retry < retryTimes; retry++ {
				connected = connToPort(host, port)
				if !connected {
					fmt.Printf("NOT able to connect %s:%s\n", host, port)
					time.Sleep(retryDelay)
				} else {
					fmt.Printf("Connect to %s:%s\n", host, port)
					retry = retryTimes + 1
				}
			}
			updateStatus <- connected

		}(status)

		peerPortReachable = <-status
		time.Sleep(time.Duration(checkInterMs) * time.Millisecond)
	}
}

func connToPort(host, port string) bool {
	if host == "" {
		return false
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), dialTimeout)
	if err != nil {
		return false
	} else {
		conn.Close()
		return true
	}
}

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

func main() {
	var host = flag.String("host", "127.0.0.1", "enter the host ip for examine the node")
	var port = flag.String("port", "2136", "enter the host port for examine the node")
	flag.Parse()
	log.Println("Detecting host port = ", *host, ":", *port)

	go CheckPortReachableRoutine(*host, *port)

	for {
		time.Sleep(30 * time.Second)
	}

}

const retryDelay = time.Duration(500 * time.Millisecond)

//CheckPortReachableRoutine is a Connection Check Routine
func CheckPortReachableRoutine(host, port string) {
	stop := make(chan bool)
	defer close(stop)
	for {
		connStat := connCheck(host, port, 1000, 3, stop)
		for {
			peerPortReachable = <-connStat
		}

	}
}

func connCheck(host, port string, checkInterMs, retryTimes int, done <-chan bool) <-chan bool {
	status := make(chan bool)

	go func() {
		for {
			connected := true
			for retry := 0; retry < retryTimes; retry++ {
				connected = connToPort(host, port)
				if !connected {
					log.Println("NOT connect to ", host, ":", port)
					retry++
					time.Sleep(retryDelay)
				} else {
					log.Println("is connected to ", host, ":", port)
					retry = retryTimes + 1
				}
			}
			select {
			case <-done:
				return
			case status <- connected:
			}
			time.Sleep(time.Duration(checkInterMs) * time.Millisecond)
		}
	}()
	return status
}

func connToPort(host, port string) bool {
	if host == "" {
		return false
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return false
	} else {
		defer conn.Close()
		return true
	}
}

// TCP server
// This server listens for text based inputs
// This inputs are commands for the crawler for process
// Useful for debugging with telnet
package main

import "net"
import "log"
import "strings"
import "bufio"

// Bind socket to...
const TEXT_PROTOCOL_BIND = ":4040"

func TextProtocolHandler(pageCrawl chan string) {
	server, err := net.Listen("tcp", TEXT_PROTOCOL_BIND)

	if err != nil {
		log.Println("Tcp server failed to listen")
		log.Fatal(err)
	}

	defer server.Close()

	log.Println("TCP server waiting for commands on", TEXT_PROTOCOL_BIND)

	for {
		conn, err := server.Accept()

		if err != nil {
			return
		}

		log.Println("Client connected")

		go TextProtocolClientHandler(conn, pageCrawl)
	}
}

func TextProtocolClientHandler(conn net.Conn, pageCrawl chan string) {
	defer conn.Close()

	for {
		cmd_str, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			return
		}

		// Trim new line
		cmd_str = strings.TrimSuffix(cmd_str, "\r\n")

		// Split cmd method from argument
		cmd := strings.SplitN(cmd_str, " ", 2)

		// Upper first part of cmd
		cmd[0] = strings.ToUpper(cmd[0])

		// Process cmd
		switch cmd[0] {
		case "CRAWL":
			// Min arguments
			if len(cmd) > 2 {
				break
			}

			// Crawl an url
			url := cmd[1]
			log.Println("Received request from", conn.RemoteAddr(), "to crawl", url)

			// Validate url
			// TODO
			// Dispatch
			// TODO
			pageCrawl <- url
			conn.Write([]byte("Queued for processing " + url + "...\n"))
		case "EXIT", "QUIT":
			return
		default:
			conn.Write([]byte("Invalid command.\n"))
		}
	}
}

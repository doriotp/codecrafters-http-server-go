package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func handler(w http.ResponseWriter, r *http.Request) {
	// Extract the URL path
	path := r.URL.Path
	fmt.Println(path)
	if path == "/" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	var (
		userAgentValue string
	)
	defer conn.Close()
	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read request:", err)
		return
	}

	parts := strings.Fields(requestLine)
	if parts[1] == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(parts[1], "/echo/") {
		text := parts[1][6:]
		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(text), text)
		conn.Write([]byte(response))
	} else if strings.HasPrefix(parts[1], "/user-agent") {
		for {
			requestLine, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("Failed to read request:", err)
				return
			}

			parts = strings.Fields(requestLine)

			if parts[0] == "User-Agent:" {
				userAgentValue = parts[1]
				break
			}
		}

		response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgentValue), userAgentValue)
		conn.Write([]byte(response))

	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}

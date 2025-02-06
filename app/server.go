package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func readFile(directory, filename string) (string, error) {
	data, err := os.ReadFile(directory + filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist")
		} else {
			fmt.Println("Error reading file:", err)
		}
		return "", err
	}

	return string(data), nil
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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	var (
		userAgentValue string
		contentLength  int
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

	} else if strings.HasPrefix(parts[1], "/files/") {
		dir := os.Args[2]
		fileName := parts[1][7:]
		fmt.Println(parts[0])
		if parts[0] == "GET" {
			fmt.Println(dir)
			data, err := readFile(dir, fileName)
			if err != nil {
				conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			} else {
				response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(data), data)
				conn.Write([]byte(response))
			}
		} else {
			fmt.Println(1)
			for {
				requestLine, err = reader.ReadString('\n')
				if err != nil {
					fmt.Println("Failed to read request:", err)
					return
				}

				parts = strings.Fields(requestLine)

				if parts[0] == "Content-Length:" {
					contentLength, _ = strconv.Atoi(parts[1])
					break
				}
			}

			// Read the body based on Content-Length
			if contentLength > 0 {
				// body := make([]byte, contentLength)
				body, err := io.ReadAll(reader)
				if err != nil {
					fmt.Println("Failed to read body:", err)
					return
				}

				fmt.Println(string(body))

				os.WriteFile(dir+fileName, body, 0755)
				conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))

			}
		}
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}

}

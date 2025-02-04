package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func handler(w http.ResponseWriter, r *http.Request) {
	// Extract the URL path
	path := r.URL.Path
	fmt.Println(path)
	if path=="/"{
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else{
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	// l, err := net.Listen("tcp", "0.0.0.0:4221")
	// if err != nil {
	// 	fmt.Println("Failed to bind to port 4221")
	// 	os.Exit(1)
	// }
	
	// conn, err := l.Accept()
	// if err != nil {
	// 	fmt.Println("Error accepting connection: ", err.Error())
	// 	os.Exit(1)
	// }



	// http.HandleFunc("/abcdefg", func(rw http.ResponseWriter, r *http.Request){
	// 	rw.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	// })

	// http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request){
	// 	rw.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	// })
	http.HandleFunc("/", handler)
	http.ListenAndServe(":4221", nil)
}

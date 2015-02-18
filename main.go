package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		log.Println("Connection received from", req.RemoteAddr)
		socketClient, _, err := res.(http.Hijacker).Hijack()
		if err != nil {
			res.WriteHeader(500)
			log.Println("Fail to hijack client", err)
			return
		}
		// Flush the options to make sure the client sets the raw mode
		socketClient.Write([]byte{})
		fmt.Fprintf(socketClient, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\n")

		log.Println("Headers sent, starting echoing client", req.RemoteAddr)

		io.Copy(socketClient, socketClient)
		socketClient.Close()

		log.Println("End of hijacking for", req.RemoteAddr)
	})

	log.Println(http.ListenAndServe(":8080", nil))
}

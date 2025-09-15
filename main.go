package main

import (
	"fmt"
	"net/http"

	"github.com/ridhwan-faris/sftp-example/handler"
)

func main() {
	sftpHandler := handler.SftpHandler{}

	http.HandleFunc("/send-file", sftpHandler.SendFile)
	fmt.Println("Starting...")

	http.ListenAndServe(":8282", nil)
}

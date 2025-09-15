package handler

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/pkg/sftp"
	"github.com/ridhwan-faris/sftp-example/common/helper"
	"github.com/ridhwan-faris/sftp-example/handler/request"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/ssh"
)

const (
	sheet       = "Sheet 1"
	sftpAddress = "localhost:2222"
)

type SftpHandler struct{}

func (c *SftpHandler) SendFile(w http.ResponseWriter, r *http.Request) {
	req := request.SftpRequest{}

	err := helper.Bind(r, &req)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("request body", req.FileName)

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	index, err := f.NewSheet(sheet)
	if err != nil {
		fmt.Println("error when new sheet", err)
		fmt.Println(err)
		http.Error(w, "error generate sheet", http.StatusInternalServerError)
		return
	}

	f.SetCellValue(sheet, "A1", "POC SFTP")
	f.SetActiveSheet(index)

	var buf bytes.Buffer
	err = f.Write(&buf)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "error at write file", http.StatusInternalServerError)
		return
	}

	sshConfig := ssh.ClientConfig{
		User:            "ridhwan",
		Auth:            []ssh.AuthMethod{ssh.Password("password123")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshConn, err := ssh.Dial("tcp", sftpAddress, &sshConfig)
	if err != nil {
		fmt.Println("error at dial sftp:", err)
		http.Error(w, "error at dial sftp", http.StatusInternalServerError)
		return
	}

	defer sshConn.Close()

	client, err := sftp.NewClient(sshConn)
	if err != nil {
		fmt.Println("error at create new client sftp:", err)
		http.Error(w, "error at create new client sftp", http.StatusInternalServerError)
		return
	}

	defer client.Close()

	fileRemote, err := client.Create("/upload/" + req.FileName + ".xlsx")
	if err != nil {
		fmt.Println("error at create file:", err)
		http.Error(w, "error at create file", http.StatusInternalServerError)
		return
	}

	defer fileRemote.Close()

	_, err = fileRemote.Write(buf.Bytes())
	if err != nil {
		fmt.Println("error at write buffer to remote sftp:", err)
		http.Error(w, "error at write buffer to remote sftp", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "File successfully upload to SFTP server")
}

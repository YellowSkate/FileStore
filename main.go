package main

import (
	"GoFileStore/handler"
	"fmt"
	"net/http"
)

func main() {
	// 设置根目录
	fs := http.FileServer(http.Dir("~/GoFileStore"))
	http.Handle("/", fs)

	//绑定方法
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSuc)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)

	http.HandleFunc("/file/update", handler.FileUpdataMetaHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to Start server,err:%s", err.Error())
	}

}

package main

import (
	"GoFileStore/handler"
	"fmt"
	"net/http"
)

func main() {
	// 设置根目录
	fs := http.FileServer(http.Dir("./www"))
	http.Handle("/www/", http.StripPrefix("/www/", fs))
	fmt.Println("231233")

	//用户相关接口
	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	//文件相关接口
	http.HandleFunc("/file/upload", handler.HTTPInterceptor(handler.UploadHandler))
	http.HandleFunc("/file/upload/suc", handler.HTTPInterceptor(handler.UploadSuc))     //todo
	http.HandleFunc("/file/meta", handler.HTTPInterceptor(handler.GetFileMetaHandler))  //摒弃
	http.HandleFunc("/file/download", handler.HTTPInterceptor(handler.DownloadHandler)) //todo

	http.HandleFunc("/file/query", handler.HTTPInterceptor(handler.FileQueryHandler))
	http.HandleFunc("/file/update", handler.HTTPInterceptor(handler.FileUpdataMetaHandler)) //todo
	http.HandleFunc("/file/delete", handler.HTTPInterceptor(handler.FileDeleteHandler))     //todo

	http.HandleFunc("/file/fastupload", handler.HTTPInterceptor(handler.TryFastUploadHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to Start server,err:%s", err.Error())
	}

}

package handler

import (
	dblayer "GoFileStore/db"
	"GoFileStore/meta"
	"GoFileStore/util"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	// "io/ioutil" //在go1.16之后 可以直接使用 io
	"net/http"
	"os"
	"time"
)

// 打印请求行和头部信息
// func display(r *http.Request) {
// 	fmt.Println("Method:", r.Method)
// 	fmt.Println("URL:", r.URL.String())
// 	fmt.Println("Headers:")
// 	for name, values := range r.Header {
// 		for _, value := range values {
// 			fmt.Printf("%s: %s\n", name, value)
// 		}
// 	}
// }

// 文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {

	// { //Debug
	// 	fmt.Println("call UploadHandler")
	// 	display(r)
	// }
	if r.Method == "GET" {
		//返回上传html页面
		data, err := os.ReadFile("./www/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		// fmt.Println(data)
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		//接收文件流及本地目录

		//获取表单提交的文件
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get FormFile :%s\n", err.Error())
			return
		}
		defer file.Close() //!!!

		//提取文件元信息
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./temp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-01 15:04:05"),
		}

		//根据文件元信息 创建 newFile
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to Create newFile :%s\n ", err.Error())
			return
		}
		defer newFile.Close()

		// newfile=file
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to Copy file,err:%s\n", err.Error())
			return

		}
		//计算文件元 key
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		// meta.UpdateFileMeta(fileMeta)

		meta.UploadFileMetaDB(fileMeta) //插入数据库
		//Todo 判断插入成功
		// fmt.Printf("Sha1:%s", fileMeta.FileSha1)//debug

		//更新用户文件表信息
		r.ParseForm()
		username := r.Form.Get("username")
		suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if suc {
			http.Redirect(w, r, "/www/view/home.html", http.StatusFound)
		} else {
			w.Write([]byte("Upload uesr_file Failed."))
		}

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound) //重定向
	}
}

// 文件上传成功重定向
func UploadSuc(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload Sucess!")
}

// 获取文件辕信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm() //解析表达数据 ，并存储在r.Form
	filehash := r.Form["filehash"][0]
	// fMeta, ok := meta.GetFileMeta(filehash)
	fMeta, ok := meta.GetFileMetaDB(filehash)

	if nil != ok {
		w.WriteHeader(http.StatusInternalServerError) //http 状态码 500
		io.WriteString(w, "Failed GetFileMeta!")

		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //http 状态码 500
		io.WriteString(w, "Failed json Marshal!")

		return
	}
	w.Write(data)
}

//文件下载

func DownloadHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	// fm, ok := meta.GetFileMeta(fsha1)
	fm, ok := meta.GetFileMetaDB(fsha1)

	if nil != ok {
		w.WriteHeader(http.StatusInternalServerError) //http 状态码 500
		io.WriteString(w, "Failed GetFileMeta!")

		return
	}
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //http 状态码 500
		io.WriteString(w, "Failed Open location!")

		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //http 状态码 500
		io.WriteString(w, "Failed Open location!")

		return
	}

	// meta.Desc(*fm)
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-disposition", `attachment;filename="`+fm.FileName+`"`)
	w.Write(data)
}

//修改文件元信息  //tod

func FileUpdataMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden) //403
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed) //405
		return
	}

	// curFileMeta, ok := meta.GetFileMeta(fileSha1)
	curFileMeta, ok := meta.GetFileMetaDB(fileSha1)

	if ok != nil {
		w.WriteHeader(http.StatusInternalServerError) //http 状态码 500
		io.WriteString(w, "Failed GetFileMeta!")

		return
	}

	//修改Meta信息
	curFileMeta.FileName = newFileName

	meta.UpdateFileMeta(*curFileMeta)
	meta.UpdateFileMetaDB(*curFileMeta) //TODO 修改文件名字

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //http 状态码 500
		io.WriteString(w, "Failed json Marshal!")

		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// 文件删除
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fMeta, ok := meta.GetFileMeta(fileSha1)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError) //http 状态码 500
		io.WriteString(w, "Failed GetFileMeta!")

		return
	}

	os.Remove(fMeta.Location)

	meta.RemoveFileMeta(fileSha1)

	w.WriteHeader(http.StatusOK)

}

// 获取用户文件信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")

	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 秒传
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	//解析请求
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	//从文件表中查询相同hash的文件记录

	fileMeta, err := meta.GetFileMetaDB(filehash)
	if nil != err {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//查不到记录返回秒传失败
	if nil == fileMeta {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败,请使用普通接口上传",
		}
		w.Write(resp.JSONBytes())
		return
	}

	//将文件信息填入用户表,进行秒传
	suc := dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		w.Write(resp.JSONBytes())
		return
	} else {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败,请重试",
		}
		w.Write(resp.JSONBytes())
		return

	}

}

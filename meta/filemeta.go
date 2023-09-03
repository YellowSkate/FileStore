package meta

import (
	mydb "GoFileStore/db"
	"fmt"
	"sync"
	"time"
)

// 描述文件元结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string //时间戳
}

func Desc(f FileMeta) {
	fmt.Println(f.FileSha1)
	fmt.Println(f.FileName)
	fmt.Println(f.FileSize)
	fmt.Println(f.Location)
	fmt.Println(f.Location)

}

var fileMetas map[string]FileMeta

var rwMutex sync.RWMutex = sync.RWMutex{}

func init() {

	fileMetas = make(map[string]FileMeta)
}

// 获取文件元信息
func GetFileMeta(k string) (f FileMeta, ok bool) {
	rwMutex.RLock()
	f, ok = fileMetas[k]
	defer rwMutex.RUnlock()
	//k 指的是文件 映射后 sha1的值
	return
}

// 更新
func UpdateFileMeta(f FileMeta) {
	fileMetas[f.FileSha1] = f
}

// 删除
func RemoveFileMeta(k string) {
	delete(fileMetas, k)
}

// 排序
const baseFormat = "2006-01-02 15:04:05"

type ByUploadTime []FileMeta

func (a ByUploadTime) Len() int {
	return len(a)
}

func (a ByUploadTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByUploadTime) Less(i, j int) bool { //按照时间排序
	iTime, _ := time.Parse(baseFormat, a[i].UploadAt)
	jTime, _ := time.Parse(baseFormat, a[j].UploadAt)
	return iTime.UnixNano() > jTime.UnixNano()
}

// DB  2.0
// 上传 数据库
func UploadFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// 更新文件
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)

}

// 获取文件
func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	tfile, err := mydb.GetFileMeta(fileSha1)
	if tfile == nil || err != nil {
		return nil, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return &fmeta, nil
}

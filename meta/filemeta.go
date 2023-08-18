package meta

import (
	"fmt"
	"sync"
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

package files

import (
	"dcontrol/server/base"
	"dcontrol/server/setting"
	"dcontrol/server/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var uploadDir = "./files" // 指定文件目录
var maxFileSize int64 = 1200 * 1024 * 1024 // 上限 1200M

type FileInfo struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Sizes string `json:"sizes"`
	Time string `json:"time"`
	Timestamp int64  `json:"timestamp"` // 使用毫秒时间戳
}

func InitDir() {
	if setting.Conf.Dir != "" {
		uploadDir = setting.Conf.Dir
	}
	fmt.Println("uploadDir:", setting.Conf.Dir, uploadDir)
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}
}

func HandleApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("Content-Type", "application/json")

	// 获取请求路径   strings.HasSuffix
	path := r.URL.Path
	fmt.Println("redis HandleApi path:", path)

	switch {
	case strings.Contains(path, "/upload"):
		uploadHandler(w, r)
	case strings.Contains(path, "/download"):
		downloadHandler(w, r)
	case strings.Contains(path, "/list"):
		listHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

// 获取唯一的文件名，如果文件已存在，递增后缀 (n)
func getUniqueFileName(filePath string) string {
	ext := filepath.Ext(filePath)
	base := strings.TrimSuffix(filepath.Base(filePath), ext)
	newFilePath := filePath

	// 递增 (n) 后缀，直到文件名唯一
	n := 1
	for {
		if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
			// 文件不存在，返回唯一的文件路径
			return newFilePath
		}
		// 文件已存在，修改文件名加上 (n)
		newFilePath = fmt.Sprintf("%s(%d)%s", filepath.Join(filepath.Dir(filePath), base), n, ext)
		n++
	}
}

// 处理文件上传
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// 限制请求方法为POST
	if r.Method != http.MethodPost {
		base.R(w).FailMsg("请求方法必须为POST")
		return
	}

	// 解析multipart表单数据，限制最大文件大小为120MB
	err := r.ParseMultipartForm(maxFileSize)
	if err != nil {
		base.R(w).FailMsg("上传文件大小不能超过1200M")
		return
	}

	// 获取文件
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		base.R(w).FailMsg("上传文件file为空")
		return
	}
	defer file.Close()

	// 创建目录如果不存在
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err := os.MkdirAll(uploadDir, 0755)
		if err != nil {
			base.R(w).FailMsg("创建目录失败")
			return
		}
	}

	// 获取文件原始文件名
	fileName := fileHeader.Filename
	// 获取文件路径（包括目录和文件名）
	filePath := filepath.Join(uploadDir, fileName)

	// 获取唯一的文件名
	uniqueFilePath := getUniqueFileName(filePath)

	// 保存文件到指定目录
	dst, err := os.Create(uniqueFilePath)
	if err != nil {
		base.R(w).FailMsg("创建文件失败")
		return
	}
	defer dst.Close()

	// 将文件内容复制到新文件中
	_, err = io.Copy(dst, file)
	if err != nil {
		base.R(w).FailMsg("保存文件失败")
		return
	}

	base.R(w).Ok(filepath.Base(uniqueFilePath))

}

// 处理文件下载
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// 获取文件名参数
	filename := r.URL.Query().Get("name")
	if filename == "" {
		base.R(w).FailMsg("文件名不能为空")
		return
	}
	fmt.Println("下载文件, " + filename)

	// 构建完整的文件路径
	filePath := filepath.Join(uploadDir, filename)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		base.R(w).FailMsg("文件不存在")
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		base.R(w).FailMsg("读取文件错误")
		return
	}

	// 设置响应头 // attachment  inline
	w.Header().Set("Content-Disposition", "inline; filename="+fileInfo.Name())
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 将文件内容写入响应
	_, err = io.Copy(w, file)
	if err != nil {
		base.R(w).FailMsg("读取文件错误")
	}
}

// 处理列出文件目录
func listHandler(w http.ResponseWriter, r *http.Request) {

	var filesInfo []FileInfo

	// 打开目录
	dir, err := os.Open(uploadDir)
	if err != nil {
		base.R(w).FailMsg("读取目录错误")
		return
	}
	defer dir.Close()

	// 读取目录中的文件列表
	files, err := dir.Readdir(-1) // -1 表示读取目录下所有文件
	if err != nil {
		base.R(w).FailMsg("读取目录文件错误")
		return
	}

	for _, file := range files {
		if !file.IsDir() { // 只处理文件
			// 格式化时间为 "YYYY-MM-DD HH:MM:SS"
			filesInfo = append(filesInfo, FileInfo{
				Name:        file.Name(),
				Size:        file.Size(),
				Sizes: utils.FormatBytes(uint64(file.Size())),
				Time: file.ModTime().Format("2006-01-02 15:04:05"),
				Timestamp: file.ModTime().UnixNano() / int64(time.Millisecond), // 转换为毫秒时间戳
			})
		}
	}

	base.R(w).Ok(filesInfo)
}
package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const FilePath = "/tmp"

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func main() {
	if !IsExist(FilePath) {
		err := os.Mkdir(FilePath, 0750)
		if err != nil {
			return
		}
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/download", downloadHandler)

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// 读取生成的文件列表
	files, err := getListOfFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 将文件列表渲染到HTML页面
	html := "<h1>Generated Files</h1>"
	html += "<ul>"
	html += "<form action=\"/generate\" method=\"get\" id=\"form1\">"
	html += "<label for=\"fileSize\">fileSize(MB):</label>"
	html += "<input type=\"text\" id=\"fileSize\" name=\"fileSize\"><br><br>"
	html += "</form>"
	html += "<button type=\"submit\" form=\"form1\" value=\"generated\">generated</button>"
	html += "</ul>"

	html += "<ul>"
	for _, file := range files {
		html += fmt.Sprintf("<li><a href=\"/download?file=%s\">%s</a></li>", file, file)
	}
	html += "</ul>"

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	// 解析页面上输入的文件大小
	fileSizeStr := r.FormValue("fileSize")
	fileSize, err := strconv.ParseInt(fileSizeStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid file size", http.StatusBadRequest)
		return
	}

	// 生成一个指定大小的临时文件
	fileName := generateFileName()
	err = generateFile(fileName, fileSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// 从查询参数中获取要下载的文件名
	fileName := r.FormValue("file")

	// 打开要下载的文件
	file, err := os.Open(FilePath + "/" + fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 设置响应头，告诉浏览器以附件形式下载文件
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Transfer-Encoding", "binary")

	// 将文件内容写入HTTP响应
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 生成一个随机的文件名
func generateFileName() string {
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(100000)
	return fmt.Sprintf("file%d.txt", randomNum)
}

// 在指定路径生成指定大小的文件
func generateFile(fileName string, fileSize int64) error {
	file, err := os.Create(FilePath + "/" + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将指定大小的随机数据写入文件
	_, err = io.CopyN(file, rand.New(rand.NewSource(time.Now().UnixNano())), fileSize*1024*1024)
	if err != nil {
		return err
	}

	return nil
}

// 获取已生成的文件列表
func getListOfFiles() ([]string, error) {
	var files []string

	dir, err := os.Open(FilePath)
	if err != nil {
		return files, err
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return files, err
	}

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			files = append(files, fileInfo.Name())
		}
	}

	return files, nil
}

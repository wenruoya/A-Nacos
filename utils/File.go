package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ReadFile(path string) []string {
	f, err := os.Open(path) //打开目录路径txt文件
	if err != nil {
		return nil
	}
	defer f.Close() //最后关闭文件

	r := bufio.NewReader(f)

	buf := make([]byte, 1024)

	chunks := []byte{}
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if 0 == n {
			break
		}
		chunks = append(chunks, buf[:n]...)
	}

	res := strings.ReplaceAll(string(chunks), "\r\n", ",")
	return strings.Split(res, ",")
}

func WriteFile(res string) string {
	file, err := os.OpenFile("res.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("File open fail! err: %v", err)
	}
	defer file.Close()

	// 创建一个写的对象
	wr := bufio.NewWriter(file)
	_, _ = wr.WriteString(res) // 将内容写入缓存
	_ = wr.Flush()             // 将缓存中的内容写入文件中

	path, err := filepath.Abs("res.txt")
	return path
}

package util

import (
	"io/ioutil"
	"log"
	"os"
)

func TruncateWrite(fileName string, data []byte) {
	OSFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer OSFile.Close()
	OSFile.Truncate(0)
	n, _ := OSFile.Seek(0, 0)
	_, err = OSFile.WriteAt(data, n)
	if err != nil {
		log.Fatal(err)
	}
}

func ReadFile(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

func IsExtraFile(filePath string) bool {
	// 文件不存在则返回error
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

package music

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
)

//go:embed mf/*.mp3
var filesFs embed.FS

var MusicsFs []*TMusicFs

type TMusicFs struct {
	Fs  *os.File
	Tag string
}

//ReloadSoundFiles  加载所有歌曲文件
func ReloadSoundFiles() {
	//打印出文件名称
	fslist, err := fs.ReadDir(filesFs, "mf")
	if err != nil {
		log.Fatal(err)
	}
	//打印出文件名称
	path := "/Users/eddiechan/go/bin"
	for _, file := range fslist {
		f, err := os.Open(path + "/music/mf/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		_, after, _ := strings.Cut(f.Name(), "-")
		after = strings.TrimRight(after, ".mp3")
		MusicsFs = append(MusicsFs, &TMusicFs{
			Fs:  f,
			Tag: after,
		})
	}
	fmt.Println("音频文件加载完成")
}

//CloseMusicFs 关闭所有音乐文件
func CloseMusicFs() {
	for _, f := range MusicsFs {
		f.Fs.Close()
	}
	fmt.Println("音频文件释放完成")

}

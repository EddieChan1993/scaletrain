package music

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"scaletrain/util"
	"strconv"
	"strings"
)

//go:embed mf/*.mp3
var filesFs embed.FS

var MusicsFs = make(map[int]*TMusicFs)

type SoundTyp = string

const (
	LowSound  SoundTyp = "低音"
	MidSound  SoundTyp = "中音"
	HighSound SoundTyp = "高音"
)

type TMusicFs struct {
	Fs        *os.File
	Tag       string
	Id        int
	SoundType SoundTyp
}

//ReloadSoundFiles  加载所有歌曲文件
func ReloadSoundFiles() {
	//打印出文件名称
	fslist, err := fs.ReadDir(filesFs, "mf")
	if err != nil {
		log.Fatal(err)
	}
	//打印出文件名称
	for _, file := range fslist {
		f, err := os.Open(util.Path + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		before, after, _ := strings.Cut(strings.TrimLeft(file.Name(), util.Path), "-")
		after = strings.TrimRight(after, ".mp3")
		soundType, _, _ := strings.Cut(after, " ")
		id, _ := strconv.Atoi(before)
		MusicsFs[id] = &TMusicFs{
			Fs:        f,
			Tag:       after,
			Id:        id,
			SoundType: soundType,
		}
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

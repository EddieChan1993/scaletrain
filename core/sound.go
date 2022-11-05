package core

import (
	_ "embed"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"log"
	"os"
	"scaletrain/music"
	"scaletrain/util"
)

type sound struct {
	fs     *os.File
	s      beep.StreamSeekCloser
	format beep.Format
	tag    string
}

//newRandSound 随机获取一歌曲文件
func newRandSound() *sound {
	musicFsIndex := util.RandInt(len(music.MusicsFs))
	musicsFsOne := music.MusicsFs[musicFsIndex-1]
	streamer, format, err := mp3.Decode(musicsFsOne.Fs)
	if err != nil {
		log.Fatal(err)
	}
	return &sound{
		fs:     musicsFsOne.Fs,
		s:      streamer,
		format: format,
		tag:    musicsFsOne.Tag,
	}
}

func (this_ *sound) resetSound() {
	this_.fs.Seek(0, 0)
}

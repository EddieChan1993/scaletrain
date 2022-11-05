package core

import (
	_ "embed"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"log"
	"os"
	"scaletrain/music"
)

type sound struct {
	fs     *os.File
	s      beep.StreamSeekCloser
	format beep.Format
	tag    string
	id     int
}

//newRandSound 随机获取一歌曲文件
func newRandSound(st music.SoundTyp) *sound {
	musicFsIndex := randIndexByScore(st)
	musicsFsOne := music.MusicsFs[musicFsIndex]
	streamer, format, err := mp3.Decode(musicsFsOne.Fs)
	if err != nil {
		log.Fatal(err)
	}
	return &sound{
		id:     musicsFsOne.Id,
		fs:     musicsFsOne.Fs,
		s:      streamer,
		format: format,
		tag:    musicsFsOne.Tag,
	}
}

func newSound(id int) *sound {
	musicsFsOne := music.MusicsFs[id]
	streamer, format, err := mp3.Decode(musicsFsOne.Fs)
	if err != nil {
		log.Fatal(err)
	}
	return &sound{
		id:     musicsFsOne.Id,
		fs:     musicsFsOne.Fs,
		s:      streamer,
		format: format,
		tag:    musicsFsOne.Tag,
	}
}
func (this_ *sound) resetSound() {
	this_.fs.Seek(0, 0)
}

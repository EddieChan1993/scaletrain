package core

import (
	"context"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"scaletrain/music"
	"time"
)

const statusStop = 1

type Player struct {
	queue    *queue
	nowSound chan *sound
	sr       beep.SampleRate
	ctx      context.Context
	cancel   context.CancelFunc
	ticker   *time.Ticker
	stop     chan struct{}
	status   int
}

func InitPlayer() *Player {
	music.ReloadSoundFiles()
	InitScore()
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(4 * time.Second)
	return &Player{
		queue:    &queue{},
		nowSound: make(chan *sound),
		sr:       beep.SampleRate(44100),
		ctx:      ctx,
		cancel:   cancel,
		ticker:   ticker,
		stop:     make(chan struct{}),
	}
}

func (this_ *Player) RunPlayer() {
	speaker.Init(this_.sr, this_.sr.N(time.Second/10))
	speaker.Play(this_.queue)
	flag := ""
	fmt.Println("=======================")
	time.Sleep(2 * time.Second)
	this_.randPlay()
	go func() {
		defer func() {
			this_.stop <- struct{}{}
			fmt.Println("stop 播放器")
		}()
		for {
			select {
			case <-this_.ctx.Done():
				return
			case s := <-this_.nowSound:
				time.Sleep(3 * time.Second)
				fmt.Println(s.tag, "正确(v) ?")
				fmt.Scanf("%s", &flag)
				if flag == "v" {
					AddScore(s.id)
				}
				if this_.status == statusStop {
					return
				}
				fmt.Println("=======================")
				time.Sleep(2 * time.Second)
				s.resetSound()
				this_.randPlay()
			}
		}
	}()
}

//randPlay 随机播放
func (this_ *Player) randPlay() {
	s := newRandSound()
	resampled := beep.Resample(4, s.format.SampleRate, this_.sr, s.s)
	volume := &effects.Volume{
		Streamer: resampled,
		Base:     2,
		Volume:   2,
	}
	speaker.Lock()
	this_.queue.Add(volume, beep.Callback(func() {
		this_.nowSound <- s
	}))
	speaker.Unlock()
}
func (this_ *Player) Stop() {
	this_.cancel()
	this_.ticker.Stop()
	this_.status = statusStop
	<-this_.stop
	saveScore()
}

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

type Player struct {
	queue    *queue
	nowSound chan *sound
	sr       beep.SampleRate
	ctx      context.Context
	cancel   context.CancelFunc
	ticker   *time.Ticker
	stop     chan struct{}
}

func InitPlayer() *Player {
	music.ReloadSoundFiles()
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
	go func() {
		for {
			select {
			case <-this_.ctx.Done():
				fmt.Println("stop 播放器")
				this_.stop <- struct{}{}
				return
			case s := <-this_.nowSound:
				fmt.Println(s.tag)
				s.resetSound()
			case <-this_.ticker.C:
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
	<-this_.stop
}

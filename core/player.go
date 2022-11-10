package core

import (
	"context"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"scaletrain/music"
	"scaletrain/util"
	"time"
)

const statusStop = 1
const waitSec = 2        //等待时间
const maxRepeatTimes = 3 //最大重复次数

type Player struct {
	queue    *queue
	nowSound chan *sound
	barCh    chan struct{}
	sr       beep.SampleRate
	ctx      context.Context
	cancel   context.CancelFunc
	ticker   *time.Ticker
	stop     chan struct{}
	status   int
	st       music.SoundTyp //音高类型
}

func InitPlayer() *Player {
	music.ReloadSoundFiles()
	InitScore()
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(4 * time.Second)

	player := &Player{
		queue:    &queue{},
		nowSound: make(chan *sound),
		barCh:    make(chan struct{}),
		sr:       beep.SampleRate(44100),
		ctx:      ctx,
		cancel:   cancel,
		ticker:   ticker,
		stop:     make(chan struct{}),
	}
	player.selectScaleST()
	return player
}

//selectScaleST 音区选择
func (this_ *Player) selectScaleST() {
	stFlg := 0
	fmt.Println("===========选择练习音区============")
	fmt.Println("1高音区  2中音区 3低音区（其他）全部")
	fmt.Scanf("%d", &stFlg)
	switch stFlg {
	case 1:
		this_.st = music.HighSound
	case 2:
		this_.st = music.MidSound
	case 3:
		this_.st = music.LowSound
	}
}

func (this_ *Player) RunPlayer() {
	speaker.Init(this_.sr, this_.sr.N(time.Second/10))
	speaker.Play(this_.queue)
	fmt.Println("===========开始准备", this_.st, "============")
	util.WaitBar(fmt.Sprintf("%d 秒后开始", waitSec), waitSec)
	this_.randPlay()
	go func() {
		defer func() {
			this_.stop <- struct{}{}
			fmt.Println("stop 播放器")
		}()
		repeatTimes := 0
		flag := ""
		for {
			select {
			case <-this_.ctx.Done():
				return
			case s := <-this_.nowSound:
				s.resetSound()
				if repeatTimes == 0 {
					util.WaitBar(fmt.Sprintf("%d 秒后公布答案", waitSec), waitSec)
					fmt.Println(s.tag, "正确(v) ?")
					flag = ""
					fmt.Scanf("%s\n", &flag)
					if flag == "v" {
						//对，加分
						AddScore(s.id)
						this_.randPlay()
						fmt.Println("=======================")
					} else {
						//错了，则重复播放
						fmt.Println(s.tag, "反复听", maxRepeatTimes, "遍")
						this_.repeatPlay(s.id)
						repeatTimes++
						SubScore(s.id)
						break
					}
				} else {
					time.Sleep(1 * time.Second)
					if repeatTimes >= maxRepeatTimes {
						//重听次数超过，正常随机
						this_.randPlay()
						repeatTimes = 0
						fmt.Println("=======================")
						break
					}
					//错了，则重复播放
					repeatTimes++
					fmt.Println(s.tag, "第", repeatTimes, "遍")
					this_.repeatPlay(s.id)
				}
			}
		}
	}()
}

//randPlay 随机播放
func (this_ *Player) randPlay() {
	s := newRandSound(this_.st)
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

func (this_ *Player) repeatPlay(id int) {
	s := newSound(id)
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
	music.CloseMusicFs()
	showEnd()
}

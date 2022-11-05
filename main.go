package main

import (
	"os"
	"os/signal"
	"scaletrain/core"
	"scaletrain/music"
	"syscall"
)

func main() {

	player := core.InitPlayer()
	player.RunPlayer()
	sigStop()
	player.Stop()
}

func sigStop() {
	c := make(chan os.Signal, syscall.SIGKILL) // 定义一个信号的通道
	signal.Notify(c, syscall.SIGINT)           // 转发键盘中断信号到c
	<-c                                        // 阻塞
	music.CloseMusicFs()
}

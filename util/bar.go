package util

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"os"
	"time"
)

func WaitBar(barName string, sec int) {
	ch := make(chan struct{})
	bar := progressbar.NewOptions(
		-1,
		progressbar.OptionSetDescription(barName),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(false),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetRenderBlankState(true),
	)
	go func() {
		defaultMilli := 10
		allMilli := sec * 1000
		for i := 0; i < allMilli/defaultMilli; i++ {
			time.Sleep(time.Duration(defaultMilli) * time.Millisecond)
			bar.Add(1)
		}
		ch <- struct{}{}
	}()
	<-ch
	fmt.Println()
}

package core

import (
	"encoding/json"
	"fmt"
	"log"
	"scaletrain/music"
	"scaletrain/util"
)

const defaultScore = 1 //默认分数500

type scoreTyp = map[int]int32

var scores scoreTyp
var totalScale int64 //总共答题
var sureScale int64  //答对

func InitScore() {
	scores = reloadScore()
	if len(scores) == 0 {
		scores = make(scoreTyp, len(music.MusicsFs))
		for i := range music.MusicsFs {
			scores[i] = defaultScore
		}
	}
	fmt.Println("训练样本加载完成")
}

//randIndexByScore 根据评分随机出一个index
func randIndexByScore(st music.SoundTyp) int {
	poolTyp := make(util.RandPoolTyp)
	total := util.DefaultInt(0)
	for id, s := range scores {
		if st == "" {
			total += s
		} else {
			if music.MusicsFs[id].SoundType == st {
				total += s
			}
		}
	}
	for id, weight := range scores {
		if st == "" {
			poolTyp[util.DefaultInt(id)] = total - weight
		} else {
			if music.MusicsFs[id].SoundType == st {
				poolTyp[util.DefaultInt(id)] = total - weight
			}
		}
	}
	return int(util.RandOne(poolTyp))
}

func AddScore(id int) {
	scores[id] += 1
	totalScale += 1
	sureScale += 1
	PrintScore()
}

func SubScore(id int) {
	if scores[id] > 1 {
		scores[id] -= 1
	}
	totalScale += 1
	PrintScore()
}

func PrintScore() {
	percent := (sureScale * 100) / totalScale
	fmt.Println(fmt.Sprintf("正确率 %d/%d（%d%%)", sureScale, totalScale, percent))
}

//saveScore 存储分数
func saveScore() {
	data, err := json.Marshal(scores)
	if err != nil {
		log.Fatal(err)
	}
	util.TruncateWrite(util.ScoreFile, data)
	fmt.Println("训练样本存储完成")
}

//reloadScore 加载分数
func reloadScore() scoreTyp {
	if !util.IsExtraFile(util.ScoreFile) {
		return nil
	}
	scores = make(scoreTyp)
	data := util.ReadFile(util.ScoreFile)
	json.Unmarshal(data, &scores)
	return scores
}

func showEnd() {
	fmt.Println("=======================")
	fmt.Println("总共答题", totalScale, "；答对", sureScale)
}

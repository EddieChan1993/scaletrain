package core

import (
	"encoding/json"
	"fmt"
	"log"
	"scaletrain/music"
	"scaletrain/util"
)

const scoreFile = util.PathScore + "score_file.json"
const defaultScore = 500 //默认分数500

type scoreTyp = map[int]int32

var scores scoreTyp

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
func randIndexByScore() int {
	poolTyp := make(util.RandPoolTyp)
	total := util.DefaultInt(0)
	for _, s := range scores {
		total += s
	}
	for id, weight := range scores {
		poolTyp[util.DefaultInt(id)] = total - weight
	}
	return int(util.RandOne(poolTyp))
}

func AddScore(id int) {
	scores[id] += 1
}

//saveScore 存储分数
func saveScore() {
	data, err := json.Marshal(scores)
	if err != nil {
		log.Fatal(err)
	}
	util.TruncateWrite(scoreFile, data)
	fmt.Println("训练样本存储完成")
}

//reloadScore 加载分数
func reloadScore() scoreTyp {
	if !util.IsExtraFile(scoreFile) {
		return nil
	}
	scores = make(scoreTyp)
	data := util.ReadFile(scoreFile)
	json.Unmarshal(data, &scores)
	return scores
}

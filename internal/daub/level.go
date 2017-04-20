package daub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/yulefox/lamp/internal/db"
)

// VIP/等级指数分布
type fixA struct {
	desire   float64
	minLevel int32
}

// 留存指数分布
type fixR struct {
	desire float64
}

var (
	fixAs = []fixA{
		0:  {26.0, 1},
		1:  {26.0, 3},
		2:  {26.0, 10},
		3:  {26.0, 20},
		4:  {26.0, 30},
		5:  {26.0, 30},
		6:  {26.0, 35},
		7:  {26.0, 40},
		8:  {26.0, 45},
		9:  {26.0, 50},
		10: {26.0, 55},
		11: {26.0, 60},
		12: {26.0, 65},
		13: {26.0, 70},
		14: {26.0, 75},
		15: {26.0, 80},
	}

	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	levels db.Levels
)

func loadLevels() {
	sql := "query_levels_active"
	days := 10000
	fileName := fmt.Sprintf("%s_%d_fixed.json", sql, days)
	buf, err := ioutil.ReadFile(fileName)

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buf, &levels)
	if err != nil {
		panic(err)
	}
}

// CalcLevel .
func CalcLevel(vip, maxTT int32) (lvl, combat, totaltime int32) {
	f := fixAs[vip]
	loop := 0
	v := vip

	for {
		lvl = int32(r.ExpFloat64()*f.desire) + f.minLevel
		if vip < 5 && lvl > 75 {
			continue
		}

		l, e := db.GetLevel(levels, lvl, v)
		if loop > 5 {
			if vip < 8 {
				l, e = db.GetLevel(levels, 40, 0)
			} else {
				l, e = db.GetLevel(levels, 75, 5)
			}
			loop = 0
		}
		if e && l.Count > 5 && l.TotalTimeAvg < maxTT {
			lvl = l.Level
			combat = int32(r.NormFloat64()*1.0 + float64(l.CombatAvg))
			totaltime = int32(r.NormFloat64()*1.0 + float64(l.TotalTimeAvg))
			return lvl, combat, totaltime
		}
		loop++
	}
}

func init() {
	loadLevels()
}

package daub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/yulefox/lamp/internal/db"
)

var (
	rn = rand.New(rand.NewSource(time.Now().UnixNano()))

	names []db.Role
)

func loadNames() {
	sql := "get_roles"
	fileName := fmt.Sprintf("%s_fixed.json", sql)
	buf, err := ioutil.ReadFile(fileName)
	names = make([]db.Role, 100000)

	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buf, &names)
	if err != nil {
		panic(err)
	}
}

// CalcName .
func CalcName(vip int32) string {
	loop := 0
	size := int32(len(names))

	for {
		n := rn.Int31n(size)
		name := names[n]

		if loop > 5 && vip > 0 && name.VIP < 5 && name.Level < 50 {
			loop++
			continue
		}
		return name.Name
	}
}

func init() {
	loadNames()
}

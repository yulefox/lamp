package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/urfave/cli"
	"github.com/yulefox/lamp/internal/db"
)

var (
	workersNum = 50
	conf       config
	lastTime   int32
	wg         sync.WaitGroup
	rolesNum   int32
	levels     []level
	ch         chan db.Role
	chExit     chan bool
	roles      db.Roles
)

type config struct {
	Opt     db.Option         `json:"db_opt"`
	Hosts   map[string]string `json:"hosts"`
	DBs     map[string]server `json:"dbs"`
	server  server
	context *db.Context
}

type server struct {
	Index        string
	Servers      [][]int32 `json:"servers"`
	Host         string    `json:"host"`
	StartTime    string    `json:"start_time"`
	maxLevel     int32
	maxTotalTime int32
}

type level struct {
	lvlA       int32
	lvlB       int32
	vipM       int32
	combatA    int32
	combatB    int32
	totalTimeA int32
	totalTimeB int32
	rate       int32
	rateVIP    int32
}

func (c *config) load(filename, dbNo string) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(buf, c)
	c.server = c.DBs[dbNo]
	c.server.Index = dbNo
	c.Opt.Name = "agame_" + dbNo
	c.Opt.Host = conf.Hosts[c.server.Host]
	c.context = db.NewContext(c.Opt)
}

func calcLevel(r *db.Role) {
	rn := rand.Int31n(10000)
	var lc, ln level

	for i, lvl := range levels {
		if r.VIP == 0 {
			if rn < lvl.rate {
				lc = levels[i]
				ln = levels[i+1]
				break
			}
		} else if r.VIP <= lvl.vipM {
			if rn < lvl.rateVIP {
				lc = levels[i]
				ln = levels[i+1]
				break
			}
		}
	}
	r.Level = lc.lvlA
	r.Combat = lc.combatA
	if lc.combatB > 0 {
		r.Combat += rand.Int31n(lc.combatB)
	}
	r.TotalTime = lc.totalTimeA
	if lc.totalTimeB > 0 {
		r.TotalTime += rand.Int31n(lc.totalTimeB)
	}

	// rand level
	d := int32(0)
	if lc.lvlB > 0 {
		d = rand.Int31n(lc.lvlB)
	}
	if d > 0 {
		r.Level += d
		if ln.combatA > lc.combatA {
			r.Combat += d * (ln.combatA - lc.combatA) / lc.lvlB
		}
		if ln.totalTimeA > lc.totalTimeA {
			r.TotalTime += d * (ln.totalTimeA - lc.totalTimeA) / lc.lvlB
		}
	}
}

func fixEnterTime(r *db.Role) {
	atomic.AddInt32(&rolesNum, 1)
	r.EnterTime = r.CreateTime + int32(float32(r.TotalTime)*(2.0+rand.Float32()*12.0))
	if r.EnterTime > lastTime {
		r.EnterTime = lastTime - 10000
	}
	if r.TotalTime > 10000 {
		r.LeaveTime = r.EnterTime + rand.Int31n(10000)
	} else if r.TotalTime > 0 {
		r.LeaveTime = r.EnterTime + rand.Int31n(r.TotalTime)
	}
}

func fixLevel(r *db.Role) {
	atomic.AddInt32(&rolesNum, 1)
	calcLevel(r)
	r.MoneyA = rand.Int31n((r.Level+r.VIP)*5 + 1)
	r.MoneyB = rand.Int31n(int32(math.Pow(float64(r.Level+r.VIP), 3)*0.5+1)) + rand.Int31n(int32(math.Pow(float64(r.Level+r.VIP), 2)*0.8+1))
}

func updateRoles(stmt *sql.Stmt) {
	for {
		select {
		case r := <-ch:
			db.UpdateRole(stmt, r)
			atomic.AddInt32(&rolesNum, -1)
			if rolesNum <= 0 {
				for i := 0; i < workersNum; i++ {
					chExit <- true
				}
			}
			fmt.Printf("%d) %s\n", rolesNum, r.String())
			break
		case <-chExit:
			wg.Done()
			return
		}
	}
}

// Check .
func check(c *sql.DB) {
	var log bytes.Buffer
	log.WriteString(fmt.Sprintf("%+v\n", conf.server))
	rows, err := c.Query(db.Queries["check"])

	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var lvl, num, minTT, maxTT, minCB, maxCB int
		rows.Scan(&lvl, &num, &minTT, &maxTT, &minCB, &maxCB)
		log.WriteString(fmt.Sprintf("%4d %5d | %10d-%10d, %10d-%10d\n", lvl, num, minTT, maxTT, minCB, maxCB))
	}

	ioutil.WriteFile(conf.server.Index+".json", log.Bytes(), 0644)
}

func fix(c *cli.Context) {
	file := c.Args().Get(0)
	dbNo := c.Args().Get(1)
	conf.load(file, dbNo)

	t, _ := time.Parse("2006-01-02 03:04", conf.server.StartTime)
	mt := int32(time.Now().Sub(t).Seconds() * 2 / 5)
	conf.server.maxTotalTime = mt
	for i, l := range levels {
		dtA := mt - l.totalTimeA
		dtB := l.totalTimeB - dtA
		if dtB > 0 {
			levels[i].totalTimeB = dtA
		}
		if dtA < 0 {
			levels[i-1].lvlB *= dtA / (levels[i].totalTimeA - levels[i-1].totalTimeA)
			levels[i-1].rate = 10000
			levels[i-1].rateVIP = 10000
			levels[i].totalTimeA = mt
			levels[i].totalTimeB = 0
			conf.server.maxLevel = levels[i-1].lvlA + levels[i-1].lvlB
			break
		}
	}
	db.GetRoles(conf.context.DB, roles, "query_roles")
	for _, r := range roles {
		calcLevel(r)
		fixLevel(r)
		fixEnterTime(r)
		ch <- *r
	}

	stmt, err := conf.context.Prepare(db.Queries["update_role"])

	if err != nil {
		panic(err)
	}

	// handle update queries
	for i := 0; i < workersNum; i++ {
		wg.Add(1)
		go updateRoles(stmt)
	}
	wg.Wait()
	check(conf.context.DB)
	conf.context.Close()
}

func dbname(r *db.Role) string {
	n := r.Server / 100
	name := strconv.Itoa(int(n))
	for name, s := range conf.DBs {
		for _, rs := range s.Servers {
			if r.Server >= rs[0] && r.Server <= rs[1] {
				return name
			}
		}
	}
	return name
}

func top(c *cli.Context) {
	file := c.Args().Get(0)
	dbName := c.Args().Get(1)
	conf.load(file, dbName)
	d := conf.context.DB

	fmt.Println(conf.Opt)
	db.GetCharges(d, roles, "query_charges", 500)
	d.Close()

	var csv bytes.Buffer
	for _, r := range roles {
		dbName = dbname(r)
		conf.load(file, dbName)
		d = conf.context.DB
		db.GetRoles(d, roles, "get_role", r.ID)
		d.Close()
		csv.WriteString(r.CSVString())
	}
	ioutil.WriteFile("top.csv", csv.Bytes(), 0644)
	os.Exit(0)
}

func init() {
	t := time.Now()
	rand.Seed(t.UnixNano())

	lastTime = int32(t.Unix())
	fmt.Println(t, lastTime)

	levels = []level{
		{1, 1, 0, 420, 0, 20, 5000, 1894, 0},
		{2, 1, 0, 499, 1021, 605, 10001, 2613, 0},
		{3, 8, 1, 577, 1801, 1032, 30030, 4924, 150},
		{11, 10, 2, 1282, 5028, 4022, 40002, 6045, 1000},
		{21, 10, 3, 3562, 20816, 8113, 100801, 8121, 3894},
		{31, 10, 4, 8954, 40501, 17129, 202013, 9253, 8613},
		{41, 10, 6, 18026, 81303, 77201, 400020, 9990, 9920},
		{51, 10, 8, 43027, 153047, 102082, 601030, 10000, 9999},
		{61, 10, 9, 78019, 200035, 271004, 802191, 10000, 10000},
		{71, 10, 11, 150401, 261012, 602021, 1200303, 10000, 10000},
		{81, 10, 15, 310301, 350203, 1501019, 1503205, 10000, 10000},
		{91, 10, 15, 450282, 601210, 2200281, 2000011, 10000, 10000},
		{101, 10, 15, 450282, 601210, 4200281, 0, 10000, 10000},
	}

	roles = make(map[int64]*db.Role)
	ch = make(chan db.Role, 65535)
	chExit = make(chan bool, workersNum)
}

// Run .
func Run(args []string) {
	app := cli.NewApp()
	app.Action = top
	app.Run(args)
}

func main() {
	Run(os.Args)
}

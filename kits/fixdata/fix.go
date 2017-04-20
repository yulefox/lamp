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

	"strings"

	"github.com/urfave/cli"
	"github.com/yulefox/lamp/internal/cmd"
	"github.com/yulefox/lamp/internal/daub"
	"github.com/yulefox/lamp/internal/db"
)

var (
	timeLayout = "2006-01-02 03:04"
	workersNum = 50
	conf       config
	wg         sync.WaitGroup
	rolesNum   int32
	levels     db.Levels
	ch         chan db.Role
	chExit     chan bool
	roles      db.Roles
	names      []Name
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
	StartTimeStr string    `json:"start_time"`
	StartTime    time.Time
	NowTime      time.Time
	Duration     time.Duration
	maxTotalTime int32
}

// Name .
type Name struct {
	Name  string `json:"name"`  // 角色名
	Class int32  `json:"class"` // 职业
	VIP   int32  `json:"vip"`   // VIP
	Level int32  `json:"level"` // 等级
}

func (c *config) load(fileName, dbNo string) {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(buf, c)
	c.server = c.DBs[dbNo]
	c.server.Index = dbNo
	c.server.StartTime, _ = time.Parse(timeLayout, conf.server.StartTimeStr)
	c.server.NowTime = time.Now()
	c.server.Duration = c.server.NowTime.Sub(c.server.StartTime)
	conf.server.maxTotalTime = int32(c.server.Duration.Seconds() / 2.5)
	if _, err := strconv.Atoi(dbNo); err != nil {
		c.Opt.Name = dbNo
	} else {
		c.Opt.Name = "agame_" + dbNo
	}
	c.Opt.Host = conf.Hosts[c.server.Host]
	c.context = db.NewContext(c.Opt)
}

func avg(v1, c1, v2, c2 int32) int32 {
	return (v1*c1 + v2*c2) / (c1 + c2)
}

func mergeLevels(src, dst db.Levels) {
	for k, vs := range src {
		if vd, exists := dst[k]; exists {
			vd.CombatAvg = avg(vs.CombatAvg, vs.Count, vd.CombatAvg, vd.Count)
			vd.CombatMin = avg(vs.CombatMin, vs.Count, vd.CombatMax, vd.Count)
			vd.CombatMax = avg(vs.CombatMax, vs.Count, vd.CombatMax, vd.Count)
			vd.TotalTimeAvg = avg(vs.TotalTimeAvg, vs.Count, vd.TotalTimeAvg, vd.Count)
			vd.TotalTimeMin = avg(vs.TotalTimeMin, vs.Count, vd.TotalTimeMin, vd.Count)
			vd.TotalTimeMax = avg(vs.TotalTimeMax, vs.Count, vd.TotalTimeMax, vd.Count)
			vd.Count += vs.Count
		} else {
			vd = &db.Level{}
			db.Merge(vs, vd, true)
			dst[k] = vd
		}
	}
}

func calcSampleLevels(c *cli.Context) {
	file := c.Args().Get(0)
	dbName := c.Args().Get(1)
	conf.load(file, dbName)

	sql := "query_levels_active"
	days := 10000
	for dbName := range conf.DBs {
		fileName := fmt.Sprintf("%s_%d_%s.json", sql, days, dbName)
		buf, err := ioutil.ReadFile(fileName)
		ls := make(db.Levels)

		if err != nil {
			continue
		}
		err = json.Unmarshal(buf, &ls)
		if err != nil {
			panic(err)
		}
		fmt.Println(fileName)
		mergeLevels(ls, levels)
	}

	var buf bytes.Buffer
	fileName := fmt.Sprintf("%s_%d_fixed.json", sql, days)
	b, _ := json.MarshalIndent(levels, "", "  ")
	buf.Write(b)
	ioutil.WriteFile(fileName, buf.Bytes(), 0644)
}

func calcSampleNames(c *cli.Context) {
	file := c.Args().Get(0)
	dbName := c.Args().Get(1)
	conf.load(file, dbName)
	names = make([]Name, 0)

	sql := "get_roles"
	for dbName := range conf.DBs {
		fileName := fmt.Sprintf("%s_%s.json", sql, dbName)
		buf, err := ioutil.ReadFile(fileName)
		ns := make([]Name, 100000)

		if err != nil {
			continue
		}
		err = json.Unmarshal(buf, &ns)
		if err != nil {
			panic(err)
		}
		fmt.Println(fileName)
		names = append(names, ns...)
	}

	var buf bytes.Buffer
	fileName := fmt.Sprintf("%s_fixed.json", sql)
	b, _ := json.MarshalIndent(names, "", "  ")
	buf.Write(b)
	ioutil.WriteFile(fileName, buf.Bytes(), 0644)
}

func sampleNames(c *cli.Context) {
	file := c.Args().Get(0)
	dbName := c.Args().Get(1)
	conf.load(file, dbName)
	d := conf.context.DB
	sql := "get_roles"

	db.GetRoles(d, roles, sql)
	names = make([]Name, 0)

	for _, r := range roles {
		n := Name{r.Name, r.Class, r.VIP, r.Level}
		names = append(names, n)
	}

	var buf bytes.Buffer
	fileName := fmt.Sprintf("%s_%s.json", sql, dbName)
	b, _ := json.MarshalIndent(names, "", "  ")
	buf.Write(b)
	ioutil.WriteFile(fileName, buf.Bytes(), 0644)
}

func sampleLevels(c *cli.Context) {
	file := c.Args().Get(0)
	dbName := c.Args().Get(1)
	conf.load(file, dbName)
	d := conf.context.DB
	sql := "query_levels_active"
	days := 10000

	db.GetLevels(d, levels, sql, days)

	var buf bytes.Buffer
	fileName := fmt.Sprintf("%s_%d_%s.json", sql, days, dbName)
	b, _ := json.MarshalIndent(levels, "", "  ")
	buf.Write(b)
	ioutil.WriteFile(fileName, buf.Bytes(), 0644)
}

func fixRole(r *db.Role) {
	r.MoneyA = rand.Int31n((r.Level+r.VIP)*5 + 1)
	r.MoneyB = rand.Int31n(int32(math.Pow(float64(r.Level+r.VIP), 3)*0.5+1)) + rand.Int31n(int32(math.Pow(float64(r.Level+r.VIP), 2)*0.8+1))
	r.Class = 101 + rand.Int31n(3)
}

func randClock() time.Duration {
	h := rand.Int63n(24)
	if h > 0 && h < 9 {
		h = rand.Int63n(24)
	}
	m := rand.Int63n(60)
	s := rand.Int63n(60)
	return time.Duration(h*3600 + m*60 + s*60)
}

/*
func fixTime(r *db.Role) {
	ct := time.Unix(int64(r.CreateTime), 0).Truncate(24 * time.Hour).Add(randClock())
	r.CreateTime = int32(ct.Unix())
	r.EnterTime = r.CreateTime + int32(float32(r.TotalTime)*(2.0+rand.Float32()*12.0))
	et := time.Unix(int64(r.EnterTime), 0)
	now := conf.server.NowTime
	if et.After(now) {
		et = now.Add(time.Duration(-rand.Int31n(10000)))
		r.EnterTime = int32(et.Unix())
	}
	if r.EnterTime > lastTime {
		r.EnterTime = lastTime - 10000
	}
	if r.TotalTime > 10000 {
		r.LeaveTime = r.EnterTime + rand.Int31n(10000)
	} else if r.TotalTime > 0 {
		r.LeaveTime = r.EnterTime + rand.Int31n(r.TotalTime)
	}
}
*/

func fixRoles(roles db.Roles) {
	sum := int32(0)
	ls := make(map[int32]int32)

	for _, v := range roles {
		v.Name = daub.CalcName(v.VIP)

		//fixRole(v)
		//fixTime(v)
		//fixOrder(v)
		ch <- *v
		//ls[v.Level]++
	}

	for i := int32(1); i <= 110; i++ {
		sum += ls[i]
		fmt.Println(i, sum)
	}
	fmt.Println(rolesNum)
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
	db.GetRoles(conf.context.DB, roles, "query_roles")
	atomic.AddInt32(&rolesNum, int32(len(roles)))
	fixRoles(roles)
	os.Exit(0)

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
	db.GetCharges(d, roles, "query_charges", 1000)
	d.Close()

	ra := make([]*db.Role, 0)
	for _, r := range roles {
		dbName = dbname(r)
		conf.load(file, dbName)
		d = conf.context.DB
		db.GetRoles(d, roles, "get_role", r.ID)
		fmt.Println(r)
		ra = append(ra, r)
		d.Close()
	}
	m := map[string]interface{}{"roles": ra}

	var buf bytes.Buffer
	fileName := "top1000.json"
	b, _ := json.MarshalIndent(m, "", "  ")
	buf.Write(b)
	ioutil.WriteFile(fileName, buf.Bytes(), 0644)
	cmd.SCP("20202", fileName, "juice@sdk.91juice.com:/juice/gmproxy/json")
	os.Exit(0)
}

func dist(c *cli.Context) {
	file := c.Args().Get(0)
	sql := c.Args().Get(1)
	dbNames := make([]string, 0)
	if c.NArg() > 2 {
		dbs := c.Args().Get(2)
		dbNames = strings.Split(dbs, ",")
	} else {
		dbNames = []string{"1001", "2001", "2070", "2136", "2154", "2161", "2173", "2174", "2175", "2176", "2177", "6001", "6030", "6036", "6039", "9001", "9011", "9020"}
	}
	dists := make(db.Dists)
	for _, dbName := range dbNames {
		conf.load(file, dbName)
		d := conf.context.DB

		db.GetDistribution(d, dists, sql)
	}
	da := make([]*db.Dist, 0)
	for _, v := range dists {
		da = append(da, v)
	}
	m := map[string]interface{}{"dist": da}

	var buf bytes.Buffer
	fileName := sql + ".json"
	b, _ := json.MarshalIndent(m, "", "  ")
	buf.Write(b)
	ioutil.WriteFile(fileName, buf.Bytes(), 0644)
	cmd.SCP("20202", fileName, "juice@sdk.91juice.com:/juice/gmproxy/json")
	os.Exit(0)
}

func init() {
	roles = make(db.Roles)
	levels = make(db.Levels)
	ch = make(chan db.Role, 65535)
	chExit = make(chan bool, workersNum)
}

// Run .
func Run(args []string) {
	app := cli.NewApp()
	app.Action = fix
	app.Run(args)
}

func main() {
	Run(os.Args)
}

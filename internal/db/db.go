package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql
)

// Option .
type Option struct {
	User     string `json:"user" opt:"user" default:"root"`
	Password string `json:"password" opt:"password" default:"root"`
	Host     string `json:"host" opt:"host" default:"localhost"`
	Port     int    `json:"port" opt:"port" default:"3306"`
	Name     string `json:"name" opt:"name" default:"test"`
	CharSet  string `json:"char_set" opt:"charset" default:"utf8"`
}

// Node .
type Node struct {
	Index        string
	Servers      [][]int32 `json:"servers"`
	Host         string    `json:"host"`
	StartTimeStr string    `json:"start_time"`
	StartTime    time.Time
	NowTime      time.Time
	Duration     time.Duration
	maxTotalTime int32
}

// Context .
type Context struct {
	*sql.DB
	Opt      Option
	QueryStr string
}

// Config .
type Config struct {
	Opt     Option            `json:"db_opt"`
	Hosts   map[string]string `json:"hosts"`
	DBs     map[string]Node   `json:"dbs"`
	Node    Node
	context *Context
}

func (o Option) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		o.User,
		o.Password,
		o.Host,
		o.Port,
		o.Name,
		o.CharSet)
}

// NewContext .
func NewContext(opt Option) *Context {
	db, err := sql.Open("mysql", opt.String())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &Context{
		DB:  db,
		Opt: opt,
	}
}

func (c *Context) check() {
	stmt, err := c.Prepare(c.QueryStr)

	if err != nil {
		panic(err)
	}
	stmt.Exec()
}

func init() {
	filepath := path.Join(os.Getenv("LAMP_CONFIG_PATH"), filename)
	buf, err := ioutil.ReadFile("db.json")
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
	c.context = NewContext(c.Opt)
}

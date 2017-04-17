package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql" // mysql
)

// Option .
type Option struct {
	User     string `opt:"user" default:"root"`
	Password string `opt:"password" default:"root"`
	Host     string `opt:"host" default:"localhost"`
	Port     int    `opt:"port" default:"3306"`
	Name     string `opt:"name" default:"test"`
	CharSet  string `opt:"charset" default:"utf8"`
}

// Context .
type Context struct {
	*sql.DB
	Opt      Option
	QueryStr string
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

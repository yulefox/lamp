package db

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql" // mysql
)

func TestNewContext(t *testing.T) {
	c := NewContext(Option{
		User:     "root",
		Password: "root",
		Host:     "localhost",
		Port:     3306,
		Name:     "test",
		CharSet:  "utf8",
	})

	fmt.Println(c.Opt)
	rows, err := c.Query("select id, vip from dat_role")

	if err != nil {
		panic(err)
	}
	for rows.Next() {
		var id int64
		var vip int
		rows.Scan(&id, &vip)
		fmt.Printf("id: %d, vip: %d\n", id, vip)
	}
	fmt.Println(c.Opt)
}

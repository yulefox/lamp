package db

import (
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql" // mysql
)

var hosts = map[string]string{
	"local":  "192.168.1.111",
	"loc130": "192.168.1.130",
	"00101":  "112.126.84.253",
	"00801":  "13.113.177.153",
	"00901":  "222.187.253.58",
	"11600":  "112.126.86.133",
	"100100": "101.201.152.55",
	"200100": "123.56.182.231",
	"200200": "112.126.87.213",
	"200300": "47.93.78.241",
	"600100": "123.56.201.33",
	"600200": "47.93.82.164",
	"300100": "123.206.103.61",
	"400200": "60.205.94.201",
	"700100": "112.126.88.162",
	"800100": "52.198.61.142",
	"900100": "222.187.253.72",
	"900300": "222.186.190.91",
}

func TestNewContext(t *testing.T) {
	c := NewContext(Option{
		User:     "lamp",
		Password: "Lamp@91juice.com",
		Host:     hosts["700100"],
		Port:     3306,
		Name:     "agame_7001",
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

package db

import (
	"database/sql"
	"fmt"
	"time"
)

var (
	// Queries .
	Queries map[string]string
)

// Roles .
type Roles map[int64]*Role

// Role .
type Role struct {
	ID               int64
	User             string
	Name             string
	Server           int32
	Channel          string
	VIP              int32
	Charge           int32
	Level            int32
	Combat           int32
	MoneyA           int32
	MoneyB           int32
	CreateTime       int32
	EnterTime        int32
	LeaveTime        int32
	TotalTime        int32
	FirstChargeLevel int32
}

func (r *Role) String() string {
	return fmt.Sprintf("%d: %d, %s\n  %d, %d\n  %s-%s", r.ID, r.Server, r.Name, r.Level, r.TotalTime,
		time.Unix(int64(r.EnterTime), 0).Format("2006-01-02 03:04"),
		time.Unix(int64(r.LeaveTime), 0).Format("2006-01-02 03:04"))
}

// CSVString .
func (r *Role) CSVString() string {
	return fmt.Sprintf("%d,%d,%d,%s,%s,%d,%d,%d,%d,%d,%s,%s,%s,",
		r.Charge, r.Server, r.ID, r.User, r.Name, r.Level,
		r.MoneyA, r.MoneyB, r.TotalTime, r.Combat,
		time.Unix(int64(r.CreateTime), 0).Format("2006-01-02 03:04"),
		time.Unix(int64(r.EnterTime), 0).Format("2006-01-02 03:04"),
		time.Unix(int64(r.LeaveTime), 0).Format("2006-01-02 03:04"))
}

// GetRoles .
func GetRoles(db *sql.DB, roles Roles, sql string, args ...interface{}) {
	rows, err := db.Query(Queries[sql], args...)

	if err != nil {
		fmt.Println(err)
		return
	}

	// traverse
	for rows.Next() {
		r := &Role{}
		if len(args) > 0 {
			r = roles[args[0].(int64)]
		}
		rows.Scan(&r.ID, &r.User, &r.Name, &r.Server, &r.Channel, &r.VIP,
			&r.Level, &r.Combat, &r.MoneyA, &r.MoneyB, &r.CreateTime, &r.EnterTime, &r.LeaveTime,
			&r.TotalTime)
		roles[r.ID] = r
		fmt.Println(r.CSVString())
	}
}

// GetCharges .
func GetCharges(db *sql.DB, roles Roles, sql string, args ...interface{}) {
	rows, err := db.Query(Queries[sql], args...)

	if err != nil {
		fmt.Println(err)
		return
	}

	// traverse
	for rows.Next() {
		r := &Role{}
		rows.Scan(&r.ID, &r.Server, &r.Charge, &r.FirstChargeLevel)
		roles[r.ID] = r
	}
}

// UpdateRole .
func UpdateRole(stmt *sql.Stmt, r Role) {
	_, err := stmt.Exec(r.Level, r.Combat, r.TotalTime, r.MoneyA, r.MoneyB,
		r.EnterTime, r.LeaveTime, r.EnterTime, r.LeaveTime, r.ID)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func init() {
	roleFields := "id, user, name, server, channel, vip, lvl, combat, money_cur_a, money_cur_b, createtime, entertime, leavetime, totaltime"
	orderFields := "rid, server, round(sum(price)/100) as charge, min(lvl)"

	Queries = map[string]string{
		"query_charges": "select " + orderFields + " from dat_order where server > 100000 group by rid order by charge desc limit ?",
		"query_roles":   "select " + roleFields + " from dat_role where channel in ('yd', 'ly', 'yljh')",
		"get_role":      "select " + roleFields + " from dat_role where id = ?",
		"check":         "select lvl, count(*) as num, min(totaltime) min_tt, max(totaltime) max_tt, min(combat) min_cb, max(combat) max_cb, max(entertime_s), max(leavetime_s) from dat_role where channel in ('ly', 'yd', 'yljh') group by lvl",
		"update_role":   "update dat_role set lvl=?, combat=?, totaltime=?, money_cur_a=?, money_cur_b=?, entertime=?, leavetime=?, entertime_s=from_unixtime(?), leavetime_s=from_unixtime(?), stamp=unix_timestamp(), stamp_s=now() where id=?",
	}
}

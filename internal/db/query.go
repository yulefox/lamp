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

// Dists .
type Dists map[string]*Dist

// Roles .
type Roles map[int64]*Role

// Levels .
type Levels map[string]*Level

// Dist Distribution
type Dist struct {
	Name  string `json:"name"`  // 名称
	Value int    `json:"value"` // 数值
}

// Role .
type Role struct {
	ID               int64  `json:"id"`                 // 角色 ID
	User             string `json:"user"`               // 用户名
	Name             string `json:"name"`               // 角色名
	Class            int32  `json:"class"`              // 职业
	Server           int32  `json:"server"`             // 区服
	Channel          string `json:"channel"`            // 渠道
	VIP              int32  `json:"vip"`                // VIP
	Charge           int32  `json:"charge"`             // 充值金额
	Level            int32  `json:"level"`              // 等级
	Combat           int32  `json:"combat"`             // 战力
	MoneyA           int32  `json:"money_a"`            // 剩余钻石
	MoneyB           int32  `json:"money_b"`            // 剩余金币
	CreateTime       int64  `json:"create_time"`        // 创建时间
	EnterTime        int64  `json:"enter_time"`         // 最后上线时间
	LeaveTime        int64  `json:"leave_time"`         // 最后下线时间
	TotalTime        int64  `json:"total_time"`         // 累计在线时长
	FirstChargeLevel int32  `json:"first_charge_level"` // 首次充值等级
	LastChargeLevel  int32  `json:"last_charge_level"`  // 最后充值等级
	FirstChargeTime  string `json:"first_charge_time"`  // 首次充值时间
	LastChargeTime   string `json:"last_charge_time"`   // 最后充值时间
	LogTimes         []int32
}

func (r *Role) String() string {
	return fmt.Sprintf("%20d: %6d, %3d, %10d, %s-%s, %s", r.ID, r.Server, r.Level, r.TotalTime,
		time.Unix(int64(r.EnterTime), 0).Format("2006-01-02 03:04"),
		time.Unix(int64(r.LeaveTime), 0).Format("2006-01-02 03:04"),
		r.Name)
}

// Level Level.
type Level struct {
	ID           string // Level_VIP
	Level        int32
	VIP          int32
	Count        int32
	CombatMin    int32
	CombatMax    int32
	CombatAvg    int32
	TotalTimeMin int32
	TotalTimeMax int32
	TotalTimeAvg int32
	CreateTime   string  // 最早创建时间
	Scale        float32 // 等级分布
}

func (l *Level) id() {
	l.ID = fmt.Sprintf("%03d_%02d", l.Level, l.VIP)
}

// GetLevel .
func GetLevel(lvls Levels, lvl, vip int32) (*Level, bool) {
	id := fmt.Sprintf("%03d_%02d", lvl, vip)
	l, e := lvls[id]

	return l, e
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
		rows.Scan(&r.ID, &r.User, &r.Name, &r.Server, &r.Channel, &r.Class, &r.VIP,
			&r.Level, &r.Combat, &r.MoneyA, &r.MoneyB, &r.CreateTime, &r.EnterTime, &r.LeaveTime,
			&r.TotalTime)
		if dest, exists := roles[r.ID]; exists {
			Merge(r, dest, true)
		} else {
			roles[r.ID] = r
		}
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
		rows.Scan(&r.ID, &r.Server, &r.Charge, &r.FirstChargeLevel, &r.LastChargeLevel, &r.FirstChargeTime, &r.LastChargeTime)
		if dest, exists := roles[r.ID]; exists {
			Merge(r, dest, true)
		} else {
			roles[r.ID] = r
		}
	}
}

// GetTotalTimes .
func GetTotalTimes(db *sql.DB, levels Levels, sql string, args ...interface{}) {
	rows, err := db.Query(Queries[sql], args...)

	if err != nil {
		fmt.Println(err)
		return
	}

	// traverse
	for rows.Next() {
		l := &Level{}
		rows.Scan(&l.Level, &l.VIP, &l.Count, &l.CombatMin, &l.CombatMax, &l.CombatAvg,
			&l.TotalTimeMin, &l.TotalTimeMax, &l.TotalTimeAvg, &l.CreateTime)
		fmt.Println(l)
		l.id()
		if dest, exists := levels[l.ID]; exists {
			Merge(l, dest, true)
		} else {
			levels[l.ID] = l
		}
	}
}

// GetLevels .
func GetLevels(db *sql.DB, levels Levels, sql string, args ...interface{}) {
	rows, err := db.Query(Queries[sql], args...)

	if err != nil {
		fmt.Println(err)
		return
	}

	// traverse
	for rows.Next() {
		l := &Level{}
		rows.Scan(&l.Level, &l.VIP, &l.Count, &l.CombatMin, &l.CombatMax, &l.CombatAvg,
			&l.TotalTimeMin, &l.TotalTimeMax, &l.TotalTimeAvg, &l.CreateTime)
		fmt.Println(l)
		l.id()
		if dest, exists := levels[l.ID]; exists {
			Merge(l, dest, true)
		} else {
			levels[l.ID] = l
		}
	}
}

// GetDistribution .
func GetDistribution(db *sql.DB, dists Dists, sql string, args ...interface{}) {
	rows, err := db.Query(Queries[sql], args...)

	if err != nil {
		fmt.Println(err)
		return
	}

	// traverse
	for rows.Next() {
		d := &Dist{}
		rows.Scan(&d.Name, &d.Value)
		if c, exists := dists[d.Name]; exists {
			c.Value += d.Value
		} else {
			dists[d.Name] = d
		}
	}
}

// UpdateRole .
func UpdateRole(stmt *sql.Stmt, r Role) {
	_, err := stmt.Exec(r.Class, r.Level, r.Combat, r.TotalTime, r.MoneyA, r.MoneyB,
		r.EnterTime, r.LeaveTime, r.EnterTime, r.LeaveTime, r.ID)

	if err != nil {
		fmt.Println(err)
		return
	}
}

// LogUser .
func LogUser(stmt *sql.Stmt, r Role) {
	_, err := stmt.Exec(r.User, r.Server, r.Channel, r.EnterTime, r.LeaveTime, r.EnterTime, r.LeaveTime)

	if err != nil {
		fmt.Println(err)
		return
	}
}

// LogRole .
func LogRole(stmt *sql.Stmt, r Role) {
	_, err := stmt.Exec(r.ID, r.Name, r.User, r.Class, r.EnterTime, r.LeaveTime, r.EnterTime, r.LeaveTime)

	if err != nil {
		fmt.Println(err)
		return
	}
}

func init() {
	roleFields := "id, user, name, server, channel, cls, vip, lvl, combat, money_cur_a, money_cur_b, createtime, entertime, leavetime, totaltime"
	chargeFields := "rid, server, round(sum(price)/100) as charge, min(lvl), max(lvl), min(req_time_s), max(req_time_s)"
	orderFields := "oid, cid, uid, rid, price, idx, num, server, channel, req_time_s"

	Queries = map[string]string{
		"query_charges":       "select " + chargeFields + " from dat_order where server > 100000 group by rid order by charge desc limit ?",
		"query_order":         "select " + orderFields + " from dat_order",
		"query_roles":         "select " + roleFields + " from dat_role where channel in ('yd', 'ly', 'yljh')",
		"get_roles":           "select " + roleFields + " from dat_role",
		"get_role":            "select " + roleFields + " from dat_role where id = ?",
		"dist_register":       "select date_format(createtime_s, '%H') as hour, count(*) as count from dat_role group by hour",
		"query_totaltimes":    "select vip, avg(totaltime) from dat_role group by vip",
		"query_levels_active": "select lvl, vip, count(*), min(combat), max(combat), round(avg(combat)), min(totaltime), max(totaltime), round(avg(totaltime)), min(createtime_s) from dat_role where TIMESTAMPDIFF(DAY, entertime_s, now()) < ? group by lvl, vip",
		"query_levels_lost":   "select lvl, vip, count(*), min(combat), max(combat), round(avg(combat)), min(totaltime), max(totaltime), round(avg(totaltime)), min(createtime_s) from dat_role where TIMESTAMPDIFF(DAY, entertime_s, now()) >= ? group by lvl, vip",
		"check":               "select lvl, count(*) as num, min(totaltime) min_tt, max(totaltime) max_tt, min(combat) min_cb, max(combat) max_cb, max(entertime_s), max(leavetime_s) from dat_role where channel in ('ly', 'yd', 'yljh') group by lvl",
		"update_role":         "update dat_role set cls=?, lvl=?, combat=?, totaltime=?, money_cur_a=?, money_cur_b=?, createtime=?, createtime_s=from_unixtime(?), entertime=?, leavetime=?, entertime_s=from_unixtime(?), leavetime_s=from_unixtime(?), stamp=unix_timestamp(), stamp_s=now() where id=?",
		"log_user":            "insert into log_user (user, server, channel, entertime, leavetime, entertime_s, leavetime_s) values(?,?,?,?,?,from_unixtime(?),from_unixtime(?))",
		"log_role":            "insert into log_role (rid, role, user, cls, entertime, leavetime, entertime_s, leavetime_s)values(?,?,?,?,?,?,?,?,from_unixtime(?),from_unixtime(?))",
	}
}

package lamp

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // mysql
)

// OpenDB open DB
func OpenDB() (err error) {
	db, err := sql.Open("mysql", "juice:Juice@1004@192.168.1.111/agame_208")
	if err != nil {
		return
	}
	id := 0
	rows, err := db.Query("SELECT name FROM dat_role WHERE id=?", id)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s is %d\n", name, id)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}

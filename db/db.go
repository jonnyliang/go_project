package db

import (
	"database/sql"
	//	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Mydb() *sql.DB {
	db, err := sql.Open("mysql", "root:root@(X.X.X.X:3306)/ipmi?charset=utf8")

	if err != nil {
		log.Fatalf("open database error %s\n", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("ping error: %s", err)
	}
	return db
}
func Get_ips() []string {
	var slice []string = []string{}
	db := Mydb()
	defer db.Close()
	rows, sel_err := db.Query(" select distinct(substring_index(ipmi_ip,'.',3)) from host_info WHERE ipmi_ip  <> '' ")
	if sel_err != nil {
		log.Fatal(sel_err)
		return slice
	}
	var ips string
	for rows.Next() {
		if rows_err := rows.Scan(&ips); rows_err != nil {
			log.Fatal(rows_err)
			continue
		}
		//fmt.Printf("%s", ips)
		slice = append(slice, ips+".1-250")
	}
	return slice

}

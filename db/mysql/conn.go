package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// 用户名:密码@tcp(IP地址:端口)/数据库名称?连接选项
const mysql_source = "root:iop890@tcp(127.0.0.1:3306)/filestore?charset=utf8"

func init() {
	var err error
	db, err = sql.Open("mysql", mysql_source)
	if err != nil {
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
	db.SetMaxOpenConns(1000) //数据库连接池最大打开数
	err = db.Ping()          //测试正常连接
	if err != nil {
		fmt.Println("Failed to ping to mysql, err:" + err.Error())
		os.Exit(2)
	}
}

// DBConn : 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

package mysql

import (
	"database/sql"
	"fmt"
	"log"
)

//
func NewMysql(option Option) (*sql.DB, error) {
	//
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		option.Username,
		option.Password,
		option.Host,
		option.Port,
		option.DbName,
	)

	log.Println(connStr)

	db, err := sql.Open("mysql", connStr)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// 最大连接
	db.SetMaxOpenConns(10)
	// 最大闲连
	db.SetMaxIdleConns(1)
	log.Println(db)
	db.Exec("set names ?", option.Charset)

	return db, nil
}

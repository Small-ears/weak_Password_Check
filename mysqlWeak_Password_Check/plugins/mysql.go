package plugins

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"golang.com/golang.com/weak_Password_Check/vars"
)

func ScanMysql(s vars.Service) (result vars.ScanResult, err error) {
	result.Server = s
	dataSourceNmae := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", s.Username, s.Password, s.Target, s.Port, "mysql")
	db, err := sql.Open("mysql", dataSourceNmae)
	if err != nil {
		return result, err
	}

	err = db.Ping() //检查与数据库的连接是否有效的操作
	if err != nil {
		return result, err
	}

	result.Result = true

	defer func() {
		if db != nil {
			_ = db.Close()
		}
	}()

	return result, err
}

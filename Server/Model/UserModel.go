package model

import (
	"database/sql"
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
	"fmt"
	"math"
)

func UserCount(db *sql.DB) int {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM User").Scan(&count)
	return count
}

func UserAdd(db *sql.DB, data entity.UserEntity) (bool, string, int) {
	sqlCom := "INSERT INTO User(Account,Name,Password,Status,Level,Createtime) VALUES(?,?,?,?,?,?)"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	data.Password = lib.MD5(lib.MD5(lib.Int64ToString(lib.TimeStamp()) + data.Password + lib.Int64ToString(lib.TimeStamp())))
	data.Createtime = int(lib.TimeStamp())
	row, err := stmt.Exec(data.Account, data.Name, data.Password, data.Status, data.Level, data.Createtime)
	if err != nil {
		return false, err.Error(), 0
	}
	id, err := row.LastInsertId()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(id)
}

func UserUpdate(db *sql.DB, id string, data entity.UserEntity) (bool, string, int) {
	sqlCom := "UPDATE User SET Account=?,Name=?,Password=?,Status=?,Level=? WHERE ID=?"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	data.Password = lib.MD5(lib.MD5(data.Password + lib.Int64ToString(lib.TimeStamp()) + data.Password))
	res, err := stmt.Exec(data.Account, data.Name, data.Password, data.Status, data.Level, id)
	if err != nil {
		return false, err.Error(), 0
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(affect)
}

func UserData(db *sql.DB, id string) (bool, string, entity.UserEntity) {
	user := entity.UserEntity{}
	sqlCom := "SELECT * FROM User WHERE ID=" + id
	rows, err := db.Query(sqlCom)
	if err != nil {
		return false, err.Error(), user
	}
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Account, &user.Name, &user.Password, &user.Status, &user.Level, &user.Createtime)
		if err != nil {
			return false, err.Error(), user
		}
	}
	return true, "", user
}

func Users(db *sql.DB, order int, account string, name string, level int, status int) []entity.UserEntity {
	users := []entity.UserEntity{}
	condition_account := "1=1"
	condition_name := "1=1"
	condition_level := "1=1"
	condition_status := "1=1"
	if account != "" {
		condition_account = "Account LIKE '%" + account + "%'"
	}
	if name != "" {
		condition_name = "Name LIKE '%" + name + "%'"
	}
	if level > 0 {
		condition_level = "Level = " + lib.IntToString(level)
	}
	if status > 0 {
		condition_status = "Status = " + lib.IntToString(status)
	}
	orderBy := ""
	if order == -1 {
		orderBy = "DESC"
	} else {
		orderBy = "ASC"
	}
	sqlCom := "SELECT * FROM User WHERE " + condition_account + " AND " + condition_name + " AND " + condition_level + " AND " + condition_status +
		" ORDER BY ID " + orderBy
	rows, err := db.Query(sqlCom)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	for rows.Next() {
		user := entity.UserEntity{}
		err := rows.Scan(&user.ID, &user.Account, &user.Name, &user.Password, &user.Status, &user.Level, &user.Createtime)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		users = append(users, user)
	}
	if len(users) > 0 {
		for i := 0; i < len(users); i++ {
			users[i].Password = ""
		}
	}
	return users
}

func UserList(db *sql.DB, page int, pageSize int, order int, account string, name string, level int, status int) (int, int, int, []entity.UserEntity) {
	users := []entity.UserEntity{}
	condition_account := "1=1"
	condition_name := "1=1"
	condition_level := "1=1"
	condition_status := "1=1"
	if account != "" {
		condition_account = "Account LIKE '%" + account + "%'"
	}
	if name != "" {
		condition_name = "Name LIKE '%" + name + "%'"
	}
	if level > 0 {
		condition_level = "Level = " + lib.IntToString(level)
	}
	if status > 0 {
		condition_status = "Status = " + lib.IntToString(status)
	}
	if page <= 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	orderBy := ""
	if order == -1 {
		orderBy = "DESC"
	} else {
		orderBy = "ASC"
	}
	totalPage := math.Ceil(float64(UserCount(db)) / float64(pageSize))
	if totalPage > 0 && page > int(totalPage) {
		page = int(totalPage)
	}
	sqlCom := "SELECT * FROM User WHERE " + condition_account + " AND " + condition_name + " AND " + condition_level + " AND " + condition_status +
		" ORDER BY ID " + orderBy + " LIMIT " + lib.IntToString(pageSize) + " OFFSET " + lib.IntToString((page-1)*pageSize)
	rows, err := db.Query(sqlCom)
	if err != nil {
		fmt.Println(err.Error())
		return 0, 0, 0, nil
	}
	for rows.Next() {
		user := entity.UserEntity{}
		err := rows.Scan(&user.ID, &user.Account, &user.Name, &user.Password, &user.Status, &user.Level, &user.Createtime)
		if err != nil {
			fmt.Println(err.Error())
			return 0, 0, 0, nil
		}
		users = append(users, user)
	}

	if len(users) > 0 {
		for i := 0; i < len(users); i++ {
			users[i].Password = ""
		}
	}
	return page, pageSize, int(totalPage), users
}

func UserDel(db *sql.DB, id string) (bool, string, int) {
	if lib.StringContains(id, ",") {
		res, err := db.Exec("DELETE FROM User WHERE ID IN (" + id + ")")
		if err != nil {
			return false, err.Error(), 0
		}
		affect, err := res.RowsAffected()
		if err != nil {
			return false, err.Error(), 0
		}
		return true, "", int(affect)
	} else {
		sqlCom := "DELETE FROM User WHERE ID=?"
		stmt, err := db.Prepare(sqlCom)
		if err != nil {
			return false, err.Error(), 0
		}
		res, err := stmt.Exec(id)
		if err != nil {
			return false, err.Error(), 0
		}
		affect, err := res.RowsAffected()
		if err != nil {
			return false, err.Error(), 0
		}
		return true, "", int(affect)
	}
}

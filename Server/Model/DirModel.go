package model

import (
	"database/sql"
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
	"fmt"
	"math"
)

func DirCount(db *sql.DB) int {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM Dir").Scan(&count)
	return count
}

func DirAdd(db *sql.DB, data entity.DirEntity) (bool, string, int) {
	sqlCom := "INSERT INTO Dir(DirName,ParentID,UserID,Createtime) VALUES(?,?,?,?)"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	data.Createtime = int(lib.TimeStamp())
	row, err := stmt.Exec(data.DirName, data.ParentID, data.UserID, data.Createtime)
	if err != nil {
		return false, err.Error(), 0
	}
	id, err := row.LastInsertId()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(id)
}

func DirUpdate(db *sql.DB, id string, data entity.DirEntity) (bool, string, int) {
	sqlCom := "UPDATE Dir SET DirName=?,ParentID=?,UserID=?,Createtime=? WHERE ID=?"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	res, err := stmt.Exec(data.DirName, data.ParentID, data.UserID, id)
	if err != nil {
		return false, err.Error(), 0
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(affect)
}

func DirData(db *sql.DB, id string) (bool, string, entity.DirEntity) {
	data := entity.DirEntity{}
	sqlCom := "SELECT * FROM Dir WHERE ID=" + id
	rows, err := db.Query(sqlCom)
	if err != nil {
		return false, err.Error(), data
	}
	for rows.Next() {
		err := rows.Scan(&data.ID, &data.DirName, &data.ParentID, &data.UserID, &data.Createtime)
		if err != nil {
			return false, err.Error(), data
		}
	}
	return true, "", data
}

func Dirs(db *sql.DB, order int, dirName string, parentID int, userID int) []entity.DirEntity {
	datas := []entity.DirEntity{}
	condition_dirName := "1=1"
	condition_parentID := "1=1"
	condition_userID := "1=1"
	if dirName != "" {
		condition_parentID = "DirName LIKE '%" + dirName + "%'"
	}
	if parentID > 0 {
		condition_parentID = "ParentID = " + lib.IntToString(parentID)
	}
	if userID > 0 {
		condition_userID = "UserID = " + lib.IntToString(userID)
	}
	orderBy := ""
	if order == -1 {
		orderBy = "DESC"
	} else {
		orderBy = "ASC"
	}
	sqlCom := "SELECT * FROM Dir WHERE " + condition_dirName + " AND " + condition_parentID + " AND " + condition_userID + " ORDER BY ID " + orderBy
	rows, err := db.Query(sqlCom)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	for rows.Next() {
		data := entity.DirEntity{}
		err := rows.Scan(&data.ID, &data.DirName, &data.ParentID, &data.UserID, &data.Createtime)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		datas = append(datas, data)
	}
	return datas
}

func DirList(db *sql.DB, page int, pageSize int, order int, dirName string, parentID int, userID int) (int, int, int, []entity.DirEntity) {
	datas := []entity.DirEntity{}
	condition_dirName := "1=1"
	condition_parentID := "1=1"
	condition_userID := "1=1"
	if dirName != "" {
		condition_parentID = "DirName LIKE '%" + dirName + "%'"
	}
	if parentID > 0 {
		condition_parentID = "ParentID = " + lib.IntToString(parentID)
	}
	if userID > 0 {
		condition_userID = "UserID = " + lib.IntToString(userID)
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
	totalPage := math.Ceil(float64(DirCount(db)) / float64(pageSize))
	if totalPage > 0 && page > int(totalPage) {
		page = int(totalPage)
	}
	sqlCom := "SELECT * FROM Dir WHERE " + condition_dirName + " AND " + condition_parentID + " AND " + condition_userID +
		" ORDER BY ID " + orderBy + " LIMIT " + lib.IntToString(pageSize) + " OFFSET " + lib.IntToString((page-1)*pageSize)
	rows, err := db.Query(sqlCom)
	if err != nil {
		fmt.Println(err.Error())
		return 0, 0, 0, nil
	}
	for rows.Next() {
		data := entity.DirEntity{}
		err := rows.Scan(&data.ID, &data.DirName, &data.ParentID, &data.UserID, &data.Createtime)
		if err != nil {
			fmt.Println(err.Error())
			return 0, 0, 0, nil
		}
		datas = append(datas, data)
	}
	return page, pageSize, int(totalPage), datas
}

func DirDel(db *sql.DB, id string) (bool, string, int) {
	if lib.StringContains(id, ",") {
		res, err := db.Exec("DELETE FROM Dir WHERE ID IN (" + id + ")")
		if err != nil {
			return false, err.Error(), 0
		}
		affect, err := res.RowsAffected()
		if err != nil {
			return false, err.Error(), 0
		}
		return true, "", int(affect)
	} else {
		sqlCom := "DELETE FROM Dir WHERE ID=?"
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
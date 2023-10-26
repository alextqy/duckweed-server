package model

import (
	"database/sql"
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
	"fmt"
	"math"
)

func DirCount(db *sql.Tx, dirName string, parentID int, userID int) int {
	var count int
	condition_dirName := "1=1"
	condition_parentID := "1=1"
	condition_userID := "1=1"
	if dirName != "" {
		condition_dirName = "DirName LIKE '%" + dirName + "%'"
	}
	if parentID > 0 {
		condition_parentID = "ParentID = " + lib.IntToString(parentID)
	} else {
		condition_parentID = "ParentID = 0"
	}
	if userID > 0 {
		condition_userID = "UserID = " + lib.IntToString(userID)
	}
	sqlCom := "SELECT COUNT(*) FROM Dir WHERE " + condition_dirName + " AND " + condition_parentID + " AND " + condition_userID
	db.QueryRow(sqlCom).Scan(&count)
	return count
}

func DirAdd(db *sql.Tx, data entity.DirEntity) (bool, string, int) {
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

func DirUpdate(db *sql.Tx, id string, data entity.DirEntity) (bool, string, int) {
	sqlCom := "UPDATE Dir SET DirName=?,ParentID=?,UserID=? WHERE ID=?"
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

func DirData(db *sql.Tx, id string) (bool, string, entity.DirEntity) {
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

func DirDataSame(db *sql.Tx, userID string, parentID string, dirName string) (bool, string, entity.DirEntity) {
	data := entity.DirEntity{}
	sqlCom := "SELECT * FROM Dir WHERE UserID=" + userID + " AND parentID=" + parentID + " AND DirName='" + dirName + "'"
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

func Dirs(db *sql.Tx, order int, dirName string, parentID int, userID int) []entity.DirEntity {
	datas := []entity.DirEntity{}
	condition_dirName := "1=1"
	condition_parentID := "1=1"
	condition_userID := "1=1"
	if dirName != "" {
		condition_dirName = "DirName LIKE '%" + dirName + "%'"
	}
	if parentID > 0 {
		condition_parentID = "ParentID = " + lib.IntToString(parentID)
	} else {
		condition_parentID = "ParentID = 0"
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

func DirList(db *sql.Tx, page int, pageSize int, order int, dirName string, parentID int, userID int) (int, int, int, []entity.DirEntity) {
	datas := []entity.DirEntity{}
	condition_dirName := "1=1"
	condition_parentID := "1=1"
	condition_userID := "1=1"
	if dirName != "" {
		condition_dirName = "DirName LIKE '%" + dirName + "%'"
	}
	if parentID > 0 {
		condition_parentID = "ParentID = " + lib.IntToString(parentID)
	} else {
		condition_parentID = "ParentID = 0"
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
	totalPage := math.Ceil(float64(DirCount(db, dirName, parentID, userID)) / float64(pageSize))
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

func DirDel(db *sql.Tx, id string) (bool, string, int) {
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

func DirDelUser(db *sql.Tx, userID string) (bool, string, int) {
	sqlCom := "DELETE FROM Dir WHERE UserID=?"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	res, err := stmt.Exec(userID)
	if err != nil {
		return false, err.Error(), 0
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(affect)
}

func DirMove(db *sql.Tx, id, ids string) (bool, string, int) {
	if lib.StringContains(ids, ",") {
		res, err := db.Exec("UPDATE Dir SET ParentID=" + id + " WHERE ID IN (" + ids + ")")
		if err != nil {
			return false, err.Error(), 0
		}
		affect, err := res.RowsAffected()
		if err != nil {
			return false, err.Error(), 0
		}
		return true, "", int(affect)
	} else {
		sqlCom := "UPDATE Dir SET ParentID=? WHERE ID=?"
		stmt, err := db.Prepare(sqlCom)
		if err != nil {
			return false, err.Error(), 0
		}
		res, err := stmt.Exec(id, ids)
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

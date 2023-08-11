package model

import (
	"database/sql"
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
	"fmt"
	"math"
)

func FileCount(db *sql.DB) int {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM File").Scan(&count)
	return count
}

func FileAdd(db *sql.DB, data entity.FileEntity) (bool, string, int) {
	sqlCom := "INSERT INTO File(FileName,FileType,FileSize,StoragePath,MD5,UserID,DirID,Createtime) VALUES(?,?,?,?,?,?,?,?)"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	data.Createtime = int(lib.TimeStamp())
	row, err := stmt.Exec(data.FileName, data.FileType, data.FileSize, data.StoragePath, data.MD5, data.UserID, data.DirID, data.Createtime)
	if err != nil {
		return false, err.Error(), 0
	}
	id, err := row.LastInsertId()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(id)
}

func FileUpdate(db *sql.DB, id string, data entity.FileEntity) (bool, string, int) {
	sqlCom := "UPDATE File SET FileName=?,FileType=?,FileSize=?,StoragePath=?,MD5=?,UserID=?,DirID=? WHERE ID=?"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	res, err := stmt.Exec(data.FileName, data.FileType, data.FileSize, data.StoragePath, data.MD5, data.UserID, data.DirID, id)
	if err != nil {
		return false, err.Error(), 0
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(affect)
}

func FileData(db *sql.DB, id string) (bool, string, entity.FileEntity) {
	data := entity.FileEntity{}
	sqlCom := "SELECT * FROM File WHERE ID=" + id
	rows, err := db.Query(sqlCom)
	if err != nil {
		return false, err.Error(), data
	}
	for rows.Next() {
		err := rows.Scan(&data.ID, &data.FileName, &data.FileType, &data.FileSize, &data.StoragePath, &data.MD5, &data.UserID, &data.DirID, &data.Createtime)
		if err != nil {
			return false, err.Error(), data
		}
	}
	return true, "", data
}

func Files(db *sql.DB, order int, fileName string, userID int, dirID int) []entity.FileEntity {
	datas := []entity.FileEntity{}
	condition_fileName := "1=1"
	condition_userID := "1=1"
	condition_dirID := "1=1"
	if fileName != "" {
		condition_fileName = "FileName LIKE '%" + fileName + "%'"
	}
	if userID > 0 {
		condition_userID = "UserID = " + lib.IntToString(userID)
	}
	if dirID > 0 {
		condition_dirID = "DirID = " + lib.IntToString(dirID)
	}
	orderBy := ""
	if order == -1 {
		orderBy = "DESC"
	} else {
		orderBy = "ASC"
	}
	sqlCom := "SELECT * FROM File WHERE " + condition_fileName + " AND " + condition_userID + " AND " + condition_dirID + " ORDER BY ID " + orderBy
	rows, err := db.Query(sqlCom)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	for rows.Next() {
		data := entity.FileEntity{}
		err := rows.Scan(&data.ID, &data.FileName, &data.FileType, &data.FileSize, &data.StoragePath, &data.MD5, &data.UserID, &data.DirID, &data.Createtime)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		datas = append(datas, data)
	}
	return datas
}

func FileList(db *sql.DB, page int, pageSize int, order int, fileName string, userID int, dirID int) (int, int, int, []entity.FileEntity) {
	datas := []entity.FileEntity{}
	condition_fileName := "1=1"
	condition_userID := "1=1"
	condition_dirID := "1=1"
	if fileName != "" {
		condition_fileName = "FileName LIKE '%" + fileName + "%'"
	}
	if userID > 0 {
		condition_userID = "UserID = " + lib.IntToString(userID)
	}
	if dirID > 0 {
		condition_dirID = "DirID = " + lib.IntToString(dirID)
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
	totalPage := math.Ceil(float64(FileCount(db)) / float64(pageSize))
	if totalPage > 0 && page > int(totalPage) {
		page = int(totalPage)
	}
	sqlCom := "SELECT * FROM File WHERE " + condition_fileName + " AND " + condition_userID + " AND " + condition_dirID +
		" ORDER BY ID " + orderBy + " LIMIT " + lib.IntToString(pageSize) + " OFFSET " + lib.IntToString((page-1)*pageSize)
	rows, err := db.Query(sqlCom)
	if err != nil {
		fmt.Println(err.Error())
		return 0, 0, 0, nil
	}
	for rows.Next() {
		data := entity.FileEntity{}
		err := rows.Scan(&data.ID, &data.FileName, &data.FileType, &data.FileSize, &data.StoragePath, &data.MD5, &data.UserID, &data.DirID, &data.Createtime)
		if err != nil {
			fmt.Println(err.Error())
			return 0, 0, 0, nil
		}
		datas = append(datas, data)
	}
	return page, pageSize, int(totalPage), datas
}

func FileDel(db *sql.DB, id string) (bool, string, int) {
	if lib.StringContains(id, ",") {
		res, err := db.Exec("DELETE FROM File WHERE ID IN (" + id + ")")
		if err != nil {
			return false, err.Error(), 0
		}
		affect, err := res.RowsAffected()
		if err != nil {
			return false, err.Error(), 0
		}
		return true, "", int(affect)
	} else {
		sqlCom := "DELETE FROM File WHERE ID=?"
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
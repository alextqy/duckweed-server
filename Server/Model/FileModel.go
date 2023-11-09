package model

import (
	"database/sql"
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
	"fmt"
	"math"
)

func FileCount(db *sql.Tx, fileName string, userID int, dirID int) int {
	var count int
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
	} else {
		condition_dirID = "DirID = 0"
	}
	sqlCom := "SELECT COUNT(*) FROM File WHERE " + condition_fileName + " AND " + condition_userID + " AND " + condition_dirID
	db.QueryRow(sqlCom).Scan(&count)
	return count
}

func FileAdd(db *sql.Tx, data entity.FileEntity) (bool, string, int) {
	sqlCom := "INSERT INTO File(FileName,FileType,FileSize,StoragePath,MD5,UserID,DirID,Createtime,Status,OutreachID) VALUES(?,?,?,?,?,?,?,?,?,?)"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	data.Createtime = int(lib.TimeStamp())
	data.Status = 1
	row, err := stmt.Exec(data.FileName, data.FileType, data.FileSize, data.StoragePath, data.MD5, data.UserID, data.DirID, data.Createtime, data.Status, data.OutreachID)
	if err != nil {
		return false, err.Error(), 0
	}
	id, err := row.LastInsertId()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(id)
}

func FileUpdate(db *sql.Tx, id string, data entity.FileEntity) (bool, string, int) {
	sqlCom := "UPDATE File SET FileName=?,FileType=?,FileSize=?,StoragePath=?,MD5=?,UserID=?,DirID=?,Status=? WHERE ID=?"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	res, err := stmt.Exec(data.FileName, data.FileType, data.FileSize, data.StoragePath, data.MD5, data.UserID, data.DirID, data.Status, id)
	if err != nil {
		return false, err.Error(), 0
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(affect)
}

func FileData(db *sql.Tx, id string) (bool, string, entity.FileEntity) {
	data := entity.FileEntity{}
	sqlCom := "SELECT * FROM File WHERE ID=" + id
	rows, err := db.Query(sqlCom)
	if err != nil {
		return false, err.Error(), data
	}
	for rows.Next() {
		err := rows.Scan(&data.ID, &data.FileName, &data.FileType, &data.FileSize, &data.StoragePath, &data.MD5, &data.UserID, &data.DirID, &data.Createtime, &data.Status, &data.OutreachID)
		if err != nil {
			return false, err.Error(), data
		}
	}
	return true, "", data
}

func FileDataSame(db *sql.Tx, dirID string, fileName string, fileType string) (bool, string, entity.FileEntity) {
	data := entity.FileEntity{}
	sqlCom := "SELECT * FROM File WHERE DirID='" + dirID + "' AND FileName='" + fileName + "' AND FileType='" + fileType + "'"
	rows, err := db.Query(sqlCom)
	if err != nil {
		return false, err.Error(), data
	}
	for rows.Next() {
		err := rows.Scan(&data.ID, &data.FileName, &data.FileType, &data.FileSize, &data.StoragePath, &data.MD5, &data.UserID, &data.DirID, &data.Createtime, &data.Status, &data.OutreachID)
		if err != nil {
			return false, err.Error(), data
		}
	}
	return true, "", data
}

func FileDataMD5(db *sql.Tx, md5 string) (bool, string, entity.FileEntity) {
	data := entity.FileEntity{}
	sqlCom := "SELECT * FROM File WHERE MD5='" + md5 + "'"
	rows, err := db.Query(sqlCom)
	if err != nil {
		return false, err.Error(), data
	}
	for rows.Next() {
		err := rows.Scan(&data.ID, &data.FileName, &data.FileType, &data.FileSize, &data.StoragePath, &data.MD5, &data.UserID, &data.DirID, &data.Createtime, &data.Status, &data.OutreachID)
		if err != nil {
			return false, err.Error(), data
		}
	}
	return true, "", data
}

func Files(db *sql.Tx, order int, fileName string, userID int, dirID int, status int) []entity.FileEntity {
	datas := []entity.FileEntity{}
	condition_fileName := "1=1"
	condition_userID := "1=1"
	condition_dirID := "1=1"
	condition_status := "1=1"
	if fileName != "" {
		condition_fileName = "FileName LIKE '%" + fileName + "%'"
	}
	if userID > 0 {
		condition_userID = "UserID = " + lib.IntToString(userID)
	}
	if dirID > 0 {
		condition_dirID = "DirID = " + lib.IntToString(dirID)
	} else {
		condition_dirID = "DirID = 0"
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
	sqlCom := "SELECT * FROM File WHERE " + condition_fileName + " AND " + condition_userID + " AND " + condition_dirID + " AND " + condition_status + " ORDER BY ID " + orderBy
	rows, err := db.Query(sqlCom)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	for rows.Next() {
		data := entity.FileEntity{}
		err := rows.Scan(&data.ID, &data.FileName, &data.FileType, &data.FileSize, &data.StoragePath, &data.MD5, &data.UserID, &data.DirID, &data.Createtime, &data.Status, &data.OutreachID, &data.SourceAddress)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		datas = append(datas, data)
	}
	return datas
}

func FileList(db *sql.Tx, page int, pageSize int, order int, fileName string, userID int, dirID int, status int) (int, int, int, []entity.FileEntity) {
	datas := []entity.FileEntity{}
	condition_fileName := "1=1"
	condition_userID := "1=1"
	condition_dirID := "1=1"
	condition_status := "1=1"
	if fileName != "" {
		condition_fileName = "FileName LIKE '%" + fileName + "%'"
	}
	if userID > 0 {
		condition_userID = "UserID = " + lib.IntToString(userID)
	}
	if dirID > 0 {
		condition_dirID = "DirID = " + lib.IntToString(dirID)
	} else {
		condition_dirID = "DirID = 0"
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
	totalPage := math.Ceil(float64(FileCount(db, fileName, userID, dirID)) / float64(pageSize))
	if totalPage > 0 && page > int(totalPage) {
		page = int(totalPage)
	}
	sqlCom := "SELECT * FROM File WHERE " + condition_fileName + " AND " + condition_userID + " AND " + condition_dirID + " AND " + condition_status + " ORDER BY ID " + orderBy + " LIMIT " + lib.IntToString(pageSize) + " OFFSET " + lib.IntToString((page-1)*pageSize)
	rows, err := db.Query(sqlCom)
	if err != nil {
		fmt.Println(err.Error())
		return 0, 0, 0, nil
	}
	for rows.Next() {
		data := entity.FileEntity{}
		err := rows.Scan(&data.ID, &data.FileName, &data.FileType, &data.FileSize, &data.StoragePath, &data.MD5, &data.UserID, &data.DirID, &data.Createtime, &data.Status, &data.OutreachID, &data.SourceAddress)
		if err != nil {
			fmt.Println(err.Error())
			return 0, 0, 0, nil
		}
		datas = append(datas, data)
	}
	return page, pageSize, int(totalPage), datas
}

func FileDel(db *sql.Tx, id string) (bool, string, int) {
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

func FileDelDir(db *sql.Tx, dirID string) (bool, string, int) {
	sqlCom := "DELETE FROM File WHERE DirID=?"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	res, err := stmt.Exec(dirID)
	if err != nil {
		return false, err.Error(), 0
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(affect)
}

func FileDelUser(db *sql.Tx, userID string) (bool, string, int) {
	sqlCom := "DELETE FROM File WHERE UserID=?"
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

func FileMove(db *sql.Tx, dirID, ids string) (bool, string, int) {
	if lib.StringContains(ids, ",") {
		res, err := db.Exec("UPDATE File SET DirID=" + dirID + " WHERE ID IN (" + ids + ")")
		if err != nil {
			return false, err.Error(), 0
		}
		affect, err := res.RowsAffected()
		if err != nil {
			return false, err.Error(), 0
		}
		return true, "", int(affect)
	} else {
		sqlCom := "UPDATE File SET DirID=? WHERE ID=?"
		stmt, err := db.Prepare(sqlCom)
		if err != nil {
			return false, err.Error(), 0
		}
		res, err := stmt.Exec(dirID, ids)
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

package model

import (
	"database/sql"
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
	"fmt"
	"math"
)

func AnnouncementCount(db *sql.Tx) int {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM Announcement").Scan(&count)
	return count
}

func AnnouncementAdd(db *sql.Tx, data entity.AnnouncementEntity) (bool, string, int) {
	sqlCom := "INSERT INTO Announcement(Content,Createtime) VALUES(?,?)"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	data.Createtime = int(lib.TimeStamp())
	row, err := stmt.Exec(data.Content, data.Createtime)
	if err != nil {
		return false, err.Error(), 0
	}
	id, err := row.LastInsertId()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(id)
}

func AnnouncementUpdate(db *sql.Tx, id string, data entity.AnnouncementEntity) (bool, string, int) {
	sqlCom := "UPDATE Announcement SET Content=? WHERE ID=?"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	res, err := stmt.Exec(data.Content)
	if err != nil {
		return false, err.Error(), 0
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(affect)
}

func AnnouncementData(db *sql.Tx, id string) (bool, string, entity.AnnouncementEntity) {
	data := entity.AnnouncementEntity{}
	sqlCom := "SELECT * FROM Announcement WHERE ID=" + id
	rows, err := db.Query(sqlCom)
	if err != nil {
		return false, err.Error(), data
	}
	for rows.Next() {
		err := rows.Scan(&data.ID, &data.Content, &data.Createtime)
		if err != nil {
			return false, err.Error(), data
		}
	}
	return true, "", data
}

func Announcements(db *sql.Tx, order int) []entity.AnnouncementEntity {
	datas := []entity.AnnouncementEntity{}
	orderBy := ""
	if order == -1 {
		orderBy = "DESC"
	} else {
		orderBy = "ASC"
	}
	sqlCom := "SELECT * FROM Announcement ORDER BY ID " + orderBy
	rows, err := db.Query(sqlCom)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	for rows.Next() {
		data := entity.AnnouncementEntity{}
		err := rows.Scan(&data.ID, &data.Content, &data.Createtime)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		datas = append(datas, data)
	}
	return datas
}

func AnnouncementList(db *sql.Tx, page int, pageSize int, order int, content string) (int, int, int, []entity.AnnouncementEntity) {
	datas := []entity.AnnouncementEntity{}
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
	totalPage := math.Ceil(float64(AnnouncementCount(db)) / float64(pageSize))
	if totalPage > 0 && page > int(totalPage) {
		page = int(totalPage)
	}
	sqlCom := "SELECT * FROM Announcement ORDER BY ID " + orderBy + " LIMIT " + lib.IntToString(pageSize) + " OFFSET " + lib.IntToString((page-1)*pageSize)
	rows, err := db.Query(sqlCom)
	if err != nil {
		fmt.Println(err.Error())
		return 0, 0, 0, nil
	}
	for rows.Next() {
		data := entity.AnnouncementEntity{}
		err := rows.Scan(&data.ID, &data.Content, &data.Createtime)
		if err != nil {
			fmt.Println(err.Error())
			return 0, 0, 0, nil
		}
		datas = append(datas, data)
	}
	return page, pageSize, int(totalPage), datas
}

func AnnouncementDel(db *sql.Tx, id string) (bool, string, int) {
	if lib.StringContains(id, ",") {
		res, err := db.Exec("DELETE FROM Announcement WHERE ID IN (" + id + ")")
		if err != nil {
			return false, err.Error(), 0
		}
		affect, err := res.RowsAffected()
		if err != nil {
			return false, err.Error(), 0
		}
		return true, "", int(affect)
	} else {
		sqlCom := "DELETE FROM Announcement WHERE ID=?"
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

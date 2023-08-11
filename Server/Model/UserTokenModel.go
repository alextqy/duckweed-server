package model

import (
	"database/sql"
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
)

func UserTokenAdd(db *sql.DB, data entity.UserTokenEntity) (bool, string, int) {
	sqlCom := "INSERT INTO UserToken(UserID,Token,Createtime) VALUES(?,?,?)"
	stmt, err := db.Prepare(sqlCom)
	if err != nil {
		return false, err.Error(), 0
	}
	data.Createtime = int(lib.TimeStamp())
	row, err := stmt.Exec(data.UserID, data.Token, data.Createtime)
	if err != nil {
		return false, err.Error(), 0
	}
	id, err := row.LastInsertId()
	if err != nil {
		return false, err.Error(), 0
	}
	return true, "", int(id)
}

func UserTokenData(db *sql.DB, token string) (bool, string, entity.UserTokenEntity) {
	data := entity.UserTokenEntity{}
	sqlCom := "SELECT * FROM UserToken WHERE Token=" + token
	rows, err := db.Query(sqlCom)
	if err != nil {
		return false, err.Error(), data
	}
	for rows.Next() {
		err := rows.Scan(&data.ID, &data.UserID, &data.Token, &data.Createtime)
		if err != nil {
			return false, err.Error(), data
		}
	}
	return true, "", data
}

func UserTokenDel(db *sql.DB, userID string) (bool, string, int) {
	sqlCom := "DELETE FROM UserToken WHERE UserID=?"
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

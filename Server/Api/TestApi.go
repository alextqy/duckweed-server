package api

import (
	entity "duckweed-server/Server/Entity"
	lib "duckweed-server/Server/Lib"
	"log"
	"math"
	"net/http"
	"strings"

	"xorm.io/xorm"
)

func XDB() (bool, string, *xorm.Session, *xorm.Engine) {
	engine, err := xorm.NewEngine("sqlite3", "../Dao.db")
	if err != nil {
		engine.Close()
		log.Fatal(err.Error())
		return false, err.Error(), nil, nil
	} else {
		session := engine.NewSession()
		defer session.Close()
		return true, "", session, engine
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(Post(r, "text"))
	// HttpWrite(w, processmodule.Test(Post(r, "text")))

	_, _, db, _ := XDB()

	file := entity.File{}
	file.FileName = lib.Int64ToString(lib.TimeStamp())
	file.FileType = "txt"
	file.FileSize = "1024"
	file.StoragePath = ""
	file.MD5 = "asdfasdfasdfasdf12324234324"
	file.UserID = 1
	file.DirID = 0
	file.Createtime = lib.TimeStamp()
	file.Status = 1
	file.OutreachID = ""
	file.SourceAddress = ""

	db.Begin()
	// fmt.Println(file_count(db, "", 0, 0))
	// _, e := file_add(db, file)
	// _, e := file_update(db, 2, file)
	// d, _ := file_data(db, 1)
	// s := file_data_same(db, 1, "", "")
	// d, _ := file_data_md5(db, "asdfasdfasdfasdf12324234324")
	// s := files(db, -1, "", 0, 1, 0)
	// _, _, _, s := file_list(db, 2, 2, -1, "", 0, 0, 0)
	// _, e := file_move(db, "1", "3,4")
	// _, e := file_del(db, "4,3")
	// _, e := file_del_dir(db, 1)
	// _, e := file_del_user(session, 2)
	// fmt.Println(e)
	db.Commit()
}

func file_count(db *xorm.Session, fileName string, userID int64, dirID int64) (int64, error) {
	file := entity.File{}
	engine := db.Where("1=1")
	if fileName != "" {
		engine = engine.And("DirName LIKE ?", "%"+fileName+"%")
	}
	if userID > 0 {
		engine = engine.And("UserID = ?", userID)
	}
	if dirID > 0 {
		engine = engine.And("DirID = ?", dirID)
	} else {
		engine = engine.And("DirID = ?", 0)
	}
	total, err := engine.Count(&file)
	return total, err
}

func file_add(db *xorm.Session, data entity.File) (int64, error) {
	affected, err := db.Insert(&data)
	return affected, err
}

func file_update(db *xorm.Session, id int64, data entity.File) (int64, error) {
	affected, err := db.ID(id).Update(&data)
	return affected, err
}

func file_data(db *xorm.Session, id int64) (entity.File, error) {
	file := entity.File{}
	_, err := db.ID(id).Get(&file)
	return file, err
}

func file_data_same(db *xorm.Session, dirID int64, fileName string, fileType string) []entity.File {
	dirs := []entity.File{}
	engine := db.Where("1=1")
	if dirID > 0 {
		engine = engine.And("DirID = ?", dirID)
	}
	if fileName != "" {
		engine = engine.And("FileName = ?", fileName)
	}
	if fileType != "" {
		engine = engine.And("FileType = ?", fileType)
	}
	engine.Desc("ID").Find(&dirs)
	return dirs
}

func file_data_md5(db *xorm.Session, md5 string) (entity.File, error) {
	file := entity.File{}
	_, err := db.Where("MD5 = ?", md5).Get(&file)
	return file, err
}

func files(db *xorm.Session, order int, fileName string, userID int, dirID int, status int) []entity.File {
	files := []entity.File{}
	engine := db.Where("1=1")
	if fileName != "" {
		engine = engine.And("FileName LIKE ?", "%"+fileName+"%")
	}
	if userID > 0 {
		engine = engine.And("UserID = ?", userID)
	}
	if dirID > 0 {
		engine = engine.And("DirID = ?", dirID)
	}
	if dirID == 0 {
		engine = engine.And("DirID = ?", 0)
	}
	if status > 0 {
		engine = engine.And("Status = ?", status)
	}
	orderBy := ""
	if order == -1 {
		orderBy = "DESC"
	} else {
		orderBy = "ASC"
	}
	engine.OrderBy("ID " + orderBy).Find(&files)
	return files
}

func file_list(db *xorm.Session, page int, pageSize int, order int, fileName string, userID int64, dirID int64, status int) (int, int, int, []entity.File) {
	files := []entity.File{}
	engine := db.Where("1=1")
	if fileName != "" {
		engine = engine.And("FileName LIKE ?", "%"+fileName+"%")
	}
	if userID > 0 {
		engine = engine.And("UserID = ?", userID)
	}
	if dirID > 0 {
		engine = engine.And("DirID = ?", dirID)
	}
	if dirID == 0 {
		engine = engine.And("DirID = ?", 0)
	}
	if status > 0 {
		engine = engine.And("Status = ?", status)
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
	fileCountInt, _ := file_count(db, fileName, userID, dirID)
	totalPage := math.Ceil(float64(fileCountInt) / float64(pageSize))
	if totalPage > 0 && page > int(totalPage) {
		page = int(totalPage)
	}
	engine.OrderBy("ID "+orderBy).Limit(pageSize, (page-1)*pageSize).Find(&files)
	return page, pageSize, int(totalPage), files
}

func file_move(db *xorm.Session, dirID, ids string) ([]map[string]string, error) {
	sql := ""
	if lib.StringContains(ids, ",") {
		sql = "UPDATE File SET DirID=" + dirID + " WHERE ID IN(" + ids + ")"
	} else {
		sql = "UPDATE File SET DirID=" + dirID + " WHERE ID=" + ids
	}
	res, err := db.QueryString(sql)
	return res, err
}

func file_del(db *xorm.Session, id string) (int, error) {
	file := entity.File{}
	if lib.StringContains(id, ",") {
		ids := strings.Split(id, ",")
		intArr := []int{}
		for i := 0; i < len(ids); i++ {
			_, _, n := lib.StringToInt(ids[i])
			intArr = append(intArr, n)
		}
		affected, err := db.In("ID", intArr).Delete(file)
		return int(affected), err
	} else {
		affected, err := db.ID(id).Delete(file)
		return int(affected), err
	}
}

func file_del_dir(db *xorm.Session, dirID int64) (int, error) {
	file := entity.File{}
	affected, err := db.Where("DirID = ?", dirID).Delete(file)
	return int(affected), err
}

func file_del_user(db *xorm.Session, userID int64) (int, error) {
	file := entity.File{}
	affected, err := db.Where("UserID = ?", userID).Delete(file)
	return int(affected), err
}

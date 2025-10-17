package api

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func FixVersionDirty(path string, version int) string {
	err := fixVersionDirty(path, version)
	if err != nil {
		return MakeJsonErrorResult2(fixVersionDirtyErr, err.Error(), nil)
	}
	return MakeJsonErrorResult2(SUCCESS_2, SUCCESS_2.Error(), nil)
}

func fixVersionDirty(path string, version int) (err error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return AppendErrorInfo(err, "sql.Open")
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Println(err)
		}
	}(db)

	sqlStmt := "update main.schema_migrations set version = ?, dirty = ?"

	_, err = db.Exec(sqlStmt, version, 0)
	if err != nil {
		return AppendErrorInfo(err, "db.Exec")
	}
	return nil
}

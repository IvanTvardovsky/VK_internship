package utils

import (
	"database/sql"
)

func IsServiceNameTaken(id int, name string, db *sql.DB) (bool, error) {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM data WHERE user_id = $1 AND service_name = $2",
		id, name).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func IsSuitableForRestrictions(lenServiceName, lenLogin, lenPassword int) bool {
	if lenServiceName > 255 || lenLogin > 255 || lenPassword > 255 {
		return false
	}
	return true
}

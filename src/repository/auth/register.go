package auth

import (
	"database/sql"
	"fmt"
)

func (o AuthenticationTable) RegisterUser(phoneNumber string, password string) (int64, *sql.DB, error) {
	var err error
	var res sql.Result
	var prepare *sql.Stmt
	db, err := o.ConnectDB()
	if err != nil {
		fmt.Println(err)
		return 0, nil, fmt.Errorf("%s", err)
	}
	defer db.Close()
	if err != nil {
		return 0, nil, fmt.Errorf("%s", err)
	}

	queryInsert := "INSERT INTO tbl_user (phoneNumber,password) values (?,?)"
	prepare, err = db.Prepare(queryInsert)
	if err != nil {
		return 0, db, fmt.Errorf("failed to insert tbl_user SQL : %v", err)
	}
	res, err = prepare.Exec(phoneNumber, password)

	if err != nil {
		return 0, db, fmt.Errorf("failed to insert error_general on tbl_user SQL : %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 0, db, fmt.Errorf("failed to populate status inserted : %v", err)
	}
	return count, db, err
}

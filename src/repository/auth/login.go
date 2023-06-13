package auth

import (
	"database/sql"
	"fmt"

	"github.com/ahfrd/grpc/auth-client/src/model"
)

func (o AuthenticationTable) SelectUsernameByEmail(email string) (model.UserStruct, *sql.DB, error) {
	var modelUser model.UserStruct
	var phoneNumber NullString
	db, err := o.ConnectDB()
	if err != nil {
		return modelUser, db, err
	}
	var query string = fmt.Sprintf(`select phoneNumber from tbl_user where email = "%s"`, email)
	db.QueryRow(query).Scan(
		&phoneNumber,
	)
	fmt.Println(query)
	modelUser.PhoneNumber = phoneNumber.String
	defer db.Close()
	if err != nil && err != sql.ErrNoRows {
		return modelUser, db, fmt.Errorf("failed Select SQL for tbl_user : %v", err)
	}

	return modelUser, db, nil
}
func (o AuthenticationTable) ValidateDataByPhoneNumber(phoneNumberReq string) (model.UserStruct, *sql.DB, error) {
	var modelUser model.UserStruct
	var id NullString
	var phoneNumber NullString
	db, err := o.ConnectDB()
	if err != nil {
		return modelUser, db, err
	}
	var query string = fmt.Sprintf(`select id,phoneNumber from tbl_user where phoneNumber = "%s"`, phoneNumberReq)
	db.QueryRow(query).Scan(
		&id,
		&phoneNumber,
	)
	fmt.Println(query)
	modelUser.Id = id.String
	modelUser.PhoneNumber = phoneNumber.String
	defer db.Close()
	if err != nil && err != sql.ErrNoRows {
		return modelUser, db, fmt.Errorf("failed Select SQL for tbl_user : %v", err)
	}

	return modelUser, db, nil
}
func (o AuthenticationTable) SelectByPhoneNumber(phoneNumberReq string) (model.UserStruct, *sql.DB, error) {
	var modelUser model.UserStruct
	var id NullString
	var phoneNumber NullString
	var password NullString
	var loginRetry NullInt
	var nextLogin NullString
	var lastLogin NullString
	var session NullString
	var status NullInt
	db, err := o.ConnectDB()
	if err != nil {
		return modelUser, db, err
	}
	var query string = fmt.Sprintf(`select id,phoneNumber,password,login_retry,next_login_date,last_login,session_id,status from tbl_user where phoneNumber = "%s"`, phoneNumberReq)
	db.QueryRow(query).Scan(
		&id,
		&phoneNumber,
		&password,
		&loginRetry,
		&nextLogin,
		&lastLogin,
		&session,
		&status,
	)
	fmt.Println(query)
	modelUser.Id = id.String
	modelUser.PhoneNumber = phoneNumber.String
	modelUser.Password = password.String
	modelUser.LoginRetry = int(loginRetry.Int64)
	modelUser.NextLoginDate = nextLogin.String
	modelUser.LastLogin = lastLogin.String
	modelUser.SessionId = session.String
	modelUser.Status = int(status.Int64)

	defer db.Close()
	if err != nil && err != sql.ErrNoRows {
		return modelUser, db, fmt.Errorf("failed Select SQL for tbl_user : %v", err)
	}

	return modelUser, db, nil
}

func (o AuthenticationTable) UpdateRetryLogin(phoneNumber string, count int) (int64, *sql.DB, error) {
	db, err := o.ConnectDB()
	if err != nil {
		return 0, nil, fmt.Errorf("%s", err)
	}

	var queryUpdate string
	queryUpdate = `UPDATE tbl_user set login_retry = ? where phoneNumber = ?`
	prepare, err := db.Prepare(queryUpdate)
	res, err := prepare.Exec(count, phoneNumber)
	defer db.Close()
	if err != nil {
		return 0, db, fmt.Errorf("failed to tbl_user SQL : %v", err)
	}

	if err != nil {
		return 0, db, fmt.Errorf("failed to tbl_user SQL : %v", err)
	}

	counter, err := res.RowsAffected()
	if err != nil {
		return 0, db, fmt.Errorf("failed to populate status updated : %v", err)
	}
	return counter, db, err
}
func (o AuthenticationTable) UpdateRetryNextLogin(phoneNumber string, next_login_date string, count int) (int64, *sql.DB, error) {
	db, err := o.ConnectDB()
	if err != nil {
		return 0, nil, fmt.Errorf("%s", err)
	}

	var queryUpdate string
	queryUpdate = `UPDATE tbl_user set login_retry = ?,next_login_date where phoneNumber = ?`
	prepare, err := db.Prepare(queryUpdate)
	res, err := prepare.Exec(count, next_login_date, phoneNumber)
	defer db.Close()
	if err != nil {
		return 0, db, fmt.Errorf("failed to tbl_user SQL : %v", err)
	}

	if err != nil {
		return 0, db, fmt.Errorf("failed to tbl_user SQL : %v", err)
	}

	counter, err := res.RowsAffected()
	if err != nil {
		return 0, db, fmt.Errorf("failed to populate status updated : %v", err)
	}
	return counter, db, err
}

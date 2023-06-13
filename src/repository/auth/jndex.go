package auth

import (
	"database/sql"

	"github.com/ahfrd/grpc/auth-client/src/db"
)

type NullString struct {
	sql.NullString
}

type NullInt struct {
	sql.NullInt64
}

type AuthenticationTable struct {
	db.Database
}

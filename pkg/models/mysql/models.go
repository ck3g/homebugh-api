package mysql

import "database/sql"

type Models struct {
	Users UserModel
}

func NewModels(db *sql.DB) interface{} {
	return Models{
		Users: UserModel{DB: db},
	}
}

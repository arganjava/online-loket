package models

import "database/sql"

type Location struct {
	ID       string
	Country  sql.NullString
	CityName sql.NullString
	Village  sql.NullString
	Address  sql.NullString
}

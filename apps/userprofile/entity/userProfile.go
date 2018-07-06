package entity

import (
	"time"
)

type UserProfile struct {
	Id int64 `xorm: "pk autoincr notnull" json: "id"`
	Password string `xorm: "varchar(255) notnull" json: "passwd"`
	Email string `xorm: "varchar(255) unique notnull" json: "email"`
	Phone int64 `xorm: "bigint" json: "phone"`
	Role string `xorm: "varchar(255) notnull" json: "role"`
	Nickname string `xorm: "varchar(255)" json: "nickname"`
	Birthday time.Time `xorm: "date" json: "birthday"`
	Address string `xorm: "varchar(255)" json: "address"`
	CreatedAt time.Time `xorm: "created" json: "created_at"`
}
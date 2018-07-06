package entity

import (
	"time"
)

type UserProfile struct {
	Id int64 `xorm: "pk autoincr notnull" json: "id"`
	Password string `xorm: "varchar(40) notnull" json: "passwd"`
	Email string `xorm: "varchar(40) unique notnull" json: "email"`
	Phone int64 `xorm: "bigint" json: "phone"`
	Nickname string `xorm: "varchar(40)" json: "nickname"`
	Birthday time.Time `xorm: "date" json: "birthday"`
	Address string `xorm: "varchar(40)" json: "address"`
	CreatedAt time.Time `xorm: "created" json: "created_at"`
}
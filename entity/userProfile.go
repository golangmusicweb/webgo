package entity

import (
	"time"
)

type UserProfile struct {
	Id int64 `xorm: "pk autoincr notnull"`
	Username string `xorm: "varchar(40) notnull unique"`
	Password string `xorm: "varchar(40) notnull"`
	Email string `xorm: "varchar(40) unique notnull"`
	Nickname string `xorm: "varchar(40)"`
	Birthday time.Time `xorm: "date"`
	Address string `xorm: "varchar(40)"`
	CreatedAt time.Time `xorm: "created"`
}
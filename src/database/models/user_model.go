package models

type User struct {
	ID       string `xorm:"pk varchar(36) unique notnull"`
	Username string `xorm:"unique"`
	Password string `xorm:"varchar(255)" private:"true"`
	Role     string `xorm:"varchar(100)"`
}

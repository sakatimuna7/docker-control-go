package models

type User struct {
	ID       int64  `xorm:"pk autoincr" json:"id"`
	Username string `xorm:"unique"`
	Password string `xorm:"varchar(255)" private:"true"`
	Role     string `xorm:"varchar(100)"`
}

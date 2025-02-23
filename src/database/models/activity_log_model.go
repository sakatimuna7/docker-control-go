package models

import (
	"time"
)

type ActivityLog struct {
	ID        int64     `xorm:"pk autoincr" json:"id"`
	UserID    int64     `xorm:"index" json:"user_id"` // Ubah varchar jadi index integer
	Action    string    `xorm:"varchar(255)" json:"action"`
	Timestamp time.Time `xorm:"created" json:"timestamp"` // Gunakan "created" agar otomatis diisi
}

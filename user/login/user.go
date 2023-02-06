package register

type User struct {
	Id       int64  `gorm:"primary_key"`
	Username string `gorm:"primary_key"`
	Password string `gorm:"size:256;notnull"`
	//UserHost用于获取用户详细信息
	UserHost int64
}

type

func (u *)

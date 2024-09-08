package user

import (
	"my_ecommerce_system/pkg/db"
	"time"
)

func addNewUser(newUser *User){
	sql := `INSERT INTO sys_user (id, name, password, created_at) VALUES (?, ?, ?, ?)`
	db.Execute(sql, db.GenId(), newUser.Name, newUser.Password, time.Now())
}

package user

import (
	"database/sql"
	"fmt"
	"my_ecommerce_system/pkg/constant"
	"my_ecommerce_system/pkg/db"

	"time"
)

type UserRepository interface {
	ExistsUserName(userName string) (bool, error)
	AddNewUser(user *User) error
}

type StdUserRepository struct {}

func (r *StdUserRepository) AddNewUser(newUser *User) error{
	sql := `INSERT INTO sys_user (id, name, password, created_at) VALUES (?, ?, ?, ?)`
	return db.Execute(sql, db.GenId(), newUser.Name, newUser.Password, time.Now())
}

func getUserByName(name string) (*User,error) {
	sqlStr := `SELECT id, name, password, created_at FROM sys_user WHERE name = ?`

	var users []User

	// 貌似当前的数据库驱动无法自动将 []uint8 转换为 time.Time
	// 故手动进行转换一下
	// 这是这样一来就破坏了代码统一性
	var createdAt sql.NullString

	err := db.ExecuteQuery(sqlStr, func(rows *sql.Rows) error {
		for rows.Next(){
			var user User
			err := rows.Scan(&user.Id, &user.Name, &user.Password, &createdAt)
			if err != nil{
				return err
			}

			if createdAt.Valid {
				// 很奇怪的时间日期格式化模板：constant.DATE_TIME_FORMAT
				user.Created_at, err = time.Parse(constant.DATE_TIME_FORMAT, createdAt.String)
				if err != nil{
					return err
				}
			}

			users = append(users, user)
		}
		return nil
	}, name)

	if err != nil{
		return nil, err
	}

	if len(users) == 0{
		return nil, nil
	}else if len(users) > 1 {
		return nil, fmt.Errorf("用户名[%s]存在多于一个的用户！", name)
	}

	return &users[0], nil
}

func (r *StdUserRepository) ExistsUserName(name string) (bool, error){
	return db.Exists("SELECT 1 FROM sys_user WHERE name = ? LIMIT 1", name)
}

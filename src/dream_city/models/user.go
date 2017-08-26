package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"sync"
	"fmt"
	"strings"
)

var (
	mUser = &MUser{users: make(map[string]*User, 0)}
)

type User struct {
	TableCommon

	Rands          string    `show:"-"`
    UserName       string    ``
	Email          string    `orm:"column(email);size(64);"`
	Password       string    `orm:"column(password_digest);" show:"-"`
	Role           string    `orm:"column(role);size(32);"`
	City           string    `orm:"column(city);size(32);"`

	CreatorId      int64     `orm:"column(creater_id);"`
	AppType        string    `orm:"column(app_type);size(255);"`
	IsActive       bool
	IsForbid       bool
}

func init() {
	orm.RegisterModel(new(User))
}

func (t *User) TableName() string {
	return "users"
}

func (t *User) LoadByName(userName string) bool {
	var err error = nil

	if strings.IndexRune(userName, '@') == -1 {
		err = FindBy("Email", userName, t)
	} else {
		err = FindBy("Email", userName, t)
	}

	return err == nil
}

func (t *User) VerifyPassword(password string) bool {
	return t.Password == password
}

func (t *User) VerifyLogin(userName, password string) bool {
	return t.LoadByName(userName) && t.VerifyPassword(password)
}

func (t *User) IsAdmin() bool {
	return true
}

func (t *User) Link() string {
	return fmt.Sprintf("/users/%d", t.Id)
}

type MUser struct{
	users map[string]*User
	sync.Mutex
}

func (m *MUser) GetUser(id interface{}) *User {
	key := ""
	if _, ok := id.(string); !ok {
		key = fmt.Sprintf("%d", id)
	} else {
		key = fmt.Sprintf("%s", id)
	}

	if key == "0" || key == "" {
		return nil
	}

	m.Lock()
	u, ok := m.users[key]
	m.Unlock()
	if ok {
		return u
	}

	i, _ := strconv.ParseInt(key, 10, 64)

	u = &User{}
	err := FindById(u, i)
	if err != nil {
		return nil
	}

	m.Lock()
	m.users[key] = u
	m.Unlock()

	return u
}

func (m *MUser)GetUserName(id interface{}) string {

	u := mUser.GetUser(id)
	if u == nil {
		return ""
	}

	return u.Email
}


func (u *User) Update(changes... string) error {
	return nil
}


func (u *User) Insert(changes... string) error {
	return nil
}
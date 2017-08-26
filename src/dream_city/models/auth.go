package models

import (
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego/context"
	"strings"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/session"
	"dream_city/library/tools"
	"dream_city/library/types"
	"dream_city/library/conf"
)



// CanRegistered checks if the username or e-mail is available.
func CanRegistered(userName string, email string) (bool, bool, error) {
	cond := orm.NewCondition()
	cond = cond.Or("UserName", userName).Or("Email", email)

	var maps []orm.Params
	o := orm.NewOrm()
	n, err := o.QueryTable("user").SetCond(cond).Values(&maps, "UserName", "Email")
	if err != nil {
		return false, false, err
	}

	e1 := true
	e2 := true

	if n > 0 {
		for _, m := range maps {
			if e1 && orm.ToStr(m["UserName"]) == userName {
				e1 = false
			}
			if e2 && orm.ToStr(m["Email"]) == email {
				e2 = false
			}
		}
	}

	return e1, e2, nil
}

// check if exist user by username or email


// register create user
func RegisterUser(user *User, username, email, password string) error {
	// use random salt encode password
	salt := types.NewGuid().String()
	pwd := tools.EncodePassword(password, salt)

	user.UserName = strings.ToLower(username)
	user.Email = strings.ToLower(email)

	// save salt and encode password, use $ as split char
	user.Password = fmt.Sprintf("%s$%s", salt, pwd)

	return user.Insert()
}

// set a new password to user
func SaveNewPassword(user *User, password string) error {
	salt := types.NewGuid().String()
	user.Password = fmt.Sprintf("%s$%s", salt, tools.EncodePassword(password, salt))
	return user.Update("Password", "Rands", "Updated")
}

// get login redirect url from cookie
func GetLoginRedirect(ctx *context.Context) string {
	loginRedirect := strings.TrimSpace(ctx.GetCookie("login_to"))
	if tools.IsMatchHost(loginRedirect) == false {
		loginRedirect = "/"
	} else {
		ctx.SetCookie("login_to", "", -1, "/")
	}
	return loginRedirect
}

// login user
func LoginUser(user *User, ctx *context.Context, remember bool) {

}

func WriteRememberCookie(user *User, ctx *context.Context) {
	secret := tools.EncodeMd5(user.Rands + user.Password)
	days := 86400 * conf.LoginRememberDays
	ctx.SetCookie(conf.CookieUserName, user.UserName, days)
	ctx.SetSecureCookie(secret, conf.CookieRememberName, user.UserName, days)
}

func DeleteRememberCookie(ctx *context.Context) {
	ctx.SetCookie(conf.CookieUserName, "", -1)
	ctx.SetCookie(conf.CookieRememberName, "", -1)
}

func LoginUserFromRememberCookie(user *User, ctx *context.Context) (success bool) {
	userName := ctx.GetCookie(conf.CookieUserName)
	if len(userName) == 0 {
		return false
	}

	defer func() {
		if !success {
			DeleteRememberCookie(ctx)
		}
	}()


	if err := FindBy("UserName", userName, user); err != nil {
		return false
	}

	secret := tools.EncodeMd5(user.Rands + user.Password)
	value, _ := ctx.GetSecureCookie(secret, conf.CookieRememberName)
	if value != userName {
		return false
	}

	LoginUser(user, ctx, true)

	return true
}

// logout user

func GetUserIdFromSession(sess session.Store) int64 {
	if id, ok := sess.Get("auth_user_id").(int64); ok && id > 0 {
		return id
	}
	return 0
}

// get user if key exist in session
func GetUserFromSession(user *User, sess session.Store) bool {
	id := GetUserIdFromSession(sess)
	if id > 0 {
		FindById(user, id)
	}

	return false
}


func VerifyUser(user *User, username, password string) (success bool) {

	if !user.LoadByName(username) {
		return
	}

	if VerifyPassword(password, user.Password) {
		// success
		success = true

		// re-save discuz password
		if len(user.Password) == 39 {
			if err := SaveNewPassword(user, password); err != nil {
				beego.Error("SaveNewPassword err: ", err.Error())
			}
		}
	}
	return
}

// compare raw password and encoded password
func VerifyPassword(rawPwd, encodedPwd string) bool {
	beego.Info("rawPwd", rawPwd, "encodedPwd", encodedPwd )
	return rawPwd == encodedPwd
	// for discuz accounts
	if len(encodedPwd) == 39 {
		salt := encodedPwd[:6]
		encoded := encodedPwd[7:]
		return encoded == tools.EncodeMd5(tools.EncodeMd5(rawPwd)+salt)
	}

	// split
	var salt, encoded string
	if len(encodedPwd) > 11 {
		salt = encodedPwd[:10]
		encoded = encodedPwd[11:]
	}

	return tools.EncodePassword(rawPwd, salt) == encoded
}

// get user by erify code
func getVerifyUser(user *User, code string) bool {
	if len(code) <= conf.TimeLimitCodeLength {
		return false
	}

	// use tail hex username query user
	hexStr := code[conf.TimeLimitCodeLength:]
	if b, err := hex.DecodeString(hexStr); err == nil {
		user.UserName = string(b)
	}

	return false
}


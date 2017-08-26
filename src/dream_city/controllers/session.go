package controllers

import (
	"dream_city/models"
	"strings"
	"dream_city/library/tools"
	"github.com/astaxie/beego"
	"dream_city/library/conf"
)

func (c *baseController) CheckSessionCookieLogin() bool {
	c.IsLogin = c.readSessionUser() || c.readCookieUser()

	if c.IsLogin {
		c.Data["User"] = &c.User
	}
	c.Data["IsLogin"] = c.IsLogin

	return c.IsLogin
}

//only when login
func (c *baseController) SessionCookieLogin(cookie bool) string {
	ctx := c.Ctx

	loginRedirect := strings.TrimSpace(ctx.GetCookie("login_to"))
	if !tools.IsMatchHost(loginRedirect) {
		loginRedirect = "/users"
	} else {
		c.Ctx.SetCookie("login_to", "", -1, "/users")
	}

	c.writeSession()

	if cookie {
		c.writeCookie()
	}

	c.setLangCookie("zh-CN")

	return loginRedirect
}


func (c *baseController) SessionCookieLogout() {
	c.clearSession()
	c.clearCookie()
}

func (c *baseController) readSessionUser() bool{
	c.IsLogin = false

	id, ok := c.Ctx.Input.CruSession.Get("auth_user_id").(int64);

	if !ok || id <= 0 {
		return false
	}

	c.IsLogin = nil == models.FindById(&c.User, id)

	return c.IsLogin
}

func (c *baseController) readCookieUser() (success bool) {
	c.IsLogin = false
	userName := c.Ctx.GetCookie(conf.CookieUserName)
	if len(userName) == 0 {
		return false
	}

	defer func() {
		if !success {
			c.clearCookie()
			c.IsLogin = false
		}
	}()


	user := &c.User
	if user.LoadByName(userName) {
		return false
	}

	secret := tools.EncodeMd5(user.Rands + user.Password)
	value, _ := c.Ctx.GetSecureCookie(secret, conf.CookieRememberName)
	if value != userName {
		return false
	}

	c.writeSession()

	return true
}

func (c *baseController) writeCookie() {
	user := &c.User

	secret := tools.EncodeMd5(user.Rands + user.Password)
	days := 86400 * conf.LoginRememberDays
	c.Ctx.SetCookie(conf.CookieUserName, user.UserName, days)
	c.Ctx.SetSecureCookie(secret, conf.CookieRememberName, user.UserName, days)
}

func (c *baseController) clearCookie() {
	c.Ctx.SetCookie(conf.CookieUserName, "", -1)
	c.Ctx.SetCookie(conf.CookieRememberName, "", -1)
}

func (c *baseController) setLangCookie(lang string) {
	c.Ctx.SetCookie("lang", lang, 60*60*24*365, "/", nil, nil, false)
}


func (c *baseController) writeSession() {
	ctx := c.Ctx

	ctx.Input.CruSession.SessionRelease(ctx.ResponseWriter)
	ctx.Input.CruSession = beego.GlobalSessions.SessionRegenerateID(ctx.ResponseWriter, ctx.Request)
	ctx.Input.CruSession.Set("auth_user_id", c.User.Id)
}

func (c *baseController) clearSession() {
	ctx := c.Ctx

	ctx.Input.CruSession.Delete("auth_user_id")
	ctx.Input.CruSession.Flush()

	beego.GlobalSessions.SessionDestroy(ctx.ResponseWriter, ctx.Request)
}
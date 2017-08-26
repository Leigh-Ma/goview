package controllers

import (
	"strings"
	"dream_city/library/conf"
	. "dream_city/controllers/form"
	"dream_city/library/tools"
)

type LoginController struct {
	baseController
}


func (c *LoginController) Index() {

	loginRedirect := strings.TrimSpace(c.GetString("to"))

	// no need login
	if c.CheckLoginRedirect(false, loginRedirect) {
		return
	}

	if len(loginRedirect) > 0 {
		c.Ctx.SetCookie("login_to", loginRedirect, 0, "/")
	}

	form := LoginForm{}
	c.SetFormSets(&form)

	c.renderView("login/login.html", &D{
		"IsLoginPage": true,
	})
}


func (c *LoginController) Login() {
	c.Data["IsLoginPage"] = true
	tpl := "login/login.html"

	if c.CheckLoginRedirect(false) {
		return
	}

	errMsg := "auth.login_error_ajax"
	success := false
	once := c.Data["once_token"]
	redirect := ""

	rsp := func() *D {
		if !success {
			c.Data["Error"] = true
		}

		return &D{
			"success" :  success,
			"message" :  c.Tr(errMsg),
			"once"    :  once,
			"redirect":  redirect,
		}
	}

	f := &LoginForm{}

	if !c.ValidFormSets(f) {
		c.render(tpl, rsp())
		return
	}

	errorTry := "auth.login." + f.UserName + c.Ctx.Input.IP()
	user := &c.User

	if times, ok := tools.TimesReachedTest(errorTry, conf.LoginMaxRetries); ok {
		c.Data["ErrorReached"] = true
		errMsg = "auth.login_error_times_reached"
		c.render(tpl, rsp())
		return

	} else if !user.VerifyLogin(f.UserName, f.Password) {
		tools.TimesReachedSet(errorTry, times, conf.LoginFailedBlocks)
		c.render(tpl, rsp())
		return
	}

	success = true
	errMsg = "auth.login_success_ajax"

	redirect = c.SessionCookieLogin(f.Remember)

	c.render(tpl, rsp(), redirect)

}



// Logout implemented user logout page.
func (c *LoginController) Logout() {
	c.SessionCookieLogout()

	// write flash message
	c.FlashWrite("HasLogout", "true")

	c.Redirect("/login", 302)
}



package controllers

import (
	"dream_city/models"
	"github.com/astaxie/beego"
	"dream_city/library/types"
	. "dream_city/controllers/form"
)
// ForgotRouter serves login page.
type ForgotRouter struct {
	baseController
}

// Get implemented Get method for ForgotRouter.
func (c *ForgotRouter) Get() {
	c.TplName = "auth/forgot.html"

	// no need login
	if c.CheckLoginRedirect(false) {
		return
	}

	form := ForgotForm{}
	c.SetFormSets(&form)
}

// Get implemented Post method for ForgotRouter.
func (c *ForgotRouter) Post() {
	c.TplName = "auth/forgot.html"

	// no need login
	if c.CheckLoginRedirect(false) {
		return
	}

	var user models.User
	form := ForgotForm{User: &user}
	// valid form and put errors to template context
	if c.ValidFormSets(&form) == false {
		return
	}

	// send reset password email
	//SendResetPwdMail(c.Locale, &user)

	c.FlashRedirect("/forgot", 302, "SuccessSend")
}

// Reset implemented user password reset.
func (c *ForgotRouter) Reset() {
	c.TplName = "auth/reset.html"

	code := c.GetString(":code")
	c.Data["Code"] = code

	//var user models.User

	if false { //auth.VerifyUserResetPwdCode(&user, code) {
		c.Data["Success"] = true
		form := ResetPwdForm{}
		c.SetFormSets(&form)
	} else {
		c.Data["Success"] = false
	}
}

// Reset implemented user password reset.
func (c *ForgotRouter) ResetPost() {
	c.TplName = "auth/reset.html"

	code := c.GetString(":code")
	c.Data["Code"] = code

	var user models.User

	if false { //auth.VerifyUserResetPwdCode(&user, code) {
		c.Data["Success"] = true

		form := ResetPwdForm{}
		if c.ValidFormSets(&form) == false {
			return
		}

		user.IsActive = true
		user.Rands = types.NewGuid().String()

		if err := models.SaveNewPassword(&user, form.Password); err != nil {
			beego.Error("ResetPost Save New Password: ", err)
		}

		if c.IsLogin {
			c.SessionCookieLogout()
		}

		c.FlashRedirect("/login", 302, "ResetSuccess")

	} else {
		c.Data["Success"] = false
	}
}


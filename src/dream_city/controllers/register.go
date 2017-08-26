package controllers

import (
	"dream_city/models"
	"github.com/astaxie/beego"
	. "dream_city/controllers/form"
	"dream_city/library/types"
)

// RegisterRouter serves register page.
type RegisterRouter struct {
	baseController
}

// Get implemented Get method for RegisterRouter.
func (c *RegisterRouter) Get() {
	// no need login
	if c.CheckLoginRedirect(false) {
		return
	}

	c.Data["IsRegister"] = true
	c.TplName = "auth/register.html"

	form := RegisterForm{Locale: c.Locale}
	c.SetFormSets(&form)
}

// Register implemented Post method for RegisterRouter.
func (c *RegisterRouter) Register() {
	c.Data["IsRegister"] = true
	c.TplName = "auth/register.html"

	// no need login
	if c.CheckLoginRedirect(false) {
		return
	}

	form := RegisterForm{Locale: c.Locale}
	// valid form and put errors to template context
	if c.ValidFormSets(&form) == false {
		return
	}

	// Create new user.
	user := new(models.User)

	if err := models.RegisterUser(user, form.UserName, form.Email, form.Password); err == nil {
		loginRedirect := "/" //c.LoginUser(user, false)
		if loginRedirect == "/" {
			c.FlashRedirect("/settings/profile", 302, "RegSuccess")
		} else {
			c.Redirect(loginRedirect, 302)
		}

		return

	} else {
		beego.Error("Register: Failed ", err)
	}
}

// Active implemented check Email actice code.
func (c *RegisterRouter) Active() {
	c.TplName = "auth/active.html"

	// no need active
	if c.CheckActiveRedirect(false) {
		return
	}

	//code := c.GetString(":code")

	var user models.User


	if false {//tools.VerifyUserActiveCode(&user, code) {
		user.IsActive = true
		user.Rands = types.NewGuid().String()
		if err := user.Update("IsActive", "Rands", "Updated"); err != nil {
			beego.Error("Active: user Update ", err)
		}
		if c.IsLogin {
			c.User = user
		}

		c.Redirect("/active/success", 302)

	} else {
		c.Data["Success"] = false
	}
}

// ActiveSuccess implemented success page when email active code verified.
func (c *RegisterRouter) ActiveSuccess() {
	c.TplName = "auth/active.html"

	c.Data["Success"] = true
}

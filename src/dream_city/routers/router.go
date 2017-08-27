// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"dream_city/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("",
		beego.NSNamespace("/users",
			beego.NSInclude(
				&controllers.UsersController{},
			),
		),
	)
	beego.AddNamespace(ns)

	beego.Router("/users", &controllers.UsersController{}, "get:Index")
	beego.Router("/login", &controllers.LoginController{}, "get:Index")
	beego.Router("/login", &controllers.LoginController{}, "post:Login")
	beego.Router("/logout", &controllers.LoginController{}, "get:Logout")
	beego.Router("/blocks", &controllers.BlockController{}, "get:Index")
	beego.Router("/blocks/test", &controllers.BlockController{}, "get:Test")
}

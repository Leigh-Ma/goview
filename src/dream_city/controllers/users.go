package controllers

import (
	"github.com/astaxie/beego/orm"
	. "dream_city/models"
)


type UsersController struct {
	baseController
}
func (c *UsersController) URLMapping() {
	c.Mapping("", c.Index)

}

func (c *UsersController) Index() {
	cond := orm.NewCondition()

	if v := c.GetString("role"); v != "" {
		cond = cond.And("role", v)
	}

	if v := c.GetString("email"); v != "" {
		cond = cond.And("email", v)
	}

	qs := NewQuery((*User)(nil)).SetCond(cond)

	page := c.paginate(qs, 20)

	instances := make([]*User, 0)

	qs.Limit(20, page.Offset()).All(&instances)

	c.renderView("common/_page.html", &D{
		"paginator": page,
		"records":   instances,
	})
}

func (c *UsersController) Debug() {
	c.renderJson(&D{"x": User{}})
}

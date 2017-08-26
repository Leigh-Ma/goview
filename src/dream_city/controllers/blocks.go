package controllers

import (
	"github.com/astaxie/beego/orm"
	. "dream_city/models"
	. "dream_city/controllers/form"
	"strings"
	"github.com/astaxie/beego"
)

type BlockController struct{
	baseController
}

func (c *BlockController) Index() {
	cond := orm.NewCondition()

	f := &FSearchBlock{}
	c.ValidFormSets(f)

	beego.Info("FSearchBlock", f)


	if f.UserName != "" {
		u := &User{}
		u.LoadByName(f.UserName)
		cond = cond.And("UserId", u.Id)
	}

	if f.CityName != "" {
		cond = cond.And("CityName", f.CityName)
	}

	if f.BlockName != "" {
		s := strings.Split(f.BlockName, ",")
		if len(s) > 1 {
			cond = cond.And("Name__in", s)
		} else {
			cond = cond.And("Name__startswith", f.BlockName)
		}
	}

	qs := NewQuery((*Block)(nil)).SetCond(cond)
	page := c.paginate(qs, 20)

	instances := make([]*Block, 0)

	qs.Limit(20, page.Offset()).All(&instances)

	c.SetFormSets(f, "FSearch")

	c.renderView("blocks/index.html", &D{
		"paginator": page,
		"records":   instances,
	})

}
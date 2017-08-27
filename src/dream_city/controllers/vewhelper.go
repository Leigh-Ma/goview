package controllers

import (
    "dream_city/library/view"
)

func (c *baseController) setColumnForm(form interface{}, localeSec string, name... string) {
    p := view.NewColumnSet(form, c , localeSec, name...)

    c.Data[p.Name] = p
}

func (c *baseController) setPaginates(records interface{}, pager *view.Paginator, name... string) {
    p := view.NewColumnSet(records, c , "column", name...)
    c.Data["Table_" + p.Name]      = p
    c.Data["paginator"] = pager
}
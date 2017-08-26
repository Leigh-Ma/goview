package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"dream_city/library/view"
	"github.com/beego/i18n"
	"dream_city/models"
	"time"
	"strings"
	"fmt"
	"net/url"
	"reflect"
	"github.com/astaxie/beego/validation"
	"dream_city/library/conf"
	"html/template"
	"dream_city/library/tools"
	. "dream_city/library/view"
)

type NestPreparer interface {
	NestPrepare()
}

type D map[string]interface{}

func (d *D)Map()map[string]interface{} {
	return map[string]interface{}(*d)
}

type baseController struct{
	beego.Controller
	i18n.Locale
	User    models.User
	IsLogin bool
}

var noNeedToLogin = map[string]bool{
	"/login": true,
	"/users": true,
}

func (c *baseController) shouldLogin() bool {
	path := c.Ctx.Request.URL.Path
	return !noNeedToLogin[path]
}

func (c *baseController) paginate(q orm.QuerySeter, limit int) *view.Paginator {
	total, err :=  q.Count()
	if err != nil {
		return nil
	}

	p := view.NewPaginator(c.Ctx.Request, limit, total)

	return p
}

func (c *baseController) isJson() bool {
	return (c.Ctx.Input.Header("X-Requested-With") == "XMLHttpRequest") || (c.Ctx.Input.Header("content") == "application/json")
}

func (c *baseController) render(template string, data *D, redirect ...string) {
	if c.isJson() {
		c.renderJson(data)
	} else if len(redirect) > 0 {
		c.Redirect(redirect[0], 302)
	} else {
		c.renderView(template, data)
	}
}

func (c *baseController) renderJson(json *D) {

	c.Data["json"] = json

	c.ServeJSON()
}

func (c *baseController) renderView(template string, data *D) {
	for key, value := range data.Map() {
		c.Data[key] = value
	}

	c.Data["locale"] = "zh-CN"

	c.Layout = "common/_layout.html"
	c.TplName  = template
}

func (c *baseController) Prepare() {

	beego.Info(c.Ctx.Request.URL, "params: ",  c.Input() )
	c.Data["PageStartTime"] = time.Now()
	c.Layout = "common/_layout.html"

	// start session
	c.StartSession()

	if match, redir := c.CheckFlashRedirect(c.Ctx.Request.RequestURI); redir {
		return
	} else if match {
		c.EndFlashRedirect()
	}

	if c.CheckSessionCookieLogin() {
		if c.User.IsForbid {
			c.SessionCookieLogout()
			c.FlashRedirect("/login", 302, "UserForbid")
			return
		}
	} else if c.shouldLogin() {
		c.Redirect("/login", 302)
		return
	}

	c.Data["AppName"]   = conf.AppName
	c.Data["AppVer"]    = conf.AppVer
	c.Data["AppUrl"]    = conf.AppUrl
	c.Data["AppLogo"]   = conf.AppLogo
	c.Data["AvatarURL"] = conf.AvatarURL
	c.Data["IsProMode"] = conf.IsProMode

	if c.setLang() {
		i := strings.Index(c.Ctx.Request.RequestURI, "?")
		c.Redirect(c.Ctx.Request.RequestURI[:i], 302)
		return
	}

	beego.ReadFromRequest(&c.Controller)

	// pass xsrf helper to template context
	xsrfToken := c.Controller.XSRFToken()
	c.Data["xsrf_token"] = xsrfToken
	c.Data["xsrf_html"] = template.HTML(c.Controller.XSRFFormHTML())

	// if method is GET then auto create a form once token
	if c.Ctx.Request.Method == "GET" {
		c.FormOnceCreate()
	}

	if app, ok := c.AppController.(NestPreparer); ok {
		app.NestPrepare()
	}
}


// on router finished
func (c *baseController) Finish() {
	beego.Info(c.Data["RouterPattern"], "IsLogin: ", c.Data["IsLogin"], c.IsLogin)
}



// check if user not active then redirect
func (c *baseController) CheckActiveRedirect(args ...interface{}) bool {
	var redirect_to string
	code := 302
	needActive := true
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			needActive = v
		case string:
			// custom redirect url
			redirect_to = v
		case int:
			code = v
		}
	}
	if needActive {
		// check login
		if c.CheckLoginRedirect() {
			return true
		}

		// redirect to active page
		if !c.User.IsActive {
			c.FlashRedirect("/settings/profile", code, "NeedActive")
			return true
		}
	} else {
		// no need active
		if c.User.IsActive {
			if redirect_to == "" {
				redirect_to = "/"
			}
			c.Redirect(redirect_to, code)
			return true
		}
	}
	return false

}

// check if not login then redirect
func (c *baseController) CheckLoginRedirect(args ...interface{}) bool {
	var redirect_to string
	code := 302
	needLogin := true
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			needLogin = v
		case string:
			redirect_to = v
		case int:
			code = v
		}
	}

	// if need login then redirect
	if needLogin && !c.IsLogin {
		if len(redirect_to) == 0 {
			req := c.Ctx.Request
			scheme := "http"
			if req.TLS != nil {
				scheme += "s"
			}
			redirect_to = fmt.Sprintf("%s://%s%s", scheme, req.Host, req.RequestURI)
		}
		redirect_to = "/login?to=" + url.QueryEscape(redirect_to)
		c.Redirect(redirect_to, code)
		return true
	}

	// if not need login then redirect
	if !needLogin && c.IsLogin {
		if len(redirect_to) == 0 {
			redirect_to = "/"
		}
		c.Redirect(redirect_to, code)
		return true
	}
	return false
}

// read beego flash message
func (c *baseController) FlashRead(key string) (string, bool) {
	if data, ok := c.Data["flash"].(map[string]string); ok {
		value, ok := data[key]
		return value, ok
	}
	return "", false
}

// write beego flash message
func (c *baseController) FlashWrite(key string, value string) {
	flash := beego.NewFlash()
	flash.Data[key] = value
	flash.Store(&c.Controller)
}

// check flash redirect, ensure browser redirect to uri and display flash message.
func (c *baseController) CheckFlashRedirect(value string) (match bool, redirect bool) {
	v := c.GetSession("on_redirect")
	if params, ok := v.([]interface{}); ok {
		if len(params) != 5 {
			c.EndFlashRedirect()
			goto end
		}
		uri := tools.ToStr(params[0])
		code := 302
		if c, ok := params[1].(int); ok {
			if c/100 == 3 {
				code = c
			}
		}
		flag := tools.ToStr(params[2])
		flagVal := tools.ToStr(params[3])
		times := 0
		if v, ok := params[4].(int); ok {
			times = v
		}

		times += 1
		if times > 3 {
			// if max retry times reached then end
			c.EndFlashRedirect()
			goto end
		}

		// match uri or flash flag
		if uri == value || flag == value {
			match = true
		} else {
			// if no match then continue redirect
			c.FlashRedirect(uri, code, flag, flagVal, times)
			redirect = true
		}
	}
	end:
	return match, redirect
}

// set flash redirect
func (c *baseController) FlashRedirect(uri string, code int, flag string, args ...interface{}) {
	flagVal := "true"
	times := 0
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			flagVal = v
		case int:
			times = v
		}
	}

	if len(uri) == 0 || uri[0] != '/' {
		panic("flash reirect only support same host redirect")
	}

	params := []interface{}{uri, code, flag, flagVal, times}
	c.SetSession("on_redirect", params)

	c.FlashWrite(flag, flagVal)
	c.Redirect(uri, code)
}

// clear flash redirect
func (c *baseController) EndFlashRedirect() {
	c.DelSession("on_redirect")
}

// check form once, void re-submit
func (c *baseController) FormOnceNotMatch() bool {
	notMatch := false
	recreat := false

	// get token from request param / header
	var value string
	if vus, ok := c.Input()["_once"]; ok && len(vus) > 0 {
		value = vus[0]
	} else {
		value = c.Ctx.Input.Header("X-Form-Once")
	}

	// exist in session
	if v, ok := c.GetSession("form_once").(string); ok && v != "" {
		// not match
		if value != v {
			notMatch = true
		} else {
			// if matched then re-creat once
			recreat = true
		}
	}

	c.FormOnceCreate(recreat)
	return notMatch
}

// create form once html
func (c *baseController) FormOnceCreate(args ...bool) {
	var value string
	var creat bool
	creat = len(args) > 0 && args[0]
	if !creat {
		if v, ok := c.GetSession("form_once").(string); ok && v != "" {
			value = v
		} else {
			creat = true
		}
	}
	if creat {
		value = tools.GetRandomString(10)
		c.SetSession("form_once", value)
	}
	c.Data["once_token"] = value
	c.Data["once_html"] = template.HTML(`<input type="hidden" name="_once" value="` + value + `">`)
}

func (c *baseController) validForm(form interface{}, names ...string) (bool, map[string]*validation.Error) {
	// parse request params to form ptr struct
	ParseForm(form, c.Input())

	// Put data back in case users input invalid data for any section.
	name := reflect.ValueOf(form).Elem().Type().Name()
	if len(names) > 0 {
		name = names[0]
	}
	c.Data[name] = form

	errName := name + "Error"

	// check form once
	if c.FormOnceNotMatch() {
		return false, nil
	}

	// Verify basic input.
	valid := validation.Validation{}
	if ok, _ := valid.Valid(form); !ok {
		errs := valid.ErrorMap()
		c.Data[errName] = &valid
		return false, errs
	}
	return true, nil
}

// valid form and put errors to tempalte context
func (c *baseController) ValidForm(form interface{}, names ...string) bool {
	valid, _ := c.validForm(form, names...)
	return valid
}

// valid form and put errors to tempalte context
func (c *baseController) ValidFormSets(form interface{}, names ...string) bool {
	valid, errs := c.validForm(form, names...)
	c.setFormSets(form, errs, names...)
	beego.Info("Param parsed form: ", form )
	return valid
}

func (c *baseController) SetFormSets(form interface{}, names ...string) *FormSets {
	return c.setFormSets(form, nil, names...)
}

func (c *baseController) setFormSets(form interface{}, errs map[string]*validation.Error, names ...string) *FormSets {
	formSets := NewFormSets(form, errs, c.Locale)
	name := reflect.ValueOf(form).Elem().Type().Name()
	if len(names) > 0 {
		name = names[0]
	}
	name += "Sets"
	c.Data[name] = formSets

	return formSets
}

// add valid error to FormError
func (c *baseController) SetFormError(form interface{}, fieldName, errMsg string, names ...string) {
	name := reflect.ValueOf(form).Elem().Type().Name()
	if len(names) > 0 {
		name = names[0]
	}
	errName := name + "Error"
	setsName := name + "Sets"

	if valid, ok := c.Data[errName].(*validation.Validation); ok {
		valid.SetError(fieldName, c.Tr(errMsg))
	}

	if fSets, ok := c.Data[setsName].(*FormSets); ok {
		fSets.SetError(fieldName, errMsg)
	}
}

// check xsrf and show a friendly page
func (c *baseController) CheckXsrfCookie() bool {
	return c.Controller.CheckXSRFCookie()
}

func (c *baseController) SystemException() {

}


func (c *baseController) SetPaginator(per int, nums int64) *view.Paginator {
	p := view.NewPaginator(c.Ctx.Request, per, nums)
	c.Data["paginator"] = p
	return p
}

func (c *baseController) JsStorage(action, key string, values ...string) {
	value := action + ":::" + key
	if len(values) > 0 {
		value += ":::" + values[0]
	}
	c.Ctx.SetCookie("JsStorage", value, 1<<31-1, "/", nil, nil, false)
}

// setLang sets site language version.
func (c *baseController) setLang() bool {
	isNeedRedir := false
	hasCookie := false

	// get all lang names from i18n
	langs := conf.Langs

	// 1. Check URL arguments.
	lang := c.GetString("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = c.Ctx.GetCookie("lang")
		hasCookie = true
	} else {
		isNeedRedir = true
	}

	// Check again in case someone modify by purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		isNeedRedir = false
		hasCookie = false
	}

	// 3. check if isLogin then use user setting
	if len(lang) == 0 && c.IsLogin {
		lang = "zh-CN"
	}

	// 4. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := c.Ctx.Input.Header("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// 4. DefaucurLang language is English.
	if len(lang) == 0 {
		lang = "en-US"
		isNeedRedir = false
	}

	// Save language information in cookies.
	if !hasCookie {
		c.setLangCookie(lang)
	}

	// Set language properties.
	c.Data["Lang"] = lang
	c.Data["Langs"] = langs

	c.Lang = lang

	return isNeedRedir
}

func Values (ptrStruct interface{}, name... string)([]interface{}) {
	values := view.ViewsM.Values(ptrStruct, name...)

	return values
}

func Keys (ptrStruct interface{}, name... string)([]string) {
	keys := view.ViewsM.Keys(ptrStruct, name...)

	return keys
}

func ModelName(ptrStruct interface{}) (string)  {
	return view.ViewsM.ModelName(ptrStruct)
}
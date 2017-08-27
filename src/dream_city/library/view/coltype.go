package view

import (
    "reflect"
    "fmt"
    "strings"
    "html/template"
)

const (
    viewTag    = "show"
    localeTag  = "sec"
    HiddenCols = "HiddenCols"
    WrapCols   = "WrapCols"
)


type IModel interface {
    HiddenCols(name... string) []string
    WrapCol(colName string) string
}


type ILocale interface {
    Tr(string, ...interface{}) string
}

type selfString struct{}

func (*selfString) Tr(text string, args ...interface{}) string {
    return text
}

type ColType struct{
    ParentName  string   //form name or modelName
    Id          string   //html element id: form_$.Name
    Name        string   //column/ struct filed name, can not set
    Type        string   //text password or other
    Title       string   //label or table head, can set by option

    LocaleSec   string   //locale key section
    OptString   string   //option.String()
    TdOption    string   //html table td element attrs

    typ         *reflect.StructField
    html        *tagOption
    locale      *tagOption
    Locale      ILocale
    SampleValue *ColValue
}

func (c *ColType) AsTableTH(options ...string) template.HTML {
    option := ""
    if len(options) > 0 {
        option = NewHtmlTagOption(options[0]).Parse().ToHtmlString()
    }

    return template.HTML(fmt.Sprintf(`<th %s> %s </th>`, option, c.Tr(c.Title)))
}

func (c *ColType) AsFromInput(options string, labelStyle... string) template.HTML {
    c.html.ParseMore(options)
    ss, placeHolder := "", ""
    style := labelStyleLabel

    if len(labelStyle) > 0 {
        style = labelStyle[0]
    }

    if style == labelStyleLabel {
        ss = `<div class="form-group">`
    }

    switch style {
    case labelStyleLabel:
        ss += fmt.Sprintf(`<label class="control-label" for="%s">%s</label>`, c.Id, c.Tr(c.Title))
    case LabelStylePlaceHolder:
        if strings.Contains(c.Type, "text") {
            placeHolder = fmt.Sprintf(`placeholder=%s`, c.Tr(c.Title))
        }
    default:
        ss += fmt.Sprintf(`<b>%s:</b>`, c.Tr(c.Title))
    }


    switch c.Type {
    case "text", "number", "checkbox", "date", "datetime", "time", "password":
        ss += fmt.Sprintf(`<input class="form-control" id="%s" type="%s" name="%s" %s %s %s />`,
            c.Id, c.Type, c.Name, placeHolder, c.valueAttr(), c.OptString)
    }

    if style == labelStyleLabel {
        ss += `</div>`
    }
    return template.HTML(ss)
}

func (c *ColType) Tr(label string, other ...interface{}) string {
    if c.LocaleSec != "" {
        return c.Locale.Tr(fmt.Sprintf("%s.%s", c.LocaleSec, label), other...)
    } else {
        return c.Locale.Tr(label, other...)
    }
}


func (ct *ColType) valueAttr() string {
    return ct.SampleValue.formValueString()
}

func (ct *ColType) setLocale(tag string) {
    ct.locale = NewTagOption(tag)
    op := ct.locale

    if sec := op.Get("sec"); sec != ""{
        ct.LocaleSec = sec
        op.Remove("sec")
    }
}

func (ct *ColType) setOption(tag string){
    ct.html = NewTagOption(tag)
    op := ct.html

    if title := op.Get("title"); title != "" {
        ct.Title = title
        op.Remove("title")
    }

    if id := op.Get("id"); id != ""{
        ct.Id = id
        op.Remove("id")
    } else {
        ct.Id = fmt.Sprintf("%s_%s", ct.ParentName, ct.Name)
    }

    ct.OptString = op.TagOptionString()

    ct.Type = func() string {
        t := ct.html.Get("type")
        if  t != "" {
            return t
        }
        return ct.htmlTagOfType(ct.typ.Type.Kind())
    }()
}

func (* ColType) htmlTagOfType(k reflect.Kind) string{
    tag := ""
    switch k {
    case reflect.Bool:
        tag = "checkbox"
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
        reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        tag = "number"
    case reflect.Array, reflect.Slice, reflect.Map:
        tag = "select"
    case reflect.String:
        tag = "text"
    case reflect.Float32, reflect.Float64:
        tag = "text"
    }
    return tag
}
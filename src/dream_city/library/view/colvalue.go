package view

import (
    "reflect"
    "fmt"
    "html/template"
)

type ColValue struct{
    Value         reflect.Value
    Error         string
}

type Tuple struct {
    Fields []*ColValue
}

func (t *Tuple) AsTableRow(options... *htmlTagOption) template.HTML {
    option := ""
    if len(options) > 0 {
        option = options[0].ToHtmlString()
    }

    row := "<tr> "
    for _, td := range t.Fields {
       row  += fmt.Sprintf(`<td %s> %s </td>`, option, td.String())
    }
    row += " </tr>"

    return template.HTML(row)
}

func (t *ColValue) AsTableTD(options ...string) template.HTML {
    option := ""
    if len(options) > 0 {
        option = NewHtmlTagOption(options[0]).Parse().ToHtmlString()
    }

    td := fmt.Sprintf(`<td %s> %s </td>`, option, t.String())
    return template.HTML(td)
}


func (v *ColValue) String() string {
    if v == nil {
        return "nil"
    }
    return AsStringV(v.Value)
}

func (v *ColValue) formValueString() string {
    if v == nil {
        return ""
    }

    switch v.Value.Type().Kind() {
    case reflect.Bool:
        return "checked"
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
        reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        if v.Value.Int() != 0 {
            return fmt.Sprintf(`value="%d"`, v.Value.Int())
        }
    case reflect.String:
        if v.Value.String() != "" {
            return fmt.Sprintf(`value="%s"`, v.Value.String())
        }
    case reflect.Float32, reflect.Float64:
        return fmt.Sprintf(`value="%f"`, v.Value.Float())
    }

    return ""
}
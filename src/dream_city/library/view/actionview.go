package view

import (
	"html/template"
	"fmt"
	"strings"
	"dream_city/library/tools"
)

type htmlTagOption struct {
	options map[string]string
	raw     string  `example: "type: submit, class: btn btn-default "`
}

func NewHtmlTagOption(option string)(*htmlTagOption){
	return &htmlTagOption{raw: option, options: make(map[string]string, 0)}
}

func(o *htmlTagOption) Parse() *htmlTagOption {
	for _, kvs := range strings.Split(o.raw, ",") {
		kv := strings.Split(kvs, "=")
		k, v := "", ""

		switch len(kv) {
		case 0:
			continue
		case 1:
			k = strings.Trim(kv[0], " ")
			v = k
		default:
			k = strings.Trim(kv[0], " ")
			v = strings.Join(kv[1:], ":" )
		}

		if k != "" {
			o.options[k] = v
		}
	}

	return o
}

func(o *htmlTagOption) ToHtmlString() string {
	s := " "
	for tag, v := range o.options {
		switch tag {
		case "label": //indicate do label
			continue
		case "hidden", "checked":
			s += tag + " "
		default:
			s += fmt.Sprintf(`%s="%s"`, tag, v)
		}
	}

	return s
}

func(o *htmlTagOption) CustomLabel() string {
	return o.options["label"]
}

func optionWrap(option string) string  {
	return  NewHtmlTagOption(option).Parse().ToHtmlString()
}

func InputTag(f *FieldSet, htmlOption string) template.HTML {
	option := NewHtmlTagOption(htmlOption).Parse()
	label := option.CustomLabel()
	switch label {
	case "-":
		return template.HTML(fmt.Sprintf(
			`<input type="%s" value="%s" name="%s" %s />`,
			f.Type, tools.ToStr(f.Value), f.Name, option.ToHtmlString()))
	case "":
		label = "column." + f.Name
	default:

	}

	return template.HTML(fmt.Sprintf(
		`<b>%s: </b> <input type="%s" value="%s" name="%s" %s />`,
		f.Locale.Tr(label),
		f.Type, tools.ToStr(f.Value), f.Name, option.ToHtmlString()))

}

func CommonTag(typ string, localedValues template.HTML, htmlOption string) template.HTML {
	switch typ{
	case "submit":
		return template.HTML(fmt.Sprintf(
			`<button type="submit" %s> %s </button>`,
			optionWrap(htmlOption), localedValues))
	}
	return template.HTML("")
}
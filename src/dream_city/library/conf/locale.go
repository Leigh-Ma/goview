package conf

import (
	"strings"
	"os"
	"github.com/beego/i18n"
	"github.com/astaxie/beego"
	"html/template"
)

const (
	localePath = "conf/locale/"
)

func settingLocales() {
	dir, err:= os.Open(localePath)
	if err !=nil {
		return
	}

	locales, err := dir.Readdirnames(0)
	if err != nil {
		return
	}

	for _, locale := range locales {
		lang := strings.Replace(strings.Split(locale, ".")[0], "locale_", "", -1)
		if err := i18n.SetMessage(lang, localePath + locale); err != nil {
			beego.Error("Fail to set message file: " + err.Error())
			os.Exit(2)
		}
	}

	fun := func(lang, format string, args ...interface{}) template.HTML {
		return template.HTML(i18n.Tr(lang, format, args...))
	}

	beego.AddFuncMap("i18n", fun)
	Langs = i18n.ListLangs()
	beego.Info("langs: ", Langs)
}
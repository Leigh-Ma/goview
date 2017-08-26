package main

import (
	_ "dream_city/routers"
	"github.com/astaxie/beego"
	"github.com/Joker/jade"
	"path/filepath"
	"io/ioutil"
	"fmt"
	"html/template"
	"os"
	_ "dream_city/models"
	"github.com/astaxie/beego/orm"
	"dream_city/controllers"
	"dream_city/library/conf"
	_ "github.com/go-sql-driver/mysql"

	"dream_city/library/view"
)

func main() {

	orm.Debug = true

	orm.RunCommand()

	conf.LoadConfig()

	beeViewSetUp()

	beego.AddTemplateEngine("jade", jadeCallback)

	beego.Run()
}

func beeViewSetUp() {
	beego.AddFuncMap("Values", controllers.Values)
	beego.AddFuncMap("Keys", controllers.Keys)
	beego.AddFuncMap("ModelName", controllers.ModelName)

	beego.AddFuncMap("CommonTag", view.CommonTag)
	if !conf.IsProMode {
		beego.SetStaticPath("static_source", "static_source")
		beego.BConfig.WebConfig.DirectoryIndex = true
	}
}

func beeLocaleSetUp(){

}

func jadeCallback(root, path string, funcs template.FuncMap) (*template.Template, error) {

	jadePath := filepath.Join(root, path)

	fi,err := os.Open(jadePath)
	if err != nil {
		return nil, fmt.Errorf("error open jade template: %v", err)
	}

	content, err := ioutil.ReadAll(fi)

	if err != nil {
		return nil, fmt.Errorf("error loading jade template: %v", err)
	}

	tpl, err := jade.Parse("name_of_tpl", string(content))
	if err != nil {
		return nil, fmt.Errorf("error loading jade template: %v", err)
	}


	tmp := template.New("Person template")
	tmp, err = tmp.Parse(tpl)
	if err != nil {
		return nil, fmt.Errorf("error loading jade template: %v", err)
	}

	return tmp, err
}



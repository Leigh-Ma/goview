package view

import (
	"sync"
	"reflect"
	"fmt"
)

const ViewTag = "show"

var (
	ViewsM = &mViews{
		views: make(map[string]*viewTitle, 0),
	}
)

type IShow interface {
	TableName() string
	Show(...string) []interface{}
	Titles(...string) []string
}

type viewTitle struct {
	ViewName string
	Titles   []string
	Fields   []string
}

type mViews struct {
	sync.Mutex
	views map[string]*viewTitle
}

func (v *viewTitle) dealModel(structType reflect.Type) {
	for i := 0; i < structType.NumField(); i++ {
		f := structType.Field(i)
		tag := f.Tag.Get(ViewTag)

		if tag == "-" {
			continue
		}

		//contains anonymous struct in model
		if f.Anonymous{
			v.dealModel(f.Type)
			continue
		}


		if tag != "" {
			v.Titles = append(v.Titles, tag)
		} else {
			v.Titles = append(v.Titles, f.Name)
		}

		v.Fields = append(v.Fields, f.Name)
	}
}

func newViewTitle(structType reflect.Type, name string) *viewTitle {
	v := &viewTitle{ViewName: name}

	v.dealModel(structType)

	return v
}

func (m *mViews) getView(ptrStruct interface{}, name... string) *viewTitle {
	rName := ""
	var structType reflect.Type

	ptrType := reflect.TypeOf(ptrStruct)
	if ptrType.Kind() != reflect.Ptr {
		if ptrType.Kind() == reflect.Struct {
			structType = ptrType
		} else {
			panic(fmt.Sprintf("Using view util with none pointor type <%s>", ptrType.Kind()))
			return nil
		}
	} else {
		structType = ptrType.Elem()
		if structType.Kind() != reflect.Struct {
			panic(fmt.Sprintf("Using view util with none struct pointor type <%s>", ptrType.Kind()))
			return nil
		}
	}

	rName = structType.Name()
	if len(name) != 0 {
		rName = name[0]
	}

	m.Lock()
	view, ok := m.views[rName]
	m.Unlock()

	if !ok {
		view = newViewTitle(structType, rName)
		m.Lock()
		m.views[rName] = view
		m.Unlock()
	}

	return view
}

func (m *mViews) AddView(viewName string, titles, fields []string)(*viewTitle, bool) {
	m.Lock()
	view, ok := m.views[viewName]
	if !ok {
		view = &viewTitle{
			ViewName: viewName,
			Titles:   titles,
			Fields:   fields,
		}
		m.views[viewName] = view
	}
	m.Unlock()

	return view, ok
}

func (m *mViews) Values (ptrStruct interface{}, name... string)([]interface{}) {

	if show, ok := ptrStruct.(IShow); ok {
		r := show.Show(name...)
		values := make([]interface{}, 0)
		for _, f := range r {
			values = append(values, AsString(f))
		}
		return values
	}

	view := m.getView(ptrStruct, name...)
	if view == nil {
		return []interface{}{}
	}

	values := make([]interface{}, 0)

	//checked to ensure ptr of struct in getView
	ptrValue := reflect.ValueOf(ptrStruct).Elem()
	structValue := ptrValue.Interface()

	v := reflect.ValueOf(structValue)

	for _, field := range view.Fields {
		fieldValue := v.FieldByName(field)
		values = append(values, AsStringV(fieldValue) )
	}

	return values
}

func (m *mViews) Keys (ptrStruct interface{}, name... string)([]string) {
	k := reflect.TypeOf(ptrStruct).Kind()
	if k == reflect.Array || k == reflect.Slice {
		ptrStruct = reflect.ValueOf(ptrStruct).Index(0).Interface()
	}

	if show, ok := ptrStruct.(IShow); ok {
		return show.Titles(name...)
	}

	view := m.getView(ptrStruct, name...)
	if view == nil {
		return []string{}
	}

	return view.Titles
}

func (m *mViews) ModelName(ptrStruct interface{}) string {
	k := reflect.TypeOf(ptrStruct).Kind()
	if k == reflect.Array || k == reflect.Slice {
		ptrStruct = reflect.ValueOf(ptrStruct).Index(0).Interface()
	}

	if show, ok := ptrStruct.(IShow); ok {
		return show.TableName()
	}

	view := m.getView(ptrStruct)

	return view.ViewName
}


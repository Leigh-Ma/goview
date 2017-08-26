package view

import (
	"reflect"
	"errors"
	"fmt"
	_"time"
)

var (
	asString = &ToString{toString: make(map[string]reflect.Value, 0)}
)

func init() {
	asString.Register("Time", asStringTime)
}

func AsStringV(v reflect.Value) string {
	return asString.AsString(v)
}

func AsString(v interface{}) string {
	return asString.AsString(reflect.ValueOf(v))
}

type ToString struct {
	toString map[string]reflect.Value
}

func (t *ToString) Register(typeName string, fun interface{}) error {
	funcType := reflect.TypeOf(fun)

	if funcType.Kind() != reflect.Func {
		err := errors.New("parse: should not register other than function")
		panic(err.Error())
		return err
	}

	if funcType.NumIn() != 1 {
		err := errors.New("parse: field parser function should have one IN param")
		panic(err.Error())
		return err
	}

	if funcType.NumOut() != 1 || funcType.Out(0).Name() != "string" {
		err := errors.New("parse: field parser function should have 1 out param with type `string`")
		panic(err.Error())
		return err
	}

	t.toString[typeName] = reflect.ValueOf(fun)

	return nil
}

func (t *ToString) AsString(fieldValue reflect.Value) (string) {
	asString, ok := t.toString[fieldValue.Type().Name()]
	if ok {
		return asString.Call([]reflect.Value{reflect.ValueOf(fieldValue)})[0].String()
	}

	m := fieldValue.MethodByName("String")
	if m.IsValid() {
		param := []reflect.Value{}
		return m.Call(param)[0].String()
	}

	return fmt.Sprint(fieldValue)
}

func asStringTime(time reflect.Value) string {
	m := time.MethodByName("Format")
	return m.Call([]reflect.Value{ reflect.ValueOf("2006-01-02 15:04:05")})[0].String()
}
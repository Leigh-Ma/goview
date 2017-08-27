package view

import (
    "reflect"
    "fmt"
    "html/template"
    "strings"
)

var (
    stringArrayTypeName  = ""
    stringTypeName       = ""
)

func init() {
    stringArrayTypeName = reflect.TypeOf([]string{}).String()
    stringTypeName      = reflect.TypeOf("").String()
}

type ColumnSet struct {
    Records       []*Tuple
    Columns       []*ColType
    ColType       map[string]*ColType
    LocaleSec     string
    Locale        ILocale
    Name          string
    modelType     reflect.Type  //原始结构信息
    modelValue    reflect.Value //原始数据信息
    hiddenCols    map[string]bool
    auth          string
    colWrapperStatus    int
}

const (
    colWrapperStatusUnknown = 0
    colWrapperStatusOk = 1
    colWrapperStatusError = 2
)

//option= `name|auth`
func NewColumnSet(structPtrOrArray interface{}, locale ILocale, localeSec string, option... string) (*ColumnSet) {
    cs := &ColumnSet{
        Columns:   make([]*ColType, 0),
        ColType:   make(map[string]*ColType, 0),
        Records:   make([]*Tuple, 0),
        hiddenCols: make(map[string]bool, 0),
        LocaleSec: localeSec,
        Locale:    locale,
    }

    cs.analysisType(structPtrOrArray)
    cs.analysisCustom(locale, option...)

    cs.dealType(cs.modelType)
    cs.dealValue(structPtrOrArray)
    fmt.Printf("locale for form %+v", cs.Locale)

    return cs
}


func (cs *ColumnSet) AsInputGroup(commonOption string, labelStyle... string) []template.HTML {
    group := make([]template.HTML, 0)
    for _, col := range cs.Columns{
        group = append(group, col.AsFromInput(commonOption, labelStyle...))
    }
    return group
}

func (cs *ColumnSet) AsFormInputs(commonOption string, labelStyle... string) template.HTML {
    inputs := ""
    for _, col := range cs.Columns{
        inputs += string(col.AsFromInput(commonOption, labelStyle...))
    }
    return template.HTML(inputs)
}

func (cs *ColumnSet) AsTable(options... string) template.HTML {
    option := ""
    if len(options) > 0 {
        option = optionWrap(options[0])
    }

    head  := cs.AsTableHead()
    body  := "<tbody> "
    for _, tr := range cs.AsTableTRs() {
        body += string(tr)
    }
    body += " </tbody>"

    return template.HTML(fmt.Sprintf(`<table%s> %s %s </table>`,option, string(head), body))
}

func (cs *ColumnSet) AsTableHead() template.HTML {
    head := "<thead> <tr> "
    for _, th := range cs.AsTableTHs() {
        head += string(th)
    }
    head += " </tr> </thead>"

    return template.HTML(head)
}

func (cs *ColumnSet) AsTableTHs() []template.HTML {
    group := make([]template.HTML, 0)
    for _, col := range cs.Columns{
        group = append(group, col.AsTableTH())
    }

    return group
}

func (cs *ColumnSet) AsTableTRs(options... string) []template.HTML {
    rows := make([]template.HTML, 0)
    option := []*htmlTagOption{}

    if len(options) > 0 {
        option = append(option, NewHtmlTagOption(options[0]))
    }

    for _, row := range cs.Records{
        rows = append(rows, row.AsTableRow(option...))
    }

    return rows
}


//private
func (cs *ColumnSet) analysisType(structPtrOrArray interface{}){
    firstElem := reflect.ValueOf(structPtrOrArray)
    modelType := reflect.TypeOf(structPtrOrArray)
    if modelType.Kind() == reflect.Array || modelType.Kind() == reflect.Slice {
        firstElem = reflect.ValueOf(structPtrOrArray).Index(0)
        modelType = reflect.TypeOf(firstElem.Interface())
    }

    if reflect.Ptr == modelType.Kind() {
        modelType = modelType.Elem()
    }

    if modelType.Kind() != reflect.Struct {
        panic(fmt.Sprintf("Using view util with none struct value type <%s>", modelType))
        return
    }

    cs.modelType = modelType
    cs.modelValue = firstElem
    cs.Name = cs.modelType.Name()
}

func (cs *ColumnSet) analysisCustom(locale ILocale, option... string) {

    if m := cs.modelValue.MethodByName(HiddenCols); m.IsValid() {
        cols := m.Call([]reflect.Value{reflect.ValueOf(cs.auth)})[0]
        if cols.Type().String() == stringArrayTypeName {
            for i:= 0; i < cols.Len(); i++{
                cs.hiddenCols[cols.Index(i).String()] = true
            }
        } else {
            fmt.Printf("%s: implement HiddenCols with return value not %s <%s>!!!",
                cs.Name, stringArrayTypeName, cols.Type().String())
        }

    }

    if len(option) == 0 {
        return
    }

    os := strings.Split(option[0], "|")
    n := len(os)

    switch {
    case n == 0:
        return
    case n >= 2:
        cs.auth = strings.Trim(os[1], " ")
    }

    name := strings.Trim(os[0], " ")
    if name != "" {
        cs.Name = name
    }
}

func (cs *ColumnSet) dealValue(structPtrOrArray interface{}) {

    valueStruct := reflect.ValueOf(structPtrOrArray)
    modelType := reflect.TypeOf(structPtrOrArray)

    if modelType.Kind() == reflect.Array || modelType.Kind() == reflect.Slice {
        for i := 0; i < valueStruct.Len(); i++ {
            cs.dealValue(valueStruct.Index(i).Interface())
        }
        return
    }

    ptrStruct := valueStruct //maybe its ptr, can also be struct
    if reflect.Ptr == modelType.Kind() {
        modelType = modelType.Elem()
        valueStruct = valueStruct.Elem()
    }

    if modelType.Kind() != reflect.Struct {
        panic(fmt.Sprintf("Using view util with none struct value type <%s>", modelType))
        return
    }

    var wraper reflect.Value

    record := &Tuple{Fields: make([]*ColValue, 0)}

    if cs.colWrapperStatus != colWrapperStatusError{
        wraper = ptrStruct.MethodByName(WrapCols)
        //USE PTR in case of not defined for instance
    }

    for _, col := range cs.Columns {
        fieldValue := valueStruct.FieldByName(col.Name)
        colV := &ColValue{Value: fieldValue}

        record.Fields = append(record.Fields, colV)
        if col.SampleValue == nil {
            col.SampleValue = colV
        }

        var wrapV reflect.Value
        switch cs.colWrapperStatus {
        case colWrapperStatusUnknown:
            if wraper.IsValid() {
                v := wraper.Call([]reflect.Value{reflect.ValueOf(col.Name)})[0]
                if v.Type().String() == stringTypeName {
                    wrapV = v
                    cs.colWrapperStatus = colWrapperStatusOk
                }
            }

            if colWrapperStatusOk != cs.colWrapperStatus {
                cs.colWrapperStatus = colWrapperStatusError
            }

        case colWrapperStatusOk:
            wrapV = wraper.Call([]reflect.Value{reflect.ValueOf(col.Name)})[0]
        }

        if cs.colWrapperStatus == colWrapperStatusOk && wrapV.IsValid() && wrapV.String() != "" {
            colV.Value = wrapV
        }
    }

    cs.Records = append(cs.Records, record)
}

func (cs *ColumnSet) dealField(f *reflect.StructField) {
    tag := f.Tag.Get(viewTag)
    if tag == "-" || cs.hiddenCols[f.Name] {
        return
    }

    colType := &ColType{typ:  f,
        Name:       f.Name,
        Title:      f.Name,
        LocaleSec:  cs.LocaleSec,
        ParentName: cs.Name,
        Locale:     cs.Locale,
    }

    colType.setOption(tag)

    colType.setLocale(f.Tag.Get(localeTag))

    cs.ColType[colType.Name] = colType

    cs.Columns = append(cs.Columns, colType)
}

func (cs *ColumnSet) dealType(structType reflect.Type) {
    for i := 0; i < structType.NumField(); i++ {
        field := structType.Field(i)
        //contains anonymous struct in model
        if field.Anonymous{
            cs.dealType(field.Type)
        } else {
            cs.dealField(&field)
        }
    }
}

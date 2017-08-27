package view

import (
    "fmt"
    "strings"
)

const (
    noValue = "$$$"
)

type tagOption struct{
    options map[string]string
    raw     string  `example: "type= submit, class= btn btn-default "`
}

func(o *tagOption) Parse() *tagOption {
    for _, kvs := range strings.Split(o.raw, ",") {
        kv := strings.Split(kvs, "=")
        k, v := "", ""

        switch len(kv) {
        case 0:
            continue
        case 1:
            k = strings.Trim(kv[0], " ")
            v = noValue
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

func(o *tagOption) ParseMore(more string) *tagOption {
    for _, kvs := range strings.Split(more, ",") {
        kv := strings.Split(kvs, "=")
        k, v := "", ""

        switch len(kv) {
        case 0:
            continue
        case 1:
            k = strings.Trim(kv[0], " ")
            v = noValue
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

func (o *tagOption) TagOptionString() string {
    s := " "
    for tag, v := range o.options {
        switch v {
        case noValue:
            s += tag + " "
        default:
            s += fmt.Sprintf(`%s="%s"`, tag, v)
        }
    }

    return s
}

func (o *tagOption) Get(key string) string {
    return o.options[key]
}

func (o *tagOption) Remove(key string) {
    delete(o.options, key)
}

func NewTagOption(option string)(*tagOption){
    return &tagOption{raw: option, options: make(map[string]string, 0)}
}

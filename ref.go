package torm

import (
	"reflect"

	"github.com/Drelf2018/TypeGo/Reflect"
)

type Parser int

func (Parser) Parse(ref *Reflect.Map[[]string], elem reflect.Type) (r []string) {
	for _, field := range Reflect.FieldOf(elem) {
		if _, ok := elem.FieldByName(field.Name + "ID"); !ok {
			continue
		}
		preloads := ref.MustGetType(field.Type)
		if len(preloads) == 0 {
			r = append(r, field.Name)
			continue
		}
		prefix := field.Name + "."
		for _, preload := range preloads {
			r = append(r, prefix+preload)
		}
	}
	return
}

var Ref = Reflect.NewMap[Parser](114514, Reflect.SLICEPTRALIAS)

func ChangeParser(p Reflect.Parser[[]string], alias ...Reflect.Alias) {
	Ref = Reflect.NewMap(p, alias...)
}

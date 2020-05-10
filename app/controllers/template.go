package controllers

import (
	"reflect"
	"strings"
)

type BaseTemplate interface {
	GetLayoutFile() string
	GetCssFile() string
}

// type ControllerTemplates []ControllerTemplate
type Template struct {
	BaseTemplate BaseTemplate
	Name         string
	Files        []string
}

// getType get name from struct implete base template
func (t *Template) getType() string {
	valueOf := reflect.ValueOf(t.BaseTemplate)

	var name string
	if valueOf.Type().Kind() == reflect.Ptr {
		name = reflect.Indirect(valueOf).Type().Name()
	} else {
		name = valueOf.Type().Name()
	}

	return strings.ToLower(name[0:strings.Index(name, "BaseTemplate")])
}

func (t *Template) GetFullLayoute() string {
	return "./views/layout/" + t.BaseTemplate.GetLayoutFile()
}

func (t *Template) GetFullCss() string {
	return "./views/css/" + t.BaseTemplate.GetCssFile()
}

func (t *Template) GetFullViews() []string {
	r := []string{}
	for _, f := range t.Files {
		r = append(r, "./views/"+t.getType()+"/"+f)
	}
	return r
}

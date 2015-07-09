package view

import "html/template"

var tpl map[string]*template.Template

type outputdata struct {
	User string
}

func init() {
	tpl = make(map[string]*template.Template)
}

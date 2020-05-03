package main

import (
	"fmt"

	controllers "github.com/user/scraping-go/app/controllers"
)

type SearchBaseTemplate struct {
}

func (t *SearchBaseTemplate) GetLayoutFile() string {
	return "base.tmpl"
}

// type SearchIndexTemplate struct {
// 	SearchBaseTemplate
// }

// func (t *SearchIndexTemplate) GetName() string {
// 	return "index"
// }

// func (t *SearchIndexTemplate) GetFiles() []string {
// 	return []string{"test.tmpl"}
// }

func main() {
	var bt controllers.BaseTemplate = &SearchBaseTemplate{}
	fmt.Println(bt)

	var t controllers.Template = controllers.Template{
		BaseTemplate: &SearchBaseTemplate{},
		Name:         "index",
		Files:        []string{"test.tmpl"},
	}
	fmt.Println(t.GetFullViews())
	// // var ct controllers.ControllerTemplate
	// ct := &SearchBaseTemplate{}
	// var _ controllers.ControllerTemplate = ct
	// // controllers.GetFullLayoutes(ct)
	// // t := controllers.ControllerTemplate{LayoutFile: "test"}
	// fmt.Println(controllers.GetFullLayoutes(ct))
	// fmt.Println(getType(ct))
}

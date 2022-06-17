package websupport

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

type Model struct {
	Map map[string]any
}

func ModelAndView(w http.ResponseWriter, resources *embed.FS, view string, data Model) error {
	views := []string{
		fmt.Sprintf("resources/templates/%v.gohtml", view),
		"resources/templates/template.gohtml",
	}
	base := filepath.Base(views[0])
	return template.Must(template.New(base).Funcs(template.FuncMap{}).ParseFS(resources, views...)).Execute(w, data)
}

package helpers

import (
	"bytes"
	"fmt"
	"github.com/markbates/pkger"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type TemplateHolder struct {
	template *template.Template
}

var Rnd = TemplateHolder{}

func init() {
	Rnd.template = parseTemplates()
}

func parseTemplates() *template.Template {
	t := template.New("")
	tmp := ""

	err := pkger.Walk("/web/templates/Components", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".html" {
			open, err := pkger.Open(path)
			fmt.Println(path)
			if err != nil {
				return err
			}
			buf := new(strings.Builder)
			if _, err := io.Copy(buf, open); err != nil {
				return err
			}
			tmp += buf.String()
		}
		return err
	})

	if err != nil {
		fmt.Println(err)
	}

	_, err = t.Parse(tmp)
	if err != nil {
		fmt.Println(err)
	}
	return t
}

func Render(w http.ResponseWriter, status int, name string, v interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)

	newBuf := new(bytes.Buffer)
	defer newBuf.Reset()
	if err := Rnd.template.ExecuteTemplate(newBuf, name, v); err != nil {
		return err
	}

	_, err := w.Write(newBuf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

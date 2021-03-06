package helpers

import (
	"bytes"
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

func RndInit() {
	Rnd.template = parseTemplates()
}

func parseHtml() (str string) {
	err := filepath.Walk(TemplatePath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".html" {
			open, err := os.Open(path)
			if err != nil {
				return err
			}

			buf := new(strings.Builder)
			if _, err := io.Copy(buf, open); err != nil {
				return err
			}
			str += buf.String()
			err = open.Close()
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		LogError(err.Error())
	}
	return
}

func parseTemplates() (t *template.Template) {
	t = template.New("")
	str := parseHtml()
	_, err := t.Parse(str)
	if err != nil {
		LogError(err.Error())
	}
	return
}

func Render(w http.ResponseWriter, status int, name string, v interface{}) error {
	RndInit()
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

package initialize

import (
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"
)

var templates = template.Must(template.ParseFS(templateFS, "template/*"))

func writeTemplates(dest string, c *Cmd) error {
	templateFS, err := fs.Sub(templateFS, "template")
	if err != nil {
		return err
	}

	return fs.WalkDir(templateFS, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return os.MkdirAll(path.Join(dest, p), 0o755)
		}

		return writeTemplate(templates, p, path.Join(dest, strings.TrimSuffix(p, ".tmpl")), c)
	})
}

func writeTemplate(t *template.Template, src, dest string, data any) error {
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.ExecuteTemplate(f, src, data)
}

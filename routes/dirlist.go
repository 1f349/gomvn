package routes

import (
	_ "embed"
	"github.com/dustin/go-humanize"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//go:embed dirlist.go.html
var dirListHtml string

var dirListTemplate = template.Must(template.New("dirlist").Parse(dirListHtml))

func (r *routeCtx) handleFiles(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if containsDotDot(req.URL.Path) {
		http.Error(rw, "invalid URL path", http.StatusBadRequest)
		return
	}
	if strings.HasSuffix(req.URL.Path, "/") {
		r.handleDirList(rw, req)
		return
	}
	open, err := os.Open(filepath.Join(r.basePath, req.URL.Path))
	if err != nil {
		http.Error(rw, "404 Not Found", http.StatusNotFound)
		return
	}
	stat, err := open.Stat()
	if err != nil {
		http.Error(rw, "500 Internal Server Error: Failed to stat file", http.StatusInternalServerError)
		return
	}
	http.ServeContent(rw, req, open.Name(), stat.ModTime(), open)
}

type fileInfo struct {
	Name    string
	URL     string
	Size    string
	ModTime string
}

func (r *routeCtx) handleDirList(rw http.ResponseWriter, req *http.Request) {
	openDir, err := os.ReadDir(filepath.Join(r.basePath, req.URL.Path))
	if err != nil {
		http.Error(rw, "404 Not Found", http.StatusNotFound)
		return
	}
	fileInfos := make([]*fileInfo, len(openDir))
	for i := range openDir {
		info, err := openDir[i].Info()
		if err != nil {
			http.Error(rw, "500 Internal Server Error: Failed to stat file", http.StatusInternalServerError)
			return
		}
		url := path.Join(req.URL.Path, info.Name())
		name := path.Base(url)
		size := ""
		if info.IsDir() {
			url += "/"
			name += "/"
		} else {
			size = humanize.IBytes(uint64(info.Size()))
		}
		fileInfos[i] = &fileInfo{
			Name:    name,
			URL:     url,
			Size:    size,
			ModTime: info.ModTime().Format("2006-01-02 15:04:05 -0700"),
		}
	}
	err = dirListTemplate.Execute(rw, map[string]any{
		"Name":  r.name,
		"Path":  req.URL.Path,
		"Files": fileInfos,
	})
	if err != nil {
		log.Println("[GoMVN] Index template error: ", err)
	}
}

func (r *routeCtx) handlePut(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	p, err := r.pathUtils.ParsePath(req)
	if err != nil {
		http.Error(rw, "404 Not Found", http.StatusNotFound)
		return
	}
	p = filepath.Join(r.basePath, p)
	err = os.MkdirAll(filepath.Dir(p), os.ModePerm)
	if err != nil {
		http.Error(rw, "500 Failed to create directory", http.StatusInternalServerError)
		return
	}
	create, err := os.Create(p)
	if err != nil {
		http.Error(rw, "500 Failed to open file", http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(create, req.Body)
	if err != nil {
		http.Error(rw, "500 Failed to write file", http.StatusInternalServerError)
		return
	}
}

func containsDotDot(v string) bool {
	if !strings.Contains(v, "..") {
		return false
	}
	for _, ent := range strings.FieldsFunc(v, isSlashRune) {
		if ent == ".." {
			return true
		}
	}
	return false
}

func isSlashRune(r rune) bool { return r == '/' || r == '\\' }

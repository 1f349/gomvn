package routes

import (
	_ "embed"
	"github.com/1f349/gomvn/database"
	"github.com/1f349/gomvn/paths"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"net/http"
	"os"
)

type routeCtx struct {
	db         *database.Queries
	pathUtils  paths.Paths
	name       string
	basePath   string
	repository []string
}

func Router(db *database.Queries, name, basePath string, repository []string) http.Handler {
	pUtils := paths.Paths{Repository: repository}
	base := routeCtx{db, pUtils, name, basePath, repository}

	r := httprouter.New()
	r.PUT("/*", base.handlePut)
	r.GET("/", base.handleIndex)
	r.GET("/*", base.handleGet)
	return r
}

//go:embed index.go.html
var indexHtml string

var indexTemplate = template.Must(template.New("index").Parse(indexHtml))

func (r *routeCtx) handleIndex(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	_ = indexTemplate.Execute(rw, map[string]any{
		"Name":         r.name,
		"Repositories": paths.GetRepositories(r.basePath, r.repository),
	})
}

func (r *routeCtx) handlePut(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	p, err := r.pathUtils.ParsePath(req)
	if err != nil {
		http.Error(rw, "404 Not Found", http.StatusNotFound)
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

func (r *routeCtx) handleGet(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {

}

package routes

import (
	_ "embed"
	"encoding/json"
	"github.com/1f349/gomvn/database"
	"github.com/1f349/gomvn/paths"
	"github.com/julienschmidt/httprouter"
	"github.com/thanhpk/randstr"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
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

	rApi := httprouter.New()
	rApi.GET("/users", func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		users, err := db.GetAllUsers(req.Context())
		if err != nil {
			http.Error(rw, "500 Database Error", http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(rw).Encode(users)
	})
	rApi.POST("/users", func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		err := req.ParseForm()
		if err != nil {
			http.Error(rw, "400 Bad Request", http.StatusBadRequest)
			return
		}
		name := req.PostForm.Get("name")
		admin := req.PostForm.Has("admin")
		hex := randstr.Hex(32)

		_, err = db.CreateUser(req.Context(), database.CreateUserParams{
			Name:      name,
			Admin:     admin,
			TokenHash: hex,
		})
		if err != nil {
			http.Error(rw, "500 Database Error", http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(rw).Encode(map[string]any{
			"token": hex,
		})
	})

	rWeb := httprouter.New()
	rWeb.PUT("/*filepath", base.repoAuth(base.handlePut))
	rWeb.GET("/", base.handleIndex)
	for _, repo := range repository {
		rWeb.ServeFiles(path.Join("/", repo, "*filepath"), http.FS(os.DirFS(filepath.Join(basePath, repo))))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(rw http.ResponseWriter, req *http.Request) {
		isAdmin, err := db.IsAdmin(req.Context(), req.Header.Get("Authorization"))
		if err != nil {
			http.Error(rw, "500 Database Error", http.StatusInternalServerError)
			return
		}
		if isAdmin != 1 {
			http.Error(rw, "403 Forbidden", http.StatusForbidden)
			return
		}
		rApi.ServeHTTP(rw, req)
	})
	mux.Handle("/", rWeb)

	return mux
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

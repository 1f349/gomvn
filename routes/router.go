package routes

import (
	_ "embed"
	"encoding/json"
	"github.com/1f349/gomvn/database"
	"github.com/1f349/gomvn/paths"
	"github.com/julienschmidt/httprouter"
	"github.com/thanhpk/randstr"
	"net/http"
	"path"
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
	rWeb.GET("/", base.handleFiles)
	for _, repo := range repository {
		rWeb.GET(path.Join("/", repo, "*filepath"), base.handleFiles)
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

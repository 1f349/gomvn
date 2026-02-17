package routes

import (
	"context"
	"encoding/base64"
	"github.com/1f349/gomvn/database"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (r *routeCtx) repoAuth(next httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
		un, pw, ok := parseBasicBearer(req)
		if !ok {
			http.Error(rw, "403 Forbidden", http.StatusForbidden)
			return
		}
		isValid, err := r.db.CheckUserDetails(context.Background(), database.CheckUserDetailsParams{
			Name:      un,
			TokenHash: pw,
		})
		if err != nil || isValid != 1 {
			http.Error(rw, "403 Forbidden", http.StatusForbidden)
			return
		}
		next(rw, req, params)
	}
}

func parseBasicBearer(req *http.Request) (string, string, bool) {
	auth := req.Header.Get("Authorization")
	details, ok := strings.CutPrefix(auth, "Basic ")
	if !ok {
		return "", "", false
	}
	decBytes, err := base64.StdEncoding.DecodeString(details)
	if err != nil {
		return "", "", false
	}
	decStr := string(decBytes)
	before, after, ok0 := strings.Cut(decStr, ":")
	if !ok0 {
		return "", "", false
	}
	return before, after, true
}

package routes

import (
	"github.com/julienschmidt/httprouter"
	"strconv"
)

func getQueryUserId(params httprouter.Params) (int64, bool) {
	if val, err := strconv.ParseInt(params.ByName("id"), 10, 64); err == nil {
		return val, true
	}
	return 0, false
}

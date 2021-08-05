package util

import (
	"net/http"

	"github.com/tnyie/journaler-api/middleware"
)

func GetUserID(r *http.Request) string {
	id := r.Context().Value(middleware.AuthCtx{})

	if id == nil {
		return ""
	}

	return id.(string)
}

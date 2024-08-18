package api

import (
	"database/sql"
	"net/http"
)

func NewRouter(db *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	return mux
}

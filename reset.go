package main

import (
	"fmt"
	"net/http"
)

func (a *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if a.platform != "dev" {
		fmt.Println("Permission denied, platform=%s", a.platform)
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	a.fileServerHits.Store(0)
	if err := a.db.DeleteAllUsers(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset the database: " + err.Error()))
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf(nil, "Count reset to %d and database reset to initial state.", a.fileServerHits.Load()))
}

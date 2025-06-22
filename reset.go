package main

import (
	"fmt"
	"net/http"
)

func (a *apiConfig) handlerReset(w http.ResponseWriter, _ *http.Request) {
	a.fileServerHits.Store(0)
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf(nil, "Count reset to %d", a.fileServerHits.Load()))
}

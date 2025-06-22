package main

import (
	"fmt"
	"net/http"
)

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		a.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (a *apiConfig) handlerHits(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf(nil, "Hits: %d", a.fileServerHits.Load()))
}

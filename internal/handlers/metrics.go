package handlers

import (
	"fmt"
	"net/http"
)

func (a *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		a.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (a *ApiConfig) HandlerHits(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf(nil, `
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>
	`, a.FileServerHits.Load()))
}

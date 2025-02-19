package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/WhileCodingDoLearn/bootdev_server_tut/docs"
)

func (apiCfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiCfg.filserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func init() {
	docs.GetDocsGenerator().AddDoc(docs.EndpointDescription{
		Name:           "/admin/metrics",
		Method:         "GET",
		InputType:      nil,
		OutputType:     nil,
		Description:    "Enpoint for Server metrics",
		Authentication: false,
		ErrorCodes:     []int{},
	})
}

func (apiCfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	clicks := &apiCfg.filserverHits
	res := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", clicks.Load())
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(res))
}

func init() {
	docs.GetDocsGenerator().AddDoc(docs.EndpointDescription{
		Name:           "/admin/reset",
		Method:         "POST",
		InputType:      nil,
		OutputType:     nil,
		Description:    "Resets Counter for Server metrics. Deletes all users (only in dev)",
		Authentication: false,
		ErrorCodes:     []int{},
	})
}

func (apiCfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	counter := &apiCfg.filserverHits
	counter.Store(0)
	if apiCfg.plattform == "Dev" {
		apiCfg.dbQueries.DeleteAllUsers(r.Context())
		responWithoutBody(w, http.StatusOK)
		return
	}
	respondWithError(w, http.StatusForbidden, apiCfg.plattform)

}

func (apiCfg *apiConfig) middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

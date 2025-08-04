package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)	
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
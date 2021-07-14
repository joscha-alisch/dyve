package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/root").Subrouter().Handle("", someOtherHandlerThatIDontControl())

	http.ListenAndServe(":9005", r)
}

func someOtherHandlerThatIDontControl() http.Handler {
	other := mux.NewRouter()
	other.HandleFunc("/sub", func(w http.ResponseWriter, request *http.Request) {
		w.WriteHeader(200)
	})
	return other
}

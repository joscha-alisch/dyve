package sdk

import "net/http"
import "github.com/gorilla/mux"

func ListenAndServeAppProvider(addr string, p AppProvider) error {
	return http.ListenAndServe(addr, NewAppProviderHandler(p))
}

func NewAppProviderHandler(p AppProvider) http.Handler {
	h := &appProviderHandler{Router: mux.NewRouter(), p: p}

	h.HandleFunc("/apps", h.listApps)
	h.HandleFunc("/apps/{id:[0-9]+}", h.getApp)

	return h
}

type appProviderHandler struct {
	*mux.Router

	p AppProvider
}

func (h *appProviderHandler) listApps(w http.ResponseWriter, r *http.Request) {

}

func (h *appProviderHandler) getApp(w http.ResponseWriter, r *http.Request) {

}

package sdk

import (
	"errors"
	"net/http"
)
import "github.com/gorilla/mux"

func ListenAndServeAppProvider(addr string, p AppProvider) error {
	return ListenAndServe(addr, ProviderConfig{Apps: p})
}

func NewAppProviderHandler(p AppProvider) http.Handler {
	h := &appProviderHandler{Router: mux.NewRouter(), p: p}

	h.HandleFunc("/apps", h.listApps)
	h.HandleFunc("/apps/", h.listApps)
	h.HandleFunc("/apps/{id:[0-9a-z-]+}", h.getApp)

	return h
}

type appProviderHandler struct {
	*mux.Router

	p AppProvider
}

func (h *appProviderHandler) listApps(w http.ResponseWriter, r *http.Request) {
	apps, err := h.p.ListApps()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, apps)
}

func (h *appProviderHandler) getApp(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	app, err := h.p.GetApp(id)
	if errors.Is(err, ErrNotFound) {
		respondErr(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, app)
}

package sdk

import (
	"errors"
	"net/http"
)
import "github.com/gorilla/mux"

func ListenAndServeGroupProvider(addr string, p GroupProvider) error {
	return ListenAndServe(addr, ProviderConfig{Groups: p})
}

func NewGroupProviderHandler(p GroupProvider) http.Handler {
	h := &groupProviderHandler{Router: mux.NewRouter(), p: p}

	h.HandleFunc("/groups", h.listGroups)
	h.HandleFunc("/groups/", h.listGroups)
	h.HandleFunc("/groups/{id:[0-9a-z-]+}", h.getGroup)

	return h
}

type groupProviderHandler struct {
	*mux.Router

	p GroupProvider
}

func (h *groupProviderHandler) listGroups(w http.ResponseWriter, r *http.Request) {
	apps, err := h.p.ListGroups()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}
	respondOk(w, apps)
}

func (h *groupProviderHandler) getGroup(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	app, err := h.p.GetGroup(id)
	if errors.Is(err, ErrNotFound) {
		respondErr(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, app)
}

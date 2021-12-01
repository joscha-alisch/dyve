package sdk

import (
	"errors"
	"net/http"
)
import "github.com/gorilla/mux"

func ListenAndServeAppInstancesProvider(addr string, p AppProvider) error {
	return ListenAndServe(addr, ProviderConfig{Apps: p})
}

func NewAppInstancesProviderHandler(p InstancesProvider) http.Handler {
	h := &appInstancesProviderHandler{Router: mux.NewRouter(), p: p}

	h.HandleFunc("/instances/{id:[0-9a-z-]+}", h.getAppInstances)

	return h
}

type appInstancesProviderHandler struct {
	*mux.Router

	p InstancesProvider
}

func (h *appInstancesProviderHandler) getAppInstances(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	app, err := h.p.GetAppInstances(id)
	if errors.Is(err, ErrNotFound) {
		respondErr(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, app)
}

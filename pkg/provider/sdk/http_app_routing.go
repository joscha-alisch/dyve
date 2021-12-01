package sdk

import (
	"errors"
	"net/http"
)
import "github.com/gorilla/mux"

func ListenAndServeAppRoutingProvider(addr string, p AppProvider) error {
	return ListenAndServe(addr, ProviderConfig{Apps: p})
}

func NewAppRoutingProviderHandler(p RoutingProvider) http.Handler {
	h := &appRoutingProviderHandler{Router: mux.NewRouter(), p: p}

	h.HandleFunc("/routing/{id:[0-9a-z-]+}", h.getAppRouting)

	return h
}

type appRoutingProviderHandler struct {
	*mux.Router

	p RoutingProvider
}

func (h *appRoutingProviderHandler) getAppRouting(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	app, err := h.p.GetAppRouting(id)
	if errors.Is(err, ErrNotFound) {
		respondErr(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, app)
}

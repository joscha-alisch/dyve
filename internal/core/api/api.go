package api

import (
	"github.com/gorilla/mux"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"net/http"
)

func New(appProviders []sdk.AppProvider) http.Handler {
	return &api{}
}

type api struct {
	*mux.Router
}

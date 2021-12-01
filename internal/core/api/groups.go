package api

import "net/http"

func (a *api) listGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := a.core.Groups.ListGroupsByProvider()
	if err != nil {
		respondErr(w, http.StatusInternalServerError, err)
		return
	}

	respondOk(w, groups)
}

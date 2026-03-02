package server

import (
	"encoding/json"
	"net/http"
)

func (s *server) handleUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	userID, err := s.ensureLoggedIn(r)
	if err != nil {
		s.writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	user, err := s.db.GetUser(r.Context(), userID)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := user.ToView()
	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		s.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

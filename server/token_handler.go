package server

import "net/http"

func (s *server) handleToken(w http.ResponseWriter, r *http.Request) {
	s.logger.Info(r.Context(), "handleToken")

}

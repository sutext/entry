package server

import (
	"net/http"

	"sutext.github.io/entry/xlog"
)

func (s *server) handleAuth(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("auth request", xlog.Ctx(r.Context()))

}

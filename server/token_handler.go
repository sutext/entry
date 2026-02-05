package server

import (
	"net/http"

	"sutext.github.io/entry/xlog"
)

func (s *server) handleToken(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("handleToken", xlog.Ctx(r.Context()))

}
func (s *server) handleTokenRefresh(w http.ResponseWriter, r *http.Request) {

	s.logger.Info("handleTokenRefresh", xlog.Ctx(r.Context()))
}

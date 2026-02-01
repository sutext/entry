package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"path"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"sutext.github.io/entry/model"
	"sutext.github.io/entry/xlog"
)

type Server interface {
	Serve() error
	Shoutdown(ctx context.Context) error
	HandleFunc(p string, h http.HandlerFunc)
	HandleCORS(p string, h http.HandlerFunc)
	HandlePrefix(p string, h http.Handler)
}
type server struct {
	db                     model.Storage
	mux                    *mux.Router
	logger                 *xlog.Logger
	dirver                 model.Driver
	issuerURL              url.URL
	allHeaders             http.Header
	realIPHeader           string
	allowedOrigins         []string
	allowedHeaders         []string
	trustedRealIPCIDRs     []*netip.Prefix
	prometheusRegistry     *prometheus.Registry
	supportedGrantTypes    []string
	supportedResponseTypes map[string]bool
}

func New(opts ...Option) Server {
	options := newOptions(opts...)
	issuerURL, err := url.Parse(options.issuerURL)
	if err != nil {
		panic(err)
	}
	s := &server{
		mux:                    mux.NewRouter(),
		logger:                 options.logger,
		dirver:                 options.dirver,
		issuerURL:              *issuerURL,
		allHeaders:             options.allHeaders,
		realIPHeader:           options.realIPHeader,
		allowedOrigins:         options.allowedOrigins,
		allowedHeaders:         options.allowedHeaders,
		trustedRealIPCIDRs:     options.trustedRealIPCIDRs,
		supportedGrantTypes:    options.supportedGrantTypes,
		supportedResponseTypes: options.supportedResponseTypes,
	}
	return s
}
func (s *server) Serve() error {
	db, err := model.Open(s.dirver)
	if err != nil {
		return err
	}
	s.db = db
	s.mux.NotFoundHandler = http.NotFoundHandler()
	s.HandleCORS("/", s.handleRoot)
	s.HandleCORS("/.well-known/openid-configuration", s.handleDiscovery)
	s.HandleFunc("/token", s.handleToken)
	s.logger.Info(context.Background(), "Server started")
	http.ListenAndServe(":8080", s.mux)
	return nil
}
func (s *server) Shoutdown(ctx context.Context) error {
	return nil
}
func (s *server) handleRoot(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w,
		`<!DOCTYPE html>
		<title>Entry</title>
		<h1>Entry IdP</h1>
		<h3>A Federated OpenID Connect Provider</h3>
		<p><a href=%q>Discovery</a></p>`, s.absURL("/.well-known/openid-configuration"))
	if err != nil {
		s.logger.Error(r.Context(), "failed to write response", xlog.Err(err))
		// s.renderError(r, w, http.StatusInternalServerError, "Handling the / path error.")
		return
	}
}
func (s *server) HandleFunc(p string, h http.HandlerFunc) {
	s.mux.Handle(path.Join(s.issuerURL.Path, p), s.handlerWithHeaders(p, h))
}
func (s *server) HandleCORS(p string, h http.HandlerFunc) {
	var handler http.Handler = h
	if len(s.allowedOrigins) > 0 {
		cors := handlers.CORS(
			handlers.AllowedOrigins(s.allowedOrigins),
			handlers.AllowedHeaders(s.allowedHeaders),
		)
		handler = cors(handler)
	}
	s.mux.Handle(path.Join(s.issuerURL.Path, p), s.handlerWithHeaders(p, handler))
}
func (s *server) HandlePrefix(p string, h http.Handler) {
	prefix := path.Join(s.issuerURL.Path, p)
	s.mux.PathPrefix(prefix).Handler(http.StripPrefix(prefix, h))
}
func (s *server) handlerWithHeaders(handlerName string, handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for k, v := range s.allHeaders {
			w.Header()[k] = v
		}
		// Context values are used for logging purposes with the log/slog logger.
		rCtx := r.Context()
		rCtx = context.WithValue(rCtx, xlog.KeyRequestID, uuid.NewString())

		if s.realIPHeader != "" {
			realIP, err := s.parseRealIP(r)
			if err == nil {
				rCtx = context.WithValue(rCtx, xlog.KeyRemoteIP, realIP)
			}
		}
		instrumentHandler := func(_ string, handler http.Handler) http.HandlerFunc {
			return handler.ServeHTTP
		}
		r = r.WithContext(rCtx)
		if s.prometheusRegistry != nil {
			requestCounter := prometheus.NewCounterVec(prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Count of all HTTP requests.",
			}, []string{"code", "method", "handler"})

			durationHist := prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name:    "request_duration_seconds",
				Help:    "A histogram of latencies for requests.",
				Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
			}, []string{"code", "method", "handler"})

			sizeHist := prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Name:    "response_size_bytes",
				Help:    "A histogram of response sizes for requests.",
				Buckets: []float64{200, 500, 900, 1500},
			}, []string{"code", "method", "handler"})

			s.prometheusRegistry.MustRegister(requestCounter, durationHist, sizeHist)

			instrumentHandler = func(handlerName string, handler http.Handler) http.HandlerFunc {
				return promhttp.InstrumentHandlerDuration(
					durationHist.MustCurryWith(prometheus.Labels{"handler": handlerName}),
					promhttp.InstrumentHandlerCounter(
						requestCounter.MustCurryWith(prometheus.Labels{"handler": handlerName}),
						promhttp.InstrumentHandlerResponseSize(
							sizeHist.MustCurryWith(prometheus.Labels{"handler": handlerName}),
							handler),
					),
				)
			}
		}
		instrumentHandler(handlerName, handler)(w, r)
	}
}
func (s *server) parseRealIP(r *http.Request) (string, error) {
	remoteAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	remoteIP, err := netip.ParseAddr(remoteAddr)
	if err != nil {
		return "", err
	}
	for _, n := range s.trustedRealIPCIDRs {
		if !n.Contains(remoteIP) {
			return remoteAddr, nil // Fallback to the address from the request if the header is provided
		}
	}
	ipVal := r.Header.Get(s.realIPHeader)
	if ipVal != "" {
		ip, err := netip.ParseAddr(ipVal)
		if err == nil {
			return ip.String(), nil
		}
	}
	return remoteAddr, nil
}
func (s *server) absURL(pathItems ...string) string {
	u := s.issuerURL
	u.Path = s.absPath(pathItems...)
	return u.String()
}
func (s *server) absPath(pathItems ...string) string {
	paths := make([]string, len(pathItems)+1)
	paths[0] = s.issuerURL.Path
	copy(paths[1:], pathItems)
	return path.Join(paths...)
}

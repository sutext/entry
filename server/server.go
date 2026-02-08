package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/netip"
	"net/url"
	"path"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"sutext.github.io/entry/model"
	"sutext.github.io/entry/xerr"
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
	db                            model.Storage
	mux                           *mux.Router
	logger                        *xlog.Logger
	dirver                        model.Driver
	issuerURL                     url.URL
	forcePKCE                     bool
	allHeaders                    http.Header
	realIPHeader                  string
	allowedOrigins                []string
	allowedHeaders                []string
	trustedRealIPCIDRs            []*netip.Prefix
	internalErrorHandler          func(error) *xerr.Response
	supportedGrantTypes           map[string]struct{}
	supportedResponseTypes        map[string]struct{}
	supportedCodeChallengeMethods map[string]struct{}
}

func New(opts ...Option) Server {
	options := newOptions(opts...)
	issuerURL, err := url.Parse(options.issuerURL)
	if err != nil {
		panic(err)
	}
	s := &server{
		mux:                           mux.NewRouter(),
		logger:                        options.logger,
		dirver:                        options.dirver,
		issuerURL:                     *issuerURL,
		allHeaders:                    options.allHeaders,
		realIPHeader:                  options.realIPHeader,
		allowedOrigins:                options.allowedOrigins,
		allowedHeaders:                options.allowedHeaders,
		trustedRealIPCIDRs:            options.trustedRealIPCIDRs,
		supportedGrantTypes:           options.supportedGrantTypes,
		supportedResponseTypes:        options.supportedResponseTypes,
		supportedCodeChallengeMethods: options.supportedCodeChallengeMethods,
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
	// s.HandleFunc("/token", s.handleToken)
	s.HandleFunc("/authorize", s.handleAuthorize)
	s.logger.Info("Server started")
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
		s.logger.Error("failed to write response", xlog.Ctx(r.Context()), xlog.Err(err))
		// s.renderError(r, w, http.StatusInternalServerError, "Handling the / path error.")
		return
	}
}
func (s *server) HandleFunc(p string, h http.HandlerFunc) {
	s.mux.Handle(path.Join(s.issuerURL.Path, p), h)
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
	s.mux.Handle(path.Join(s.issuerURL.Path, p), handler)
}
func (s *server) HandlePrefix(p string, h http.Handler) {
	prefix := path.Join(s.issuerURL.Path, p)
	s.mux.PathPrefix(prefix).Handler(http.StripPrefix(prefix, h))
}
func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
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
func (s *server) checkResponseType(rt ResponseType) bool {
	for art := range s.supportedResponseTypes {
		if art == rt.String() {
			return true
		}
	}
	return false
}

// CheckCodeChallengeMethod checks for allowed code challenge method
func (s *server) checkCodeChallengeMethod(ccm CodeChallengeMethod) bool {
	for c := range s.supportedCodeChallengeMethods {
		if c == ccm.String() {
			return true
		}
	}
	return false
}
func (s *server) checkGrantType(gt GrantType) bool {
	for agt := range s.supportedGrantTypes {
		if agt == gt.String() {
			return true
		}
	}
	return false
}
func (s *server) getErrorData(err error) (map[string]any, int, http.Header) {
	var re xerr.Response
	if v, ok := xerr.Descriptions[err]; ok {
		re.Error = err
		re.Description = v
		re.StatusCode = xerr.StatusCodes[err]
	} else {
		if fn := s.internalErrorHandler; fn != nil {
			if v := fn(err); v != nil {
				re = *v
			}
		}

		if re.Error == nil {
			re.Error = xerr.ErrServerError
			re.Description = xerr.Descriptions[xerr.ErrServerError]
			re.StatusCode = xerr.StatusCodes[xerr.ErrServerError]
		}
	}

	data := make(map[string]interface{})
	if err := re.Error; err != nil {
		data["error"] = err.Error()
	}

	if v := re.ErrorCode; v != 0 {
		data["error_code"] = v
	}

	if v := re.Description; v != "" {
		data["error_description"] = v
	}

	if v := re.URI; v != "" {
		data["error_uri"] = v
	}

	statusCode := http.StatusInternalServerError
	if v := re.StatusCode; v > 0 {
		statusCode = v
	}

	return data, statusCode, re.Header
}

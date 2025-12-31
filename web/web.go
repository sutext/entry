package web

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
)

//go:embed static/* templates/* themes/* robots.txt
var files embed.FS

// FS returns a filesystem with the default web assets.
func FS() fs.FS {
	return files
}

type WebSite struct {
	Static    http.Handler
	Themes    http.Handler
	Robots    http.HandlerFunc
	templates map[string]*template.Template
}

const (
	tmplApproval      = "approval.html"
	tmplLogin         = "login.html"
	tmplPassword      = "password.html"
	tmplOOB           = "oob.html"
	tmplError         = "error.html"
	tmplDevice        = "device.html"
	tmplDeviceSuccess = "device_success.html"
)

var requiredTmpls = []string{
	tmplApproval,
	tmplLogin,
	tmplPassword,
	tmplOOB,
	tmplError,
	tmplDevice,
	tmplDeviceSuccess,
}

type Config struct {
	FS        fs.FS
	LogoURL   string
	Issuer    string
	Theme     string
	IssuerURL string
	Extra     map[string]string
}

func funcMap(c Config) (template.FuncMap, error) {
	issuerURL, err := url.Parse(c.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing issuerURL: %v", err)
	}

	funcs := map[string]any{
		"extra":  func(k string) string { return c.Extra[k] },
		"issuer": func() string { return c.Issuer },
		"logo":   func() string { return c.LogoURL },
		"lower":  strings.ToLower,
		"trim":   strings.TrimSpace,
		"upper":  strings.ToUpper,
		"url": func(reqPath, assetPath string) string {
			return relativeURL(issuerURL.Path, reqPath, assetPath)
		},
	}

	return funcs, nil
}

// LoadConfig returns static assets, theme assets, and templates used by the frontend by
// reading the dir specified in the webConfig. If directory is not specified it will
// use the file system specified by webFS.
//
// The directory layout is expected to be:
//
//	( web directory )
//	|- static
//	|- themes
//	|  |- (theme name)
//	|- templates
func NewWebSite(c Config) (*WebSite, error) {
	// fallback to the default theme if the legacy theme name is provided
	if c.Theme == "coreos" || c.Theme == "tectonic" {
		c.Theme = ""
	}
	if c.Theme == "" {
		c.Theme = "light"
	}
	if c.Issuer == "" {
		c.Issuer = "dex"
	}
	if c.LogoURL == "" {
		c.LogoURL = "theme/logo.png"
	}

	staticFiles, err := fs.Sub(c.FS, "static")
	if err != nil {
		return nil, fmt.Errorf("read static dir: %v", err)
	}
	themeFiles, err := fs.Sub(c.FS, path.Join("themes", c.Theme))
	if err != nil {
		return nil, fmt.Errorf("read themes dir: %v", err)
	}
	robotsContent, err := fs.ReadFile(c.FS, "robots.txt")
	if err != nil {
		return nil, fmt.Errorf("read robots.txt dir: %v", err)
	}

	static := http.FileServer(http.FS(staticFiles))
	theme := http.FileServer(http.FS(themeFiles))
	robots := func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, string(robotsContent)) }

	templates, err := loadTemplates(c, "templates")

	return &WebSite{static, theme, robots, templates}, err
}

// loadTemplates parses the expected templates from the provided directory.
func loadTemplates(c Config, templatesDir string) (map[string]*template.Template, error) {
	files, err := fs.ReadDir(c.FS, templatesDir)
	if err != nil {
		return nil, fmt.Errorf("read dir: %v", err)
	}

	filenames := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filenames = append(filenames, path.Join(templatesDir, file.Name()))
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("no files in template dir %q", templatesDir)
	}

	funcs, err := funcMap(c)
	if err != nil {
		return nil, err
	}

	tmpls, err := template.New("").Funcs(funcs).ParseFS(c.FS, filenames...)
	if err != nil {
		return nil, fmt.Errorf("parse files: %v", err)
	}
	templates := make(map[string]*template.Template)
	missingTmpls := []string{}
	for _, tmplName := range requiredTmpls {
		tmpl := tmpls.Lookup(tmplName)
		if tmpl == nil {
			missingTmpls = append(missingTmpls, tmplName)
		} else {
			templates[tmplName] = tmpl
		}
	}
	if len(missingTmpls) > 0 {
		return nil, fmt.Errorf("missing template(s): %s", missingTmpls)
	}
	return templates, nil
}

// relativeURL returns the URL of the asset relative to the URL of the request path.
// The serverPath is consulted to trim any prefix due in case it is not listening
// to the root path.
//
// Algorithm:
// 1. Remove common prefix of serverPath and reqPath
// 2. Remove common prefix of assetPath and reqPath
// 3. For each part of reqPath remaining(minus one), go up one level (..)
// 4. For each part of assetPath remaining, append it to result
//
// eg
// server listens at localhost/dex so serverPath is dex
// reqPath is /dex/auth
// assetPath is static/main.css
// relativeURL("/dex", "/dex/auth", "static/main.css") = "../static/main.css"
func relativeURL(serverPath, reqPath, assetPath string) string {
	if u, err := url.ParseRequestURI(assetPath); err == nil && u.Scheme != "" {
		// assetPath points to the external URL, no changes needed
		return assetPath
	}

	splitPath := func(p string) []string {
		res := []string{}
		parts := strings.Split(path.Clean(p), "/")
		for _, part := range parts {
			if part != "" {
				res = append(res, part)
			}
		}
		return res
	}

	stripCommonParts := func(s1, s2 []string) ([]string, []string) {
		min := len(s1)
		if len(s2) < min {
			min = len(s2)
		}

		splitIndex := min
		for i := 0; i < min; i++ {
			if s1[i] != s2[i] {
				splitIndex = i
				break
			}
		}
		return s1[splitIndex:], s2[splitIndex:]
	}

	server, req, asset := splitPath(serverPath), splitPath(reqPath), splitPath(assetPath)

	// Remove common prefix of request path with server path
	_, req = stripCommonParts(server, req)

	// Remove common prefix of request path with asset path
	asset, req = stripCommonParts(asset, req)

	// For each part of the request remaining (minus one) -> go up one level (..)
	// For each part of the asset remaining               -> append it
	var relativeURL string
	for i := 0; i < len(req)-1; i++ {
		relativeURL = path.Join("..", relativeURL)
	}
	relativeURL = path.Join(relativeURL, path.Join(asset...))

	return relativeURL
}

var scopeDescriptions = map[string]string{
	"offline_access": "Have offline access",
	"profile":        "View basic profile information",
	"email":          "View your email address",
	// 'groups' is not a standard OIDC scope, and Dex only returns groups only if the upstream provider does too.
	// This warning is added for convenience to show that the user may expose some sensitive data to the application.
	"groups": "View your groups",
}

type ConnectorInfo struct {
	ID   string
	Name string
	URL  template.URL
	Type string
}

type byName []ConnectorInfo

func (n byName) Len() int           { return len(n) }
func (n byName) Less(i, j int) bool { return n[i].Name < n[j].Name }
func (n byName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }

func (s *WebSite) RenderDevice(r *http.Request, w http.ResponseWriter, postURL string, userCode string, lastWasInvalid bool) error {
	if lastWasInvalid {
		w.WriteHeader(http.StatusBadRequest)
	}
	data := struct {
		PostURL  string
		UserCode string
		Invalid  bool
		ReqPath  string
	}{postURL, userCode, lastWasInvalid, r.URL.Path}
	return renderTemplate(w, s.templates[tmplDevice], data)
}

func (s *WebSite) RenderDeviceSuccess(r *http.Request, w http.ResponseWriter, clientName string) error {
	data := struct {
		ClientName string
		ReqPath    string
	}{clientName, r.URL.Path}
	return renderTemplate(w, s.templates[tmplDeviceSuccess], data)
}

func (s *WebSite) RenderLogin(r *http.Request, w http.ResponseWriter, connectors []ConnectorInfo) error {
	sort.Sort(byName(connectors))
	data := struct {
		Connectors []ConnectorInfo
		ReqPath    string
	}{connectors, r.URL.Path}
	return renderTemplate(w, s.templates[tmplLogin], data)
}

func (s *WebSite) RenderPassword(r *http.Request, w http.ResponseWriter, postURL, lastUsername, usernamePrompt string, lastWasInvalid bool, backLink string) error {
	if lastWasInvalid {
		w.WriteHeader(http.StatusUnauthorized)
	}
	data := struct {
		PostURL        string
		BackLink       string
		Username       string
		UsernamePrompt string
		Invalid        bool
		ReqPath        string
	}{postURL, backLink, lastUsername, usernamePrompt, lastWasInvalid, r.URL.Path}
	return renderTemplate(w, s.templates[tmplPassword], data)
}

func (s *WebSite) RenderApproval(r *http.Request, w http.ResponseWriter, authReqID, username, clientName string, scopes []string) error {
	accesses := []string{}
	for _, scope := range scopes {
		access, ok := scopeDescriptions[scope]
		if ok {
			accesses = append(accesses, access)
		}
	}
	sort.Strings(accesses)
	data := struct {
		User      string
		Client    string
		AuthReqID string
		Scopes    []string
		ReqPath   string
	}{username, clientName, authReqID, accesses, r.URL.Path}
	return renderTemplate(w, s.templates[tmplApproval], data)
}

func (s *WebSite) RenderOOB(r *http.Request, w http.ResponseWriter, code string) error {
	data := struct {
		Code    string
		ReqPath string
	}{code, r.URL.Path}
	return renderTemplate(w, s.templates[tmplOOB], data)
}

func (s *WebSite) RenderError(r *http.Request, w http.ResponseWriter, errCode int, errMsg string) error {
	w.WriteHeader(errCode)
	data := struct {
		ErrType string
		ErrMsg  string
		ReqPath string
	}{http.StatusText(errCode), errMsg, r.URL.Path}
	if err := s.templates[tmplError].Execute(w, data); err != nil {
		return fmt.Errorf("rendering template %s failed: %s", s.templates[tmplError].Name(), err)
	}
	return nil
}

// small io.Writer utility to determine if executing the template wrote to the underlying response writer.
type writeRecorder struct {
	wrote bool
	w     io.Writer
}

func (w *writeRecorder) Write(p []byte) (n int, err error) {
	w.wrote = true
	return w.w.Write(p)
}

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, data any) error {
	wr := &writeRecorder{w: w}
	if err := tmpl.Execute(wr, data); err != nil {
		if !wr.wrote {
			// TODO(ericchiang): replace with better internal server error.
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return fmt.Errorf("rendering template %s failed: %s", tmpl.Name(), err)
	}
	return nil
}

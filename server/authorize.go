package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sutext.github.io/entry/model"
	"sutext.github.io/entry/xerr"
	"sutext.github.io/suid/guid"
)

func (s *server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	respType := ResponseType(r.FormValue("response_type"))
	if respType.String() == "" {
		http.Error(w, "response_type is empty", http.StatusBadRequest)
		return
	}
	if respType != ResponseTypeCode {
		http.Error(w, "only support response type code", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/#/approve?"+r.Form.Encode(), http.StatusFound)
}

type PreviewResponse struct {
	Scopes     []string `json:"scopes"`
	ClientID   string   `json:"client_id"`
	ClientName string   `json:"client_name"`
	ClientLogo string   `json:"client_logo"`
}

func (s *server) handleAuthorizePreview(w http.ResponseWriter, r *http.Request) {
	_, err := s.ensureLoggedIn(r)
	if err != nil {
		http.Error(w, "failed to ensure logged in: "+err.Error(), http.StatusUnauthorized)
		return
	}
	req, err := s.validateAuthorizeRequest(r)
	if err != nil {
		http.Error(w, "failed to validate authorize request: "+err.Error(), http.StatusBadRequest)
		return
	}
	client, err := s.db.GetClient(r.Context(), req.ClientID)
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to get client %s: %s", req.ClientID, err.Error()),
			http.StatusInternalServerError,
		)
		return
	}
	if err := s.validateClientSettings(client, req, r); err != nil {
		http.Error(
			w,
			fmt.Sprintf("failed to validate client settings: %s", err.Error()),
			http.StatusBadRequest,
		)
		return
	}
	resp := PreviewResponse{
		Scopes:     strings.Split(req.Scope, " "),
		ClientID:   req.ClientID,
		ClientName: client.Name,
		ClientLogo: client.LogoURL,
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	if err := e.Encode(resp); err != nil {
		http.Error(w, "failed to encode preview response: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
func (s *server) handleAuthorizeApprove(w http.ResponseWriter, r *http.Request) {
	userID, err := s.ensureLoggedIn(r)
	if err != nil {
		http.Redirect(w, r, "/#/login?"+r.Form.Encode(), http.StatusFound)
		return
	}
	req, err := s.validateAuthorizeRequest(r)
	if err != nil {
		http.Error(w, "failed to validate authorize request: "+err.Error(), http.StatusBadRequest)
		return
	}

	req.UserID = userID
	code := guid.New().String()
	s.codeCache.Set(code, req, time.Minute*10)
	u, err := url.Parse(req.RedirectURI)
	if err != nil {
		s.redirectError(w, req, err)
		return
	}
	q := u.Query()
	if req.State != "" {
		q.Set("state", req.State)
	}
	q.Set("code", code)
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusFound)
}

func (s *server) validateClientSettings(client *model.Client, req *AuthorizeRequest, r *http.Request) error {
	if !client.Public {
		if !s.checkTrustedPeer(client.TrustedPeers, r.RemoteAddr) {
			return xerr.ErrUnauthorizedClient
		}
	}
	if !client.Scopes.Contains(req.Scope) {
		return xerr.ErrInvalidScope
	}
	if !client.RedirectURIs.Contains(req.RedirectURI) {
		return xerr.ErrInvalidRedirectURI
	}
	return nil
}
func (s *server) checkTrustedPeer(trustedPeers []string, remoteAddr string) bool {
	for _, tp := range trustedPeers {
		if tp == remoteAddr {
			return true
		}
	}
	return false
}
func (s *server) validateAuthorizeRequest(r *http.Request) (*AuthorizeRequest, error) {
	redirectURI := r.FormValue("redirect_uri")
	clientID := r.FormValue("client_id")
	if !(r.Method == "GET" || r.Method == "POST") ||
		clientID == "" {
		return nil, xerr.ErrInvalidRequest
	}

	resType := ResponseType(r.FormValue("response_type"))
	if resType.String() == "" {
		return nil, xerr.ErrUnsupportedResponseType
	} else if allowed := s.checkResponseType(resType); !allowed {
		return nil, xerr.ErrUnauthorizedClient
	}

	cc := r.FormValue("code_challenge")
	if cc == "" {
		return nil, xerr.ErrCodeChallengeRquired
	}
	if cc != "" && (len(cc) < 43 || len(cc) > 128) {
		return nil, xerr.ErrInvalidCodeChallengeLen
	}

	ccm := CodeChallengeMethod(r.FormValue("code_challenge_method"))
	// set default
	if ccm == "" {
		ccm = CodeChallengePlain
	}
	if ccm != "" && !s.checkCodeChallengeMethod(ccm) {
		return nil, xerr.ErrUnsupportedCodeChallengeMethod
	}

	req := &AuthorizeRequest{
		RedirectURI:         redirectURI,
		ResponseType:        resType,
		ClientID:            clientID,
		State:               r.FormValue("state"),
		Scope:               r.FormValue("scope"),
		CodeChallenge:       cc,
		CodeChallengeMethod: ccm,
	}
	return req, nil
}
func (s *server) getRedirectURI(req *AuthorizeRequest, data map[string]any) (string, error) {
	u, err := url.Parse(req.RedirectURI)
	if err != nil {
		return "", err
	}
	q := u.Query()
	if req.State != "" {
		q.Set("state", req.State)
	}
	for k, v := range data {
		q.Set(k, fmt.Sprint(v))
	}
	switch req.ResponseType {
	case ResponseTypeCode:
		u.RawQuery = q.Encode()
	case ResponseTypeToken:
		u.RawQuery = ""
		fragment, err := url.QueryUnescape(q.Encode())
		if err != nil {
			return "", err
		}
		u.Fragment = fragment
	}

	return u.String(), nil
}
func (s *server) getAuthorizeToken(ctx context.Context, req *AuthorizeRequest) (*model.TokenInfo, error) {
	// check the client allows the grant type
	// if fn := s.ClientAuthorizedHandler; fn != nil {
	// 	gt := ResponseTypeCode
	// 	if req.ResponseType == ResponseTypeToken {
	// 		gt = ResponseTypeToken
	// 	}

	// 	allowed, err := fn(req.ClientID, gt)
	// 	if err != nil {
	// 		return nil, err
	// 	} else if !allowed {
	// 		return nil, err.ErrUnauthorizedClient
	// 	}
	// }

	tgr := &TokenGenerateRequest{
		// ClientID:       req.ClientID,
		// UserID:         req.UserID,
		// RedirectURI:    req.RedirectURI,
		// Scope:          req.Scope,
		// AccessTokenExp: req.AccessTokenExp,
		// Request:        req.Request,
	}

	// check the client allows the authorized scope
	// if fn := s.ClientScopeHandler; fn != nil {
	// 	allowed, err := fn(tgr)
	// 	if err != nil {
	// 		return nil, err
	// 	} else if !allowed {
	// 		return nil, err.ErrInvalidScope
	// 	}
	// }

	tgr.CodeChallenge = req.CodeChallenge
	tgr.CodeChallengeMethod = req.CodeChallengeMethod

	return s.generateAuthToken(ctx, req.ResponseType, tgr)
}
func (s *server) generateAuthToken(ctx context.Context, rt ResponseType, tgr *TokenGenerateRequest) (ti *model.TokenInfo, err error) {
	cli, err := s.db.GetClient(ctx, tgr.ClientID)
	if err != nil {
		return ti, err
	} else if tgr.RedirectURI != "" {
		if !cli.RedirectURIs.Contains(tgr.RedirectURI) {
			return ti, xerr.ErrInvalidRedirectURI
		}
	}
	return s.db.CreateTokenInfo(ctx)
	// ti := s.db.CreateTokenInfo(ctx)
	// if m.extractExtension != nil {
	// 	m.extractExtension(tgr, ti)
	// }
	// ti.SetClientID(tgr.ClientID)
	// ti.SetUserID(tgr.UserID)
	// ti.SetRedirectURI(tgr.RedirectURI)
	// ti.SetScope(tgr.Scope)

	// createAt := time.Now()
	// td := &oauth2.GenerateBasic{
	// 	Client:    cli,
	// 	UserID:    tgr.UserID,
	// 	CreateAt:  createAt,
	// 	TokenInfo: ti,
	// 	Request:   tgr.Request,
	// }
	// switch rt {
	// case oauth2.Code:
	// 	codeExp := m.codeExp
	// 	if codeExp == 0 {
	// 		codeExp = DefaultCodeExp
	// 	}
	// 	ti.SetCodeCreateAt(createAt)
	// 	ti.SetCodeExpiresIn(codeExp)
	// 	if exp := tgr.AccessTokenExp; exp > 0 {
	// 		ti.SetAccessExpiresIn(exp)
	// 	}
	// 	if tgr.CodeChallenge != "" {
	// 		ti.SetCodeChallenge(tgr.CodeChallenge)
	// 		ti.SetCodeChallengeMethod(tgr.CodeChallengeMethod)
	// 	}

	// 	tv, err := m.authorizeGenerate.Token(ctx, td)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	ti.SetCode(tv)
	// case oauth2.Token:
	// 	// set access token expires
	// 	icfg := m.grantConfig(oauth2.Implicit)
	// 	aexp := icfg.AccessTokenExp
	// 	if exp := tgr.AccessTokenExp; exp > 0 {
	// 		aexp = exp
	// 	}
	// 	ti.SetAccessCreateAt(createAt)
	// 	ti.SetAccessExpiresIn(aexp)

	// 	if icfg.IsGenerateRefresh {
	// 		ti.SetRefreshCreateAt(createAt)
	// 		ti.SetRefreshExpiresIn(icfg.RefreshTokenExp)
	// 	}

	// 	tv, rv, err := m.accessGenerate.Token(ctx, td, icfg.IsGenerateRefresh)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	ti.SetAccess(tv)

	// 	if rv != "" {
	// 		ti.SetRefresh(rv)
	// 	}
	// }

	// err = m.tokenStore.Create(ctx, ti)
	// if err != nil {
	// 	return nil, err
	// }
	// return ti, nil
}

// GetAuthorizeData get authorization response data
func (s *server) getAuthorizeData(rt ResponseType, ti *model.TokenInfo) map[string]interface{} {
	if rt == ResponseTypeCode {
		return map[string]interface{}{
			"code": ti.Code,
		}
	}
	return s.getTokenData(ti)
}
func (s *server) getTokenData(ti *model.TokenInfo) map[string]interface{} {
	data := map[string]interface{}{
		"access_token": ti.Access,
		"token_type":   "access_token", //s.Config.TokenType,
		"expires_in":   int64(ti.AccessExpiresIn / time.Second),
	}

	if scope := ti.Scope; scope != "" {
		data["scope"] = scope
	}

	if refresh := ti.Refresh; refresh != "" {
		data["refresh_token"] = refresh
	}

	// if fn := s.ExtensionFieldsHandler; fn != nil {
	// 	ext := fn(ti)
	// 	for k, v := range ext {
	// 		if _, ok := data[k]; ok {
	// 			continue
	// 		}
	// 		data[k] = v
	// 	}
	// }
	return data
}

func (s *server) redirectError(w http.ResponseWriter, req *AuthorizeRequest, err error) {
	data, _, _ := s.getErrorData(err)
	s.redirect(w, req, data)
}

func (s *server) redirect(w http.ResponseWriter, req *AuthorizeRequest, data map[string]any) {
	uri, err := s.getRedirectURI(req, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Location", uri)
	w.WriteHeader(302)
}

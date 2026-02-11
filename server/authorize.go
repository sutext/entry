package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-session/session/v3"
	"sutext.github.io/entry/model"
	"sutext.github.io/entry/xerr"
)

func (s *server) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	err := dumpRequest(w, "Request", r)
	if err != nil {
		http.Error(w, "failed to dump request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := r.Context()
	store, err := session.Start(ctx, w, r)
	if err != nil {
		http.Error(w, "failed to start session: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Form == nil {
		r.ParseForm()
	}
	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		store.Set("ReturnUri", r.Form)
		store.Save()
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	userID := uid.(string)
	store.Delete("LoggedInUserID")
	if v, ok := store.Get("ReturnUri"); ok {
		form := v.(url.Values)
		store.Delete("ReturnUri")
		r.Form = form
	}
	store.Save()

	req, err := s.validateAuthorizeRequest(r)
	if err != nil {
		s.redirectError(w, req, err)
		return
	}
	req.UserID = userID

	// // specify the scope of authorization
	// if fn := s.AuthorizeScopeHandler; fn != nil {
	// 	scope, err := fn(w, r)
	// 	if err != nil {
	// 		return err
	// 	} else if scope != "" {
	// 		req.Scope = scope
	// 	}
	// }

	// // specify the expiration time of access token
	// if fn := s.AccessTokenExpHandler; fn != nil {
	// 	exp, err := fn(w, r)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	req.AccessTokenExp = exp
	// }

	ti, err := s.getAuthorizeToken(ctx, req)
	if err != nil {
		s.redirectError(w, req, err)
		return
	}

	// If the redirect URI is empty, the default domain provided by the client is used.
	if req.RedirectURI == "" {
		client, err := s.db.GetClient(ctx, req.ClientID)
		if err != nil {
			s.redirectError(w, req, err)
			return
		}
		req.RedirectURI = client.RedirectURIs[0]
	}
	s.redirect(w, req, s.getAuthorizeData(req.ResponseType, ti))
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
	if cc == "" && s.forcePKCE {
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
		Request:             r,
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
		ClientID:       req.ClientID,
		UserID:         req.UserID,
		RedirectURI:    req.RedirectURI,
		Scope:          req.Scope,
		AccessTokenExp: req.AccessTokenExp,
		Request:        req.Request,
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

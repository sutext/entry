package server

import "net/http"

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/url"
// 	"time"

//	"sutext.github.io/entry/model"
//	"sutext.github.io/entry/xerr"
//
// )
func (s *server) handleToken(w http.ResponseWriter, r *http.Request) {

}

// func (s *server) handleTokenRequest(w http.ResponseWriter, r *http.Request) error {
// 	ctx := r.Context()

// 	gt, tgr, err := s.validationTokenRequest(r)
// 	if err != nil {
// 		return s.tokenError(w, err)
// 	}

// 	ti, err := s.getAccessToken(ctx, gt, tgr)
// 	if err != nil {
// 		return s.tokenError(w, err)
// 	}

// 	return s.token(w, s.getTokenData(ti), nil)
// }
// func (s *server) getAccessToken(ctx context.Context, gt GrantType, tgr *TokenGenerateRequest) (*model.TokenInfo,
// 	error) {
// 	if allowed := s.checkGrantType(gt); !allowed {
// 		return nil, xerr.ErrUnauthorizedClient
// 	}

// 	// if fn := s.ClientAuthorizedHandler; fn != nil {
// 	// 	allowed, err := fn(tgr.ClientID, gt)
// 	// 	if err != nil {
// 	// 		return nil, err
// 	// 	} else if !allowed {
// 	// 		return nil, xerr.ErrUnauthorizedClient
// 	// 	}
// 	// }

// 	switch gt {
// 	case AuthorizationCode:
// 		ti, err := s.GenerateAccessToken(ctx, gt, tgr)
// 		if err != nil {
// 			switch err {
// 			case xerr.ErrInvalidAuthorizeCode, xerr.ErrInvalidCodeChallenge, xerr.ErrMissingCodeChallenge:
// 				return nil, xerr.ErrInvalidGrant
// 			case xerr.ErrInvalidClient:
// 				return nil, xerr.ErrInvalidClient
// 			default:
// 				return nil, err
// 			}
// 		}
// 		return ti, nil
// 	case PasswordCredentials, ClientCredentials:
// 		// if fn := s.ClientScopeHandler; fn != nil {
// 		// 	allowed, err := fn(tgr)
// 		// 	if err != nil {
// 		// 		return nil, err
// 		// 	} else if !allowed {
// 		// 		return nil, xerr.ErrInvalidScope
// 		// 	}
// 		// }
// 		return s.GenerateAccessToken(ctx, gt, tgr)
// 	case Refreshing:
// 		// check scope
// 		if scopeFn := s.RefreshingScopeHandler; tgr.Scope != "" && scopeFn != nil {
// 			rti, err := s.LoadRefreshToken(ctx, tgr.Refresh)
// 			if err != nil {
// 				if err == xerr.ErrInvalidRefreshToken || err == xerr.ErrExpiredRefreshToken {
// 					return nil, xerr.ErrInvalidGrant
// 				}
// 				return nil, err
// 			}

// 			allowed, err := scopeFn(tgr, rti.GetScope())
// 			if err != nil {
// 				return nil, err
// 			} else if !allowed {
// 				return nil, xerr.ErrInvalidScope
// 			}
// 		}

// 		if validationFn := s.RefreshingValidationHandler; validationFn != nil {
// 			rti, err := s.LoadRefreshToken(ctx, tgr.Refresh)
// 			if err != nil {
// 				if err == xerr.ErrInvalidRefreshToken || err == xerr.ErrExpiredRefreshToken {
// 					return nil, xerr.ErrInvalidGrant
// 				}
// 				return nil, err
// 			}
// 			allowed, err := validationFn(rti)
// 			if err != nil {
// 				return nil, err
// 			} else if !allowed {
// 				return nil, xerr.ErrInvalidScope
// 			}
// 		}

// 		ti, err := s.RefreshAccessToken(ctx, tgr)
// 		if err != nil {
// 			if err == xerr.ErrInvalidRefreshToken || err == xerr.ErrExpiredRefreshToken {
// 				return nil, xerr.ErrInvalidGrant
// 			}
// 			return nil, err
// 		}
// 		return ti, nil
// 	}

// 	return nil, xerr.ErrUnsupportedGrantType
// }
// func (s *server) validationTokenRequest(r *http.Request) (GrantType, *TokenGenerateRequest, error) {
// 	// if v := r.Method; !(v == "POST" ||
// 	// 	(s.AllowGetAccessRequest && v == "GET")) {
// 	// 	return "", nil, xerr.ErrInvalidRequest
// 	// }

// 	gt := GrantType(r.FormValue("grant_type"))
// 	if gt.String() == "" {
// 		return "", nil, xerr.ErrUnsupportedGrantType
// 	}

// 	clientID, clientSecret, ok := r.BasicAuth()
// 	if !ok {
// 		return "", nil, xerr.ErrInvalidClient
// 	}

// 	tgr := &TokenGenerateRequest{
// 		ClientID:     clientID,
// 		ClientSecret: clientSecret,
// 		Request:      r,
// 	}

// 	switch gt {
// 	case AuthorizationCode:
// 		tgr.RedirectURI = r.FormValue("redirect_uri")
// 		tgr.Code = r.FormValue("code")
// 		if tgr.RedirectURI == "" ||
// 			tgr.Code == "" {
// 			return "", nil, xerr.ErrInvalidRequest
// 		}
// 		tgr.CodeVerifier = r.FormValue("code_verifier")
// 		if s.forcePKCE && tgr.CodeVerifier == "" {
// 			return "", nil, xerr.ErrInvalidRequest
// 		}
// 	case PasswordCredentials:
// 		tgr.Scope = r.FormValue("scope")
// 		username, password := r.FormValue("username"), r.FormValue("password")
// 		if username == "" || password == "" {
// 			return "", nil, xerr.ErrInvalidRequest
// 		}

// 		userID, err := s.passwordAuthorizationHandler(r.Context(), clientID, username, password)
// 		if err != nil {
// 			return "", nil, err
// 		} else if userID == "" {
// 			return "", nil, xerr.ErrInvalidGrant
// 		}
// 		tgr.UserID = userID
// 	case ClientCredentials:
// 		tgr.Scope = r.FormValue("scope")
// 	case Refreshing:
// 		refresh, err := s.refreshTokenResolveHandler(r)
// 		if err != nil {
// 			return "", nil, err
// 		}
// 		tgr.Scope = r.FormValue("scope")
// 		tgr.Refresh = refresh
// 	}
// 	return gt, tgr, nil
// }

// func (s *server) tokenError(w http.ResponseWriter, err error) error {
// 	data, statusCode, header := s.getErrorData(err)
// 	return s.token(w, data, header, statusCode)
// }

// func (s *server) token(w http.ResponseWriter, data map[string]any, header http.Header, statusCode ...int) error {
// 	// if fn := s.ResponseTokenHandler; fn != nil {
// 	// 	return fn(w, data, header, statusCode...)
// 	// }
// 	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
// 	w.Header().Set("Cache-Control", "no-store")
// 	w.Header().Set("Pragma", "no-cache")

// 	for key := range header {
// 		w.Header().Set(key, header.Get(key))
// 	}

// 	status := http.StatusOK
// 	if len(statusCode) > 0 && statusCode[0] > 0 {
// 		status = statusCode[0]
// 	}

// 	w.WriteHeader(status)
// 	return json.NewEncoder(w).Encode(data)
// }
// func (m *server) GenerateAuthToken(ctx context.Context, rt ResponseType, tgr *TokenGenerateRequest) (*model.TokenInfo, error) {
// 	cli, err := m.db.GetClient(ctx, tgr.ClientID)
// 	if err != nil {
// 		return nil, err
// 	} else if tgr.RedirectURI != "" {
// 		if err := m.validateURI(cli.GetDomain(), tgr.RedirectURI); err != nil {
// 			return nil, err
// 		}
// 	}

// 	ti := models.NewToken()
// 	if m.extractExtension != nil {
// 		m.extractExtension(tgr, ti)
// 	}
// 	ti.SetClientID(tgr.ClientID)
// 	ti.SetUserID(tgr.UserID)
// 	ti.SetRedirectURI(tgr.RedirectURI)
// 	ti.SetScope(tgr.Scope)

// 	createAt := time.Now()
// 	td := &oauth2.GenerateBasic{
// 		Client:    cli,
// 		UserID:    tgr.UserID,
// 		CreateAt:  createAt,
// 		TokenInfo: ti,
// 		Request:   tgr.Request,
// 	}
// 	switch rt {
// 	case oauth2.Code:
// 		codeExp := m.codeExp
// 		if codeExp == 0 {
// 			codeExp = DefaultCodeExp
// 		}
// 		ti.SetCodeCreateAt(createAt)
// 		ti.SetCodeExpiresIn(codeExp)
// 		if exp := tgr.AccessTokenExp; exp > 0 {
// 			ti.SetAccessExpiresIn(exp)
// 		}
// 		if tgr.CodeChallenge != "" {
// 			ti.SetCodeChallenge(tgr.CodeChallenge)
// 			ti.SetCodeChallengeMethod(tgr.CodeChallengeMethod)
// 		}

// 		tv, err := m.authorizeGenerate.Token(ctx, td)
// 		if err != nil {
// 			return nil, err
// 		}
// 		ti.SetCode(tv)
// 	case oauth2.Token:
// 		// set access token expires
// 		icfg := m.grantConfig(oauth2.Implicit)
// 		aexp := icfg.AccessTokenExp
// 		if exp := tgr.AccessTokenExp; exp > 0 {
// 			aexp = exp
// 		}
// 		ti.SetAccessCreateAt(createAt)
// 		ti.SetAccessExpiresIn(aexp)

// 		if icfg.IsGenerateRefresh {
// 			ti.SetRefreshCreateAt(createAt)
// 			ti.SetRefreshExpiresIn(icfg.RefreshTokenExp)
// 		}

// 		tv, rv, err := m.accessGenerate.Token(ctx, td, icfg.IsGenerateRefresh)
// 		if err != nil {
// 			return nil, err
// 		}
// 		ti.SetAccess(tv)

// 		if rv != "" {
// 			ti.SetRefresh(rv)
// 		}
// 	}

// 	err = m.tokenStore.Create(ctx, ti)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ti, nil
// }

// // get authorization code data
// func (m *server) getAuthorizationCode(ctx context.Context, code string) (*model.TokenInfo, error) {
// 	ti, err := m.tokenStore.GetByCode(ctx, code)
// 	if err != nil {
// 		return nil, err
// 	} else if ti == nil || ti.GetCode() != code || ti.GetCodeCreateAt().Add(ti.GetCodeExpiresIn()).Before(time.Now()) {
// 		err = xerr.ErrInvalidAuthorizeCode
// 		return nil, xerr.ErrInvalidAuthorizeCode
// 	}
// 	return ti, nil
// }

// // delete authorization code data
// func (m *server) delAuthorizationCode(ctx context.Context, code string) error {
// 	return m.tokenStore.RemoveByCode(ctx, code)
// }

// // get and delete authorization code data
// func (m *server) getAndDelAuthorizationCode(ctx context.Context, tgr *TokenGenerateRequest) (*model.TokenInfo, error) {
// 	code := tgr.Code
// 	ti, err := m.getAuthorizationCode(ctx, code)
// 	if err != nil {
// 		return nil, err
// 	} else if ti.GetClientID() != tgr.ClientID {
// 		return nil, xerr.ErrInvalidAuthorizeCode
// 	} else if codeURI := ti.GetRedirectURI(); codeURI != "" && codeURI != tgr.RedirectURI {
// 		return nil, xerr.ErrInvalidAuthorizeCode
// 	}

// 	err = m.delAuthorizationCode(ctx, code)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ti, nil
// }

// func (m *server) validateCodeChallenge(ti model.TokenInfo, ver string) error {
// 	cc := ti.GetCodeChallenge()
// 	// early return
// 	if cc == "" && ver == "" {
// 		return nil
// 	}
// 	if cc == "" {
// 		return xerr.ErrMissingCodeVerifier
// 	}
// 	if ver == "" {
// 		return xerr.ErrMissingCodeVerifier
// 	}
// 	ccm := ti.GetCodeChallengeMethod()
// 	if ccm.String() == "" {
// 		ccm = oauth2.CodeChallengePlain
// 	}
// 	if !ccm.Validate(cc, ver) {
// 		return xerr.ErrInvalidCodeChallenge
// 	}
// 	return nil
// }

// // GenerateAccessToken generate the access token
// func (m *server) GenerateAccessToken(ctx context.Context, gt GrantType, tgr *TokenGenerateRequest) (*model.TokenInfo, error) {
// 	cli, err := m.GetClient(ctx, tgr.ClientID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if cliPass, ok := cli.(oauth2.ClientPasswordVerifier); ok {
// 		if !cliPass.VerifyPassword(tgr.ClientSecret) {
// 			return nil, xerr.ErrInvalidClient
// 		}
// 	} else if len(cli.GetSecret()) > 0 && tgr.ClientSecret != cli.GetSecret() {
// 		return nil, xerr.ErrInvalidClient
// 	}
// 	if tgr.RedirectURI != "" {
// 		if err := m.validateURI(cli.GetDomain(), tgr.RedirectURI); err != nil {
// 			return nil, err
// 		}
// 	}

// 	if gt == ClientCredentials && cli.IsPublic() == true {
// 		return nil, xerr.ErrInvalidClient
// 	}

// 	var extension url.Values

// 	if gt == AuthorizationCode {
// 		ti, err := m.getAndDelAuthorizationCode(ctx, tgr)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if err := m.validateCodeChallenge(ti, tgr.CodeVerifier); err != nil {
// 			return nil, err
// 		}
// 		tgr.UserID = ti.GetUserID()
// 		tgr.Scope = ti.GetScope()
// 		if exp := ti.GetAccessExpiresIn(); exp > 0 {
// 			tgr.AccessTokenExp = exp
// 		}
// 		if eti, ok := ti.(oauth2.ExtendableTokenInfo); ok {
// 			extension = eti.GetExtension()
// 		}
// 	}

// 	ti := models.NewToken()
// 	ti.SetExtension(extension)
// 	if m.extractExtension != nil {
// 		m.extractExtension(tgr, ti)
// 	}
// 	ti.SetClientID(tgr.ClientID)
// 	ti.SetUserID(tgr.UserID)
// 	ti.SetRedirectURI(tgr.RedirectURI)
// 	ti.SetScope(tgr.Scope)

// 	createAt := time.Now()
// 	ti.SetAccessCreateAt(createAt)

// 	// set access token expires
// 	gcfg := m.grantConfig(gt)
// 	aexp := gcfg.AccessTokenExp
// 	if exp := tgr.AccessTokenExp; exp > 0 {
// 		aexp = exp
// 	}
// 	ti.SetAccessExpiresIn(aexp)
// 	if gcfg.IsGenerateRefresh {
// 		ti.SetRefreshCreateAt(createAt)
// 		ti.SetRefreshExpiresIn(gcfg.RefreshTokenExp)
// 	}

// 	td := &oauth2.GenerateBasic{
// 		Client:    cli,
// 		UserID:    tgr.UserID,
// 		CreateAt:  createAt,
// 		TokenInfo: ti,
// 		Request:   tgr.Request,
// 	}

// 	av, rv, err := m.accessGenerate.Token(ctx, td, gcfg.IsGenerateRefresh)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ti.SetAccess(av)

// 	if rv != "" {
// 		ti.SetRefresh(rv)
// 	}

// 	err = m.tokenStore.Create(ctx, ti)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return ti, nil
// }

// // RefreshAccessToken refreshing an access token
// func (m *server) RefreshAccessToken(ctx context.Context, tgr *TokenGenerateRequest) (*model.TokenInfo, error) {
// 	ti, err := m.LoadRefreshToken(ctx, tgr.Refresh)
// 	if err != nil {
// 		return nil, err
// 	}

// 	cli, err := m.GetClient(ctx, ti.GetClientID())
// 	if err != nil {
// 		return nil, err
// 	}

// 	oldAccess, oldRefresh := ti.GetAccess(), ti.GetRefresh()

// 	td := &oauth2.GenerateBasic{
// 		Client:    cli,
// 		UserID:    ti.GetUserID(),
// 		CreateAt:  time.Now(),
// 		TokenInfo: ti,
// 		Request:   tgr.Request,
// 	}

// 	rcfg := DefaultRefreshTokenCfg
// 	if v := m.rcfg; v != nil {
// 		rcfg = v
// 	}

// 	ti.SetAccessCreateAt(td.CreateAt)
// 	if v := rcfg.AccessTokenExp; v > 0 {
// 		ti.SetAccessExpiresIn(v)
// 	}

// 	if v := rcfg.RefreshTokenExp; v > 0 {
// 		ti.SetRefreshExpiresIn(v)
// 	}

// 	if rcfg.IsResetRefreshTime {
// 		ti.SetRefreshCreateAt(td.CreateAt)
// 	}

// 	if scope := tgr.Scope; scope != "" {
// 		ti.SetScope(scope)
// 	}

// 	tv, rv, err := m.accessGenerate.Token(ctx, td, rcfg.IsGenerateRefresh)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ti.SetAccess(tv)
// 	if rv != "" {
// 		ti.SetRefresh(rv)
// 	}

// 	if err := m.tokenStore.Create(ctx, ti); err != nil {
// 		return nil, err
// 	}

// 	if rcfg.IsRemoveAccess {
// 		// remove the old access token
// 		if err := m.tokenStore.RemoveByAccess(ctx, oldAccess); err != nil {
// 			return nil, err
// 		}
// 	}

// 	if rcfg.IsRemoveRefreshing && rv != "" {
// 		// remove the old refresh token
// 		if err := m.tokenStore.RemoveByRefresh(ctx, oldRefresh); err != nil {
// 			return nil, err
// 		}
// 	}

// 	if rv == "" {
// 		ti.SetRefresh("")
// 		ti.SetRefreshCreateAt(time.Now())
// 		ti.SetRefreshExpiresIn(0)
// 	}

// 	return ti, nil
// }

// // RemoveAccessToken use the access token to delete the token information
// func (m *server) RemoveAccessToken(ctx context.Context, access string) error {
// 	if access == "" {
// 		return xerr.ErrInvalidAccessToken
// 	}
// 	return m.tokenStore.RemoveByAccess(ctx, access)
// }

// // RemoveRefreshToken use the refresh token to delete the token information
// func (m *server) RemoveRefreshToken(ctx context.Context, refresh string) error {
// 	if refresh == "" {
// 		return xerr.ErrInvalidAccessToken
// 	}
// 	return m.tokenStore.RemoveByRefresh(ctx, refresh)
// }

// // LoadAccessToken according to the access token for corresponding token information
// func (m *server) LoadAccessToken(ctx context.Context, access string) (*model.TokenInfo, error) {
// 	if access == "" {
// 		return nil, xerr.ErrInvalidAccessToken
// 	}

// 	ct := time.Now()
// 	ti, err := m.tokenStore.GetByAccess(ctx, access)
// 	if err != nil {
// 		return nil, err
// 	} else if ti == nil || ti.GetAccess() != access {
// 		return nil, xerr.ErrInvalidAccessToken
// 	} else if ti.GetRefresh() != "" && ti.GetRefreshExpiresIn() != 0 &&
// 		ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(ct) {
// 		return nil, xerr.ErrExpiredRefreshToken
// 	} else if ti.GetAccessExpiresIn() != 0 &&
// 		ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Before(ct) {
// 		return nil, xerr.ErrExpiredAccessToken
// 	}
// 	return ti, nil
// }

// // LoadRefreshToken according to the refresh token for corresponding token information
// func (m *server) LoadRefreshToken(ctx context.Context, refresh string) (*model.TokenInfo, error) {
// 	if refresh == "" {
// 		return nil, xerr.ErrInvalidRefreshToken
// 	}

// 	ti, err := m.tokenStore.GetByRefresh(ctx, refresh)
// 	if err != nil {
// 		return nil, err
// 	} else if ti == nil || ti.GetRefresh() != refresh {
// 		return nil, xerr.ErrInvalidRefreshToken
// 	} else if ti.GetRefreshExpiresIn() != 0 && // refresh token set to not expire
// 		ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(time.Now()) {
// 		return nil, xerr.ErrExpiredRefreshToken
// 	}
// 	return ti, nil
// }

package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-jose/go-jose/v4/jwt"
	"sutext.github.io/entry/model"
	"sutext.github.io/entry/xerr"
	"sutext.github.io/suid/guid"
)

func (s *server) handleToken(w http.ResponseWriter, r *http.Request) {
	gtype := GrantType(r.FormValue("grant_type"))
	if gtype.String() == "" {
		http.Error(w, "grant_type is required", http.StatusBadRequest)
		return
	}
	if allowed := s.checkGrantType(gtype); !allowed {
		http.Error(w, "grant_type not supported", http.StatusBadRequest)
		return
	}
	switch gtype {
	case AuthorizationCode:
		data, err := s.validateCodeGrant(r)
		if err != nil {
			s.tokenError(w, err)
			return
		}
		s.token(w, data, nil, http.StatusOK)
	case Refreshing:
		data, err := s.validateRefreshGrant(r)
		if err != nil {
			s.tokenError(w, err)
			return
		}
		s.token(w, data, nil, http.StatusOK)
	case ClientCredentials:
		data, err := s.validateClientCredentialsGrant(r)
		if err != nil {
			s.tokenError(w, err)
			return
		}
		s.token(w, data, nil, http.StatusOK)
	default:
		http.Error(w, "grant_type not supported", http.StatusBadRequest)
		return
	}
}
func (s *server) validateCodeGrant(r *http.Request) (data map[string]any, err error) {
	ctx := r.Context()
	redirectURI := r.FormValue("redirect_uri")
	if redirectURI == "" {
		return data, xerr.ErrInvalidRedirectURI
	}
	code := r.FormValue("code")
	if code == "" {
		return data, xerr.ErrInvalidAuthorizeCode
	}
	codeVerifier := r.FormValue("code_verifier")
	if codeVerifier == "" {
		return data, xerr.ErrMissingCodeChallenge
	}
	clientID := r.FormValue("client_id")
	if clientID == "" {
		return data, xerr.ErrInvalidClient
	}
	client, err := s.db.GetClient(ctx, clientID)
	if err != nil {
		return data, err
	}
	if !client.RedirectURIs.Contains(redirectURI) {
		return data, xerr.ErrInvalidRedirectURI
	}
	if !client.Public {
		clientSecret := r.FormValue("client_secret")
		if clientSecret == "" {
			return data, xerr.ErrUnauthorizedClient
		}
		if client.Secret != clientSecret {
			return data, xerr.ErrUnauthorizedClient
		}
		if !s.checkTrustedPeer(client.TrustedPeers, r.RemoteAddr) {
			return data, xerr.ErrUnauthorizedClient
		}
	}
	codeReq, err := s.codeCache.Get(code)
	if err != nil {
		return data, err
	}
	s.codeCache.Delete(code)
	if !codeReq.CodeChallengeMethod.Validate(codeReq.CodeChallenge, codeVerifier) {
		return data, xerr.ErrInvalidCodeChallenge
	}
	refreshToken := model.RefreshToken{
		ID:        guid.New(),
		ClientID:  clientID,
		UserID:    codeReq.UserID,
		Scope:     codeReq.Scope,
		ExpiryIn:  time.Now().Add(s.refreshTokenDuration),
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
	}
	if err = s.db.CreateRefresh(ctx, refreshToken); err != nil {
		return data, err
	}
	claims := jwt.Claims{
		Subject:  codeReq.UserID.String(),
		Audience: []string{clientID},
		Expiry:   jwt.NewNumericDate(time.Now().Add(s.accessTokenDuration)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	accessToken, err := jwt.Signed(s.signer).Claims(claims).Serialize()
	if err != nil {
		return data, err
	}
	data = map[string]any{
		"refresh_token": refreshToken.ID.String(),
		"expires_in":    claims.Expiry,
		"token_type":    "Bearer",
		"access_token":  accessToken,
	}
	return data, nil
}
func (s *server) validateRefreshGrant(r *http.Request) (data map[string]any, err error) {
	ctx := r.Context()
	refreshToken, err := guid.Parse(r.FormValue("refresh_token"))
	if err != nil {
		return data, xerr.ErrInvalidRefreshToken
	}
	clientID := r.FormValue("client_id")
	if clientID == "" {
		return data, xerr.ErrInvalidClient
	}
	client, err := s.db.GetClient(ctx, clientID)
	if err != nil {
		return data, err
	}
	if !client.Public {
		clientSecret := r.FormValue("client_secret")
		if clientSecret == "" {
			return data, xerr.ErrUnauthorizedClient
		}
		if client.Secret != clientSecret {
			return data, xerr.ErrUnauthorizedClient
		}
		if !s.checkTrustedPeer(client.TrustedPeers, r.RemoteAddr) {
			return data, xerr.ErrUnauthorizedClient
		}
	}
	rt, err := s.db.GetRefresh(ctx, refreshToken)
	if err != nil {
		return data, err
	}
	if rt.ClientID != clientID {
		return data, xerr.ErrUnauthorizedClient
	}
	if rt.ExpiryIn.Before(time.Now()) {
		return data, xerr.ErrExpiredRefreshToken
	}
	claims := jwt.Claims{
		Subject:  rt.UserID.String(),
		Audience: []string{clientID},
		Expiry:   jwt.NewNumericDate(time.Now().Add(s.accessTokenDuration)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	accessToken, err := jwt.Signed(s.signer).Claims(claims).Serialize()
	if err != nil {
		return data, err
	}
	data = map[string]any{
		"refresh_token": rt.ID.String(),
		"expires_in":    claims.Expiry,
		"token_type":    "Bearer",
		"access_token":  accessToken,
	}
	return data, nil
}
func (s *server) validateClientCredentialsGrant(r *http.Request) (data map[string]any, err error) {
	ctx := r.Context()
	clientID := r.FormValue("client_id")
	if clientID == "" {
		return data, xerr.ErrInvalidClient
	}
	client, err := s.db.GetClient(ctx, clientID)
	if err != nil {
		return data, err
	}
	clientSecret := r.FormValue("client_secret")
	if clientSecret == "" {
		return data, xerr.ErrUnauthorizedClient
	}
	if client.Secret != clientSecret {
		return data, xerr.ErrUnauthorizedClient
	}
	if !client.Public {
		if !s.checkTrustedPeer(client.TrustedPeers, r.RemoteAddr) {
			return data, xerr.ErrUnauthorizedClient
		}
	}
	claims := jwt.Claims{
		Subject:  clientID,
		Expiry:   jwt.NewNumericDate(time.Now().Add(s.accessTokenDuration)),
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	accessToken, err := jwt.Signed(s.signer).Claims(claims).Serialize()
	if err != nil {
		return data, err
	}
	data = map[string]any{
		"token_type":   "Bearer",
		"access_token": accessToken,
	}
	return data, nil
}

func (s *server) tokenError(w http.ResponseWriter, err error) error {
	data, statusCode, header := s.getErrorData(err)
	return s.token(w, data, header, statusCode)
}

func (s *server) token(w http.ResponseWriter, data map[string]any, header http.Header, statusCode ...int) error {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	for key := range header {
		w.Header().Set(key, header.Get(key))
	}
	status := http.StatusOK
	if len(statusCode) > 0 && statusCode[0] > 0 {
		status = statusCode[0]
	}
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

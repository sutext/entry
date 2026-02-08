package server

import (
	"context"
	"encoding/json"
	"net/http"

	"sutext.github.io/entry/model"
	"sutext.github.io/entry/xerr"
)

func (s *server) handleTokenRequest(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	gt, tgr, err := s.validationTokenRequest(r)
	if err != nil {
		return s.tokenError(w, err)
	}

	ti, err := s.getAccessToken(ctx, gt, tgr)
	if err != nil {
		return s.tokenError(w, err)
	}

	return s.token(w, s.getTokenData(ti), nil)
}
func (s *server) getAccessToken(ctx context.Context, gt GrantType, tgr *TokenGenerateRequest) (*model.TokenInfo,
	error) {
	if allowed := s.checkGrantType(gt); !allowed {
		return nil, xerr.ErrUnauthorizedClient
	}

	if fn := s.ClientAuthorizedHandler; fn != nil {
		allowed, err := fn(tgr.ClientID, gt)
		if err != nil {
			return nil, err
		} else if !allowed {
			return nil, xerr.ErrUnauthorizedClient
		}
	}

	switch gt {
	case AuthorizationCode:
		ti, err := s.Manager.GenerateAccessToken(ctx, gt, tgr)
		if err != nil {
			switch err {
			case xerr.ErrInvalidAuthorizeCode, xerr.ErrInvalidCodeChallenge, xerr.ErrMissingCodeChallenge:
				return nil, xerr.ErrInvalidGrant
			case xerr.ErrInvalidClient:
				return nil, xerr.ErrInvalidClient
			default:
				return nil, err
			}
		}
		return ti, nil
	case PasswordCredentials, ClientCredentials:
		if fn := s.ClientScopeHandler; fn != nil {
			allowed, err := fn(tgr)
			if err != nil {
				return nil, err
			} else if !allowed {
				return nil, xerr.ErrInvalidScope
			}
		}
		return s.Manager.GenerateAccessToken(ctx, gt, tgr)
	case Refreshing:
		// check scope
		if scopeFn := s.RefreshingScopeHandler; tgr.Scope != "" && scopeFn != nil {
			rti, err := s.Manager.LoadRefreshToken(ctx, tgr.Refresh)
			if err != nil {
				if err == xerr.ErrInvalidRefreshToken || err == xerr.ErrExpiredRefreshToken {
					return nil, xerr.ErrInvalidGrant
				}
				return nil, err
			}

			allowed, err := scopeFn(tgr, rti.GetScope())
			if err != nil {
				return nil, err
			} else if !allowed {
				return nil, xerr.ErrInvalidScope
			}
		}

		if validationFn := s.RefreshingValidationHandler; validationFn != nil {
			rti, err := s.Manager.LoadRefreshToken(ctx, tgr.Refresh)
			if err != nil {
				if err == xerr.ErrInvalidRefreshToken || err == xerr.ErrExpiredRefreshToken {
					return nil, xerr.ErrInvalidGrant
				}
				return nil, err
			}
			allowed, err := validationFn(rti)
			if err != nil {
				return nil, err
			} else if !allowed {
				return nil, xerr.ErrInvalidScope
			}
		}

		ti, err := s.Manager.RefreshAccessToken(ctx, tgr)
		if err != nil {
			if err == xerr.ErrInvalidRefreshToken || err == xerr.ErrExpiredRefreshToken {
				return nil, xerr.ErrInvalidGrant
			}
			return nil, err
		}
		return ti, nil
	}

	return nil, xerr.ErrUnsupportedGrantType
}
func (s *server) validationTokenRequest(r *http.Request) (GrantType, *TokenGenerateRequest, error) {
	if v := r.Method; !(v == "POST" ||
		(s.AllowGetAccessRequest && v == "GET")) {
		return "", nil, xerr.ErrInvalidRequest
	}

	gt := GrantType(r.FormValue("grant_type"))
	if gt.String() == "" {
		return "", nil, xerr.ErrUnsupportedGrantType
	}

	clientID, clientSecret, err := s.ClientInfoHandler(r)
	if err != nil {
		return "", nil, err
	}

	tgr := TokenGenerateRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Request:      r,
	}

	switch gt {
	case AuthorizationCode:
		tgr.RedirectURI = r.FormValue("redirect_uri")
		tgr.Code = r.FormValue("code")
		if tgr.RedirectURI == "" ||
			tgr.Code == "" {
			return "", nil, xerr.ErrInvalidRequest
		}
		tgr.CodeVerifier = r.FormValue("code_verifier")
		if s.forcePKCE && tgr.CodeVerifier == "" {
			return "", nil, xerr.ErrInvalidRequest
		}
	case PasswordCredentials:
		tgr.Scope = r.FormValue("scope")
		username, password := r.FormValue("username"), r.FormValue("password")
		if username == "" || password == "" {
			return "", nil, xerr.ErrInvalidRequest
		}

		userID, err := s.PasswordAuthorizationHandler(r.Context(), clientID, username, password)
		if err != nil {
			return "", nil, err
		} else if userID == "" {
			return "", nil, xerr.ErrInvalidGrant
		}
		tgr.UserID = userID
	case ClientCredentials:
		tgr.Scope = r.FormValue("scope")
	case Refreshing:
		tgr.Refresh, err = s.RefreshTokenResolveHandler(r)
		tgr.Scope = r.FormValue("scope")
		if err != nil {
			return "", nil, err
		}
	}
	return gt, tgr, nil
}

func (s *server) tokenError(w http.ResponseWriter, err error) error {
	data, statusCode, header := s.getErrorData(err)
	return s.token(w, data, header, statusCode)
}

func (s *server) token(w http.ResponseWriter, data map[string]any, header http.Header, statusCode ...int) error {
	if fn := s.ResponseTokenHandler; fn != nil {
		return fn(w, data, header, statusCode...)
	}
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

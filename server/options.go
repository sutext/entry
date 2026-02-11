package server

import (
	"context"
	"net/http"
	"net/netip"

	"sutext.github.io/entry/model"
	"sutext.github.io/entry/xlog"
)

type RefreshTokenResolveHandler func(r *http.Request) (string, error)
type PasswordAuthorizationHandler func(ctx context.Context, clientID, username, password string) (userID string, err error)

type options struct {
	addr                          string
	dirver                        model.Driver
	logger                        *xlog.Logger
	issuerURL                     string
	allHeaders                    http.Header
	realIPHeader                  string
	allowedOrigins                []string
	allowedHeaders                []string
	trustedRealIPCIDRs            []*netip.Prefix
	supportedGrantTypes           map[string]struct{}
	supportedResponseTypes        map[string]struct{}
	supportedCodeChallengeMethods map[string]struct{}
}

func newOptions(opts ...Option) *options {
	os := &options{
		addr:   ":8080",
		dirver: nil,
		logger: xlog.NewText(xlog.LevelInfo),
	}
	for _, o := range opts {
		o.apply(os)
	}
	return os
}

type Option struct {
	apply func(*options)
}

func option(apply func(*options)) Option {
	return Option{apply: apply}
}
func WithAddr(addr string) Option {
	return option(func(o *options) {
		o.addr = addr
	})
}
func WithDriver(d model.Driver) Option {
	return option(func(o *options) {
		o.dirver = d
	})
}
func WithLogger(logger *xlog.Logger) Option {
	return option(func(o *options) {
		o.logger = logger
	})
}
func WithAllHeaders(headers http.Header) Option {
	return option(func(o *options) {
		o.allHeaders = headers
	})
}
func WithIssuerURL(url string) Option {
	return option(func(o *options) {
		o.issuerURL = url
	})
}
func WithRealIPHeader(header string) Option {
	return option(func(o *options) {
		o.realIPHeader = header
	})
}
func WithCORS(origins, headers []string) Option {
	return option(func(o *options) {
		o.allowedOrigins = origins
		o.allowedHeaders = headers
	})
}
func WithTrustedRealIPCIDRs(cidrs []*netip.Prefix) Option {
	return option(func(o *options) {
		o.trustedRealIPCIDRs = cidrs
	})
}
func WithSupportedGrantTypes(grantTypes []string) Option {
	return option(func(o *options) {
		o.supportedGrantTypes = make(map[string]struct{}, len(grantTypes))
		for _, gt := range grantTypes {
			o.supportedGrantTypes[gt] = struct{}{}
		}
	})
}
func WithSupportedResponseTypes(responseTypes []string) Option {
	return option(func(o *options) {
		o.supportedResponseTypes = make(map[string]struct{}, len(responseTypes))
		for _, rt := range responseTypes {
			o.supportedResponseTypes[rt] = struct{}{}
		}
	})
}
func WithSupportedCodeChallengeMethods(methods []CodeChallengeMethod) Option {
	return option(func(o *options) {
		o.supportedCodeChallengeMethods = make(map[string]struct{}, len(methods))
		for _, m := range methods {
			o.supportedCodeChallengeMethods[m.String()] = struct{}{}
		}
	})
}

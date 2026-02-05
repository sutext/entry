package server

import (
	"net/http"
	"net/netip"

	"sutext.github.io/entry/model"
	"sutext.github.io/entry/xlog"
)

type options struct {
	addr                   string
	dirver                 model.Driver
	logger                 *xlog.Logger
	issuerURL              string
	allHeaders             http.Header
	realIPHeader           string
	allowedOrigins         []string
	allowedHeaders         []string
	trustedRealIPCIDRs     []*netip.Prefix
	supportedGrantTypes    []string
	supportedResponseTypes map[string]bool
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
		o.supportedGrantTypes = grantTypes
	})
}
func WithSupportedResponseTypes(responseTypes map[string]bool) Option {
	return option(func(o *options) {
		o.supportedResponseTypes = responseTypes
	})
}

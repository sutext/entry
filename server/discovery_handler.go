package server

import (
	"encoding/json"
	"net/http"
	"sort"
)

type discovery struct {
	Issuer             string `json:"issuer"`
	JwksURI            string `json:"jwks_uri"`
	AuthEndpoint       string `json:"authorization_endpoint"`
	TokenEndpoint      string `json:"token_endpoint"`
	UserInfoEndpoint   string `json:"userinfo_endpoint"`
	DeviceEndpoint     string `json:"device_authorization_endpoint"`
	IntrospectEndpoint string `json:"introspection_endpoint"`
	// RevocationEndpoint string   `json:"revocation_endpoint"`
	GrantTypes        []string `json:"grant_types_supported"`
	ResponseTypes     []string `json:"response_types_supported"`
	Subjects          []string `json:"subject_types_supported"`
	IDTokenAlgs       []string `json:"id_token_signing_alg_values_supported"`
	CodeChallengeAlgs []string `json:"code_challenge_methods_supported"`
	Scopes            []string `json:"scopes_supported"`
	AuthMethods       []string `json:"token_endpoint_auth_methods_supported"`
	Claims            []string `json:"claims_supported"`
}

func (s *server) handleDiscovery(w http.ResponseWriter, r *http.Request) {
	d := s.constructDiscovery()
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		http.Error(w, "failed to marshal discovery data: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
func (s *server) constructDiscovery() discovery {
	d := discovery{
		Issuer:             s.issuerURL.String(),
		AuthEndpoint:       s.absURL("/auth"),
		TokenEndpoint:      s.absURL("/token"),
		JwksURI:            s.absURL("/keys"),
		UserInfoEndpoint:   s.absURL("/userinfo"),
		DeviceEndpoint:     s.absURL("/device/code"),
		IntrospectEndpoint: s.absURL("/token/introspect"),
		Subjects:           []string{"public"},
		IDTokenAlgs:        []string{"RS256"},
		CodeChallengeAlgs:  []string{"plain", "S256"},
		Scopes:             []string{"openid", "email", "profile"},
		AuthMethods:        []string{"client_secret_basic", "client_secret_post"},
		Claims: []string{
			"iss", "sub", "aud", "iat", "exp", "email", "phone",
		},
	}

	for responseType := range s.supportedResponseTypes {
		d.ResponseTypes = append(d.ResponseTypes, responseType)
	}
	sort.Strings(d.ResponseTypes)

	d.GrantTypes = s.supportedGrantTypes
	return d
}

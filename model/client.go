package model

type Client struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Status       uint8   `json:"status"`
	Secret       string  `json:"secret"`
	Scopes       Strings `json:"scopes"`
	Public       bool    `json:"public,omitempty"`
	LogoURL      string  `json:"logo_url,omitempty"`
	Description  string  `json:"description"`
	RedirectURIs Strings `json:"redirect_uris"`
	TrustedPeers Strings `json:"trusted_peers,omitempty"` // TrustedPeers holds the value of the "trusted_peers" field.
}

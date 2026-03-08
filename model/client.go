package model

import (
	"sutext.github.io/suid/guid"
)

type ClientType string

const (
	ClientTypeOfficial     ClientType = "official"
	ClientTypePublic       ClientType = "public"
	ClientTypeConfidential ClientType = "confidential"
)

type ClientStatus uint8

const (
	ClientStatusNormal  ClientStatus = 0
	ClientStatusBanned  ClientStatus = 1
	ClientStatusDeleted ClientStatus = 2
)

type Client struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         ClientType   `json:"type"`
	Status       ClientStatus `json:"status"`
	Secret       string       `json:"secret"`
	Scopes       Strings      `json:"scopes"`
	Public       bool         `json:"public,omitempty"`
	LogoURL      string       `json:"logo_url,omitempty"`
	Description  string       `json:"description,omitempty"`
	RedirectURIs Strings      `json:"redirect_uris,omitempty"`
	TrustedPeers Strings      `json:"trusted_peers,omitempty"`
}

func NewClient() *Client {
	return &Client{
		ID:     guid.New().String(),
		Type:   ClientTypePublic,
		Status: ClientStatusNormal,
		Secret: guid.New().String(),
	}
}
